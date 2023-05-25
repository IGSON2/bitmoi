package api

import (
	db "bitmoi/backend/db/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

type RegisterRank struct {
	User        string `json:"user"`
	Displayname string `json:"displayname"`
	Scoreid     string `json:"scoreid"`
	Comment     string `json:"comment"`
}

type MoreInfoData struct {
	AvgLev     float64       `json:"avglev"`
	AvgPnl     float64       `json:"avgpnl"`
	AvgRoe     float64       `json:"avgroe"`
	StageArray []StageSimple `json:"stagearray"`
}

type StageSimple struct {
	Name string  `json:"name"`
	Date string  `json:"date"`
	Roe  float64 `json:"roe"`
}

const (
	long                   = "LONG"
	short                  = "SHORT"
	defaultBalance float64 = 1000
)

var (
	ErrNotUpdatedScore    = errors.New("failed to update rank due to low score")
	ErrLiquidation        = errors.New("you were liquidated due to insufficient balance")
	ErrInvalidStageLength = errors.New("ㅑnsufficient number of stages cleared")
)

func (s *Server) insertUserScore(o *OrderRequest, r *ResultScore, c *fiber.Ctx) error {
	var position string
	if o.IsLong {
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

func (s *Server) getScoreToStage(o *OrderRequest, c *fiber.Ctx) error {
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

func (s *Server) getScoresByUserID(userId string, limit int32) ([]db.Score, error) {
	return s.store.GetScoresByUserID(context.Background(), db.GetScoresByUserIDParams{
		UserID: userId,
		Limit:  limit,
	})
}

func (s *Server) getScoresByScoreID(scoreId, userId string, limit int32) ([]db.Score, error) {
	return s.store.GetScoresByScoreID(context.Background(), db.GetScoresByScoreIDParams{
		ScoreID: scoreId,
		UserID:  userId,
		Limit:   limit,
	})
}

func (s *Server) getAllRanks(limit int32) ([]db.RankingBoard, error) {
	return s.store.GetAllRanks(context.Background(), limit)
}

func (s *Server) getRankByUserID(userId string) (db.RankingBoard, error) {
	return s.store.GetRankByUserID(context.Background(), userId)
}

func (s *Server) insertScoreToRankBoard(params db.InsertRankParams, c *fiber.Ctx) error {
	length, err := s.store.GetStageLenByScoreID(c.Context(), db.GetStageLenByScoreIDParams{
		ScoreID: params.ScoreID,
		UserID:  params.UserID,
	})
	if err != nil {
		return err
	} else if length != finalstage {
		return ErrInvalidStageLength
	}

	r, err := s.getRankByUserID(params.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = s.store.InsertRank(c.Context(), params)
		}
		return err
	}

	if r.FinalBalance > params.FinalBalance {
		return ErrNotUpdatedScore
	} else {
		_, err = s.store.InsertRank(c.Context(), params)
	}
	return err
}
