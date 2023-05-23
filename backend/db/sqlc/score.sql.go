// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: score.sql

package db

import (
	"context"
	"database/sql"
)

const getScore = `-- name: GetScore :one
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, endprice, pnl, roe FROM score
WHERE score_id = ? AND stage = ?
`

type GetScoreParams struct {
	ScoreID string `json:"score_id"`
	Stage   int32  `json:"stage"`
}

func (q *Queries) GetScore(ctx context.Context, arg GetScoreParams) (Score, error) {
	row := q.db.QueryRowContext(ctx, getScore, arg.ScoreID, arg.Stage)
	var i Score
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
		&i.Endprice,
		&i.Pnl,
		&i.Roe,
	)
	return i, err
}

const getScoresByScoreID = `-- name: GetScoresByScoreID :many
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, endprice, pnl, roe FROM score
WHERE score_id = ? AND user_id = ?
LIMIT ?
`

type GetScoresByScoreIDParams struct {
	ScoreID string `json:"score_id"`
	UserID  string `json:"user_id"`
	Limit   int32  `json:"limit"`
}

func (q *Queries) GetScoresByScoreID(ctx context.Context, arg GetScoresByScoreIDParams) ([]Score, error) {
	rows, err := q.db.QueryContext(ctx, getScoresByScoreID, arg.ScoreID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Score{}
	for rows.Next() {
		var i Score
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
			&i.Endprice,
			&i.Pnl,
			&i.Roe,
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

const getScoresByUserID = `-- name: GetScoresByUserID :many
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, endprice, pnl, roe FROM score
WHERE user_id = ?
ORDER BY score_id DESC 
LIMIT ?
`

type GetScoresByUserIDParams struct {
	UserID string `json:"user_id"`
	Limit  int32  `json:"limit"`
}

func (q *Queries) GetScoresByUserID(ctx context.Context, arg GetScoresByUserIDParams) ([]Score, error) {
	rows, err := q.db.QueryContext(ctx, getScoresByUserID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Score{}
	for rows.Next() {
		var i Score
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
			&i.Endprice,
			&i.Pnl,
			&i.Roe,
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

const insertScore = `-- name: InsertScore :execresult
INSERT INTO score (
    score_id,
    user_id,
    stage,
    pairname,
    entrytime,
    position,
    leverage,
    outtime,
    entryprice,
    endprice,
    pnl,
    roe
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

type InsertScoreParams struct {
	ScoreID    string  `json:"score_id"`
	UserID     string  `json:"user_id"`
	Stage      int32   `json:"stage"`
	Pairname   string  `json:"pairname"`
	Entrytime  string  `json:"entrytime"`
	Position   string  `json:"position"`
	Leverage   int32   `json:"leverage"`
	Outtime    int32   `json:"outtime"`
	Entryprice float64 `json:"entryprice"`
	Endprice   float64 `json:"endprice"`
	Pnl        float64 `json:"pnl"`
	Roe        float64 `json:"roe"`
}

func (q *Queries) InsertScore(ctx context.Context, arg InsertScoreParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertScore,
		arg.ScoreID,
		arg.UserID,
		arg.Stage,
		arg.Pairname,
		arg.Entrytime,
		arg.Position,
		arg.Leverage,
		arg.Outtime,
		arg.Entryprice,
		arg.Endprice,
		arg.Pnl,
		arg.Roe,
	)
}
