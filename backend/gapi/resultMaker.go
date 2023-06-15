package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"context"
	"encoding/json"
	"fmt"
	"math"
)

const (
	commissionRate = 0.0002
)

func (s *Server) createPracResult(order *pb.OrderRequest, c context.Context) (*pb.OrderResponse, error) {
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
	var result = pb.OrderResponse{
		ResultChart: resultchart,
		Score:       calculateResult(resultchart, order, practice, nil),
	}
	return &result, nil
}

func (s *Server) createCompResult(compOrder *pb.OrderRequest, c context.Context) (*pb.OrderResponse, error) {

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
	var result = pb.OrderResponse{
		ResultChart: resultchart,
		Score:       calculateResult(resultchart, compOrder, competition, compInfo),
	}

	originchart, err := s.selectStageChart(compInfo.Name, compInfo.Interval, compInfo.RefTimestamp, c)
	if err != nil {
		return nil, fmt.Errorf("cannot select origin competition chart. err : %w", err)
	}

	result.OriginChart = originchart

	result.Score.Name = result.ResultChart.PData[0].Name
	result.Score.Entrytime = utilities.EntryTimeFormatter(originchart.PData[len(originchart.PData)-1].Time)
	return &result, nil
}

func calculateResult(resultchart *pb.CandleData, order *pb.OrderRequest, mode string, info *utilities.IdentificationData) *pb.Score {
	var (
		roe      float64
		pnl      float64
		endIdx   int
		endPrice float64
	)

	if mode == competition {
		order.EntryPrice = math.Floor(order.EntryPrice/info.PriceFactor*common.Decimal) / common.Decimal
		order.ProfitPrice = math.Floor(order.ProfitPrice/info.PriceFactor*common.Decimal) / common.Decimal
		order.LossPrice = math.Floor(order.LossPrice/info.PriceFactor*common.Decimal) / common.Decimal
		order.Quantity = math.Floor(order.Quantity*info.PriceFactor*common.Decimal) / common.Decimal
	}

	for idx, candle := range resultchart.PData {
		if order.IsLong {
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

	score := pb.Score{
		Stage:      order.Stage,
		Name:       order.Name,
		Leverage:   order.Leverage,
		EntryPrice: order.EntryPrice,
		EndPrice:   (math.Floor(endPrice*common.Decimal) / common.Decimal),
		OutTime:    int32(endIdx),
		Roe:        (math.Floor(roe*common.Decimal*float64(order.Leverage)) / 100),
		Pnl:        math.Floor(pnl*common.Decimal) / common.Decimal,
		Commission: math.Floor(commissionRate*order.EntryPrice*order.Quantity*common.Decimal) / common.Decimal,
	}
	if order.Balance+score.Pnl-score.Commission < 1 {
		score.IsLiquidated = true
	}
	return &score
}

func (s *Server) selectResultChart(info *utilities.IdentificationData, waitingTerm int, c context.Context) (*pb.CandleData, error) {
	cdd := new(pb.CandleData)

	switch info.Interval {
	case db.OneD:
		candles, err := s.store.Get1dResult(c, db.Get1dResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateTerm(db.OneD, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1dSlice(candles)
		cdd = cs.InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hResult(c, db.Get4hResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateTerm(db.FourH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles4hSlice(candles)
		cdd = cs.InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hResult(c, db.Get1hResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateTerm(db.OneH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1hSlice(candles)
		cdd = cs.InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mResult(c, db.Get15mResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateTerm(db.FifM, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles15mSlice(candles)
		cdd = cs.InitCandleData()
	}
	if cdd.PData == nil || cdd.VData == nil {
		return nil, ErrGetResultChart
	}

	return cdd, nil
}
