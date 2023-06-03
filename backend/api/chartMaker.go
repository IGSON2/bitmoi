package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
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
	Name         string      `json:"name"`
	OneChart     *CandleData `json:"onechart"`
	BtcRatio     float64     `json:"btcratio"`
	EntryTime    string      `json:"entrytime"`
	Identifier   string      `json:"identifier"`
	interval     string
	refTimestamp int64
	priceFactor  float64
	volumeFactor float64
	timeFactor   int64
}

type Charts struct {
	Charts OnePairChart `json:"charts"`
}

func (s *Server) SelectMinMaxTime(interval, name string, c *fiber.Ctx) (int64, int64, error) {
	switch interval {
	case db.OneD:
		r, err := s.store.Get1dMinMaxTime(c.Context(), name)
		return r.Min.(int64), r.Max.(int64), err
	case db.FourH:
		r, err := s.store.Get4hMinMaxTime(c.Context(), name)
		return r.Min.(int64), r.Max.(int64), err
	case db.OneH:
		r, err := s.store.Get1hMinMaxTime(c.Context(), name)
		return r.Min.(int64), r.Max.(int64), err
	case db.FifM:
		r, err := s.store.Get15mMinMaxTime(c.Context(), name)
		return r.Min.(int64), r.Max.(int64), err
	case db.FiveM:
		r, err := s.store.Get5mMinMaxTime(c.Context(), name)
		return r.Min.(int64), r.Max.(int64), err
	}
	return 0, 0, fmt.Errorf("invalid interval %s", interval)
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

func (s *Server) selectStageChart(name, interval string, refTimestamp int64, c *fiber.Ctx) (*CandleData, error) {
	var cdd CandleData

	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandles(c.Context(), db.Get1dCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles1dSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hCandles(c.Context(), db.Get4hCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles4hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hCandles(c.Context(), db.Get1hCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles1hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mCandles(c.Context(), db.Get15mCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles15mSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FiveM:
		candles, err := s.store.Get5mCandles(c.Context(), db.Get5mCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles5mSlice(candles)
		cdd = (&cs).InitCandleData()
	}
	if cdd.PData == nil || cdd.VData == nil {
		return nil, ErrGetStageChart
	}
	return &cdd, nil
}

func calculateRefTimestamp(section int64, name, interval string) int64 {
	fiveMonth, waitingTime := 150*24*time.Hour.Seconds(), 24*time.Hour.Seconds()
	if fiveMonth > float64(section) {
		fmt.Printf("%s is Shorter than fiveMonth.\n", name)
	}
	return int64(utilities.MakeRanInt(int(waitingTime), int(section)))
}

func (s *Server) makeChartToRef(interval, name string, mode string, prevStage int, c *fiber.Ctx) (*OnePairChart, error) {

	min, max, err := s.SelectMinMaxTime(interval, name, c)
	if err != nil {
		return nil, fmt.Errorf("cannot count all rows. name : %s, interval : %s, err : %w", name, interval, err)
	}

	refTimestamp := max - calculateRefTimestamp(max-min, name, interval)
	cdd, err := s.selectStageChart(name, interval, refTimestamp, c)
	if err != nil {
		return nil, fmt.Errorf("cannot make chart to reference timestamp. name : %s, interval : %s, err : %w", name, interval, err)
	}

	ratio, err := s.calcBtcRatio(interval, name, refTimestamp, c)
	if err != nil {
		return nil, fmt.Errorf("cannot calculate btc ratio. name : %s, interval : %s, err : %w", name, interval, err)
	}

	var oc = &OnePairChart{
		Name:         name,
		OneChart:     cdd,
		EntryTime:    utilities.EntryTimeFormatter(cdd.PData[len(cdd.PData)-1].Time),
		interval:     interval,
		refTimestamp: refTimestamp,
		BtcRatio:     ratio,
	}

	if mode == competition {
		if err := oc.setFactors(); err != nil {
			return nil, fmt.Errorf("cannot set factors. name : %s, interval : %s, err : %w", name, interval, err)
		}
		oc.anonymization(prevStage)
	} else {
		oc.addIdentifier()
	}
	return oc, nil
}

func (o *OnePairChart) setFactors() error {
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
	rf, err := utilities.MakeRanFloat(0, 100)
	if err != nil {
		return err
	}

	o.priceFactor = math.Floor(rf*decimal/mins[0]) / decimal
	o.volumeFactor = math.Floor(rf*decimal/vols[0]) / decimal
	o.timeFactor = timeFactor

	return nil
}

func (o *OnePairChart) anonymization(stage int) {

	o.OneChart.encodeChart(o.priceFactor, o.volumeFactor, o.timeFactor)

	o.addIdentifier()
	o.EntryTime = "Sometime"
	o.Name = fmt.Sprintf("STAGE %02d", stage+1)
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
			Name:  "",
			Open:  math.Floor(onebar.Open*pFactor*decimal) / decimal,
			Close: math.Floor(onebar.Close*pFactor*decimal) / decimal,
			High:  math.Floor(onebar.High*pFactor*decimal) / decimal,
			Low:   math.Floor(onebar.Low*pFactor*decimal) / decimal,
			Time:  onebar.Time - tFactor,
		}
		tempPData = append(tempPData, newPbar)
	}
	for _, onebar := range c.VData {
		newVbar := VolumeData{
			Value: math.Floor(onebar.Value*vFactor*decimal) / decimal,
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
	return math.Round(decimal*(rp*rv)/(bp*bv)) / 100, nil
}
