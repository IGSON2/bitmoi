// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	Get15mCandles(ctx context.Context, arg Get15mCandlesParams) ([]Candles15m, error)
	Get1dCandles(ctx context.Context, arg Get1dCandlesParams) ([]Candles1d, error)
	Get1hCandles(ctx context.Context, arg Get1hCandlesParams) ([]Candles1h, error)
	Get4hCandles(ctx context.Context, arg Get4hCandlesParams) ([]Candles4h, error)
	GetOne15mCandle(ctx context.Context, arg GetOne15mCandleParams) (Candles15m, error)
	GetOne1dCandle(ctx context.Context, arg GetOne1dCandleParams) (Candles1d, error)
	GetOne1hCandle(ctx context.Context, arg GetOne1hCandleParams) (Candles1h, error)
	GetOne4hCandle(ctx context.Context, arg GetOne4hCandleParams) (Candles4h, error)
	Insert15mCandles(ctx context.Context, arg Insert15mCandlesParams) (sql.Result, error)
	Insert1dCandles(ctx context.Context, arg Insert1dCandlesParams) (sql.Result, error)
	Insert1hCandles(ctx context.Context, arg Insert1hCandlesParams) (sql.Result, error)
	Insert4hCandles(ctx context.Context, arg Insert4hCandlesParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
