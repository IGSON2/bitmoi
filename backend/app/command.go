package app

import (
	"bitmoi/backend/futureclient"
	"bitmoi/backend/utilities"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	StoreCommand = &cli.Command{
		Action: GetCandleData,
		Name:   "store",
		Usage:  "Store candles data form binance",
		Flags:  []cli.Flag{IntervalFlag, TimestampFlag, GetAllFlag, BackwardFlag, PairListFlag},
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
		if pairsFlag := ctx.String("pairs"); pairsFlag == "" {
			return fmt.Errorf("require at least one pair")
		} else {
			for _, n := range utilities.SplitAndTrim(pairsFlag) {
				names = append(names, n+"USDT")
			}
		}
	}

	var intervals []string
	if intervalsFlag := ctx.String("interval"); intervalsFlag == "" {
		return fmt.Errorf("require at least one interval")
	} else {
		intervals = utilities.SplitAndTrim(intervalsFlag)
	}

	var cnt int
	for _, name := range names {
		for _, interval := range intervals {
			err = f.StoreCandles(interval, name, ctx.Int64("timestamp"), ctx.Bool("backward"), &cnt)
			if err != nil {
				return fmt.Errorf("cannot store candles, err : %w", err)
			}
		}
	}
	log.Info().Msg("All pairs are stored completely")
	return nil
}
