package db

import (
	"bitmoi/backend/utilities"
	"fmt"
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

func SendCompResult(compOrder OrderStruct) ResultData {
	compInfo := utilities.DecryptByASE(compOrder.Identifier)
	// TODO: decode compInfo
	resultchart := compResult(compInfo)
	var compResult = ResultData{
		OriginChart: compOrigin(compInfo),
		ResultChart: resultchart,
		ResultScore: calculateResult(resultchart, compOrder),
	}
	compResult.ResultScore.Name = compResult.ResultChart.PData[0].Name
	compResult.ResultScore.Entrytime = utilities.EntryTimeFormatter(compResult.ResultChart.PData[0].Time*1000 - 3600000)
	return compResult
}

func SendPracResult(pracOrder OrderStruct) ResultData {
	pracInfo := moisha.DecodeInfo(pracOrder.Identifier)
	resultchart := pracResult(pracInfo)
	var pracResult = ResultData{
		OriginChart: CandleData{},
		ResultChart: resultchart,
		ResultScore: calculateResult(resultchart, pracOrder),
	}
	return pracResult
}

func compOrigin(info *moisha.OriginInfo) CandleData {
	var loadedChart = CandleData{
		PData: (*AC.InitAllchart(OneH))[info.BmName].PData[:len((*AC.InitAllchart(OneH))[info.BmName].PData)-info.BmBacksteps],
		VData: (*AC.InitAllchart(OneH))[info.BmName].VData[:len((*AC.InitAllchart(OneH))[info.BmName].VData)-info.BmBacksteps],
	}
	return decodeChartWithoutPrice(loadedChart, info)
}

func compResult(info *moisha.OriginInfo) CandleData {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("PANIC : ", info.BmName)
		}
	}()
	strIdx := len((*AC.InitAllchart(OneH))[info.BmName].PData) - info.BmBacksteps
	loadedChart := CandleData{
		PData: (*AC.InitAllchart(OneH))[info.BmName].PData[strIdx : strIdx+24],
		VData: (*AC.InitAllchart(OneH))[info.BmName].VData[strIdx : strIdx+24],
	}
	return decodeChartWithoutPrice(loadedChart, info)
}

func pracResult(info *moisha.OriginInfo) CandleData {

	strIdx := len((*AC.InitAllchart(OneH))[info.BmName].PData) - info.BmBacksteps
	loadedChart := CandleData{
		PData: (*AC.InitAllchart(OneH))[info.BmName].PData[strIdx : strIdx+24],
		VData: (*AC.InitAllchart(OneH))[info.BmName].VData[strIdx : strIdx+24],
	}
	loadedChart.transformTime()
	return loadedChart
}

func decodeChartWithoutPrice(loadedChart CandleData, info *moisha.OriginInfo) CandleData {
	var tempPdata []PriceData
	var tempVdata []VolumeData
	for _, onebar := range loadedChart.PData {
		var newPbar = PriceData{
			Name:  info.BmName,
			Open:  onebar.Open * info.BmPriceFactor,
			Close: onebar.Close * info.BmPriceFactor,
			High:  onebar.High * info.BmPriceFactor,
			Low:   onebar.Low * info.BmPriceFactor,
			Time:  onebar.Time/1000 + 32400,
		}
		tempPdata = append(tempPdata, newPbar)
	}
	for _, onebar := range loadedChart.VData {
		var newVbar = VolumeData{
			Value: onebar.Value * info.BmVolumeFactor,
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
