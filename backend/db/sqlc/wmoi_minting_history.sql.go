// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: wmoi_minting_history.sql

package db

import (
	"context"
	"database/sql"
)

const createWmoiMintingHist = `-- name: CreateWmoiMintingHist :execresult
INSERT INTO wmoi_minting_history (
    to_user,
    amount,
    title
) VALUES (
    ?, ?, ?
)
`

type CreateWmoiMintingHistParams struct {
	ToUser string `json:"to_user"`
	Amount int64  `json:"amount"`
	Title  string `json:"title"`
}

func (q *Queries) CreateWmoiMintingHist(ctx context.Context, arg CreateWmoiMintingHistParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createWmoiMintingHist, arg.ToUser, arg.Amount, arg.Title)
}

const getWmoiMintingHist = `-- name: GetWmoiMintingHist :many
SELECT id, to_user, amount, title, created_at FROM wmoi_minting_history
WHERE to_user = ? AND title LIKE ?
ORDER BY created_at DESC 
LIMIT ?
OFFSET ?
`

type GetWmoiMintingHistParams struct {
	ToUser string `json:"to_user"`
	Title  string `json:"title"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetWmoiMintingHist(ctx context.Context, arg GetWmoiMintingHistParams) ([]WmoiMintingHistory, error) {
	rows, err := q.db.QueryContext(ctx, getWmoiMintingHist,
		arg.ToUser,
		arg.Title,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []WmoiMintingHistory{}
	for rows.Next() {
		var i WmoiMintingHistory
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
