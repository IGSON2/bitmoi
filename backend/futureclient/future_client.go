package futureclient

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"bitmoi/backend/config"

	"github.com/adshao/go-binance/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	// "github.com/adshao/go-binance/v2/futures"
)

const (
	LimitCandlesNum = 1000
)

type FutureClient struct {
	// Client    *futures.Client
	Client    *binance.Client
	Store     db.Store
	Pairs     []string
	Yesterday int64
	Logger    *zerolog.Logger
}

func NewFutureClient(c *config.Config) (*FutureClient, error) {
	dbConn, err := sql.Open(c.DBDriver, c.DBSource)

	if err != nil {
		return nil, fmt.Errorf("cannot open db store %w", err)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger().Level(c.LogLevel)
	f := &FutureClient{
		// Client:    futures.NewClient("", ""),
		Client:    binance.NewClient("", ""),
		Store:     db.NewStore(dbConn),
		Yesterday: utilities.Yesterday9AM(),
		Logger:    &logger,
	}
	return f, nil
}

func (f *FutureClient) InitPairs(isBinance bool) error {
	if isBinance {
		return f.GetAllPairsFromBinance()
	}
	return f.GetAllPairsFromStore()
}

// GetAllPairsFromBinance is storing pairnames from binance server but "deprecated" now.
func (f *FutureClient) GetAllPairsFromBinance() error {
	f.Logger.Info().Msg("start to get all pair names from binance server")
	info, err := f.Client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get allpairs %w", err)
	}
	for _, s := range info.Symbols {
		if strings.HasSuffix(s.Symbol, "USDT") {
			f.Pairs = append(f.Pairs, s.Symbol)
		}
	}
	f.Logger.Info().Msgf("init all pair names completely, total %d pairs.", len(f.Pairs))
	return nil
}

func (f *FutureClient) GetAllPairsFromStore() error {
	f.Logger.Info().Msg("start to get all pair names from database store")
	var err error
	f.Pairs, err = f.Store.GetAllPairsInDB1D(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get allpairs %w", err)
	}
	f.Logger.Info().Msgf("init all pair names completely. total %d pairs.", len(f.Pairs))
	return nil
}

// StoreCandles retrieves the candle data from the binance and stores it in db.
func (f *FutureClient) StoreCandles(interval, name string, timestamp int64, backward bool, cnt *int) error {
	f.Logger.Info().Any("pair", name).Any("interval", interval).Msg("Start to store")
	c := context.Background()
	min, max, err := f.Store.SelectMinMaxTime(interval, name, c)
	min, max = (min-32400)*1000, (max-32400)*1000
	if err != nil {
		f.Logger.Error().Err(err).Msgf("cannot get min max timestamp name:%s interval:%s", name, interval)
		return err
	}

	var (
		startTime int64
		endTime   int64
	)

	if backward {
		endTime = f.Yesterday
		if min < timestamp {
			if min <= 0 {
				startTime = timestamp
				f.Logger.Info().Msgf("there's no candles. start time : %s, end time :%s", utilities.TransMilli(startTime), utilities.TransMilli(endTime))
			} else {
				f.Logger.Info().Str("Given", utilities.TransMilli(timestamp)).Str("Min", utilities.TransMilli(min)).Msg("given timestamp is more futuristic than minimum timestamp.")
				return nil
			}
		} else if min == timestamp {
			startTime = max + 1
			f.Logger.Info().Str("Given", utilities.TransMilli(timestamp)).Str("Max", utilities.TransMilli(max)).Str("Min", utilities.TransMilli(min)).
				Msg("given timestamp is equal with minimum timestamp. set to maximum timestamp + 1")
		} else {
			startTime = timestamp
			endTime = min - 1
			f.Logger.Info().Str("Given", utilities.TransMilli(timestamp)).Str("Max", utilities.TransMilli(max)).Str("Min", utilities.TransMilli(min)).
				Msgf("start time has been set to given timestamp start:%s - end:%s", utilities.TransMilli(startTime), utilities.TransMilli(endTime))
		}
	} else {
		startTime = getStartTime(max, interval, 1)
		endTime = f.Yesterday
		f.Logger.Info().Str("Given", utilities.TransMilli(timestamp)).Str("Max", utilities.TransMilli(max)).Str("Min", utilities.TransMilli(min)).Str("EndTime", utilities.TransMilli(endTime)).Str("StartTime", utilities.TransMilli(startTime)).Msg("Get candles from max to yesterday")
		if startTime > f.Yesterday {
			f.Logger.Info().Str("Start time", utilities.TransMilli(startTime)).Str("Yesterday", utilities.TransMilli(f.Yesterday)).
				Msg("start time is more futureistic than yesterday 9am")
			return nil
		}
	}

	info := initInsertInfo(interval)

	for startTime <= endTime {
		f.Logger.Info().Str("End", utilities.TransMilli(endTime)).Str("Start", utilities.TransMilli(startTime)).Msg("get klines")

		klines, err := f.Client.NewKlinesService().Symbol(name).StartTime(startTime).EndTime(endTime).Limit(LimitCandlesNum).
			Interval(interval).Do(c)
		if err != nil {
			if strings.Contains(err.Error(), "connection reset by peer") {
				time.Sleep(1 * time.Minute)
				continue
			}
			return fmt.Errorf("cannot get kilnes, err :%w", err)
		}
		if len(klines) == 0 {
			startTime = getStartTime(startTime, interval, LimitCandlesNum)
			continue
		}

		for _, k := range klines {
			var color string

			if utilities.StrToFloat(k.Close) >= utilities.StrToFloat(k.Open) {
				color = "rgba(38,166,154,0.5)"
			} else {
				color = "rgba(239,83,80,0.5)"
			}

			info.SetCandleInfo(
				utilities.StrToFloat(k.Open),
				utilities.StrToFloat(k.Close),
				utilities.StrToFloat(k.High),
				utilities.StrToFloat(k.Low),
				utilities.StrToFloat(k.Volume),
				(k.OpenTime/1000)+32400,
				name,
				color,
			)

			switch interval {
			case db.OneD:
				param, _ := info.(*db.Insert1dCandlesParams)
				_, err = f.Store.Insert1dCandles(c, *param)
			case db.FourH:
				param, _ := info.(*db.Insert4hCandlesParams)
				_, err = f.Store.Insert4hCandles(c, *param)
			case db.OneH:
				param, _ := info.(*db.Insert1hCandlesParams)
				_, err = f.Store.Insert1hCandles(c, *param)
			case db.FifM:
				param, _ := info.(*db.Insert15mCandlesParams)
				_, err = f.Store.Insert15mCandles(c, *param)
			case db.FiveM:
				param, _ := info.(*db.Insert5mCandlesParams)
				_, err = f.Store.Insert5mCandles(c, *param)
			default:
				err := fmt.Errorf("unsupported interval: %s", interval)
				f.Logger.Error().Err(err)
				return err
			}

			if err != nil {
				if strings.Contains(err.Error(), "Duplicate entry") {
					continue
				}
				return fmt.Errorf("cannot insert candle into db err : %w", err)
			}
		}

		startTime = klines[len(klines)-1].CloseTime + 1

		*cnt++
		if *cnt%900 == 0 {
			f.Logger.Info().Int("count", *cnt).Msg("Every 900th request takes a minute off")
			time.Sleep(1 * time.Minute)
		}
	}
	f.Logger.Info().Str("pair", name).Msg("stored complete.")
	return nil
}

