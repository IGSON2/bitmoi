package api

import (
	db "bitmoi/backend/db/sqlc"
	"context"

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

func (s *Server) getScoresByUserID()  {}
func (s *Server) getScoresByScoreID() {}
func (s *Server) getRankByUserID()    {}
func (s *Server) insertScoreToRankBoard() {
	//if firt rank
	//or already rankded
}
func (s *Server) getRankedStagesByScoreID() {}
