package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	oneTimeStageLoad = 2000
	BTCUSDT          = "BTCUSDT"
)

var (
	ErrGetStageChart  = errors.New("server cannot get stage chart")
	ErrGetResultChart = errors.New("server cannot get result chart")
	ErrShortRange     = errors.New("chart range is too short")
)

type PriceData struct {
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
	Name         string      `json:"name"`
	OneChart     *CandleData `json:"onechart"`
	BtcRatio     float64     `json:"btcratio"`
	EntryTime    string      `json:"entrytime"`
	EntryPrice   float64     `json:"entry_price"`
	Identifier   string      `json:"identifier"`
	interval     string
	refTimestamp int64
	priceFactor  float64
	volumeFactor float64
	timeFactor   int64
}

type PriceVolumeData struct {
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
	Time   int64   `json:"time"`
}

func (s *Server) calcBtcRatio(interval, name string, refTimestamp int64, c *fiber.Ctx) (float64, error) {
	switch interval {
	case db.OneD:
		b, err1 := s.store.Get1dVolSumPriceAVG(c.Context(), db.Get1dVolSumPriceAVGParams{Name: BTCUSDT, Time: refTimestamp})
		r, err2 := s.store.Get1dVolSumPriceAVG(c.Context(), db.Get1dVolSumPriceAVGParams{Name: name, Time: refTimestamp})
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("getBTC err : %w, getREQ err : %w", err1, err2)
		}
		return convTypeAndCalcRatio(b.Priceavg, b.Volsum, r.Priceavg, r.Volsum)
	case db.FourH:
		b, err1 := s.store.Get4hVolSumPriceAVG(c.Context(), db.Get4hVolSumPriceAVGParams{Name: BTCUSDT, Time: refTimestamp})
		r, err2 := s.store.Get4hVolSumPriceAVG(c.Context(), db.Get4hVolSumPriceAVGParams{Name: name, Time: refTimestamp})
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("getBTC err : %w, getREQ err : %w", err1, err2)
		}
		return convTypeAndCalcRatio(b.Priceavg, b.Volsum, r.Priceavg, r.Volsum)
	case db.OneH:
		b, err1 := s.store.Get1hVolSumPriceAVG(c.Context(), db.Get1hVolSumPriceAVGParams{Name: BTCUSDT, Time: refTimestamp})
		r, err2 := s.store.Get1hVolSumPriceAVG(c.Context(), db.Get1hVolSumPriceAVGParams{Name: name, Time: refTimestamp})
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("getBTC err : %w, getREQ err : %w", err1, err2)
		}
		return convTypeAndCalcRatio(b.Priceavg, b.Volsum, r.Priceavg, r.Volsum)
	case db.FifM:
		b, err1 := s.store.Get15mVolSumPriceAVG(c.Context(), db.Get15mVolSumPriceAVGParams{Name: BTCUSDT, Time: refTimestamp})
		r, err2 := s.store.Get15mVolSumPriceAVG(c.Context(), db.Get15mVolSumPriceAVGParams{Name: name, Time: refTimestamp})
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("getBTC err : %w, getREQ err : %w", err1, err2)
		}
		return convTypeAndCalcRatio(b.Priceavg, b.Volsum, r.Priceavg, r.Volsum)
	case db.FiveM:
		b, err1 := s.store.Get5mVolSumPriceAVG(c.Context(), db.Get5mVolSumPriceAVGParams{Name: BTCUSDT, Time: refTimestamp})
		r, err2 := s.store.Get5mVolSumPriceAVG(c.Context(), db.Get5mVolSumPriceAVGParams{Name: name, Time: refTimestamp})
		if err1 != nil || err2 != nil {
			return -1, fmt.Errorf("getBTC err : %w, getREQ err : %w", err1, err2)
		}
		return convTypeAndCalcRatio(b.Priceavg, b.Volsum, r.Priceavg, r.Volsum)
	}
	return -1, fmt.Errorf("invalid interval %s", interval)
}

