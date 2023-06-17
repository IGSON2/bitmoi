package gapi

import (
	"bitmoi/backend/api"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"fmt"
	"math"
	"strings"
)

const (
	long  = "LONG"
	short = "SHORT"
)

var modeValidations = []string{practice, competition}

func validateOrderRequest(o *pb.OrderRequest) error {
	limit := math.Pow(float64(o.Leverage), float64(-1))
	if o.IsLong {
		if !(o.EntryPrice < o.ProfitPrice && o.LossPrice < o.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon: %s, entry price: %f", long, o.EntryPrice)
		}
		if (o.EntryPrice-o.LossPrice)/o.EntryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				long, o.EntryPrice, o.LossPrice, o.Leverage, common.CeilDecimal(o.EntryPrice*(1-limit)))
		}
	} else {
		if !(o.EntryPrice > o.ProfitPrice && o.LossPrice > o.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon : %s, entry price : %f", short, o.EntryPrice)
		}
		if (o.LossPrice-o.EntryPrice)/o.EntryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				short, o.EntryPrice, o.LossPrice, o.Leverage, common.FloorDecimal(o.EntryPrice*(1+limit)))
		}
	}

	if (o.Balance * float64(o.Leverage)) < (o.Quantity * o.EntryPrice) {
		return fmt.Errorf("invalid order. check your balance. order amount : %.5f, limit amount : %.5f ", o.Quantity*o.EntryPrice, o.Balance*float64(o.Leverage))
	}
	if errs := utilities.ValidateStruct(convertOrderRequest(o)); errs != nil {
		return fmt.Errorf("request validation error: %s", errs.Error())
	}
	return nil
}

func validateAndGetNextPair(g *pb.GetCandlesRequest, pairs []string) (string, int, error) {
	if !strings.EqualFold(modeValidations[0], g.Mode) && !strings.EqualFold(modeValidations[1], g.Mode) {
		return "", 10, fmt.Errorf("mode must be specified between practice or competition")
	}
	history := utilities.Splitnames(g.Names)
	prevStage := len(history)
	if prevStage >= finalstage {
		return "", 10, fmt.Errorf("invalid current stage")
	}
	return utilities.FindDiffPair(pairs, history), prevStage, nil
}

func convertOrderRequest(r *pb.OrderRequest) *api.OrderRequest {
	return &api.OrderRequest{
		Mode:        r.Mode,
		UserId:      r.UserId,
		Name:        r.Name,
		Stage:       r.Stage,
		IsLong:      &r.IsLong,
		EntryPrice:  r.EntryPrice,
		Quantity:    r.Quantity,
		ProfitPrice: r.ProfitPrice,
		LossPrice:   r.LossPrice,
		Leverage:    r.Leverage,
		Balance:     r.Balance,
		Identifier:  r.Identifier,
		ScoreId:     r.ScoreId,
		WaitingTerm: r.WaitingTerm,
	}
}
