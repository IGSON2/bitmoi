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
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return fmt.Errorf("cannot connect db %w", err)
	}
	dbStore := db.NewStore(conn)
	go runGateWayServer(config, dbStore)
	go runGrpcServer(config, dbStore)
	runHttpServer(config, dbStore)
	return nil
}

func runGrpcServer(config *utilities.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create gRPC server")
	}
	go server.ListenGRPC()
}

func runHttpServer(config *utilities.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create server")
	}

	log.Panic().Err(server.Listen()).Msg("cannot start http server")
}

func runGateWayServer(config *utilities.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create gateway server")
	}
	go server.ListenGRPCGateWay()
}
