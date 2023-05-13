package app

import (
	"github.com/urfave/cli/v2"
)

const (
	version = "v1.0.0"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Version = version
	app.Copyright = "Copyright 2023 IGSON"
	return app
}
