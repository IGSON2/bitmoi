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

const (
	oneTimeLoad = 2000
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
	refTimestamp int
	priceFactor  float64
	volumeFactor float64
	timeFactor   int64
}

type Charts struct {
	Charts OnePairChart `json:"charts"`
}

// func getAllchart(start, end int64, intN int, intU string) map[string]CandleData {
// 	var allChart = make(map[string]CandleData)

// 	cycleNum := cycleByCase(start, end, intN, intU)
// 	for i := 1; i <= cycleNum; i++ {
// 		if i%3 == 0 {
// 			time.Sleep(1 * time.Minute)
// 		}
// 		Xcandles = i * 1000
// 		getOneIntvChart(intN, intU, allPairs)
// 		for j := 0; j < len(allPairs); j++ {
// 			var oneCoinOfOneIntv CandleData
// 			infoByCase(intN, intU, &oneCoinOfOneIntv)
// 			if oneCoinOfOneIntv.PData == nil || oneCoinOfOneIntv.VData == nil {
// 				continue
// 			}
// 			name := oneCoinOfOneIntv.PData[0].Name
// 			var tempPdatas []PriceData = allChart[name].PData
// 			var tempVdatas []VolumeData = allChart[name].VData
// 			tempPdatas = append(oneCoinOfOneIntv.PData, tempPdatas...)
// 			tempVdatas = append(oneCoinOfOneIntv.VData, tempVdatas...)
// 			allChart[name] = CandleData{tempPdatas, tempVdatas}
// 			fmt.Println(name, " Done.")
// 		}
// 		fmt.Printf("%s - %d / %d Done.\n", MakeInterval(intN, intU), i, cycleNum)
// 	}
// 	return allChart
// }

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

func (s *Server) SelectMinMaxTime(interval, name string) (int64, int64, error) {
	switch interval {
	case db.OneD:
		r, err := s.store.Get1dMinMaxTime(context.Background(), name)
		return r.Min.(int64), r.Max.(int64), err
	case db.FourH:
		r, err := s.store.Get4hMinMaxTime(context.Background(), name)
		return r.Min.(int64), r.Max.(int64), err
	case db.OneH:
		r, err := s.store.Get1hMinMaxTime(context.Background(), name)
		return r.Min.(int64), r.Max.(int64), err
	case db.FifM:
		r, err := s.store.Get15mMinMaxTime(context.Background(), name)
		return r.Min.(int64), r.Max.(int64), err
	}
	return 0, 0, fmt.Errorf("invalid interval %s", interval)
}

func calculateRefTimestamp(section int64, name, interval string) int64 {
	fiveMonth, waitingTime := 150*24*time.Hour.Seconds(), 24*time.Hour.Seconds()
	if fiveMonth > float64(section) {
		fmt.Printf("%s is Shorter than fiveMonth.\n", name)
	}
	return int64(utilities.MakeRanInt(int(waitingTime), int(section)))
}

func (s *Server) selectCandlesToRef(interval, name string) (CandlesInterface, int64, error) {
	min, max, err := s.SelectMinMaxTime(interval, name)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot count all rows. name : %s, interval : %s, err : %w", name, interval, err)
	}

	refTimestamp := max - calculateRefTimestamp(max-min, name, interval)

	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandles(context.Background(), db.Get1dCandlesParams{name, refTimestamp, oneTimeLoad})
		cs := Candles1dSlice(candles)
		return &cs, refTimestamp, err
	case db.FourH:
		candles, err := s.store.Get4hCandles(context.Background(), db.Get4hCandlesParams{name, refTimestamp, oneTimeLoad})
		cs := Candles4hSlice(candles)
		return &cs, refTimestamp, err
	case db.OneH:
		candles, err := s.store.Get1hCandles(context.Background(), db.Get1hCandlesParams{name, refTimestamp, oneTimeLoad})
		cs := Candles1hSlice(candles)
		return &cs, refTimestamp, err
	case db.FifM:
		candles, err := s.store.Get15mCandles(context.Background(), db.Get15mCandlesParams{name, refTimestamp, oneTimeLoad})
		cs := Candles15mSlice(candles)
		return &cs, refTimestamp, err
	}
	return nil, 0, fmt.Errorf("invalid interval %s", interval)
}

