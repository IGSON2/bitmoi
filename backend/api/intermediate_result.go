package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
)

func calculateInterResult(resultchart *CandleData, order *InterScoreRequest, info *utilities.IdentificationData) *InterMediateResult {
	var (
		roe          float64
		pnl          float64
		endTimestamp int64
		endPrice     float64
	)

	if !info.IsPracticeMode() {
		order.EntryPrice = common.FloorDecimal(order.EntryPrice / info.PriceFactor)
		order.ProfitPrice = common.FloorDecimal(order.ProfitPrice / info.PriceFactor)
		order.LossPrice = common.FloorDecimal(order.LossPrice / info.PriceFactor)
		order.Quantity = common.FloorDecimal(order.Quantity * info.PriceFactor)
	}

	maxQuantity := (float64(order.Leverage) * order.Balance) / order.EntryPrice
	levQuanRate := float64(order.Leverage) * (order.Quantity / maxQuantity)
	for idx, candle := range resultchart.PData {
		if *order.IsLong {
			if candle.High >= order.ProfitPrice {
				roe = levQuanRate * ((order.ProfitPrice - order.EntryPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.High
				break
			}
			if candle.Low <= order.LossPrice {
				roe = levQuanRate * ((order.LossPrice - order.EntryPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.Low
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = levQuanRate * ((candle.Close - order.EntryPrice) / order.EntryPrice)
				break
			}
		} else {
			if candle.Low <= order.ProfitPrice {
				roe = levQuanRate * ((order.EntryPrice - order.ProfitPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.Low
				break
			}
			if candle.High >= order.LossPrice {
				roe = levQuanRate * ((order.EntryPrice - order.LossPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.High
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = levQuanRate * ((order.EntryPrice - candle.Close) / order.EntryPrice)
				break
			}
		}
	}
	pnl = (roe * order.Balance)

	resultInfo := InterMediateResult{
		Name:       order.Name,
		Leverage:   order.Leverage,
		EndPrice:   common.FloorDecimal(endPrice),
		OutTime:    endTimestamp,
		Roe:        common.FloorDecimal(roe * 100),
		Pnl:        common.FloorDecimal(pnl),
		Commission: common.FloorDecimal(commissionRate * order.EntryPrice * order.Quantity),
	}
	if order.Balance+resultInfo.Pnl-resultInfo.Commission < 1 {
		resultInfo.Isliquidated = true
	}
	return &resultInfo
}

func scoreToInterResult(s *db.PracScore) InterMediateResult {
	return InterMediateResult{
		Name:         s.Pairname,
		Entrytime:    s.Entrytime,
		Leverage:     s.Leverage,
		EndPrice:     s.Endprice,
		OutTime:      s.Outtime,
		Roe:          s.Roe,
		Pnl:          s.Pnl,
		Commission:   common.FloorDecimal(commissionRate * s.Entryprice * s.Quantity),
		Isliquidated: s.RemainBalance < 0,
	}
}
