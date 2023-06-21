package gapi

import (
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"fmt"
	"math"
)

const (
	long  = "LONG"
	short = "SHORT"
)

func validateOrderRequest(r *pb.ScoreRequest) error {
	if err := r.ValidateAll(); err != nil {
		return fmt.Errorf("request validation error: %s", err.Error())
	}

	limit := math.Pow(float64(r.Leverage), float64(-1))
	if r.IsLong {
		if !(r.EntryPrice < r.ProfitPrice && r.LossPrice < r.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon: %s, entry price: %f", long, r.EntryPrice)
		}
		if (r.EntryPrice-r.LossPrice)/r.EntryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				long, r.EntryPrice, r.LossPrice, r.Leverage, common.CeilDecimal(r.EntryPrice*(1-limit)))
		}
	} else {
		if !(r.EntryPrice > r.ProfitPrice && r.LossPrice > r.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon : %s, entry price : %f", short, r.EntryPrice)
		}
		if (r.LossPrice-r.EntryPrice)/r.EntryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				short, r.EntryPrice, r.LossPrice, r.Leverage, common.FloorDecimal(r.EntryPrice*(1+limit)))
		}
	}

	if (r.Balance * float64(r.Leverage)) < (r.Quantity * r.EntryPrice) {
		return fmt.Errorf("invalid order. check your balance. order amount : %.5f, limit amount : %.5f ", r.Quantity*r.EntryPrice, r.Balance*float64(r.Leverage))
	}
	return nil
}

func validateGetCandlesRequest(r *pb.CandlesRequest, pairs []string) (next string, prevStage int, err error) {
	if err := r.ValidateAll(); err != nil {
		return "", 10, err
	}
	history := utilities.SplitPairnames(r.Names)
	prevStage = len(history)
	if prevStage >= finalstage {
		return "", 10, fmt.Errorf("invalid current stage")
	}
	next = utilities.FindDiffPair(pairs, history)
	return next, prevStage, nil
}

func validateAnotherIntervalRequest(r *pb.AnotherIntervalRequest) error {
	if err := r.ValidateAll(); err != nil {
		return err
	}
	return nil
}
