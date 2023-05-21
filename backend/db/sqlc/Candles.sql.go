// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: candles.sql

package db

import (
	"context"
	"database/sql"
)

const get15mCandles = `-- name: Get15mCandles :many
SELECT name, open, close, high, low, time, volume, color FROM candles_15m 
WHERE name = ?  AND time <= ?
ORDER BY time ASC 
LIMIT ?
`

type Get15mCandlesParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get15mCandles(ctx context.Context, arg Get15mCandlesParams) ([]Candles15m, error) {
	rows, err := q.db.QueryContext(ctx, get15mCandles, arg.Name, arg.Time, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Candles15m{}
	for rows.Next() {
		var i Candles15m
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

const get15mMinMaxTime = `-- name: Get15mMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_15m
WHERE name = ?
`

type Get15mMinMaxTimeRow struct {
	Min interface{} `json:"min"`
	Max interface{} `json:"max"`
}

func (q *Queries) Get15mMinMaxTime(ctx context.Context, name string) (Get15mMinMaxTimeRow, error) {
	row := q.db.QueryRowContext(ctx, get15mMinMaxTime, name)
	var i Get15mMinMaxTimeRow
	err := row.Scan(&i.Min, &i.Max)
	return i, err
}

const get15mResult = `-- name: Get15mResult :many
SELECT name, open, close, high, low, time, volume, color FROM candles_15m 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?
`

type Get15mResultParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get15mResult(ctx context.Context, arg Get15mResultParams) ([]Candles15m, error) {
	rows, err := q.db.QueryContext(ctx, get15mResult, arg.Name, arg.Time, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Candles15m{}
	for rows.Next() {
		var i Candles15m
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

const get1dCandles = `-- name: Get1dCandles :many
SELECT name, open, close, high, low, time, volume, color FROM candles_1d 
WHERE name = ? AND time <= ?
ORDER BY time ASC 
LIMIT ?
`

type Get1dCandlesParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get1dCandles(ctx context.Context, arg Get1dCandlesParams) ([]Candles1d, error) {
	rows, err := q.db.QueryContext(ctx, get1dCandles, arg.Name, arg.Time, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Candles1d{}
	for rows.Next() {
		var i Candles1d
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

const get1dMinMaxTime = `-- name: Get1dMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_1d
WHERE name = ?
`

type Get1dMinMaxTimeRow struct {
	Min interface{} `json:"min"`
	Max interface{} `json:"max"`
}

func (q *Queries) Get1dMinMaxTime(ctx context.Context, name string) (Get1dMinMaxTimeRow, error) {
	row := q.db.QueryRowContext(ctx, get1dMinMaxTime, name)
	var i Get1dMinMaxTimeRow
	err := row.Scan(&i.Min, &i.Max)
	return i, err
}

const get1dResult = `-- name: Get1dResult :many
SELECT name, open, close, high, low, time, volume, color FROM candles_1d 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?
`

type Get1dResultParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get1dResult(ctx context.Context, arg Get1dResultParams) ([]Candles1d, error) {
	rows, err := q.db.QueryContext(ctx, get1dResult, arg.Name, arg.Time, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Candles1d{}
	for rows.Next() {
		var i Candles1d
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

const get1hCandles = `-- name: Get1hCandles :many
SELECT name, open, close, high, low, time, volume, color FROM candles_1h 
WHERE name = ?  AND time <= ?
ORDER BY time ASC 
LIMIT ?
`

type Get1hCandlesParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get1hCandles(ctx context.Context, arg Get1hCandlesParams) ([]Candles1h, error) {
	rows, err := q.db.QueryContext(ctx, get1hCandles, arg.Name, arg.Time, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Candles1h{}
	for rows.Next() {
		var i Candles1h
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

const get1hMinMaxTime = `-- name: Get1hMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_1h
WHERE name = ?
`

type Get1hMinMaxTimeRow struct {
	Min interface{} `json:"min"`
	Max interface{} `json:"max"`
}

func (q *Queries) Get1hMinMaxTime(ctx context.Context, name string) (Get1hMinMaxTimeRow, error) {
	row := q.db.QueryRowContext(ctx, get1hMinMaxTime, name)
	var i Get1hMinMaxTimeRow
	err := row.Scan(&i.Min, &i.Max)
	return i, err
}

const get1hResult = `-- name: Get1hResult :many
SELECT name, open, close, high, low, time, volume, color FROM candles_1h 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?
`

type Get1hResultParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get1hResult(ctx context.Context, arg Get1hResultParams) ([]Candles1h, error) {
	rows, err := q.db.QueryContext(ctx, get1hResult, arg.Name, arg.Time, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Candles1h{}
	for rows.Next() {
		var i Candles1h
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

const get4hCandles = `-- name: Get4hCandles :many
SELECT name, open, close, high, low, time, volume, color FROM candles_4h 
WHERE name = ?  AND time <= ?
ORDER BY time ASC 
LIMIT ?
`

type Get4hCandlesParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get4hCandles(ctx context.Context, arg Get4hCandlesParams) ([]Candles4h, error) {
	rows, err := q.db.QueryContext(ctx, get4hCandles, arg.Name, arg.Time, arg.Limit)
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

const get4hMinMaxTime = `-- name: Get4hMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_4h
WHERE name = ?
`

type Get4hMinMaxTimeRow struct {
	Min interface{} `json:"min"`
	Max interface{} `json:"max"`
}

func (q *Queries) Get4hMinMaxTime(ctx context.Context, name string) (Get4hMinMaxTimeRow, error) {
	row := q.db.QueryRowContext(ctx, get4hMinMaxTime, name)
	var i Get4hMinMaxTimeRow
	err := row.Scan(&i.Min, &i.Max)
	return i, err
}

const get4hResult = `-- name: Get4hResult :many
SELECT name, open, close, high, low, time, volume, color FROM candles_4h 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?
`

type Get4hResultParams struct {
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Limit int32  `json:"limit"`
}

func (q *Queries) Get4hResult(ctx context.Context, arg Get4hResultParams) ([]Candles4h, error) {
	rows, err := q.db.QueryContext(ctx, get4hResult, arg.Name, arg.Time, arg.Limit)
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

const getOne15mCandle = `-- name: GetOne15mCandle :one
SELECT name, open, close, high, low, time, volume, color FROM candles_15m
WHERE name = ? AND time = ?
`

type GetOne15mCandleParams struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

func (q *Queries) GetOne15mCandle(ctx context.Context, arg GetOne15mCandleParams) (Candles15m, error) {
	row := q.db.QueryRowContext(ctx, getOne15mCandle, arg.Name, arg.Time)
	var i Candles15m
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

const getOne1dCandle = `-- name: GetOne1dCandle :one
SELECT name, open, close, high, low, time, volume, color FROM candles_1d
WHERE name = ? AND time = ?
`

type GetOne1dCandleParams struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

func (q *Queries) GetOne1dCandle(ctx context.Context, arg GetOne1dCandleParams) (Candles1d, error) {
	row := q.db.QueryRowContext(ctx, getOne1dCandle, arg.Name, arg.Time)
	var i Candles1d
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

const getOne1hCandle = `-- name: GetOne1hCandle :one
SELECT name, open, close, high, low, time, volume, color FROM candles_1h
WHERE name = ? AND time = ?
`

type GetOne1hCandleParams struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

func (q *Queries) GetOne1hCandle(ctx context.Context, arg GetOne1hCandleParams) (Candles1h, error) {
	row := q.db.QueryRowContext(ctx, getOne1hCandle, arg.Name, arg.Time)
	var i Candles1h
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

const getOne4hCandle = `-- name: GetOne4hCandle :one
SELECT name, open, close, high, low, time, volume, color FROM candles_4h
WHERE name = ? AND time = ?
`

type GetOne4hCandleParams struct {
	Name string `json:"name"`
	Time int64  `json:"time"`
}

func (q *Queries) GetOne4hCandle(ctx context.Context, arg GetOne4hCandleParams) (Candles4h, error) {
	row := q.db.QueryRowContext(ctx, getOne4hCandle, arg.Name, arg.Time)
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

const insert15mCandles = `-- name: Insert15mCandles :execresult
INSERT INTO candles_15m (
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

type Insert15mCandlesParams struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

func (q *Queries) Insert15mCandles(ctx context.Context, arg Insert15mCandlesParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insert15mCandles,
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

const insert1dCandles = `-- name: Insert1dCandles :execresult
INSERT INTO candles_1d (
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

type Insert1dCandlesParams struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

func (q *Queries) Insert1dCandles(ctx context.Context, arg Insert1dCandlesParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insert1dCandles,
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

const insert1hCandles = `-- name: Insert1hCandles :execresult
INSERT INTO candles_1h (
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

type Insert1hCandlesParams struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

func (q *Queries) Insert1hCandles(ctx context.Context, arg Insert1hCandlesParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insert1hCandles,
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

const insert4hCandles = `-- name: Insert4hCandles :execresult
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

type Insert4hCandlesParams struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

func (q *Queries) Insert4hCandles(ctx context.Context, arg Insert4hCandlesParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insert4hCandles,
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
