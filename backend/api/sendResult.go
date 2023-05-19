package api

import (
	"bitmoi/backend/utilities"
	"math"
)

type OrderStruct struct {
	Mode         string  `json:"mode"`
	Uid          string  `json:"uid"`
	Name         string  `json:"name"`
	Entrytime    string  `json:"entrytime"`
	Stage        int     `json:"stage"`
	IsLong       bool    `json:"islong"`
	EntryPrice   float64 `json:"entryprice"`
	Quantity     float64 `json:"quantity"`
	QuantityRate float64 `json:"quantityrate"`
	ProfitPrice  float64 `json:"profitprice"`
	LossPrice    float64 `json:"lossprice"`
	Leverage     int     `json:"leverage"`
	Balance      float64 `json:"balance"`
	Identifier   string  `json:"identifier,omitempty"`
	ScoreId      string  `json:"scoreid"`
}

type ResultScore struct {
	Stage        int     `json:"stage"`
	Name         string  `json:"name"`
	Entrytime    string  `json:"entrytime"`
	Leverage     int     `json:"leverage"`
	EntryPrice   float64 `json:"entryprice"`
	EndPrice     float64 `json:"-"`
	OutHour      int     `json:"outhour"`
	Roe          float64 `json:"roe"`
	Pnl          float64 `json:"pnl"`
	Commission   float64 `json:"commission"`
	Isliquidated bool    `json:"isliquidated"`
}

type ResultData struct {
	OriginChart    CandleData  `json:"originchart"`
	ResultChart    CandleData  `json:"resultchart"`
	ResultScore    ResultScore `json:"resultscore"`
	CompOriginName string      `json:"comporiginname,omitempty"`
}

// func SendCompResult(compOrder OrderStruct) (*ResultData, error) {
// 	compInfoByte := utilities.DecryptByASE(compOrder.Identifier)
// 	var compInfo utilities.ChartInfo
// 	err := json.Unmarshal(compInfoByte, &compInfo)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot unmarshal compition chart identifier err : %w", err)
// 	}
// 	resultchart := compResult(&compInfo)
// 	var compResult = ResultData{
// 		OriginChart: compOrigin(&compInfo),
// 		ResultChart: resultchart,
// 		ResultScore: calculateResult(resultchart, compOrder),
// 	}
// 	compResult.ResultScore.Name = compResult.ResultChart.PData[0].Name
// 	compResult.ResultScore.Entrytime = utilities.EntryTimeFormatter(compResult.ResultChart.PData[0].Time*1000 - 3600000)
// 	return &compResult, nil
// }

// func SendPracResult(pracOrder OrderStruct) (*ResultData, error) {
// 	pracInfoByte := utilities.DecryptByASE(pracOrder.Identifier)
// 	var pracInfo utilities.ChartInfo
// 	err := json.Unmarshal(pracInfoByte, &pracInfo)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot unmarshal compition chart identifier err : %w", err)
// 	}
// 	resultchart := pracResult(&pracInfo)
// 	var pracResult = ResultData{
// 		OriginChart: CandleData{},
// 		ResultChart: resultchart,
// 		ResultScore: calculateResult(resultchart, pracOrder),
// 	}
// 	return &pracResult, nil
// }

// func compOrigin(info *utilities.ChartInfo) CandleData {
// 	var loadedChart = CandleData{
// 		PData: (*AC.InitAllchart(OneH))[info.Name].PData[:len((*AC.InitAllchart(OneH))[info.Name].PData)-info.Backsteps],
// 		VData: (*AC.InitAllchart(OneH))[info.Name].VData[:len((*AC.InitAllchart(OneH))[info.Name].VData)-info.Backsteps],
// 	}
// 	return decodeChartWithoutPrice(loadedChart, info)
// }

// func compResult(info *utilities.ChartInfo) CandleData {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Println("PANIC : ", info.Name)
// 		}
// 	}()
// 	strIdx := len((*AC.InitAllchart(OneH))[info.Name].PData) - info.Backsteps
// 	loadedChart := CandleData{
// 		PData: (*AC.InitAllchart(OneH))[info.Name].PData[strIdx : strIdx+24],
// 		VData: (*AC.InitAllchart(OneH))[info.Name].VData[strIdx : strIdx+24],
// 	}
// 	return decodeChartWithoutPrice(loadedChart, info)
// }

// func pracResult(info *utilities.ChartInfo) CandleData {

