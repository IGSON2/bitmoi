package db

import (
	"bitmoi/backend/utilities"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

var (
	apikey, scrkey       = getkeys()
	client               = futures.NewClient(apikey, scrkey)
	yesterday      int64 = yesterday9AM()
	wg             sync.WaitGroup
	OneDayCh           = make(chan CandleData, 150)
	FourHourCh         = make(chan CandleData, 150)
	OneHourCh          = make(chan CandleData, 150)
	FifMinuteCh        = make(chan CandleData, 150)
	FiveMinuteCh       = make(chan CandleData, 150)
	Xcandles       int = 1000
)

func getkeys() (string, string) {
	keys := utilities.ReadText("./backend/db/chartData/keys.TXT")

	apikey := (strings.Split(keys[0], ":"))
	scrkey := (strings.Split(keys[1], ":"))
	return apikey[1], scrkey[1]
}

//모든 페어들을 대상으로 15m, 1h, 1d 단위 각각 최대 1000개의 캔들 정보들을 수집합니다.
func getOneIntvChart(intN int, intU string, coinnames []string) {
	wg.Add(len(coinnames))
	for _, name := range coinnames {
		go getInfos(intN, intU, name)
	}
	wg.Wait()
}

//Go routine을 이용, 전역변수로 선언된 각 시간 단위별 채널에 대해 name 페어의 intN + intU 단위 최대 1000개의 캔들 정보를 수집하고 그에 맞는 채널로 전송합니다.
func getInfos(intN int, intU string, name string) {
	startTimemilli := howcandles(intN, intU, Xcandles)
	endTimeMilli := howcandles(intN, intU, Xcandles-1000)
	var infos CandleData
	var pDatas []PriceData
	var vDatas []VolumeData
	klines, err := client.NewKlinesService().Symbol(name).StartTime(startTimemilli).EndTime(endTimeMilli).Limit(1000).
		Interval(utilities.MakeInterval(intN, intU)).Do(context.Background())
	utilities.Errchk(err)
	for _, k := range klines {
		volume := utilities.StrToFloat(k.Volume)
		open := utilities.StrToFloat(k.Open)
		close := utilities.StrToFloat(k.Close)
		high := utilities.StrToFloat(k.High)
		low := utilities.StrToFloat(k.Low)
		date := k.OpenTime
		pData := PriceData{name, open, close, high, low, date}
		vData := VolumeData{volume, date, ""}
		if close >= open {
			vData.Color = "rgba(38,166,154,0.5)"
		} else {
			vData.Color = "rgba(239,83,80,0.5)"
		}
		pDatas = append(pDatas, pData)
		vDatas = append(vDatas, vData)
	}
	infos = CandleData{pDatas, vDatas}
	switch intU {
	case "m":
		switch intN {
		case 5:
			FiveMinuteCh <- infos
		case 15:
			FifMinuteCh <- infos
		}
	case "h":
		switch intN {
		case 1:
			OneHourCh <- infos
		case 4:
			FourHourCh <- infos
		}
	case "d":
		OneDayCh <- infos
	}

	wg.Done()
}

//현재로부터 intN + intU 단위의 캔들을 candles개 만큼 가져올 수 있는 일자를 Millisecond로 반환합니다.
func howcandles(intN int, intU string, xcandles int) int64 {
	var start int64

	switch intU {
	case "m":
		start = yesterday - int64(time.Minute.Milliseconds()*int64(intN)*int64(xcandles))
	case "h":
		start = yesterday - int64(time.Hour.Milliseconds()*int64(intN)*int64(xcandles))
	case "d":
		start = yesterday - int64(time.Hour.Milliseconds()*int64(intN)*24*int64(xcandles))
	}

	return start
}
