package app

import (
	"bitmoi/backend/futureclient"
	"bitmoi/backend/utilities"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	StoreCommand = &cli.Command{
		Action: GetCandleData,
		Name:   "store",
		Usage:  "Store candles data form binance",
		Flags:  []cli.Flag{IntervalFlag, TimestampFlag, GetAllFlag, PairListFlage},
	}
)

func GetCandleData(ctx *cli.Context) error {
	var names []string
	f, err := futureclient.NewFutureClient(utilities.GetConfig("./"))
	if err != nil {
		return fmt.Errorf("cannot create future client, err : %w", err)
	}

	if ctx.Bool("all") {
		names = f.Pairs
	} else {
		if pairsflag := ctx.String("pairs"); pairsflag == "" {
			return fmt.Errorf("require at least one pair")
		} else {
			for _, n := range strings.Split(pairsflag, ",") {
				names = append(names, n+"USDT")
			}
		}
	}

	var cnt int
	for _, name := range names {
		err = f.StoreCandles(ctx.String("interval"), name, ctx.Int64("timestamp"), &cnt)
		if err != nil {
			return fmt.Errorf("cannot store candles, err : %w", err)
		}
	}
	f.Logger.Info().Msg("All pairs are stored completely")
	return nil
}