// getStartTime은 주어진 interval에 맞는 다음 요청에 필요한 starttime을 구합니다.
func getStartTime(root int64, interval string, candles int) int64 {
	intN, intU := db.ParseInterval(interval)

	var start int64

	switch intU {
	case "m":
		start = root + int64(time.Minute.Milliseconds()*int64(intN)*int64(candles))
	case "h":
		start = root + int64(time.Hour.Milliseconds()*int64(intN)*int64(candles))
	case "d":
		start = root + int64(time.Hour.Milliseconds()*int64(intN)*24*int64(candles))
	}

	return start
}

func initInsertInfo(interval string) db.InsertQueryInterface {
	switch interval {
	case db.OneD:
		return &db.Insert1dCandlesParams{}
	case db.FourH:
		return &db.Insert4hCandlesParams{}
	case db.OneH:
		return &db.Insert1hCandlesParams{}
	case db.FifM:
		return &db.Insert15mCandlesParams{}
	case db.FiveM:
		return &db.Insert5mCandlesParams{}
	}
	return nil
}

func (f *FutureClient) PruneCandles(term string) error {
	pairs, err := f.Store.GetUnder1YPairs(context.Background(), term)
	if err != nil {
		return fmt.Errorf("cannot get pairs, err : %w", err)
	}

	f.Logger.Info().Msgf("start to prune candles, total %d pairs", len(pairs))
	for _, p := range pairs {
		_, err = f.Store.DeletePairs1d(context.Background(), p)
		if err != nil {
			return fmt.Errorf("cannot prune candles, err : %w", err)
		}

		_, err = f.Store.DeletePairs4h(context.Background(), p)
		if err != nil {
			return fmt.Errorf("cannot prune candles, err : %w", err)
		}

		_, err = f.Store.DeletePairs1h(context.Background(), p)
		if err != nil {
			return fmt.Errorf("cannot prune candles, err : %w", err)
		}

		_, err = f.Store.DeletePairs15m(context.Background(), p)
		if err != nil {
			return fmt.Errorf("cannot prune candles, err : %w", err)
		}

		_, err = f.Store.DeletePairs5m(context.Background(), p)
		if err != nil {
			return fmt.Errorf("cannot prune candles, err : %w", err)
		}
		f.Logger.Info().Msgf("prune candles complete, pair : %s", p)
	}

	return nil
}
