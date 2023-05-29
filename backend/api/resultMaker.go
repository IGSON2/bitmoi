package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"encoding/json"
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"
)

const (
	commissionRate = 0.0002
	decimal        = 10000
)

type ResultScore struct {
	Stage        int32   `json:"stage"`
	Name         string  `json:"name"`
	Entrytime    string  `json:"entrytime"`
	Leverage     int32   `json:"leverage"`
	EntryPrice   float64 `json:"entryprice"`
	EndPrice     float64 `json:"-"`
	OutTime      int32   `json:"outtime"`
	Roe          float64 `json:"roe"`
	Pnl          float64 `json:"pnl"`
	Commission   float64 `json:"commission"`
	Isliquidated bool    `json:"isliquidated"`
}

type ResultData struct {
	OriginChart *CandleData  `json:"originchart"`
	ResultChart *CandleData  `json:"resultchart"`
	ResultScore *ResultScore `json:"resultscore"`
}

func (s *Server) createPracResult(order *OrderRequest, c *fiber.Ctx) (*ResultData, error) {
	pracInfo := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(order.Identifier)
	err := json.Unmarshal(infoByte, pracInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal chart identifier. err : %w", err)
	}
	resultchart, err := s.selectResultChart(pracInfo, int(order.WaitingTerm), c)
	if err != nil {
		return nil, fmt.Errorf("cannot select result chart. err : %w", err)
	}
	var result = ResultData{
		ResultChart: resultchart,
		ResultScore: calculateResult(resultchart, order, PracticeMode, nil),
	}
	return &result, nil
}

func (s *Server) createCompResult(compOrder *OrderRequest, c *fiber.Ctx) (*ResultData, error) {

	compInfo := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(compOrder.Identifier)
	err := json.Unmarshal(infoByte, compInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal chart identifier. err : %w", err)
	}
	resultchart, err := s.selectResultChart(compInfo, int(compOrder.WaitingTerm), c)
	if err != nil {
		return nil, fmt.Errorf("cannot select result chart. err : %w", err)
	}
	var result = ResultData{
		ResultChart: resultchart,
		ResultScore: calculateResult(resultchart, compOrder, CompetitionMode, compInfo),
	}

	originchart, err := s.selectStageChart(compInfo.Name, compInfo.Interval, compInfo.RefTimestamp, c)
	if err != nil {
		return nil, fmt.Errorf("cannot select origin competition chart. err : %w", err)
	}

	result.OriginChart = originchart

	result.ResultScore.Name = result.ResultChart.PData[0].Name
	result.ResultScore.Entrytime = utilities.EntryTimeFormatter(originchart.PData[len(originchart.PData)-1].Time)
	return &result, nil
}

func calculateResult(resultchart *CandleData, order *OrderRequest, mode int8, info *utilities.IdentificationData) *ResultScore {
	var (
		roe      float64
		pnl      float64
		endIdx   int
		endPrice float64
	)

	if mode == CompetitionMode {
		order.EntryPrice = math.Floor(order.EntryPrice/info.PriceFactor*decimal) / decimal
		order.ProfitPrice = math.Floor(order.ProfitPrice/info.PriceFactor*decimal) / decimal
		order.LossPrice = math.Floor(order.LossPrice/info.PriceFactor*decimal) / decimal
		order.Quantity = math.Floor(order.Quantity*info.PriceFactor*decimal) / decimal
	}

	for idx, candle := range resultchart.PData {
		if *order.IsLong {
			if candle.High >= order.ProfitPrice {
				roe = float64(order.QuantityRate/100) * (order.ProfitPrice - order.EntryPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.High
				break
			}
			if candle.Low <= order.LossPrice {
				roe = float64(order.QuantityRate/100) * (order.LossPrice - order.EntryPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Low
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = float64(order.QuantityRate/100) * (candle.Close - order.EntryPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Close
				break
			}
		} else {
			if candle.Low <= order.ProfitPrice {
				roe = float64(order.QuantityRate/100) * (order.EntryPrice - order.ProfitPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Low
				break
			}
			if candle.High >= order.LossPrice {
				roe = float64(order.QuantityRate/100) * (order.EntryPrice - order.LossPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.High
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = float64(order.QuantityRate/100) * (order.EntryPrice - candle.Close) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Close
				break
			}
		}
	}
	pnl = roe * float64(100/order.QuantityRate) * order.EntryPrice * order.Quantity

	resultInfo := ResultScore{
		Stage:      order.Stage,
		Name:       order.Name,
		Entrytime:  order.Entrytime,
		Leverage:   order.Leverage,
		EntryPrice: order.EntryPrice,
		EndPrice:   (math.Floor(endPrice*decimal) / decimal),
		OutTime:    int32(endIdx),
		Roe:        (math.Floor(roe*decimal*float64(order.Leverage)) / 100),
		Pnl:        math.Floor(pnl*decimal) / decimal,
		Commission: math.Floor(commissionRate*order.EntryPrice*order.Quantity*decimal) / decimal,
	}
	if order.Balance+resultInfo.Pnl-resultInfo.Commission < 1 {
		resultInfo.Isliquidated = true
	}
	return &resultInfo
}

func (s *Server) selectResultChart(info *utilities.IdentificationData, waitingTerm int, c *fiber.Ctx) (*CandleData, error) {
	var cdd CandleData

	switch info.Interval {
	case db.OneD:
		candles, err := s.store.Get1dResult(c.Context(), db.Get1dResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.OneD, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1dSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hResult(c.Context(), db.Get4hResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.FourH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles4hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hResult(c.Context(), db.Get1hResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.OneH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mResult(c.Context(), db.Get15mResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.FifM, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles15mSlice(candles)
		cdd = (&cs).InitCandleData()
	}
	if cdd.PData == nil || cdd.VData == nil {
		return nil, ErrGetResultChart
	}

	return &cdd, nil
}
