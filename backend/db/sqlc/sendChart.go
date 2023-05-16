package db

import (
	"bitmoi/backend/utilities"
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

const (
	PracticeMode int = iota
	CompetitionMode
)

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

var allPairs = []string{"BTCUSDT"}

func randomMinMax(interval string) (fivMonth int, waitingTime int) {
	switch interval {
	case FiveM:
		return 43200, 288
	case FifM:
		return 14400, 96
	case OneH:
		return 3600, 24
	case FourH:
		return 900, 6
	case OneD:
		return 150, 1
	default:
		return 0, 0
	}
}

func SendCharts(mode int, interval string, names []string) Charts {
	var tenCharts Charts
	var ranName string
Outer:
	for {
		ranName = allPairs[utilities.MakeRanNum(len(allPairs), 0)]
		var sameHere bool = false
		for _, name := range names {
			if ranName == name {
				sameHere = true
			}
		}
		if !sameHere {
			break Outer
		}
	}
	chartBymode := randomizeChart(ranName, interval)
	if mode == CompetitionMode {
		(&chartBymode).anonymization(len(names))
	} else {
		chartBymode.addIdentifier()
	}
	tenCharts.Charts = chartBymode
	return tenCharts
}

func randomizeChart(name, interval string) OnePairChart {
	intervalChart := AC.InitAllchart(interval)
	tempchart := OnePairChart{Name: name, OneChart: (*intervalChart)[name], interval: interval}
	tempchart.setRandomBackSteps()
	tempchart.compareBTCvalue((*intervalChart)["BTCUSDT"])
	tempchart.EntryTime = utilities.EntryTimeFormatter(tempchart.OneChart.PData[len(tempchart.OneChart.PData)-1].Time + 32400000)
	tempchart.OneChart.transformTime()
	return tempchart
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

	o.OneChart.encodeValue(priceFactor, volumeFactor, ranPastDate)
	o.priceFactor = priceFactor
	o.volumeFactor = volumeFactor
	o.ranPastDate = ranPastDate
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

func (c *CandleData) transformTime() {
	var tempPdata []PriceData
	var tempVdata []VolumeData
	for _, oneBar := range c.PData {
		oneBar.Time = (oneBar.Time / 1000) + 32400
		tempPdata = append(tempPdata, oneBar)
	}
	for _, oneBar := range c.VData {
		oneBar.Time = (oneBar.Time / 1000) + 32400
		tempVdata = append(tempVdata, oneBar)
	}
	*c = CandleData{tempPdata, tempVdata}
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

func SendOtherInterval(identifier, reqInterval, mode string) CandleData {
	intervalChart := AC.InitAllchart(reqInterval)
	unmarshalOriginInfo := utilities.DecryptByASE(identifier)
	var originInfo = struct {
		Name         string  `json:"name"`
		Interval     string  `json:"interval"`
		Backsteps    int     `json:"backsteps"`
		PriceFactor  float64 `json:"pricefactor"`
		VolumeFactor float64 `json:"volumefactor"`
		RanPastDate  int64   `json:"ranpastdate"`
	}{}
	json.Unmarshal(unmarshalOriginInfo, &originInfo)
	tempchart := OnePairChart{Name: originInfo.Name, OneChart: (*intervalChart)[originInfo.Name], interval: reqInterval}
	tempchart.calculateBacksteps(originInfo.Backsteps, reqInterval)
	tempchart.OneChart.transformTime()
	if mode == "competition" {
		tempchart.OneChart.encodeValue(originInfo.PriceFactor, originInfo.VolumeFactor, originInfo.RanPastDate)
	}
	return tempchart.OneChart
}
