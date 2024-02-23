package app

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	IntervalFlag = &cli.StringFlag{
		Name:  "interval",
		Usage: "Set intervals to get. e.b : 5m,15m,1h,4h,1d",
		Value: db.FourH,
	}
	TimestampFlag = &cli.Int64Flag{
		Name:  "timestamp",
		Usage: "Set the unixmilli timestamp to retrieve data by the specified date.",
		Value: time.Now().Add(-1 * 24 * 7 * time.Hour).UnixMilli(),
	}
	GetAllFlag = &cli.BoolFlag{
		Name:  "all",
		Usage: "If it's true, get all pairs",
		Value: false,
	}
	FromBinanceFlag = &cli.BoolFlag{
		Name:  "binance",
		Usage: "If it's true, init future client's pair list from binance. if it's false, init from db.",
		Value: false,
	}
	BackwardFlag = &cli.BoolFlag{
		Name:  "backward",
		Usage: "If it's true, store candles before minimum timestamp otherwise, store candles after maximum timestamp",
		Value: true,
	}
	PairListFlag = &cli.StringFlag{
		Name:  "pairs",
		Usage: "Specify pairs to get, type symbal and seperate by comma e.b : BTC,ETH",
		Value: "",
	}
)
var (
	DatadirFlag = &cli.StringFlag{
		Name:  "datadir",
		Usage: "Specify pairs to datadir path.",
		Value: utilities.DefaultDataDir(),
	}
	GRPCFlag = &cli.BoolFlag{
		Name:  "grpc",
		Usage: "If it's true, run grpc and gateway server.",
		Value: false,
	}
	HTTPFlag = &cli.BoolFlag{
		Name:  "http",
		Usage: "If it's true, run http server.",
		Value: false,
	}
	LogLevelFlag = &cli.IntFlag{
		Name:  "loglevel",
		Usage: "Set log level. 0 : Debug, 1 : Info, 2 : Warn, 3 : Error, 4 : Fatal, 5 : Panic",
		Value: 0,
	}
)
