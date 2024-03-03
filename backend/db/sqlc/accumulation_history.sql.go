// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: accumulation_history.sql

package db

import (
	"context"
	"database/sql"
)

const createAccumulationHist = `-- name: CreateAccumulationHist :execresult
INSERT INTO accumulation_history (
    to_user,
    amount,
    title
) VALUES (
    ?, ?, ?
)
`

type CreateAccumulationHistParams struct {
	ToUser string  `json:"to_user"`
	Amount float64 `json:"amount"`
	Title  string  `json:"title"`
}

func (q *Queries) CreateAccumulationHist(ctx context.Context, arg CreateAccumulationHistParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createAccumulationHist, arg.ToUser, arg.Amount, arg.Title)
}

const getAccumulationHist = `-- name: GetAccumulationHist :many
SELECT id, to_user, amount, title, created_at FROM accumulation_history
WHERE to_user = ? AND title LIKE ?
ORDER BY created_at DESC 
LIMIT ?
OFFSET ?
`

type GetAccumulationHistParams struct {
	ToUser string `json:"to_user"`
	Title  string `json:"title"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetAccumulationHist(ctx context.Context, arg GetAccumulationHistParams) ([]AccumulationHistory, error) {
	rows, err := q.db.QueryContext(ctx, getAccumulationHist,
		arg.ToUser,
		arg.Title,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AccumulationHistory{}
	for rows.Next() {
		var i AccumulationHistory
		if err := rows.Scan(
			&i.ID,
			&i.ToUser,
			&i.Amount,
			&i.Title,
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
