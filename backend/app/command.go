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
	timeStamp, err := strconv.Atoi(ctx.Args().Get(2))
	if err != nil {
		return fmt.Errorf("cannot parse string timestamp flag : %w", err)
	}
	return db.SaveCandles(ctx.Args().First(), int64(timeStamp))
}
