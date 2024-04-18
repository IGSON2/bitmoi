package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"context"
	"fmt"
)

func (s *Server) insertUserScore(o *pb.ScoreRequest, r *pb.Score, c context.Context) error {
	var position string
	if o.IsLong {
		position = "long"
	} else {
		position = "short"
	}

	_, err := s.store.InsertPracScore(c, db.InsertPracScoreParams{
		ScoreID:    o.ScoreId,
		UserID:     o.UserId,
		Stage:      int8(o.Stage),
		Pairname:   r.Name,
		Entrytime:  r.Entrytime,
		Position:   position,
		Leverage:   int8(o.Leverage),
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
