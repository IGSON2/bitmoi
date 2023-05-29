// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	Get15mCandles(ctx context.Context, arg Get15mCandlesParams) ([]Candles15m, error)
	Get15mMinMaxTime(ctx context.Context, name string) (Get15mMinMaxTimeRow, error)
	Get15mResult(ctx context.Context, arg Get15mResultParams) ([]Candles15m, error)
	Get15mVolSumPriceAVG(ctx context.Context, arg Get15mVolSumPriceAVGParams) (Get15mVolSumPriceAVGRow, error)
	Get1dCandles(ctx context.Context, arg Get1dCandlesParams) ([]Candles1d, error)
	Get1dMinMaxTime(ctx context.Context, name string) (Get1dMinMaxTimeRow, error)
	Get1dResult(ctx context.Context, arg Get1dResultParams) ([]Candles1d, error)
	Get1dVolSumPriceAVG(ctx context.Context, arg Get1dVolSumPriceAVGParams) (Get1dVolSumPriceAVGRow, error)
	Get1hCandles(ctx context.Context, arg Get1hCandlesParams) ([]Candles1h, error)
	Get1hMinMaxTime(ctx context.Context, name string) (Get1hMinMaxTimeRow, error)
	Get1hResult(ctx context.Context, arg Get1hResultParams) ([]Candles1h, error)
	Get1hVolSumPriceAVG(ctx context.Context, arg Get1hVolSumPriceAVGParams) (Get1hVolSumPriceAVGRow, error)
	Get4hCandles(ctx context.Context, arg Get4hCandlesParams) ([]Candles4h, error)
	Get4hMinMaxTime(ctx context.Context, name string) (Get4hMinMaxTimeRow, error)
	Get4hResult(ctx context.Context, arg Get4hResultParams) ([]Candles4h, error)
	Get4hVolSumPriceAVG(ctx context.Context, arg Get4hVolSumPriceAVGParams) (Get4hVolSumPriceAVGRow, error)
	Get5mCandles(ctx context.Context, arg Get5mCandlesParams) ([]Candles5m, error)
	Get5mMinMaxTime(ctx context.Context, name string) (Get5mMinMaxTimeRow, error)
	Get5mResult(ctx context.Context, arg Get5mResultParams) ([]Candles5m, error)
	Get5mVolSumPriceAVG(ctx context.Context, arg Get5mVolSumPriceAVGParams) (Get5mVolSumPriceAVGRow, error)
	GetAllParisInDB(ctx context.Context) ([]string, error)
	GetAllRanks(ctx context.Context, arg GetAllRanksParams) ([]RankingBoard, error)
	GetLastUser(ctx context.Context) (User, error)
	GetOne15mCandle(ctx context.Context, arg GetOne15mCandleParams) (Candles15m, error)
	GetOne1dCandle(ctx context.Context, arg GetOne1dCandleParams) (Candles1d, error)
	GetOne1hCandle(ctx context.Context, arg GetOne1hCandleParams) (Candles1h, error)
	GetOne4hCandle(ctx context.Context, arg GetOne4hCandleParams) (Candles4h, error)
	GetOne5mCandle(ctx context.Context, arg GetOne5mCandleParams) (Candles5m, error)
	GetRandomUser(ctx context.Context) (User, error)
	GetRankByUserID(ctx context.Context, userID string) (RankingBoard, error)
	GetScore(ctx context.Context, arg GetScoreParams) (Score, error)
	GetScoreToStage(ctx context.Context, arg GetScoreToStageParams) (interface{}, error)
	GetScoresByScoreID(ctx context.Context, arg GetScoresByScoreIDParams) ([]Score, error)
	GetScoresByUserID(ctx context.Context, arg GetScoresByUserIDParams) ([]Score, error)
	GetStageLenByScoreID(ctx context.Context, arg GetStageLenByScoreIDParams) (int64, error)
	GetUser(ctx context.Context, userID string) (User, error)
	Insert15mCandles(ctx context.Context, arg Insert15mCandlesParams) (sql.Result, error)
	Insert1dCandles(ctx context.Context, arg Insert1dCandlesParams) (sql.Result, error)
	Insert1hCandles(ctx context.Context, arg Insert1hCandlesParams) (sql.Result, error)
	Insert4hCandles(ctx context.Context, arg Insert4hCandlesParams) (sql.Result, error)
	Insert5mCandles(ctx context.Context, arg Insert5mCandlesParams) (sql.Result, error)
	InsertRank(ctx context.Context, arg InsertRankParams) (sql.Result, error)
	InsertScore(ctx context.Context, arg InsertScoreParams) (sql.Result, error)
	UpdatePhotoURL(ctx context.Context, arg UpdatePhotoURLParams) (sql.Result, error)
	UpdateUserRank(ctx context.Context, arg UpdateUserRankParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
