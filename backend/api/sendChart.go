package api

const (
	PracticeMode int8 = iota
	CompetitionMode
)

var allPairs = []string{"BTCUSDT"}

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
