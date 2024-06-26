// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: prac_after_score.sql

package db

import (
	"context"
	"database/sql"
)

const insertPracAfterScore = `-- name: InsertPracAfterScore :execresult
INSERT INTO prac_after_score (
    score_id,
    user_id,
    max_roe,
    min_roe,
    after_outtime
) VALUES (
    ?, ?, ?, ?, ?
)
`

type InsertPracAfterScoreParams struct {
	ScoreID      string  `json:"score_id"`
	UserID       string  `json:"user_id"`
	MaxRoe       float64 `json:"max_roe"`
	MinRoe       float64 `json:"min_roe"`
	AfterOuttime int64   `json:"after_outtime"`
}

func (q *Queries) InsertPracAfterScore(ctx context.Context, arg InsertPracAfterScoreParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertPracAfterScore,
		arg.ScoreID,
		arg.UserID,
		arg.MaxRoe,
		arg.MinRoe,
		arg.AfterOuttime,
	)
}
