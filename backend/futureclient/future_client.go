package futureclient

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"

	"github.com/adshao/go-binance/v2/futures"
)

const (
	LimitCandlesNum = 1000
)

type FutureClient struct {
	Client        *futures.Client
	Store         db.Store
	Pairs         []string
	Yesterday     int64
	CandlesCh     chan db.InsertQueryInterface
	LimitCalndles int
	Logger        *zerolog.Logger
}

func NewFutureClient(c *utilities.Config) (*FutureClient, error) {
	dbConn, err := sql.Open(c.DBDriver, c.DBSource)

	if err != nil {
		return nil, fmt.Errorf("cannot open db store %w", err)
	}
	clientlogger := zerolog.New(os.Stdout)
	zerolog.TimeFieldFormat = zerolog.TimestampFunc().Format("2006-01-02 15:04:05")

	f := &FutureClient{
		Client:        futures.NewClient(utilities.GetAPIKeys()),
		Store:         db.NewStore(dbConn),
		Yesterday:     utilities.Yesterday9AM(),
		LimitCalndles: LimitCandlesNum,
		Logger:        &clientlogger,
	}
	if getErr := f.GetAllPairs(); getErr != nil {
		return nil, getErr
	}
	f.CandlesCh = make(chan db.InsertQueryInterface, len(f.Pairs))
	return f, nil
}

func (f *FutureClient) GetAllPairs() error {
	info, err := f.Client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get allpairs %w", err)
	}
	for _, s := range info.Symbols {
		f.Pairs = append(f.Pairs, s.Symbol)
	}
	return nil
}

func (f *FutureClient) StoreCandles(interval, name string, timestamp int64) error {
	startTimemilli := f.howcandles(interval, 0)
	var endTimeMilli int64
	if timestamp <= 0 {
		endTimeMilli = f.howcandles(interval, 1000)
		f.Logger.Info().Any("endtime", endTimeMilli).Msg("Timestamp dosen't specified. Set the endtime 1000 candles before.")
	} else {
		endTimeMilli = timestamp
		f.Logger.Info().Any("endtime", endTimeMilli).Msg("EndTime set by received timestamp.")
	}

	info := *new(db.InsertQueryInterface) // TODO: Nil

	klines, err := f.Client.NewKlinesService().Symbol(name).StartTime(startTimemilli).EndTime(endTimeMilli).Limit(1000).
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
	return nil
}

// 현재로부터 intN + intU 단위의 캔들을 candles개 만큼 가져올 수 있는 일자를 Millisecond로 반환합니다.
func (f *FutureClient) howcandles(interval string, sub int) int64 {
	intN, intU := db.ParseInterval(interval)

	var start int64

	switch intU {
	case "m":
		start = f.Yesterday - int64(time.Minute.Milliseconds()*int64(intN)*int64(f.LimitCalndles-sub))
	case "h":
		start = f.Yesterday - int64(time.Hour.Milliseconds()*int64(intN)*int64(f.LimitCalndles-sub))
	case "d":
		start = f.Yesterday - int64(time.Hour.Milliseconds()*int64(intN)*24*int64(f.LimitCalndles-sub))
	}

	return start
}
