package app

import (
	db "bitmoi/backend/db/chartData"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	APIPortFlag = &cli.IntFlag{
		Name:  "port",
		Usage: "Set API port number",
		Value: 4000,
	}
	IntervalFlag = &cli.StringFlag{
		Name:  "interval",
		Usage: "Set interval to get",
		Value: db.FourH,
	}
	TimestampFlag = &cli.Int64Flag{
		Name:  "timestamp",
		Usage: "Set timestamp to retrieve data by the specified date.",
		Value: time.Now().Add(-1 * 24 * 7 * time.Hour).UnixMilli(),
	}
)
