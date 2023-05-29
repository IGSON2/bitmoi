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

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"

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
	Logger    *zerolog.Logger
}

func NewFutureClient(c *utilities.Config) (*FutureClient, error) {
	dbConn, err := sql.Open(c.DBDriver, c.DBSource)

	if err != nil {
		return nil, fmt.Errorf("cannot open db store %w", err)
	}
	clientlogger := zerolog.New(os.Stdout)
	zerolog.TimeFieldFormat = zerolog.TimestampFunc().Format("2006-01-02 15:04:05")

	f := &FutureClient{
		Client:    futures.NewClient(utilities.GetAPIKeys()),
		Store:     db.NewStore(dbConn),
		Yesterday: utilities.Yesterday9AM(),
		Logger:    &clientlogger,
	}
	if getErr := f.GetAllPairs(); getErr != nil {
		return nil, getErr
	}
	return f, nil
}

func (f *FutureClient) GetAllPairs() error {
	f.Logger.Info().Msg("start to store all pair names")
	info, err := f.Client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get allpairs %w", err)
	}
	for _, s := range info.Symbols {
		if strings.HasSuffix(s.Symbol, "USDT") {
			f.Pairs = append(f.Pairs, s.Symbol)
		}
	}
	f.Logger.Info().Msg("all pair names ars stored completely")
	return nil
}

// StoreCandles retrieves the candle data from the binance and stores it in db.
// starttime -> endtime (current)
func (f *FutureClient) StoreCandles(interval, name string, timestamp int64, cnt *int) error {
	f.Logger.Info().Any("pair", name).Msg("Start to store")
	endTime := f.Yesterday
	var startTime int64
	if timestamp <= 0 {
		startTime = howcandles(f.Yesterday, interval, LimitCandlesNum)
		f.Logger.Info().Any("start time", startTime).Msg(fmt.Sprintf("Timestamp dosen't specified. Set the start time %d candles before.", LimitCandlesNum))
	} else {
		startTime = timestamp
		f.Logger.Info().Any("start time", startTime).Msg("EndTime set by received timestamp.")
	}

	info := initInsertInfo(interval)

	for startTime <= endTime {
		f.Logger.Info().Any("Start", utilities.TransMilli(startTime)).Any("End", utilities.TransMilli(endTime)).Msg("get klines")

		klines, err := f.Client.NewKlinesService().Symbol(name).StartTime(startTime).EndTime(endTime).Limit(LimitCandlesNum).
			Interval(interval).Do(context.Background())
		if err != nil {
			return fmt.Errorf("cannot get kilnes, err :%w", err)
		}

		for _, k := range klines {
			var color string

			if utilities.StrToFloat(k.Close) >= utilities.StrToFloat(k.Open) {
				color = "rgba(38,166,154,0.5)"
			} else {
				color = "rgba(239,83,80,0.5)"
			}

			// TODO : 자릿수 변환으로 인한 연속성 오류 해결 필요
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
				_, err = f.Store.Insert1dCandles(context.Background(), *param)
			case db.FourH:
				param, _ := info.(*db.Insert4hCandlesParams)
				_, err = f.Store.Insert4hCandles(context.Background(), *param)
			case db.OneH:
				param, _ := info.(*db.Insert1hCandlesParams)
				_, err = f.Store.Insert1hCandles(context.Background(), *param)
			case db.FifM:
				param, _ := info.(*db.Insert15mCandlesParams)
				_, err = f.Store.Insert15mCandles(context.Background(), *param)
			}

			if err != nil {
				return fmt.Errorf("cannot insert candle into db err : %w", err)
			}
		}
		endTime = howcandles(endTime, interval, LimitCandlesNum)

		*cnt++
		if *cnt%900 == 0 {
			f.Logger.Info().Any("count", *cnt).Msg("Every 900th request takes a minute off")
			time.Sleep(1 * time.Minute)
		}
	}
	f.Logger.Info().Any("pair", name).Msg("stored complete.")
	return nil
}

// 주어진 시간으로부터 intN + intU 단위의 캔들을 candles개 만큼 가져올 수 있는 일자를 Millisecond로 반환합니다.
func howcandles(root int64, interval string, candles int) int64 {
	intN, intU := db.ParseInterval(interval)

	var start int64

	switch intU {
	case "m":
		start = root - int64(time.Minute.Milliseconds()*int64(intN)*int64(candles))
	case "h":
		start = root - int64(time.Hour.Milliseconds()*int64(intN)*int64(candles))
	case "d":
		start = root - int64(time.Hour.Milliseconds()*int64(intN)*24*int64(candles))
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
