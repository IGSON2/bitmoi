package api

import (
	db "bitmoi/backend/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	long                   = "LONG"
	short                  = "SHORT"
	defaultBalance float64 = 1000
	rankRows               = 10
	myscoreRows            = 15
)

var (
	ErrNotUpdatedScore    = errors.New("failed to update rank due to low score")
	ErrLiquidation        = errors.New("you were liquidated due to insufficient balance")
	ErrInvalidStageLength = errors.New("insufficient number of stages cleared")
)

func (s *Server) insertUserScore(o *ScoreRequest, r *OrderResult, c *fiber.Ctx) error {
	var position string
	if *o.IsLong {
		position = long
	} else {
		position = short
	}

	_, err := s.store.InsertScore(c.Context(), db.InsertScoreParams{
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

func (s *Server) getScoreToStage(o *ScoreRequest, c *fiber.Ctx) error {
	i, err := s.store.GetScoreToStage(c.Context(), db.GetScoreToStageParams{
		ScoreID: o.ScoreId,
		UserID:  o.UserId,
		Stage:   o.Stage,
	})
	if err != nil {
		return fmt.Errorf("cannot get score to current stage, err: %w", err)
	}
	totalScore, ok := i.(float64)
	if !ok {
		return fmt.Errorf("cannot assertion totalScore to float, err: %w", err)
	}

	if defaultBalance+totalScore <= 0 {
		return ErrLiquidation
	}
	return nil
}

func (s *Server) getMyscores(userId string, pages int32, c *fiber.Ctx) ([]db.Score, error) {
	return s.store.GetScoresByUserID(c.Context(), db.GetScoresByUserIDParams{
		UserID: userId,
		Limit:  myscoreRows,
		Offset: (pages - 1) * myscoreRows,
	})
}

func (s *Server) sendMoreInfo(scoreId, userId string, c *fiber.Ctx) ([]db.Score, error) {
	return s.store.GetScoresByScoreID(c.Context(), db.GetScoresByScoreIDParams{
		ScoreID: scoreId,
		UserID:  userId,
	})
}

func (s *Server) getAllRanks(pages int32, c *fiber.Ctx) ([]db.RankingBoard, error) {
	return s.store.GetAllRanks(context.Background(), db.GetAllRanksParams{
		Limit:  rankRows,
		Offset: (pages - 1) * rankRows,
	})
}

func (s *Server) getRankByUserID(userId string) (db.RankingBoard, error) {
	return s.store.GetRankByUserID(context.Background(), userId)
}

func (s *Server) insertScoreToRankBoard(req *RankInsertRequest, c *fiber.Ctx) error {
	length, err := s.store.GetStageLenByScoreID(c.Context(), db.GetStageLenByScoreIDParams{
		ScoreID: req.ScoreId,
		UserID:  req.UserId,
	})
	if err != nil {
		return err
	} else if length != finalstage {
		return ErrInvalidStageLength
	}

	t, err := s.store.GetScoreToStage(c.Context(), db.GetScoreToStageParams{
		ScoreID: req.ScoreId,
		UserID:  req.UserId,
		Stage:   finalstage,
	})
	if err != nil {
		return err
	}
	totalScore, ok := t.(float64)
	if !ok {
		return fmt.Errorf("cannot assign totalscore")
	}

	r, err := s.getRankByUserID(req.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = s.store.InsertRank(c.Context(), db.InsertRankParams{
				UserID:       r.UserID,
				ScoreID:      r.ScoreID,
				DisplayName:  r.DisplayName,
				Comment:      r.Comment,
				FinalBalance: totalScore,
			})
		}
		return err
	}

	if r.FinalBalance > totalScore {
		return ErrNotUpdatedScore
	} else {
		_, err = s.store.UpdateUserRank(c.Context(), db.UpdateUserRankParams{
			UserID:       r.UserID,
			ScoreID:      r.ScoreID,
			DisplayName:  r.DisplayName,
			Comment:      r.Comment,
			FinalBalance: totalScore,
		})
	}
	return err
}
