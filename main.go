package main

import (
	"bitmoi/backend/app"
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
	applog.Info().Any("port", ctx.Args().First()).Msg("Start bitmoi")
	return nil
}
