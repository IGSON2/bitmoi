package main

import (
	"bitmoi/backend/api"
	"bitmoi/backend/app"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"database/sql"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var bApp = app.NewApp()
var applog = zerolog.New(os.Stdout)

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
	zerolog.TimeFieldFormat = zerolog.TimestampFunc().Format("2006-01-02 15:04:05")
	applog.Info().Msg("Start bitmoi")
	config := utilities.GetConfig("./")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return fmt.Errorf("cannot connect db %w", err)
	}
	dbStore := db.NewStore(conn)
	server, err := api.NewServer(config, dbStore)
	if err != nil {
		return fmt.Errorf("cannot create server %w", err)
	}

	return server.Listen()
}
