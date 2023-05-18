package db

import (
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"testing"

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

	iq := Insert4hCandlesParams{
		Name:   "TWELVETHPAIR",
		Open:   93.241,
		Close:  95.882,
		High:   100.7732,
		Low:    90.001,
		Time:   10000000,
		Volume: 3.05,
	}

	ctx := context.Background()
	_, err = s.Insert4hCandles(ctx, iq)
	require.NoError(t, err)

	gq := GetOne4hCandleParams{
		Name: iq.Name,
		Time: iq.Time,
	}
	r, err := s.GetOne4hCandle(ctx, gq)
	require.NoError(t, err)
	require.Equal(t, iq.Name, r.Name)
	require.Equal(t, iq.Open, r.Open)
	require.Equal(t, iq.Close, r.Close)
	require.Equal(t, iq.High, r.High)
	require.Equal(t, iq.Low, r.Low)
	require.Equal(t, iq.Time, r.Time)
	require.Equal(t, iq.Volume, r.Volume)

	for i := 0; i < 10; i++ {
		iq.Time = 10000000 + int64((i+1)*5000)
		_, err = s.Insert4hCandles(ctx, iq)
		require.NoError(t, err)
	}

	rs, err := s.Get4hCandles(ctx, Get4hCandlesParams{
		Name:  iq.Name,
		Limit: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, rs)
	require.Equal(t, 10, len(rs))

	for i, r := range rs {
		require.Equal(t, iq.Name, r.Name)
		require.Equal(t, iq.Open, r.Open)
		require.Equal(t, iq.Close, r.Close)
		require.Equal(t, iq.High, r.High)
		require.Equal(t, iq.Low, r.Low)
		require.Equal(t, 10000000+int64((i)*5000), r.Time)
		require.Equal(t, iq.Volume, r.Volume)
	}
}
