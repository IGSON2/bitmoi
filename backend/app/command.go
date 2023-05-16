package app

import (
	db "bitmoi/backend/db"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

var (
	GetCommand = &cli.Command{
		Action:    getCandleData,
		Name:      "get",
		Usage:     "Get candles data form binance",
		ArgsUsage: "<Set_interval , Set_target_date_of_timestamp>",
		Flags:     []cli.Flag{IntervalFlag, TimestampFlag},
	}
)

func getCandleData(ctx *cli.Context) error {
	var timestamp int
	var err error

	if timestampArg := ctx.Args().Get(2); timestampArg != "" {
		timestamp, err = strconv.Atoi(timestampArg)
		if err != nil {
			return fmt.Errorf("cannot parse string timestamp flag : %w", err)
		}
	}

	return db.SaveCandles(ctx.Args().First(), int64(timestamp))
}
