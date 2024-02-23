package futureclient

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"

	"github.com/adshao/go-binance/v2/futures"
)

const (
	LimitCandlesNum = 1000
)

type FutureClient struct {
	Client    *futures.Client
	Store     db.Store
	Pairs     []string
	Yesterday int64
}

func NewFutureClient(c *utilities.Config) (*FutureClient, error) {
	dbConn, err := sql.Open(c.DBDriver, c.DBSource)

	if err != nil {
		return nil, fmt.Errorf("cannot open db store %w", err)
	}

	f := &FutureClient{
		Client:    futures.NewClient("", ""),
		Store:     db.NewStore(dbConn),
		Yesterday: utilities.Yesterday9AM(),
	}
	if getErr := f.GetAllPairsFromBinance(); getErr != nil {
		return nil, getErr
	}
	return f, nil
}

// GetAllPairsFromBinance is storing pairnames from binance server but "deprecated" now.
func (f *FutureClient) GetAllPairsFromBinance() error {
	log.Info().Msg("start to get all pair names from binance server")
	info, err := f.Client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get allpairs %w", err)
	}
	for _, s := range info.Symbols {
		if strings.HasSuffix(s.Symbol, "USDT") {
			f.Pairs = append(f.Pairs, s.Symbol)
		}
	}
	log.Info().Msgf("init all pair names completely, total %d pairs.\n %v", len(f.Pairs), f.Pairs)
	return nil
}

func (f *FutureClient) GetAllPairsFromStore() error {
	log.Info().Msg("start to get all pair names from database store")
	var err error
	f.Pairs, err = f.Store.GetAllParisInDB(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get allpairs %w", err)
	}
	log.Info().Msg("init all pair names completely")
	return nil
}

// StoreCandles retrieves the candle data from the binance and stores it in db.
func (f *FutureClient) StoreCandles(interval, name string, timestamp int64, backward bool, cnt *int) error {
	log.Info().Any("pair", name).Any("interval", interval).Msg("Start to store")
	c := context.Background()
	min, max, err := f.Store.SelectMinMaxTime(interval, name, c)
	min, max = (min-32400)*1000, (max-32400)*1000
	if err != nil {
		log.Err(err).Msgf("cannot get min max timestamp name:%s interval:%s", name, interval)
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
				log.Info().Msgf("there's no candles. start time : %s, end time :%s", utilities.TransMilli(startTime), utilities.TransMilli(endTime))
			} else {
				log.Info().Any("Given", utilities.TransMilli(timestamp)).Any("Min", utilities.TransMilli(min)).Msg("given timestamp is more futuristic than minimum timestamp.")
				return nil
			}
		} else if min == timestamp {
			startTime = max + 1
			log.Info().Any("Given", utilities.TransMilli(timestamp)).Any("Min", utilities.TransMilli(min)).Any("Max", utilities.TransMilli(max)).
				Msg("given timestamp is equal with minimum timestamp. set to maximum timestamp + 1")
		} else {
			startTime = timestamp
			endTime = min - 1
			log.Info().Any("Given", utilities.TransMilli(timestamp)).Any("Min", utilities.TransMilli(min)).Any("Max", utilities.TransMilli(max)).
				Msgf("start time has been set to given timestamp start:%s - end:%s", utilities.TransMilli(startTime), utilities.TransMilli(endTime))
		}
	} else {
		startTime = getStartTime(max, interval, 1)
		if timestamp <= max {
			log.Info().Any("Given", utilities.TransMilli(timestamp)).Any("Max", utilities.TransMilli(max)).Msg("given timestamp is equal or more past than maximum timestamp.")
			return nil
		} else {
			if max <= 0 {
				startTime = getStartTime(f.Yesterday, interval, -1*LimitCandlesNum)
				endTime = f.Yesterday
				log.Info().Msgf("there's no candles. start time : %s, end time :%s", utilities.TransMilli(startTime), utilities.TransMilli(endTime))
			} else {
				if startTime > f.Yesterday {
					log.Info().Any("Start time", utilities.TransMilli(startTime)).Any("Yesterday", utilities.TransMilli(f.Yesterday)).
						Msg("start time is more futureistic than yesterday 9am")
					return nil
				}
				endTime = f.Yesterday
				log.Info().Any("Given", utilities.TransMilli(timestamp)).Any("Min", utilities.TransMilli(min)).Any("Max", utilities.TransMilli(max)).
					Msgf("start time has been set to given timestamp start:%s - end:%s", utilities.TransMilli(startTime), utilities.TransMilli(endTime))
			}
		}
	}

	info := initInsertInfo(interval)

	for startTime <= endTime {
		log.Info().Any("Start", utilities.TransMilli(startTime)).Any("End", utilities.TransMilli(endTime)).Msg("get klines")

		klines, err := f.Client.NewKlinesService().Symbol(name).StartTime(startTime).EndTime(endTime).Limit(LimitCandlesNum).
			Interval(interval).Do(c)
		if err != nil {
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
				log.Err(err)
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
			log.Info().Any("count", *cnt).Msg("Every 900th request takes a minute off")
			time.Sleep(1 * time.Minute)
		}
	}
	log.Info().Any("pair", name).Msg("stored complete.")
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
