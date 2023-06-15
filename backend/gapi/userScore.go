package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"context"
	"fmt"
)

func (s *Server) insertUserScore(o *pb.OrderRequest, r *pb.Score, c context.Context) error {
	var position string
	if o.IsLong {
		position = "long"
	} else {
		position = "short"
	}

	_, err := s.store.InsertScore(c, db.InsertScoreParams{
		ScoreID:    o.ScoreId,
		UserID:     o.UserId,
		Stage:      o.Stage,
		Pairname:   r.Name,
		Entrytime:  r.Entrytime,
		Position:   position,
		Leverage:   o.Leverage,
		Outtime:    r.OutTime,
		Entryprice: r.EntryPrice,
		Endprice:   r.EndPrice,
		Pnl:        r.Pnl,
		Roe:        r.Roe,
	})
	if err != nil {
		return fmt.Errorf("cannot insert score, err: %w", err)
	}

	return err
}
