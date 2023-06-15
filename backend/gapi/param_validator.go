package gapi

import (
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"fmt"
)

func validateOrderRequest(o *pb.OrderRequest) error {
	if o.IsLong {
		if !(o.EntryPrice < o.ProfitPrice && o.LossPrice < o.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon : long, entry price : %f", o.EntryPrice)
		}
	} else {
		if !(o.EntryPrice > o.ProfitPrice && o.LossPrice > o.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon : short, entry price : %f", o.EntryPrice)
		}
	}

	if (o.Balance * float64(o.Leverage)) < (o.Quantity * o.EntryPrice) {
		return fmt.Errorf("invalid order. check your balance. order amound : %.5f, limit amount : %.5f ", o.Quantity*o.EntryPrice, o.Balance*float64(o.Leverage))
	}
	if err := o.ValidateAll(); err != nil {
		return err
	}
	return nil
}

func validateAndGetNextPair(g *pb.GetCandlesRequest, pairs []string) (error, string, int) {
	if err := g.ValidateAll(); err != nil {
		return err, "", 10
	}
	history := utilities.Splitnames(g.Names)
	prevStage := len(history)
	if prevStage >= finalstage {
		return fmt.Errorf("invalid current stage"), "", 10
	}
	return nil, utilities.FindDiffPair(pairs, history), prevStage
}
