// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: comp_score.sql

package db

import (
	"context"
	"database/sql"
)

const getCompScore = `-- name: GetCompScore :one
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, remain_balance, created_at FROM comp_score
WHERE user_id = ? AND score_id = ? AND stage = ?
`

type GetCompScoreParams struct {
	UserID  string `json:"user_id"`
	ScoreID string `json:"score_id"`
	Stage   int32  `json:"stage"`
}

func (q *Queries) GetCompScore(ctx context.Context, arg GetCompScoreParams) (CompScore, error) {
	row := q.db.QueryRowContext(ctx, getCompScore, arg.UserID, arg.ScoreID, arg.Stage)
	var i CompScore
	err := row.Scan(
		&i.ScoreID,
		&i.UserID,
		&i.Stage,
		&i.Pairname,
		&i.Entrytime,
		&i.Position,
		&i.Leverage,
		&i.Outtime,
		&i.Entryprice,
		&i.Quantity,
		&i.Endprice,
		&i.Pnl,
		&i.Roe,
		&i.RemainBalance,
		&i.CreatedAt,
	)
	return i, err
}

const getCompScoreToStage = `-- name: GetCompScoreToStage :one
SELECT SUM(pnl) FROM comp_score
WHERE score_id = ? AND user_id = ? AND stage <= ?
`

type GetCompScoreToStageParams struct {
	ScoreID string `json:"score_id"`
	UserID  string `json:"user_id"`
	Stage   int32  `json:"stage"`
}

func (q *Queries) GetCompScoreToStage(ctx context.Context, arg GetCompScoreToStageParams) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getCompScoreToStage, arg.ScoreID, arg.UserID, arg.Stage)
	var sum interface{}
	err := row.Scan(&sum)
	return sum, err
}

const getCompScoresByScoreID = `-- name: GetCompScoresByScoreID :many
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, remain_balance, created_at FROM comp_score
WHERE score_id = ? AND user_id = ?
`

type GetCompScoresByScoreIDParams struct {
	ScoreID string `json:"score_id"`
	UserID  string `json:"user_id"`
}

func (q *Queries) GetCompScoresByScoreID(ctx context.Context, arg GetCompScoresByScoreIDParams) ([]CompScore, error) {
	rows, err := q.db.QueryContext(ctx, getCompScoresByScoreID, arg.ScoreID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CompScore{}
	for rows.Next() {
		var i CompScore
		if err := rows.Scan(
			&i.ScoreID,
			&i.UserID,
			&i.Stage,
			&i.Pairname,
			&i.Entrytime,
			&i.Position,
			&i.Leverage,
			&i.Outtime,
			&i.Entryprice,
			&i.Quantity,
			&i.Endprice,
			&i.Pnl,
			&i.Roe,
			&i.RemainBalance,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCompScoresByUserID = `-- name: GetCompScoresByUserID :many
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, remain_balance, created_at FROM comp_score
WHERE user_id = ?
ORDER BY score_id DESC 
LIMIT ?
OFFSET ?
`

type GetCompScoresByUserIDParams struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetCompScoresByUserID(ctx context.Context, arg GetCompScoresByUserIDParams) ([]CompScore, error) {
	rows, err := q.db.QueryContext(ctx, getCompScoresByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []CompScore{}
	for rows.Next() {
		var i CompScore
		if err := rows.Scan(
			&i.ScoreID,
			&i.UserID,
			&i.Stage,
			&i.Pairname,
			&i.Entrytime,
			&i.Position,
			&i.Leverage,
			&i.Outtime,
			&i.Entryprice,
			&i.Quantity,
			&i.Endprice,
			&i.Pnl,
			&i.Roe,
			&i.RemainBalance,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCompStageLenByScoreID = `-- name: GetCompStageLenByScoreID :one
SELECT COUNT(stage) FROM comp_score
WHERE score_id = ? AND user_id = ?
`

type GetCompStageLenByScoreIDParams struct {
	ScoreID string `json:"score_id"`
	UserID  string `json:"user_id"`
}

func (q *Queries) GetCompStageLenByScoreID(ctx context.Context, arg GetCompStageLenByScoreIDParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCompStageLenByScoreID, arg.ScoreID, arg.UserID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const insertCompScore = `-- name: InsertCompScore :execresult
INSERT INTO comp_score (
    score_id,
    user_id,
    stage,
    pairname,
    entrytime,
    position,
    leverage,
    outtime,
    entryprice,
    quantity,
    endprice,
    pnl,
    roe,
    remain_balance
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

type InsertCompScoreParams struct {
	ScoreID       string  `json:"score_id"`
	UserID        string  `json:"user_id"`
	Stage         int32   `json:"stage"`
	Pairname      string  `json:"pairname"`
	Entrytime     string  `json:"entrytime"`
	Position      string  `json:"position"`
	Leverage      int32   `json:"leverage"`
	Outtime       int64   `json:"outtime"`
	Entryprice    float64 `json:"entryprice"`
	Quantity      float64 `json:"quantity"`
	Endprice      float64 `json:"endprice"`
	Pnl           float64 `json:"pnl"`
	Roe           float64 `json:"roe"`
	RemainBalance float64 `json:"remain_balance"`
}

func (q *Queries) InsertCompScore(ctx context.Context, arg InsertCompScoreParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertCompScore,
		arg.ScoreID,
		arg.UserID,
		arg.Stage,
		arg.Pairname,
		arg.Entrytime,
		arg.Position,
		arg.Leverage,
		arg.Outtime,
		arg.Entryprice,
		arg.Quantity,
		arg.Endprice,
		arg.Pnl,
		arg.Roe,
		arg.RemainBalance,
	)
}

const updateCompcScore = `-- name: UpdateCompcScore :execresult
UPDATE comp_score SET pairname = ?, entrytime = ?, outtime = ?, entryprice = ?, endprice = ?, pnl = ?, roe = ?, remain_balance = ?
WHERE user_id = ? AND score_id = ? AND stage = ?
`

type UpdateCompcScoreParams struct {
	Pairname      string  `json:"pairname"`
	Entrytime     string  `json:"entrytime"`
	Outtime       int64   `json:"outtime"`
	Entryprice    float64 `json:"entryprice"`
	Endprice      float64 `json:"endprice"`
	Pnl           float64 `json:"pnl"`
	Roe           float64 `json:"roe"`
	RemainBalance float64 `json:"remain_balance"`
	UserID        string  `json:"user_id"`
	ScoreID       string  `json:"score_id"`
	Stage         int32   `json:"stage"`
}

func (q *Queries) UpdateCompcScore(ctx context.Context, arg UpdateCompcScoreParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateCompcScore,
		arg.Pairname,
		arg.Entrytime,
		arg.Outtime,
		arg.Entryprice,
		arg.Endprice,
		arg.Pnl,
		arg.Roe,
		arg.RemainBalance,
		arg.UserID,
		arg.ScoreID,
		arg.Stage,
	)
}