func makeChart(candles CandlesInterface, mode int8, refTimestamp int64) (*OnePairChart, error) {
	var oc = OnePairChart{
		Name:         candles.Name(),
		interval:     candles.Interval(),
		EntryTime:    candles.EntryTime(),
		refTimestamp: int(refTimestamp),
		OneChart:     candles.InitCandleData(),
	}

	if mode == CompetitionMode {
		oc.setFactors()
		oc.anonymization()
	} else {
		oc.addIdentifier()
	}
	return &oc, nil
}

func (o *OnePairChart) setFactors() {
	var mins []float64
	var vols []float64
	var timeFactor int64 = int64(86400 * (utilities.MakeRanInt(10950, 19000)))
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
	basisPrice := mins[utilities.MakeRanInt(0, int(len(mins)/2))]
	priceFactor := (100 / (mins[len(mins)-1] / mins[0])) / basisPrice
	basisVolume := vols[utilities.MakeRanInt(0, int(len(vols)/2))]
	volumeFactor := ((mins[len(mins)-1] / mins[0]) / 20) / basisVolume

	o.priceFactor = priceFactor
	o.volumeFactor = volumeFactor
	o.timeFactor = timeFactor
}

func (o *OnePairChart) anonymization() {

	o.OneChart.encodeChart(o.priceFactor, o.volumeFactor, o.timeFactor)

	o.addIdentifier()
	o.refTimestamp = 0
	o.EntryTime = "Sometime"
	o.Name = "SomePair"
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

	btcEndIdx := len(btcChart.PData) - o.refTimestamp - 1
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
		Backsteps:    o.refTimestamp,
		PriceFactor:  o.priceFactor,
		VolumeFactor: o.volumeFactor,
		TimeFactor:   o.timeFactor,
	}
	o.Identifier = utilities.EncrtpByASE(&uniqueInfo)
}

func (o *OnePairChart) calculateBacksteps(oneHourBacksteps int, reqInterval string) {
	switch reqInterval {
	case "5m":
		o.refTimestamp = ((oneHourBacksteps) * 12) - 11 + 576
	case "15m":
		o.refTimestamp = ((oneHourBacksteps) * 4) - 3 + 192
	case "4h":
		if oneHourBacksteps%4 == 0 {
			o.refTimestamp = int(oneHourBacksteps/4) + 1
		} else {
			o.refTimestamp = int(oneHourBacksteps/4) + 2
		}
	}
	if len(o.OneChart.PData)-o.refTimestamp < 2000 {
		o.OneChart = CandleData{
			PData: o.OneChart.PData[:len(o.OneChart.PData)-o.refTimestamp],
			VData: o.OneChart.VData[:len(o.OneChart.VData)-o.refTimestamp],
		}
	} else {
		o.OneChart = CandleData{
			PData: o.OneChart.PData[len(o.OneChart.PData)-o.refTimestamp-2000 : len(o.OneChart.PData)-o.refTimestamp],
			VData: o.OneChart.VData[len(o.OneChart.PData)-o.refTimestamp-2000 : len(o.OneChart.VData)-o.refTimestamp],
		}
	}
}

func (c *CandleData) encodeChart(pFactor, vFactor float64, tFactor int64) {
	var tempPData []PriceData
	var tempVData []VolumeData
	for _, onebar := range c.PData {
		newPbar := PriceData{
			Name:  "Anonymous",
			Open:  onebar.Open * pFactor,
			Close: onebar.Close * pFactor,
			High:  onebar.High * pFactor,
			Low:   onebar.Low * pFactor,
			Time:  onebar.Time - tFactor,
		}
		tempPData = append(tempPData, newPbar)
	}
	for _, onebar := range c.VData {
		newVbar := VolumeData{
			Value: onebar.Value * vFactor,
			Time:  onebar.Time - tFactor,
			Color: onebar.Color,
		}
		tempVData = append(tempVData, newVbar)
	}
	*c = CandleData{tempPData, tempVData}
}
