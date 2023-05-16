package futureclient

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

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
	CandlesCh     chan db.InsertCandlesParams
	LimitCalndles int
}

func NewFutureClient(c utilities.Config) (*FutureClient, error) {
	dbConn, err := sql.Open(c.DBDriver, c.DBSource)

	if err != nil {
		return nil, fmt.Errorf("cannot open db store %w", err)
	}

	f := &FutureClient{
		Client:        futures.NewClient(utilities.GetAPIKeys()),
		Store:         db.NewStore(dbConn),
		Yesterday:     utilities.Yesterday9AM(),
		LimitCalndles: LimitCandlesNum,
	}
	if getErr := f.getAllPairs(); getErr != nil {
		return nil, getErr
	}
	f.CandlesCh = make(chan db.InsertCandlesParams, len(f.Pairs))
	return f, nil
}

func (f *FutureClient) getAllPairs() error {
	info, err := f.Client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get allpairs %w", err)
	}
	for _, s := range info.Symbols {
		f.Pairs = append(f.Pairs, s.Symbol)
	}
	return nil
}

// Go routine을 이용, 전역변수로 선언된 각 시간 단위별 채널에 대해 name 페어의 intN + intU 단위 최대 1000개의 캔들 정보를 수집하고 그에 맞는 채널로 전송합니다.
func (f *FutureClient) getInfos(intN int, intU string, name string) ([]db.InsertCandlesParams, error) {
	startTimemilli := f.howcandles(intN, intU, 0)
	endTimeMilli := f.howcandles(intN, intU, 1000)
	var infos []db.InsertCandlesParams

	klines, err := f.Client.NewKlinesService().Symbol(name).StartTime(startTimemilli).EndTime(endTimeMilli).Limit(10).
		Interval(utilities.MakeInterval(intN, intU)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	for _, k := range klines {
		i := db.InsertCandlesParams{
			Name:   name,
			Open:   utilities.StrToFloat(k.Open),
			Close:  utilities.StrToFloat(k.Close),
			High:   utilities.StrToFloat(k.High),
			Low:    utilities.StrToFloat(k.Low),
			Time:   k.OpenTime,
			Volume: utilities.StrToFloat(k.Volume),
		}
		if i.Close >= i.Open {
			i.Color = "rgba(38,166,154,0.5)"
		} else {
			i.Color = "rgba(239,83,80,0.5)"
		}
		_, err = f.Store.InsertCandles(context.Background(), i)
		if err != nil {
			return nil, err
		}
		infos = append(infos, i)
	}
	return infos, nil
}

// 현재로부터 intN + intU 단위의 캔들을 candles개 만큼 가져올 수 있는 일자를 Millisecond로 반환합니다.
func (f *FutureClient) howcandles(intN int, intU string, sub int) int64 {
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
