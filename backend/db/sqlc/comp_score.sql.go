// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: comp_score.sql

package db

import (
	"context"
	"database/sql"
)

const getCompScore = `-- name: GetCompScore :one
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, settled_at, created_at FROM comp_score
WHERE user_id = ? AND score_id = ? AND pairname = ?
`

type GetCompScoreParams struct {
	UserID   string `json:"user_id"`
	ScoreID  string `json:"score_id"`
	Pairname string `json:"pairname"`
}

func (q *Queries) GetCompScore(ctx context.Context, arg GetCompScoreParams) (CompScore, error) {
	row := q.db.QueryRowContext(ctx, getCompScore, arg.UserID, arg.ScoreID, arg.Pairname)
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
		&i.SettledAt,
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
	Stage   int8   `json:"stage"`
}

func (q *Queries) GetCompScoreToStage(ctx context.Context, arg GetCompScoreToStageParams) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getCompScoreToStage, arg.ScoreID, arg.UserID, arg.Stage)
	var sum interface{}
	err := row.Scan(&sum)
	return sum, err
}

const getCompScoresByScoreID = `-- name: GetCompScoresByScoreID :many
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, settled_at, created_at FROM comp_score
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
			&i.SettledAt,
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

const getCompScoresByStage = `-- name: GetCompScoresByStage :one
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, settled_at, created_at FROM comp_score
WHERE score_id = ? AND user_id = ? AND stage = ?
`

type GetCompScoresByStageParams struct {
	ScoreID string `json:"score_id"`
	UserID  string `json:"user_id"`
	Stage   int8   `json:"stage"`
}

func (q *Queries) GetCompScoresByStage(ctx context.Context, arg GetCompScoresByStageParams) (CompScore, error) {
	row := q.db.QueryRowContext(ctx, getCompScoresByStage, arg.ScoreID, arg.UserID, arg.Stage)
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
		&i.SettledAt,
		&i.CreatedAt,
	)
	return i, err
}

const getCompScoresByUserID = `-- name: GetCompScoresByUserID :many
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, settled_at, created_at FROM comp_score
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
			&i.SettledAt,
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

const getUnsettledCompScores = `-- name: GetUnsettledCompScores :many
SELECT score_id, user_id, stage, pairname, entrytime, position, leverage, outtime, entryprice, quantity, endprice, pnl, roe, settled_at, created_at FROM comp_score
WHERE user_id = ? AND pnl <> 0 AND outtime = 0 AND settled_at IS NULL
`

func (q *Queries) GetUnsettledCompScores(ctx context.Context, userID string) ([]CompScore, error) {
	rows, err := q.db.QueryContext(ctx, getUnsettledCompScores, userID)
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
			&i.SettledAt,
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

const getUserCompScoreSummary = `-- name: GetUserCompScoreSummary :one
SELECT 
  SUM(pnl) AS total_pnl,
  COUNT(CASE WHEN  pnl > 0 THEN 1 END) AS total_win,
  COUNT(CASE WHEN s.pnl < 0 THEN 1 END) AS total_lose,
  SUM(CASE WHEN s.created_at >= CURDATE() - INTERVAL 1 MONTH THEN s.pnl ELSE 0 END) AS monthly_pnl,
  COUNT(CASE WHEN s.created_at >= CURDATE() - INTERVAL 1 MONTH AND s.pnl > 0 THEN 1 END) AS monthly_win,
  COUNT(CASE WHEN s.created_at >= CURDATE() - INTERVAL 1 MONTH AND s.pnl < 0 THEN 1 END) AS monthly_lose
FROM comp_score s
JOIN users u ON s.user_id = u.user_id
WHERE u.nickname = ?
`

type GetUserCompScoreSummaryRow struct {
	TotalPnl    interface{} `json:"total_pnl"`
	TotalWin    int64       `json:"total_win"`
	TotalLose   int64       `json:"total_lose"`
	MonthlyPnl  interface{} `json:"monthly_pnl"`
	MonthlyWin  int64       `json:"monthly_win"`
	MonthlyLose int64       `json:"monthly_lose"`
}

func (q *Queries) GetUserCompScoreSummary(ctx context.Context, nickname string) (GetUserCompScoreSummaryRow, error) {
	row := q.db.QueryRowContext(ctx, getUserCompScoreSummary, nickname)
	var i GetUserCompScoreSummaryRow
	err := row.Scan(
		&i.TotalPnl,
		&i.TotalWin,
		&i.TotalLose,
		&i.MonthlyPnl,
		&i.MonthlyWin,
		&i.MonthlyLose,
	)
	return i, err
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
    roe
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

type InsertCompScoreParams struct {
	ScoreID    string  `json:"score_id"`
	UserID     string  `json:"user_id"`
	Stage      int8    `json:"stage"`
	Pairname   string  `json:"pairname"`
	Entrytime  string  `json:"entrytime"`
	Position   string  `json:"position"`
	Leverage   int8    `json:"leverage"`
	Outtime    int64   `json:"outtime"`
	Entryprice float64 `json:"entryprice"`
	Quantity   float64 `json:"quantity"`
	Endprice   float64 `json:"endprice"`
	Pnl        float64 `json:"pnl"`
	Roe        float64 `json:"roe"`
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
	)
}

const updateCompScoreSettledAt = `-- name: UpdateCompScoreSettledAt :execresult
UPDATE comp_score SET settled_at = ?
WHERE user_id = ? AND score_id = ?
`

type UpdateCompScoreSettledAtParams struct {
	SettledAt sql.NullTime `json:"settled_at"`
	UserID    string       `json:"user_id"`
	ScoreID   string       `json:"score_id"`
}

func (q *Queries) UpdateCompScoreSettledAt(ctx context.Context, arg UpdateCompScoreSettledAtParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateCompScoreSettledAt, arg.SettledAt, arg.UserID, arg.ScoreID)
}

const updateCompcScore = `-- name: UpdateCompcScore :execresult
UPDATE comp_score SET pairname = ?, entrytime = ?, outtime = ?, entryprice = ?, endprice = ?, pnl = ?, roe = ?
WHERE user_id = ? AND score_id = ? AND pairname = ?
`

type UpdateCompcScoreParams struct {
	Pairname   string  `json:"pairname"`
	Entrytime  string  `json:"entrytime"`
	Outtime    int64   `json:"outtime"`
	Entryprice float64 `json:"entryprice"`
	Endprice   float64 `json:"endprice"`
	Pnl        float64 `json:"pnl"`
	Roe        float64 `json:"roe"`
	UserID     string  `json:"user_id"`
	ScoreID    string  `json:"score_id"`
	Pairname_2 string  `json:"pairname_2"`
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
		arg.UserID,
		arg.ScoreID,
		arg.Pairname_2,
	)
}
