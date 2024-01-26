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
		endTimestamp int64
		endPrice     float64
	)

	if !info.IsPracticeMode() {
		order.EntryPrice = common.FloorDecimal(order.EntryPrice / info.PriceFactor)
		order.ProfitPrice = common.FloorDecimal(order.ProfitPrice / info.PriceFactor)
		order.LossPrice = common.FloorDecimal(order.LossPrice / info.PriceFactor)
		order.Quantity = common.FloorDecimal(order.Quantity * info.PriceFactor)
	}

	deposit := order.Quantity * order.EntryPrice / float64(order.Leverage)

	// 주문 수량에 관계 없이 손익률을 제공하기 위해 변경
	// maxQuantity := (float64(order.Leverage) * order.Balance) / order.EntryPrice
	// levQuanRate := float64(order.Leverage) * (order.Quantity / maxQuantity)

	for idx, candle := range resultchart.PData {
		if *order.IsLong {
			if candle.High >= order.ProfitPrice {
				roe = float64(order.Leverage) * ((order.ProfitPrice - order.EntryPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.High
				break
			}
			if candle.Low <= order.LossPrice {
				roe = float64(order.Leverage) * ((order.LossPrice - order.EntryPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.Low
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = float64(order.Leverage) * ((candle.Close - order.EntryPrice) / order.EntryPrice)
				break
			}
		} else {
			if candle.Low <= order.ProfitPrice {
				roe = float64(order.Leverage) * ((order.EntryPrice - order.ProfitPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.Low
				break
			}
			if candle.High >= order.LossPrice {
				roe = float64(order.Leverage) * ((order.EntryPrice - order.LossPrice) / order.EntryPrice)
				endTimestamp = candle.Time
				endPrice = candle.High
				break
			}
			if idx == len(resultchart.PData)-1 {
				roe = float64(order.Leverage) * ((order.EntryPrice - candle.Close) / order.EntryPrice)
				break
			}
		}
	}
	pnl := (roe * deposit)
	commission := common.FloorDecimal(commissionRate * order.EntryPrice * order.Quantity)
	if info.IsPracticeMode() {
		commission = 0
	}

	resultInfo := InterMediateResult{
		Name:       order.Name,
		IsLong:     *order.IsLong,
		Entrytime:  utilities.EntryTimeFormatter(info.RefTimestamp),
		Leverage:   order.Leverage,
		EndPrice:   common.FloorDecimal(endPrice),
		OutTime:    endTimestamp,
		Roe:        common.FloorDecimal(roe * 100),
		Pnl:        common.FloorDecimal(pnl),
		Commission: commission,
	}
	if order.Balance+resultInfo.Pnl-resultInfo.Commission < 1 {
		resultInfo.Isliquidated = true
	}
	return &resultInfo
}

func (s *Server) calculateAfterInterResult(resultchart *CandleData, order *InterScoreRequest, info *utilities.IdentificationData) *AfterScore {
	var (
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
			// max roe는 high, min은 row값
			_maxroe := levQuanRate * ((candle.High - order.EntryPrice) / order.EntryPrice)
			_minroe := levQuanRate * ((candle.Low - order.EntryPrice) / order.EntryPrice)
			if _maxroe > maxRoe {
				maxRoe = _maxroe
			}
			if _minroe < minRoe {
				minRoe = _minroe
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
			_maxroe := levQuanRate * ((order.EntryPrice - candle.Low) / order.EntryPrice)
			_minroe := levQuanRate * ((order.EntryPrice - candle.High) / order.EntryPrice)
			if _maxroe > maxRoe {
				maxRoe = _maxroe
			}
			if _minroe < minRoe {
				minRoe = _minroe
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
