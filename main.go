package main

import (
	"bitmoi/backend/api"
	"bitmoi/backend/app"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi"
	"bitmoi/backend/mail"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	"database/sql"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var bApp = app.NewApp()

func init() {
	bApp.Commands = []*cli.Command{
		app.StoreCommand,
	}
	bApp.Flags = []cli.Flag{
		app.DatadirFlag,
		app.GRPCFlag,
		app.HTTPFlag,
	}
	bApp.Action = bitmoi
}

func main() {
	if err := bApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bitmoi(ctx *cli.Context) error {
	config := utilities.GetConfig("./")

	if path := ctx.String(app.DatadirFlag.Name); path != "" {
		config.SetDataDir(path)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return fmt.Errorf("cannot connect db %w", err)
	}
	dbStore := db.NewStore(conn)

	errCh := make(chan error)

	if isGrpcRun := ctx.Bool(app.GRPCFlag.Name); isGrpcRun {
		server, err := gapi.NewServer(config, dbStore)
		if err != nil {
			log.Panic().Err(err).Msg("cannot create gRPC server")
		}
		go server.ListenGRPC(errCh)
		go server.ListenGRPCGateWay(errCh)
	}

	if isHTTPRun := ctx.Bool(app.HTTPFlag.Name); isHTTPRun {
		go runTaskProcessor(config, dbStore)
		go runHttpServer(config, dbStore, errCh)
	}

	return <-errCh
}

func runHttpServer(config *utilities.Config, store db.Store, errCh chan error) {
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create server")
	}

	errCh <- server.Listen()
}

func runTaskProcessor(config *utilities.Config, store db.Store) {
	gmailSender := mail.NewGmailSender(config)
	taskProcessor := worker.NewRedisTaskProcessor(asynq.RedisClientOpt{Addr: config.RedisAddress}, store, gmailSender, config.AccessTokenDuration)
	log.Info().Msg("start task processor")
	if err := taskProcessor.Start(); err != nil {
		log.Panic().Err(err).Msg("failed to start task processor")
	}
}
