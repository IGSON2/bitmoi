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
	c := utilities.GetConfig()
	conn, err := sql.Open(c.DBDriver, c.DBSource)
	require.NoError(t, err)
	require.NotNil(t, conn)

	s := NewStore(conn)
	require.NotNil(t, s)

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

	for i := 0; i < 100; i++ {
		iq := Insert4hCandlesParams{
			Name:   name,
			Open:   (price + 10) - float64(10*i),
			Close:  (price - 10) - float64(10*i),
			High:   (price + 15) - float64(10*i),
			Low:    (price - 15) - float64(10*i),
			Time:   timestamp - int64(i*14400),
			Volume: volume - float64(i)*0.01,
		}

		_, err = s.Insert4hCandles(ctx, iq)
		require.NoError(t, err)
	}

	minmax, err := s.Get4hMinMaxTime(ctx, name)
	require.NoError(t, err)
	require.NotEmpty(t, minmax)

	min, ok := minmax.Min.(int64)
	require.Equal(t, true, ok)
	require.Greater(t, min, timestamp-int64(1440000))

	max, ok := minmax.Max.(int64)
	require.Equal(t, true, ok)
	require.Greater(t, max, min)

	waitingTime := 24 * time.Hour.Seconds()

	refTime := max - int64(utilities.MakeRanInt(int(waitingTime), int(max-min)))
	require.Greater(t, refTime, min)
	require.Less(t, refTime, max)

	gq := Get4hCandlesParams{
		Name:  name,
		Time:  refTime,
		Limit: 1000,
	}
	candles4h, err := s.Get4hCandles(ctx, gq)
	require.NoError(t, err)
	require.Greater(t, len(candles4h), 0)
}
