package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

const (
	oneTimeLoad = 2000
)

var (
	ErrGetCandleData = errors.New("server cannot get candle data")
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
	BtcRatio     float64    `json:"btcratio"`
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

func (s *Server) makeChartUpToRef(interval, name string, mode, prevStage int8) (*OnePairChart, error) {
	var cdd CandleData
	logErr := func(err error) {
		s.logger.Error().Err(err).Msg(ErrGetCandleData.Error())
	}
	min, max, err := s.SelectMinMaxTime(interval, name)
	if err != nil {
		return nil, fmt.Errorf("cannot count all rows. name : %s, interval : %s, err : %w", name, interval, err)
	}

	refTimestamp := max - calculateRefTimestamp(max-min, name, interval)

	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandles(context.Background(), db.Get1dCandlesParams{name, refTimestamp, oneTimeLoad})
		logErr(err)
		cs := Candles1dSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hCandles(context.Background(), db.Get4hCandlesParams{name, refTimestamp, oneTimeLoad})
		logErr(err)
		cs := Candles4hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hCandles(context.Background(), db.Get1hCandlesParams{name, refTimestamp, oneTimeLoad})
		logErr(err)
		cs := Candles1hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mCandles(context.Background(), db.Get15mCandlesParams{name, refTimestamp, oneTimeLoad})
		logErr(err)
		cs := Candles15mSlice(candles)
		cdd = (&cs).InitCandleData()
	}
	if cdd.PData == nil || cdd.VData == nil {
		return nil, ErrGetCandleData
	}

	var oc = &OnePairChart{
		Name:         name,
		OneChart:     cdd,
		EntryTime:    utilities.EntryTimeFormatter(cdd.PData[len(cdd.PData)-1].Time),
		interval:     interval,
		refTimestamp: int(refTimestamp),
	}

	//TODO : SetBtcRatio

	if mode == CompetitionMode {
		oc.setFactors()
		oc.anonymization(prevStage)
	} else {
		oc.addIdentifier()
	}
	return oc, nil
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

func (o *OnePairChart) anonymization(stage int8) {

	o.OneChart.encodeChart(o.priceFactor, o.volumeFactor, o.timeFactor)

	o.addIdentifier()
	o.refTimestamp = 0
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
	thisBtcRatio := closeSum * volSum

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
	btcBtcRatio := btcCloseSum * btcVolSum

	o.BtcRatio = math.Round((thisBtcRatio/btcBtcRatio)*10000) / 100

}

func (o *OnePairChart) addIdentifier() {
	uniqueInfo := utilities.IdentificationData{
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
			Name:  "",
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
