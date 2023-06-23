package main

import (
	"bitmoi/backend/api"
	"bitmoi/backend/app"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi"
	"bitmoi/backend/utilities"
	"database/sql"
	"fmt"
	"os"

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
		go runHttpServer(config, dbStore, errCh)
	}

	return <-errCh
}

func runHttpServer(config *utilities.Config, store db.Store, errCh chan error) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create server")
	}

	errCh <- server.Listen()
}
