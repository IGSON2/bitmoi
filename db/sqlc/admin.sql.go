// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: admin.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const getAdminScores = `-- name: GetAdminScores :many
SELECT u.id, u.nickname, p.score_id, p.user_id, p.stage, p.pairname, p.entrytime, p.position, p.leverage, p.outtime, p.entryprice, p.quantity, p.endprice, p.pnl, p.roe, p.settled_at, p.created_at ,a.min_roe, a.max_roe, a.after_outtime from prac_score p
 INNER JOIN users u ON p.user_id = u.user_id
 LEFT JOIN prac_after_score a ON p.user_id = a.user_id AND p.score_id = a.score_id
 LIMIT ? OFFSET ?
`

type GetAdminScoresParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetAdminScoresRow struct {
	ID           int64           `json:"id"`
	Nickname     string          `json:"nickname"`
	ScoreID      string          `json:"score_id"`
	UserID       string          `json:"user_id"`
	Stage        int8            `json:"stage"`
	Pairname     string          `json:"pairname"`
	Entrytime    string          `json:"entrytime"`
	Position     string          `json:"position"`
	Leverage     int8            `json:"leverage"`
	Outtime      string          `json:"outtime"`
	Entryprice   float64         `json:"entryprice"`
	Quantity     float64         `json:"quantity"`
	Endprice     float64         `json:"endprice"`
	Pnl          float64         `json:"pnl"`
	Roe          float64         `json:"roe"`
	SettledAt    sql.NullTime    `json:"settled_at"`
	CreatedAt    time.Time       `json:"created_at"`
	MinRoe       sql.NullFloat64 `json:"min_roe"`
	MaxRoe       sql.NullFloat64 `json:"max_roe"`
	AfterOuttime sql.NullInt64   `json:"after_outtime"`
}

func (q *Queries) GetAdminScores(ctx context.Context, arg GetAdminScoresParams) ([]GetAdminScoresRow, error) {
	rows, err := q.db.QueryContext(ctx, getAdminScores, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAdminScoresRow{}
	for rows.Next() {
		var i GetAdminScoresRow
		if err := rows.Scan(
			&i.ID,
			&i.Nickname,
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
			&i.MinRoe,
			&i.MaxRoe,
			&i.AfterOuttime,
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

const getAdminUsdpInfo = `-- name: GetAdminUsdpInfo :many
SELECT u.id, u.nickname, a.to_user, a.amount, a.title, a.created_at, a.method, a.giver FROM accumulation_history a 
INNER JOIN users u ON a.to_user = u.user_id
LIMIT ? OFFSET ?
`

type GetAdminUsdpInfoParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetAdminUsdpInfoRow struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	ToUser    string    `json:"to_user"`
	Amount    float64   `json:"amount"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Method    string    `json:"method"`
	Giver     string    `json:"giver"`
}

func (q *Queries) GetAdminUsdpInfo(ctx context.Context, arg GetAdminUsdpInfoParams) ([]GetAdminUsdpInfoRow, error) {
	rows, err := q.db.QueryContext(ctx, getAdminUsdpInfo, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAdminUsdpInfoRow{}
	for rows.Next() {
		var i GetAdminUsdpInfoRow
		if err := rows.Scan(
			&i.ID,
			&i.Nickname,
			&i.ToUser,
			&i.Amount,
			&i.Title,
			&i.CreatedAt,
			&i.Method,
			&i.Giver,
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

const getAdminUsers = `-- name: GetAdminUsers :many
SELECT id, user_id, nickname, prac_balance, wmoi_balance, recommender_code, created_at, last_accessed_at,
(SELECT COUNT(*) FROM accumulation_history WHERE users.user_id = accumulation_history.to_user AND accumulation_history.title='출석 체크 보상') AS attendance,
(SELECT COUNT(*) FROM recommend_history WHERE users.user_id = recommend_history.recommender) AS referral,
(SELECT COUNT(*) FROM prac_score WHERE users.user_id = prac_score.user_id AND prac_score.pnl >= 0) AS prac_win,
(SELECT COUNT(*) FROM prac_score WHERE users.user_id = prac_score.user_id AND prac_score.pnl < 0) AS prac_lose,
(SELECT COUNT(*) FROM comp_score WHERE users.user_id = comp_score.user_id AND comp_score.pnl >= 0) AS comp_win,
(SELECT COUNT(*) FROM comp_score WHERE users.user_id = comp_score.user_id AND comp_score.pnl < 0) AS comp_lose
FROM users
LIMIT ? OFFSET ?
`

type GetAdminUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetAdminUsersRow struct {
	ID              int64        `json:"id"`
	UserID          string       `json:"user_id"`
	Nickname        string       `json:"nickname"`
	PracBalance     float64      `json:"prac_balance"`
	WmoiBalance     float64      `json:"wmoi_balance"`
	RecommenderCode string       `json:"recommender_code"`
	CreatedAt       time.Time    `json:"created_at"`
	LastAccessedAt  sql.NullTime `json:"last_accessed_at"`
	Attendance      int64        `json:"attendance"`
	Referral        int64        `json:"referral"`
	PracWin         int64        `json:"prac_win"`
	PracLose        int64        `json:"prac_lose"`
	CompWin         int64        `json:"comp_win"`
	CompLose        int64        `json:"comp_lose"`
}

func (q *Queries) GetAdminUsers(ctx context.Context, arg GetAdminUsersParams) ([]GetAdminUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAdminUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAdminUsersRow{}
	for rows.Next() {
		var i GetAdminUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Nickname,
			&i.PracBalance,
			&i.WmoiBalance,
			&i.RecommenderCode,
			&i.CreatedAt,
			&i.LastAccessedAt,
			&i.Attendance,
			&i.Referral,
			&i.PracWin,
			&i.PracLose,
			&i.CompWin,
			&i.CompLose,
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