func (s *Server) selectStageChart(name, interval string, refTimestamp int64, ctx context.Context) (*CandleData, error) {
	cdd := new(CandleData)

	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandles(ctx, db.Get1dCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles1dSlice(candles)
		cdd = (cs).InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hCandles(ctx, db.Get4hCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles4hSlice(candles)
		cdd = (cs).InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hCandles(ctx, db.Get1hCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles1hSlice(candles)
		cdd = (cs).InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mCandles(ctx, db.Get15mCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles15mSlice(candles)
		cdd = (cs).InitCandleData()
	case db.FiveM:
		candles, err := s.store.Get5mCandles(ctx, db.Get5mCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles5mSlice(candles)
		cdd = (cs).InitCandleData()
	}
	if cdd.PData == nil || cdd.VData == nil {
		return nil, ErrGetStageChart
	}
	return cdd, nil
}

func calculateRefTimestamp(section int64) int64 {
	oneMonth, waitingTime := 30*24*time.Hour.Seconds(), 30*24*time.Hour.Seconds()
	return int64(utilities.MakeRanInt(int(waitingTime), int(section-int64(oneMonth))))
}

func (s *Server) makeChartToRef(interval, name string, mode string, ctx context.Context) (*OnePairChart, error) {

	min, max, err := s.store.SelectMinMaxTime(interval, name, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot count all rows. name : %s, interval : %s, err : %w", name, interval, err)
	}

	refTimestamp := max - calculateRefTimestamp(max-min)
	if refTimestamp == max {
		return nil, ErrShortRange
	}
	cdd, err := s.selectStageChart(name, interval, refTimestamp, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot make chart to reference timestamp. name : %s, interval : %s, err : %w", name, interval, err)
	}

	// ratio, err := s.calcBtcRatio(interval, name, refTimestamp, c)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot calculate btc ratio. name : %s, interval : %s, refTime : %d, err : %w",
	// 		name, interval, refTimestamp, err)
	// }

	var oc = &OnePairChart{
		Name:         name,
		OneChart:     cdd,
		EntryTime:    utilities.EntryTimeFormatter(cdd.PData[0].Time),
		interval:     interval,
		refTimestamp: cdd.PData[0].Time,
		// BtcRatio:     common.CeilDecimal(ratio) * 100,
		BtcRatio: 0,
	}

	if mode == competition {
		if err := oc.setFactors(); err != nil {
			return nil, fmt.Errorf("cannot set factors. name : %s, interval : %s, err : %w", name, interval, err)
		}
		oc.anonymization()
	} else {
		oc.addIdentifier()
		oc.EntryPrice = oc.OneChart.PData[0].Close
	}
	return oc, nil
}

func (o *OnePairChart) setFactors() error {
	head := int(float64(len(o.OneChart.PData)) * 0.1)
	pd := make([]PriceData, head)
	vd := make([]VolumeData, head)

	var timeFactor int64 = int64(86400 * (utilities.MakeRanInt(10950, 19000)))
	copy(pd, o.OneChart.PData[:head])
	copy(vd, o.OneChart.VData[:head])

	sort.Slice(pd, func(i, j int) bool {
		return pd[i].Low < pd[j].Low
	})
	sort.Slice(vd, func(i, j int) bool {
		return vd[i].Value < vd[j].Value
	})

	rf, err := utilities.MakeRanFloat(1, 100)
	if err != nil {
		return err
	}

	vrf, err := utilities.MakeRanFloat(10000, 10000000)
	if err != nil {
		return err
	}

	o.timeFactor = timeFactor
	o.priceFactor = common.FloorDecimal(rf / pd[0].Low)
	o.volumeFactor = math.Floor(vrf/vd[0].Value*100000000000) / 100000000000

	return nil
}

func (o *OnePairChart) anonymization() {

	o.OneChart.encodeChart(o.priceFactor, o.volumeFactor, o.timeFactor)
	o.EntryPrice = o.OneChart.PData[0].Close
	o.addIdentifier()
	o.EntryTime = "Sometime"
	o.Name = "STAGE"
}

func (o *OnePairChart) addIdentifier() {
	uniqueInfo := utilities.IdentificationData{
		Name:         o.Name,
		Interval:     o.interval,
		RefTimestamp: o.refTimestamp,
		PriceFactor:  o.priceFactor,
		VolumeFactor: o.volumeFactor,
		TimeFactor:   o.timeFactor,
	}

	o.Identifier = utilities.EncrtpByASE(&uniqueInfo)
}

func (c *CandleData) encodeChart(pFactor, vFactor float64, tFactor int64) {
	var tempPData []PriceData
	var tempVData []VolumeData

	for _, onebar := range c.PData {
		newPbar := PriceData{
			Open:  common.FloorDecimal(onebar.Open * pFactor),
			Close: common.FloorDecimal(onebar.Close * pFactor),
			High:  common.FloorDecimal(onebar.High * pFactor),
			Low:   common.FloorDecimal(onebar.Low * pFactor),
			Time:  onebar.Time - tFactor,
		}
		tempPData = append(tempPData, newPbar)
	}
	for _, onebar := range c.VData {
		newVbar := VolumeData{
			Value: common.FloorDecimal(onebar.Value * vFactor),
			Time:  onebar.Time - tFactor,
			Color: onebar.Color,
		}
		tempVData = append(tempVData, newVbar)
	}
	*c = CandleData{tempPData, tempVData}
}

func convTypeAndCalcRatio(btcP, btcV, reqP, reqV interface{}) (float64, error) {
	bp, ok1 := btcP.(float64)
	bv, ok2 := btcV.(float64)
	rp, ok3 := reqP.(float64)
	rv, ok4 := reqV.(float64)
	if bp == 0 || bv == 0 {
		return -1, fmt.Errorf("cannot select btcusdt data")
	}
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return -1, fmt.Errorf("cannot conver type into float64, bp,pv,rp,rv : %t,%t,%t,%t", ok1, ok2, ok3, ok4)
	}

	return (rp * rv) / (bp * bv), nil
}
