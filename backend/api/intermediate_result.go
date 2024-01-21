package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"context"
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
		Entrytime:  utilities.EntryTimeFormatter(info.RefTimestamp),
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

func (s *Server) calculateAfterInterResult(resultchart *CandleData, order *InterScoreRequest, info *utilities.IdentificationData) *AfterScore {
	var (
		roe          float64
		maxRoe       float64 = -100
		minRoe       float64 = 100
		endTimestamp int64
		aInfo        = &AfterScore{0, 100, -100}
	)

	if !info.IsPracticeMode() {
		order.EntryPrice = common.FloorDecimal(order.EntryPrice / info.PriceFactor)
		order.ProfitPrice = common.FloorDecimal(order.ProfitPrice / info.PriceFactor)
		order.LossPrice = common.FloorDecimal(order.LossPrice / info.PriceFactor)
		order.Quantity = common.FloorDecimal(order.Quantity * info.PriceFactor)
	}

	// 상위 단위에서 검색에 성공했을 경우, 하위 단위에서 축소 검색을 시도한다.
	getDetailedChart := func(from, to int64) *CandleData {
		intv := db.GetIntervalFromRange(from, to-1)
		if intv == "" {
			return nil
		}
		cdd, err := s.getResultByInterval(intv, info.Name, to-1, 10, context.Background())
		if err != nil {
			s.logger.Error().Err(err).Msgf("cannot get detailed chart. name : %s, interval : %s, from : %d, to : %d", info.Name, intv, from, to)
			return nil
		}
		return cdd
	}

	maxQuantity := (float64(order.Leverage) * order.Balance) / order.EntryPrice
	levQuanRate := float64(order.Leverage) * (order.Quantity / maxQuantity)
	for _, candle := range resultchart.PData {
		timeGap := resultchart.PData[1].Time - resultchart.PData[0].Time
		if *order.IsLong {
			roe = levQuanRate * ((candle.Close - order.EntryPrice) / order.EntryPrice)
			if roe > maxRoe {
				maxRoe = roe
			}
			if roe < minRoe {
				minRoe = roe
			}

			if candle.High >= order.ProfitPrice {
				mcdd := getDetailedChart(candle.Time-timeGap, candle.Time)
				if mcdd == nil {
					endTimestamp = candle.Time
					break
				}
				aInfo = s.calculateAfterInterResult(mcdd, order, info)
				break
			}
			if candle.Low <= order.LossPrice {
				mcdd := getDetailedChart(candle.Time-timeGap, candle.Time)
				if mcdd == nil {
					endTimestamp = candle.Time
					break
				}
				aInfo = s.calculateAfterInterResult(mcdd, order, info)
				break
			}
		} else {
			roe = levQuanRate * ((order.EntryPrice - candle.Close) / order.EntryPrice)
			if roe > maxRoe {
				maxRoe = roe
			}
			if roe < minRoe {
				minRoe = roe
			}

			if candle.Low <= order.ProfitPrice {
				mcdd := getDetailedChart(candle.Time-timeGap, candle.Time)
				if mcdd == nil {
					endTimestamp = candle.Time
					break
				}
				aInfo = s.calculateAfterInterResult(mcdd, order, info)
				break
			}
			if candle.High >= order.LossPrice {
				mcdd := getDetailedChart(candle.Time-timeGap, candle.Time)
				if mcdd == nil {
					endTimestamp = candle.Time
					break
				}
				aInfo = s.calculateAfterInterResult(mcdd, order, info)
				break
			}
		}
	}
	if aInfo.MaxRoe > maxRoe {
		maxRoe = aInfo.MaxRoe
	}
	if aInfo.MinRoe < minRoe {
		minRoe = aInfo.MinRoe
	}
	if aInfo.ClosedTime > endTimestamp {
		endTimestamp = aInfo.ClosedTime
	}

	afterResultInfo := AfterScore{
		ClosedTime: endTimestamp,
		MaxRoe:     common.FloorDecimal(maxRoe),
		MinRoe:     common.FloorDecimal(minRoe),
	}
	return &afterResultInfo
}

func scoreToInterResult(s *db.PracScore) *InterMediateResult {
	return &InterMediateResult{
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
