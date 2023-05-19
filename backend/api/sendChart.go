package api

import (
	db "bitmoi/backend/db/sqlc"
)

const (
	PracticeMode int8 = iota
	CompetitionMode
)

var allPairs = []string{"BTCUSDT"}

func fiveMonthAndOneDay(interval string) (fivMonth int, waitingTime int) {
	switch interval {
	case db.FiveM:
		return 43200, 288
	case db.FifM:
		return 14400, 96
	case db.OneH:
		return 3600, 24
	case db.FourH:
		return 900, 6
	case db.OneD:
		return 150, 1
	default:
		return 0, 0
	}
}

// func SendCharts(mode int, interval string, names []string) Charts {
// 	var tenCharts Charts
// 	var ranName string
// Outer:
// 	for {
// 		ranName = allPairs[utilities.MakeRanInt(0, len(allPairs))]
// 		var sameHere bool = false
// 		for _, name := range names {
// 			if ranName == name {
// 				sameHere = true
// 			}
// 		}
// 		if !sameHere {
// 			break Outer
// 		}
// 	}
// 	chartBymode := makeChart(ranName, interval)
// 	if mode == CompetitionMode {
// 		(&chartBymode).anonymization(len(names))
// 	} else {
// 		chartBymode.addIdentifier()
// 	}
// 	tenCharts.Charts = chartBymode
// 	return tenCharts
// }

// func SendOtherInterval(identifier, reqInterval, mode string) CandleData {
// 	intervalChart := AC.InitAllchart(reqInterval)
// 	unmarshalOriginInfo := utilities.DecryptByASE(identifier)
// 	var originInfo = struct {
// 		Name         string  `json:"name"`
// 		Interval     string  `json:"interval"`
// 		Backsteps    int     `json:"backsteps"`
// 		PriceFactor  float64 `json:"pricefactor"`
// 		VolumeFactor float64 `json:"volumefactor"`
// 		TimeFactor   int64   `json:"timefactor"`
// 	}{}
// 	json.Unmarshal(unmarshalOriginInfo, &originInfo)
// 	tempchart := OnePairChart{Name: originInfo.Name, OneChart: (*intervalChart)[originInfo.Name], interval: reqInterval}
// 	tempchart.calculateBacksteps(originInfo.Backsteps, reqInterval)
// 	tempchart.OneChart.transformTime()
// 	if mode == "competition" {
// 		tempchart.OneChart.encodeChart(originInfo.PriceFactor, originInfo.VolumeFactor, originInfo.TimeFactor)
// 	}
// 	return tempchart.OneChart
// }