// 	strIdx := len((*AC.InitAllchart(OneH))[info.Name].PData) - info.Backsteps
// 	loadedChart := CandleData{
// 		PData: (*AC.InitAllchart(OneH))[info.Name].PData[strIdx : strIdx+24],
// 		VData: (*AC.InitAllchart(OneH))[info.Name].VData[strIdx : strIdx+24],
// 	}
// 	loadedChart.transformTime()
// 	return loadedChart
// }

func decodeChartWithoutPrice(loadedChart CandleData, info *utilities.ChartInfo) CandleData {
	var tempPdata []PriceData
	var tempVdata []VolumeData
	for _, onebar := range loadedChart.PData {
		var newPbar = PriceData{
			Name:  info.Name,
			Open:  onebar.Open * info.PriceFactor,
			Close: onebar.Close * info.PriceFactor,
			High:  onebar.High * info.PriceFactor,
			Low:   onebar.Low * info.PriceFactor,
			Time:  onebar.Time/1000 + 32400,
		}
		tempPdata = append(tempPdata, newPbar)
	}
	for _, onebar := range loadedChart.VData {
		var newVbar = VolumeData{
			Value: onebar.Value * info.VolumeFactor,
			Time:  onebar.Time/1000 + 32400,
			Color: onebar.Color,
		}
		tempVdata = append(tempVdata, newVbar)
	}

	return CandleData{tempPdata, tempVdata}
}

func decodeChartTimeOnly(loadedChart CandleData) CandleData {
	var tempPdata []PriceData
	var tempVdata []VolumeData
	for _, onebar := range loadedChart.PData {
		newPbar := onebar
		newPbar.Time = onebar.Time/1000 + 32400
		tempPdata = append(tempPdata, newPbar)
	}
	for _, onebar := range loadedChart.VData {
		newVbar := onebar
		newVbar.Time = onebar.Time/1000 + 32400
		tempVdata = append(tempVdata, newVbar)
	}
	return CandleData{tempPdata, tempVdata}
}

func calculateResult(resultchart CandleData, order OrderStruct) ResultScore {
	var roe float64
	var pnl float64
	var endIdx int
	var endPrice float64
	var out int
	for idx, candle := range resultchart.PData {
		if order.IsLong {
			if candle.High >= order.ProfitPrice {
				roe = float64(order.QuantityRate/100) * (order.ProfitPrice - order.EntryPrice) / order.EntryPrice
				endIdx = idx
				endPrice = candle.High
				break
			}
			if candle.Low <= order.LossPrice {
				roe = float64(order.QuantityRate/100) * (order.LossPrice - order.EntryPrice) / order.EntryPrice
				endIdx = idx
				endPrice = candle.Low
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = float64(order.QuantityRate/100) * (candle.Close - order.EntryPrice) / order.EntryPrice
				endIdx = idx
				endPrice = candle.Close
				break
			}
		} else {
			if candle.Low <= order.ProfitPrice {
				roe = float64(order.QuantityRate/100) * (order.EntryPrice - order.ProfitPrice) / order.EntryPrice
				endIdx = idx
				endPrice = candle.Low
				break
			}
			if candle.High >= order.LossPrice {
				roe = float64(order.QuantityRate/100) * (order.EntryPrice - order.LossPrice) / order.EntryPrice
				endIdx = idx
				endPrice = candle.High
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = float64(order.QuantityRate/100) * (order.EntryPrice - candle.Close) / order.EntryPrice
				endIdx = idx
				endPrice = candle.Close
				break
			}
		}
	}
	out = endIdx + 1
	pnl = roe * float64(100/order.QuantityRate) * order.EntryPrice * order.Quantity

	tempInfo := ResultScore{
		Stage:      order.Stage,
		Name:       order.Name,
		Entrytime:  order.Entrytime,
		Leverage:   order.Leverage,
		EntryPrice: order.EntryPrice,
		EndPrice:   (math.Floor(endPrice*10000) / 10000),
		OutHour:    out,
		Roe:        (math.Floor(roe*10000*float64(order.Leverage)) / 100),
		Pnl:        math.Floor(pnl*10000) / 10000,
		Commission: math.Floor(0.0002*order.EntryPrice*order.Quantity*10000) / 10000,
	}
	if order.Balance+tempInfo.Pnl-tempInfo.Commission < 1 {
		tempInfo.Isliquidated = true
	}
	return tempInfo
}
