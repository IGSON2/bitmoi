package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"fmt"
	"math"
	"sort"
	"time"
)

type PriceData struct {
	Name  string  `json:"-"`
	Open  float64 `json:"open"`
	Close float64 `json:"close"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Time  int64   `json:"time"`
}

type VolumeData struct {
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
	Color string  `json:"color"`
}

type CandleData struct {
	PData []PriceData  `json:"pdata"`
	VData []VolumeData `json:"vdata"`
}

type OnePairChart struct {
	Name         string     `json:"name"`
	OneChart     CandleData `json:"onechart"`
	TradingValue float64    `json:"tradingvalue"`
	EntryTime    string     `json:"entrytime"`
	Identifier   string     `json:"identifier"`
	interval     string
	backSteps    int
	priceFactor  float64
	volumeFactor float64
	ranPastDate  int64
}

type Charts struct {
	Charts OnePairChart `json:"charts"`
}

func getAllchart(start, end int64, intN int, intU string) map[string]CandleData {
	var allChart = make(map[string]CandleData)

	cycleNum := cycleByCase(start, end, intN, intU)
	for i := 1; i <= cycleNum; i++ {
		if i%3 == 0 {
			time.Sleep(1 * time.Minute)
		}
		Xcandles = i * 1000
		getOneIntvChart(intN, intU, allPairs)
		for j := 0; j < len(allPairs); j++ {
			var oneCoinOfOneIntv CandleData
			infoByCase(intN, intU, &oneCoinOfOneIntv)
			if oneCoinOfOneIntv.PData == nil || oneCoinOfOneIntv.VData == nil {
				continue
			}
			name := oneCoinOfOneIntv.PData[0].Name
			var tempPdatas []PriceData = allChart[name].PData
			var tempVdatas []VolumeData = allChart[name].VData
			tempPdatas = append(oneCoinOfOneIntv.PData, tempPdatas...)
			tempVdatas = append(oneCoinOfOneIntv.VData, tempVdatas...)
			allChart[name] = CandleData{tempPdatas, tempVdatas}
			fmt.Println(name, " Done.")
		}
		fmt.Printf("%s - %d / %d Done.\n", MakeInterval(intN, intU), i, cycleNum)
	}
	return allChart
}

func cycleByCase(start, end int64, intN int, intU string) int {
	var howHours int
	var cyclenum int
	termsByHours := (end - start) / (1000 * 60 * 60)
	howHours = int(termsByHours/250) + 1

	fmt.Println("TermHours : ", termsByHours)
	switch intU {
	case "m":
		switch intN {
		case 5:
			cyclenum = 3 * howHours
		case 15:
			cyclenum = howHours
		}
	case "h":
		switch intN {
		case 1:
			cyclenum = int(howHours/4) + 1
		case 4:
			cyclenum = int(howHours/16) + 1
		}
	case "d":
		cyclenum = int(howHours/96) + 1
	}
	fmt.Println("CycleNum : ", cyclenum)
	return cyclenum
}

func (s *Server) selectCandles(interval, name string, limit int) (CandlesInterface, error) {
	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandles(context.Background(), db.Get1dCandlesParams{name, int32(limit)})
		cs := Candles1dSlice(candles)
		return &cs, err
	case db.FourH:
		candles, err := s.store.Get4hCandles(context.Background(), db.Get4hCandlesParams{name, int32(limit)})
		cs := Candles4hSlice(candles)
		return &cs, err
	case db.OneH:
		candles, err := s.store.Get1hCandles(context.Background(), db.Get1hCandlesParams{name, int32(limit)})
		cs := Candles1hSlice(candles)
		return &cs, err
	case db.FifM:
		candles, err := s.store.Get15mCandles(context.Background(), db.Get15mCandlesParams{name, int32(limit)})
		cs := Candles15mSlice(candles)
		return &cs, err
	}
	return nil, fmt.Errorf("invalid interval %s", interval)
}

func makeChart(candles CandlesInterface, mode int) (*OnePairChart, error) {
	cd := candles.InitCandleData()
	var oc = OnePairChart{
		Name:      candles.Name(),
		interval:  candles.Interval(),
		EntryTime: candles.EntryTime(),
		OneChart:  cd,
	}

	oc.setRandomBackSteps()
	oc.setFactors()
	oc.addIdentifier()

}

func (o *OnePairChart) setFactors() {
	var mins []float64
	var vols []float64
	var ranPastDate int64 = int64(86400 * (utilities.MakeRanNum(19000, 10950)))
	for _, onebar := range o.OneChart.PData {
		mins = append(mins, onebar.Low)
	}
	for _, onebar := range o.OneChart.VData {
		vols = append(vols, onebar.Value)
	}
	sort.Slice(mins, func(i, j int) bool {
		return mins[i] < mins[j]
	})
	sort.Slice(vols, func(i, j int) bool {
		return vols[i] < vols[j]
	})
	basisPrice := mins[utilities.MakeRanNum(int(len(mins)/2), 0)]
	priceFactor := (100 / (mins[len(mins)-1] / mins[0])) / basisPrice
	basisVolume := vols[utilities.MakeRanNum(int(len(vols)/2), 0)]
	volumeFactor := ((mins[len(mins)-1] / mins[0]) / 20) / basisVolume

	o.priceFactor = priceFactor
	o.volumeFactor = volumeFactor
	o.ranPastDate = ranPastDate
}

func (o *OnePairChart) setRandomBackSteps() {
	fiveMonth, waitingTime := randomMinMax(o.interval)
	if fiveMonth > len(o.OneChart.PData) {
		fmt.Printf("%s is Shorter than fiveMonth.\n", o.Name)
		fiveMonth = 0
	}
	o.backSteps = utilities.MakeRanNum(len(o.OneChart.PData)-fiveMonth, waitingTime)
	if len(o.OneChart.PData)-o.backSteps < 2000 {
		o.OneChart = CandleData{
			PData: o.OneChart.PData[:len(o.OneChart.PData)-o.backSteps],
			VData: o.OneChart.VData[:len(o.OneChart.VData)-o.backSteps],
		}
	} else {
		o.OneChart = CandleData{
			PData: o.OneChart.PData[len(o.OneChart.PData)-o.backSteps-2000 : len(o.OneChart.PData)-o.backSteps],
			VData: o.OneChart.VData[len(o.OneChart.PData)-o.backSteps-2000 : len(o.OneChart.VData)-o.backSteps],
		}
	}
}

func (o *OnePairChart) anonymization(stage int) {

	o.OneChart.encodeValue(priceFactor, volumeFactor, ranPastDate)

	o.addIdentifier()
	o.backSteps = 0
	o.EntryTime = "Sometime"
	o.Name = fmt.Sprintf("STAGE %02d", stage+1)
}

func (o *OnePairChart) compareBTCvalue(btcChart CandleData) {

	var closeSum float64
	var volSum float64

	for _, onebar := range o.OneChart.PData[len(o.OneChart.PData)-721 : len(o.OneChart.PData)-1] {
		closeSum += onebar.Close
	}
	for _, onebar := range o.OneChart.VData[len(o.OneChart.VData)-721 : len(o.OneChart.VData)-1] {
		volSum += onebar.Value
	}
	thisTradingValue := closeSum * volSum

	btcEndIdx := len(btcChart.PData) - o.backSteps - 1
	btcStartIdx := btcEndIdx - 720

	var btcCloseSum float64
	var btcVolSum float64
	for _, oneBTCbar := range btcChart.PData[btcStartIdx:btcEndIdx] {
		btcCloseSum += oneBTCbar.Close
	}
	for _, oneBTCbar := range btcChart.VData[btcStartIdx:btcEndIdx] {
		btcVolSum += oneBTCbar.Value
	}
	btcTradingValue := btcCloseSum * btcVolSum

	o.TradingValue = math.Round((thisTradingValue/btcTradingValue)*10000) / 100

}

func (o *OnePairChart) addIdentifier() {
	uniqueInfo := utilities.ChartInfo{
		Name:         o.Name,
		Interval:     o.interval,
		Backsteps:    o.backSteps,
		PriceFactor:  o.priceFactor,
		VolumeFactor: o.volumeFactor,
		RanPastDate:  o.ranPastDate,
	}
	o.Identifier = utilities.EncrtpByASE(&uniqueInfo)
}

func (o *OnePairChart) calculateBacksteps(oneHourBacksteps int, reqInterval string) {
	switch reqInterval {
	case "5m":
		o.backSteps = ((oneHourBacksteps) * 12) - 11 + 576
	case "15m":
		o.backSteps = ((oneHourBacksteps) * 4) - 3 + 192
	case "4h":
		if oneHourBacksteps%4 == 0 {
			o.backSteps = int(oneHourBacksteps/4) + 1
		} else {
			o.backSteps = int(oneHourBacksteps/4) + 2
		}
	}
	if len(o.OneChart.PData)-o.backSteps < 2000 {
		o.OneChart = CandleData{
			PData: o.OneChart.PData[:len(o.OneChart.PData)-o.backSteps],
			VData: o.OneChart.VData[:len(o.OneChart.VData)-o.backSteps],
		}
	} else {
		o.OneChart = CandleData{
			PData: o.OneChart.PData[len(o.OneChart.PData)-o.backSteps-2000 : len(o.OneChart.PData)-o.backSteps],
			VData: o.OneChart.VData[len(o.OneChart.PData)-o.backSteps-2000 : len(o.OneChart.VData)-o.backSteps],
		}
	}
}

func (c *CandleData) encodeValue(pFactor, vFactor float64, pastDays int64) {
	var tempPData []PriceData
	var tempVData []VolumeData
	for _, onebar := range c.PData {
		newPbar := PriceData{
			Name:  "Anonymous",
			Open:  onebar.Open * pFactor,
			Close: onebar.Close * pFactor,
			High:  onebar.High * pFactor,
			Low:   onebar.Low * pFactor,
			Time:  onebar.Time - pastDays,
		}
		tempPData = append(tempPData, newPbar)
	}
	for _, onebar := range c.VData {
		newVbar := VolumeData{
			Value: onebar.Value * vFactor,
			Time:  onebar.Time - pastDays,
			Color: onebar.Color,
		}
		tempVData = append(tempVData, newVbar)
	}
	*c = CandleData{tempPData, tempVData}
}
