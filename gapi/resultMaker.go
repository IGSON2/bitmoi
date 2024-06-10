package gapi

import (
	db "bitmoi/db/sqlc"
	"bitmoi/gapi/pb"
	"bitmoi/utilities"
	"bitmoi/utilities/common"
	"context"
	"encoding/json"
	"fmt"
)

const (
	commissionRate = 0.0002
)

func (s *Server) createPracResult(order *pb.ScoreRequest, c context.Context) (*pb.ScoreResponse, error) {
	pracInfo := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(order.Identifier)
	err := json.Unmarshal(infoByte, pracInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal chart identifier. err : %w", err)
	}
	if !pracInfo.IsPracticeMode() {
		return nil, fmt.Errorf("factor is not set to 0, mode: practice")
	}
	resultchart, err := s.selectResultChart(pracInfo, int(order.WaitingTerm), c)
	if err != nil {
		return nil, fmt.Errorf("cannot select result chart. err : %w", err)
	}
	var result = pb.ScoreResponse{
		ResultChart: resultchart,
		Score:       calculateResult(resultchart, order, practice, nil),
	}
	return &result, nil
}

func (s *Server) createCompResult(compOrder *pb.ScoreRequest, c context.Context) (*pb.ScoreResponse, error) {

	compInfo := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(compOrder.Identifier)
	err := json.Unmarshal(infoByte, compInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal chart identifier. err : %w", err)
	}
	if compInfo.IsPracticeMode() {
		return nil, fmt.Errorf("factor is set to 0, mode: competition")
	}
	resultchart, err := s.selectResultChart(compInfo, int(compOrder.WaitingTerm), c)
	if err != nil {
		return nil, fmt.Errorf("cannot select result chart. err : %w", err)
	}
	var result = pb.ScoreResponse{
		ResultChart: resultchart,
		Score:       calculateResult(resultchart, compOrder, competition, compInfo),
	}

	originchart, err := s.selectStageChart(compInfo.Name, compInfo.Interval, compInfo.RefTimestamp, c)
	if err != nil {
		return nil, fmt.Errorf("cannot select origin competition chart. err : %w", err)
	}

	result.OriginChart = originchart

	result.Score.Name = compInfo.Name
	result.Score.Entrytime = utilities.EntryTimeFormatter(originchart.PData[0].Time)
	return &result, nil
}

func calculateResult(resultchart *pb.CandleData, order *pb.ScoreRequest, mode string, info *utilities.IdentificationData) *pb.Score {
	var (
		roe      float64
		pnl      float64
		endIdx   int
		endPrice float64
	)

	if mode == competition {
		order.EntryPrice = common.FloorDecimal(order.EntryPrice / info.PriceFactor)
		order.ProfitPrice = common.FloorDecimal(order.ProfitPrice / info.PriceFactor)
		order.LossPrice = common.FloorDecimal(order.LossPrice / info.PriceFactor)
		order.Quantity = common.FloorDecimal(order.Quantity * info.PriceFactor)
	}

	maxQuantity := (float64(order.Leverage) * order.Balance) / order.EntryPrice
	levQuanRate := float64(order.Leverage) * (order.Quantity / maxQuantity)
	for idx, candle := range resultchart.PData {
		if order.IsLong {
			if candle.High >= order.ProfitPrice {
				roe = levQuanRate * (order.ProfitPrice - order.EntryPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.High
				break
			}
			if candle.Low <= order.LossPrice {
				roe = levQuanRate * (order.LossPrice - order.EntryPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Low
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = levQuanRate * (candle.Close - order.EntryPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Close
				break
			}
		} else {
			if candle.Low <= order.ProfitPrice {
				roe = levQuanRate * (order.EntryPrice - order.ProfitPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Low
				break
			}
			if candle.High >= order.LossPrice {
				roe = levQuanRate * (order.EntryPrice - order.LossPrice) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.High
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = levQuanRate * (order.EntryPrice - candle.Close) / order.EntryPrice
				endIdx = idx + 1
				endPrice = candle.Close
				break
			}
		}
	}
	pnl = (roe * order.Balance)

	score := pb.Score{
		Stage:      order.Stage,
		Name:       order.Name,
		Leverage:   order.Leverage,
		EntryPrice: order.EntryPrice,
		Entrytime:  utilities.EntryTimeFormatter(resultchart.PData[0].Time - (resultchart.PData[1].Time - resultchart.PData[0].Time)),
		EndPrice:   common.FloorDecimal(endPrice),
		OutTime:    int64(endIdx),
		Roe:        common.FloorDecimal(roe * 100),
		Pnl:        common.FloorDecimal(pnl),
		Commission: common.FloorDecimal(commissionRate * order.EntryPrice * order.Quantity),
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
		candles, err := s.store.Get1dResult(c, db.Get1dResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateWaitingTerm(db.OneD, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1dSlice(candles)
		cdd = cs.InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hResult(c, db.Get4hResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateWaitingTerm(db.FourH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles4hSlice(candles)
		cdd = cs.InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hResult(c, db.Get1hResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateWaitingTerm(db.OneH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1hSlice(candles)
		cdd = cs.InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mResult(c, db.Get15mResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: int32(db.CalculateWaitingTerm(db.FifM, waitingTerm))})
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
