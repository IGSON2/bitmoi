package api

import (
	db "bitmoi/backend/db/sqlc"
	"context"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Score struct {
	User       string  `json:"user"`
	Scoreid    string  `json:"scoreid"`
	Stage      int     `json:"stage"`
	Pairname   string  `json:"pairname"`
	Entrytime  string  `json:"entrytime"`
	Position   string  `json:"position"`
	Leverage   int     `json:"leverage"`
	Outtime    int     `json:"outtime"`
	Entryprice float64 `json:"entryprice"`
	Endprice   float64 `json:"endprice"`
	Pnl        float64 `json:"pnl"`
	Roe        float64 `json:"roe"`
	Pkey       string  `json:"-"`
}

type ScoreList struct {
	ScoreList []Score `json:"scorelist"`
}

type RankedData struct {
	User         string  `json:"user"`
	Displayname  string  `json:"displayname"`
	PhotoUrl     string  `json:"photourl"`
	Scoreid      string  `json:"scoreid"`
	FinalBalance float64 `json:"balance"`
	Comment      string  `json:"comment"`
	scoreId      string  `json:"-"`
}

type RankingBoard struct {
	RankingBoard []RankedData `json:"rankingBoard"`
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
	long  = "LONG"
	short = "SHORT"
)

var (
	ErrNotUpdatedScore = errors.New("failed to update rank due to low score")
)

func (s *Server) insertUserScore(o *OrderStruct, r *ResultScore) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var position string
	if o.IsLong {
		position = long
	} else {
		position = short
	}

	_, err := s.store.InsertScore(ctx, db.InsertScoreParams{
		ScoreID:    o.ScoreId,
		UserID:     o.UserId,
		Stage:      int32(o.Stage),
		Pairname:   r.Name,
		Entrytime:  r.Entrytime,
		Position:   position,
		Leverage:   int32(o.Leverage),
		Outtime:    int32(r.OutTime),
		Entryprice: r.EntryPrice,
		Endprice:   r.EndPrice,
		Pnl:        r.Pnl,
		Roe:        r.Roe,
	})

	return err
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

func (s *Server) insertScoreToRankBoard(params db.InsertRankParams) error {
	r, err := s.getRankByUserID(params.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = s.store.InsertRank(context.Background(), params)
		}
		return err
	}
	if r.FinalBalance > params.FinalBalance {
		return ErrNotUpdatedScore
	} else {
		_, err = s.store.InsertRank(context.Background(), params)
	}
	return err
}
