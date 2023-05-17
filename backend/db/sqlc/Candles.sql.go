// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: Candles.sql

package db

import (
	"context"
	"database/sql"
)

const getCandles = `-- name: GetCandles :many
SELECT name, open, close, high, low, time, volume, color FROM candles_4h 
WHERE name = ? 
ORDER BY time DESC 
LIMIT ?
`

type GetCandlesParams struct {
	Name  string `json:"name"`
	Limit int32  `json:"limit"`
}

func (q *Queries) GetCandles(ctx context.Context, arg GetCandlesParams) ([]Candles4h, error) {
	rows, err := q.db.QueryContext(ctx, getCandles, arg.Name, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Candles4h{}
	for rows.Next() {
		var i Candles4h
		if err := rows.Scan(
			&i.Name,
			&i.Open,
			&i.Close,
			&i.High,
			&i.Low,
			&i.Time,
			&i.Volume,
			&i.Color,
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

const getOneCandle = `-- name: GetOneCandle :one
SELECT name, open, close, high, low, time, volume, color FROM candles_4h
WHERE name = ? AND time = ?
`

type GetOneCandleParams struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

func (q *Queries) GetOneCandle(ctx context.Context, arg GetOneCandleParams) (Candles4h, error) {
	row := q.db.QueryRowContext(ctx, getOneCandle, arg.Name, arg.Time)
	var i Candles4h
	err := row.Scan(
		&i.Name,
		&i.Open,
		&i.Close,
		&i.High,
		&i.Low,
		&i.Time,
		&i.Volume,
		&i.Color,
	)
	return i, err
}

const insertCandles = `-- name: InsertCandles :execresult
INSERT INTO candles_4h (
    name,
    open,
    close,
    high,
    low,
    time,
    volume,
    color
) VALUES (
  ?,?,?,?,?,?,?,?
)
`

type InsertCandlesParams struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

func (q *Queries) InsertCandles(ctx context.Context, arg InsertCandlesParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertCandles,
		arg.Name,
		arg.Open,
		arg.Close,
		arg.High,
		arg.Low,
		arg.Time,
		arg.Volume,
		arg.Color,
	)
}
