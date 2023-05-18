package app

import (
	db "bitmoi/backend/db/sqlc"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	IntervalFlag = &cli.StringFlag{
		Name:  "interval",
		Usage: "Set interval to get",
		Value: db.FourH,
	}
	TimestampFlag = &cli.Int64Flag{
		Name:  "timestamp",
		Usage: "Set the unixmilli timestamp to retrieve data by the specified date. 15m 1h 4h 1d",
		Value: time.Now().Add(-1 * 24 * 7 * time.Hour).UnixMilli(),
	}
	GetAllFlag = &cli.BoolFlag{
		Name:  "all",
		Usage: "If it's true, get all pairs",
		Value: false,
	}
	PairListFlage = &cli.StringFlag{
		Name:  "pairs",
		Usage: "Specify pairs to get, type symbal and seperate by space e.b : BTC ETH",
		Value: "",
	}
)
