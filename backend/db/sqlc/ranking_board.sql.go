// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: ranking_board.sql

package db

import (
	"context"
	"database/sql"
)

const getAllRanks = `-- name: GetAllRanks :many
SELECT user_id, photo_url, score_id, display_name, final_balance, comment FROM ranking_board
ORDER BY balance DESC
LIMIT ?
OFFSET ?
`

type GetAllRanksParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllRanks(ctx context.Context, arg GetAllRanksParams) ([]RankingBoard, error) {
	rows, err := q.db.QueryContext(ctx, getAllRanks, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []RankingBoard{}
	for rows.Next() {
		var i RankingBoard
		if err := rows.Scan(
			&i.UserID,
			&i.PhotoUrl,
			&i.ScoreID,
			&i.DisplayName,
			&i.FinalBalance,
			&i.Comment,
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

const getRankByUserID = `-- name: GetRankByUserID :one
SELECT user_id, photo_url, score_id, display_name, final_balance, comment FROM ranking_board
WHERE user_id = ?
`

func (q *Queries) GetRankByUserID(ctx context.Context, userID string) (RankingBoard, error) {
	row := q.db.QueryRowContext(ctx, getRankByUserID, userID)
	var i RankingBoard
	err := row.Scan(
		&i.UserID,
		&i.PhotoUrl,
		&i.ScoreID,
		&i.DisplayName,
		&i.FinalBalance,
		&i.Comment,
	)
	return i, err
}

const insertRank = `-- name: InsertRank :execresult
INSERT INTO ranking_board (
    user_id,
    score_id,
    display_name,
    photo_url,
    comment,
    final_balance
) VALUES (
   ?, ?, ?, ?, ?, ?
)
`

type InsertRankParams struct {
	UserID       string  `json:"user_id"`
	ScoreID      string  `json:"score_id"`
	DisplayName  string  `json:"display_name"`
	PhotoUrl     string  `json:"photo_url"`
	Comment      string  `json:"comment"`
	FinalBalance float64 `json:"final_balance"`
}

func (q *Queries) InsertRank(ctx context.Context, arg InsertRankParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertRank,
		arg.UserID,
		arg.ScoreID,
		arg.DisplayName,
		arg.PhotoUrl,
		arg.Comment,
		arg.FinalBalance,
	)
}

const updateUserRank = `-- name: UpdateUserRank :execresult
UPDATE ranking_board 
SET score_id = ?, final_balance = ?, comment = ?, display_name =?
WHERE user_id = ?
`

type UpdateUserRankParams struct {
	ScoreID      string  `json:"score_id"`
	FinalBalance float64 `json:"final_balance"`
	Comment      string  `json:"comment"`
	DisplayName  string  `json:"display_name"`
	UserID       string  `json:"user_id"`
}

func (q *Queries) UpdateUserRank(ctx context.Context, arg UpdateUserRankParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserRank,
		arg.ScoreID,
		arg.FinalBalance,
		arg.Comment,
		arg.DisplayName,
		arg.UserID,
	)
}
