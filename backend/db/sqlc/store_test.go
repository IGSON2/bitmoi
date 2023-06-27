package db

import (
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
)

func TestGetStore(t *testing.T) {

	// name := "BTCUSDT"
	name := strings.ToUpper(utilities.MakeRanString(3)) + "USDT"
	require.Equal(t, 7, len(name))

	price, err := utilities.MakeRanFloat(20000, 30000)
	require.NoError(t, err)
	require.Greater(t, price, float64(20000))
	require.Less(t, price, float64(30000))

	volume, err := utilities.MakeRanFloat(5, 100)
	require.NoError(t, err)
	require.Greater(t, volume, float64(5))
	require.Less(t, volume, float64(100))

	var timestamp int64 = 1684482575

	ctx := context.Background()

	for i := 0; i < 200; i++ {
		iq1 := Insert1hCandlesParams{
			Name:   name,
			Open:   (price + 10) - float64(10*i),
			Close:  (price - 10) - float64(10*i),
			High:   (price + 15) - float64(10*i),
			Low:    (price - 15) - float64(10*i),
			Time:   timestamp - int64(i*3600),
			Volume: volume - float64(i)*0.01,
		}

		_, err = testQueries.Insert1hCandles(ctx, iq1)
		require.NoError(t, err)

		iq15 := Insert15mCandlesParams{
			Name:   name,
			Open:   (price + 10) - float64(10*i),
			Close:  (price - 10) - float64(10*i),
			High:   (price + 15) - float64(10*i),
			Low:    (price - 15) - float64(10*i),
			Time:   timestamp - int64(i*900),
			Volume: volume - float64(i)*0.01,
		}

		_, err = testQueries.Insert15mCandles(ctx, iq15)
		require.NoError(t, err)

		iq5 := Insert5mCandlesParams{
			Name:   name,
			Open:   (price + 10) - float64(10*i),
			Close:  (price - 10) - float64(10*i),
			High:   (price + 15) - float64(10*i),
			Low:    (price - 15) - float64(10*i),
			Time:   timestamp - int64(i*300),
			Volume: volume - float64(i)*0.01,
		}

		_, err = testQueries.Insert5mCandles(ctx, iq5)
		require.NoError(t, err)
	}

	//1h
	minmax1, err := testQueries.Get1hMinMaxTime(ctx, name)
	require.NoError(t, err)
	require.NotEmpty(t, minmax1)

	min, max := minmax1.Convert()
	require.Greater(t, min, timestamp-int64(3600*200))
	require.Greater(t, max, min)

	waitingTime := 24 * time.Hour.Seconds()

	refTime := max - int64(utilities.MakeRanInt(int(waitingTime), int(max-min)))
	require.Greater(t, refTime, min)
	require.Less(t, refTime, max)

	gq1 := Get1hCandlesParams{
		Name:  name,
		Time:  refTime,
		Limit: 1000,
	}
	candles1h, err := testQueries.Get1hCandles(ctx, gq1)
	require.NoError(t, err)
	require.Greater(t, len(candles1h), 0)

	//15m
	minmax15, err := testQueries.Get15mMinMaxTime(ctx, name)
	require.NoError(t, err)
	require.NotEmpty(t, minmax15)

	min, max = minmax15.Convert()
	require.Greater(t, min, timestamp-int64(900*200))
	require.Greater(t, max, min)

	waitingTime = 6 * time.Hour.Seconds()
	refTime = max - int64(utilities.MakeRanInt(int(waitingTime), int(max-min)))
	require.Greater(t, refTime, min)
	require.Less(t, refTime, max)

	gq15 := Get15mCandlesParams{
		Name:  name,
		Time:  refTime,
		Limit: 1000,
	}
	candles15m, err := testQueries.Get15mCandles(ctx, gq15)
	require.NoError(t, err)
	require.Greater(t, len(candles15m), 0)

	//5m
	minmax5, err := testQueries.Get5mMinMaxTime(ctx, name)
	require.NoError(t, err)
	require.NotEmpty(t, minmax5)

	min, max = minmax5.Convert()
	require.Greater(t, min, timestamp-int64(300*200))
	require.Greater(t, max, min)

	waitingTime = 2 * time.Hour.Seconds()
	refTime = max - int64(utilities.MakeRanInt(int(waitingTime), int(max-min)))
	require.Greater(t, refTime, min)
	require.Less(t, refTime, max)

	gq5 := Get5mCandlesParams{
		Name:  name,
		Time:  refTime,
		Limit: 1000,
	}
	candles5m, err := testQueries.Get5mCandles(ctx, gq5)
	require.NoError(t, err)
	require.Greater(t, len(candles5m), 0)
}

func TestInsertUser(t *testing.T) {
	ctx := context.Background()
	testQueries.CreateUser(ctx, CreateUserParams{
		UserID:         "user",
		OauthUid:       sql.NullString{String: "1234", Valid: true},
		FullName:       "user_full",
		HashedPassword: "392cdf",
		Email:          "example@gmail.com",
		PhotoUrl:       sql.NullString{String: "photh.url", Valid: true},
	})
}

func TestInsertScore(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_, err := testQueries.InsertScore(ctx, InsertScoreParams{
			ScoreID:    utilities.MakeRanString(3),
			UserID:     "igson",
			Stage:      int32(i + 1),
			Pairname:   "SOMPAIR",
			Entrytime:  "2023",
			Position:   "LONG",
			Leverage:   1,
			Outtime:    2,
			Entryprice: 10.01,
			Endprice:   11.01,
			Pnl:        -240.24,
			Roe:        -15.1,
		})
		require.NoError(t, err)
	}
}

func TestGetTotalScore(t *testing.T) {
	ctx := context.Background()
	i, err := testQueries.GetScoreToStage(ctx, GetScoreToStageParams{
		ScoreID: "1",
		UserID:  "user",
		Stage:   5,
	})
	require.NoError(t, err)

	totalScore, ok := i.(float64)
	require.Equal(t, true, ok)
	require.Equal(t, -1201.2, totalScore)

}
