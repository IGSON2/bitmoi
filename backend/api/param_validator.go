package api

import db "bitmoi/backend/db/sqlc"

type validInsertRankParams struct {
	UserId       string  `json:"user_id" validate:"required,alpha"`
	ScoreId      string  `json:"score_id" validate:"required,numeric"`
	DisplayName  string  `json:"display_name" validate:"required"`
	PhotoUrl     string  `json:"photo_url"`
	Comment      string  `json:"comment" validate:"required"`
	FinalBalance float64 `json:"final_balance" validate:"required,number"`
}

func convertInsertRankParams(p db.InsertRankParams) *validInsertRankParams {
	return &validInsertRankParams{
		UserId:       p.UserID,
		ScoreId:      p.ScoreID,
		DisplayName:  p.DisplayName,
		PhotoUrl:     p.PhotoUrl,
		Comment:      p.Comment,
		FinalBalance: p.FinalBalance,
	}
}

type OrderStruct struct {
	Mode         string  `json:"mode" validate:"required,eq_ignore_case=competition,eq_ignore_case=practice"`
	UserId       string  `json:"userid" validate:"required,alpha"`
	Name         string  `json:"name"`
	Entrytime    string  `json:"entrytime"`
	Stage        int32   `json:"stage"`
	IsLong       bool    `json:"islong"`
	EntryPrice   float64 `json:"entryprice" validate:"required,number"`
	Quantity     float64 `json:"quantity" validate:"required,number"`
	QuantityRate float64 `json:"quantityrate" validate:"required,number"`
	ProfitPrice  float64 `json:"profitprice" validate:"required,number"`
	LossPrice    float64 `json:"lossprice" validate:"required,number"`
	Leverage     int32   `json:"leverage"`
	Balance      float64 `json:"balance" validate:"required,number"`
	Identifier   string  `json:"identifier,omitempty"`
	ScoreId      string  `json:"scoreid" validate:"required,numeric"`
	ResultTerm   int32   `json:"resultterm" validate:"required,number"`
}

type ResultScore struct {
	Stage        int32   `json:"stage"`
	Name         string  `json:"name"`
	Entrytime    string  `json:"entrytime"`
	Leverage     int32   `json:"leverage"`
	EntryPrice   float64 `json:"entryprice"`
	EndPrice     float64 `json:"-"`
	OutTime      int32   `json:"outtime"`
	Roe          float64 `json:"roe"`
	Pnl          float64 `json:"pnl"`
	Commission   float64 `json:"commission"`
	Isliquidated bool    `json:"isliquidated"`
}
