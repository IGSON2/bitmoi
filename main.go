package main

import (
	"bitmoi/backend/app"
	"fmt"
	"os"

	log "github.com/inconshreveable/log15"
	"github.com/urfave/cli/v2"
)

var (
	flags = []cli.Flag{
		app.APIPortFlag,
	}
)

var bApp = app.NewApp()
var applog = log.New("module", "app")

func init() {
	bApp.Commands = []*cli.Command{
		app.GetCommand,
	}
	bApp.Action = bitmoi
	bApp.Flags = flags
}

func main() {
	if err := bApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bitmoi(ctx *cli.Context) error {
	applog.Info("Start bitmoi", "port", ctx.Args().First())
	return nil
}
