package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

type TestName struct {
	Names string
}

func (t *TestName) append(pair string) {
	if n := strings.Count(t.Names, ","); n == 9 {
		t.Names += pair
		return
	} else if n > 9 {
		log.Error().Err(fmt.Errorf("cannot append anymore: n=%d", n))
		return
	}
	t.Names += pair + ","
}

func TestSomePairs(t *testing.T) {
	store := newTestStore(t)
	server := newTestServer(t, store)
	client := newGRPCClient(t)
	pairs, err := store.GetAllParisInDB(context.Background())

	require.NoError(t, err)

	for i := 0; i < len(pairs); i++ {
		testAnotherInterval(t, store, server, client)
	}

}

func testAnotherInterval(t *testing.T, store db.Store, s *Server, client pb.BitmoiClient) {
	tn := new(TestName)

	testCases := []struct {
		Name      string
		CandleReq *pb.GetCandlesRequest
		Req       *pb.AnotherIntervalRequest
		SetUpAuth func(t *testing.T, tm *token.PasetoMaker) context.Context
		CheckResp func(t *testing.T, candleRes, res *pb.GetCandlesResponse, req *pb.AnotherIntervalRequest, err error)
	}{
		{
			Name: "OK_Practice_5m",
			CandleReq: &pb.GetCandlesRequest{
				Mode:   practice,
				UserId: "",
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FiveM,
				Mode:        practice,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: testResponseWithRequest,
		},
		{
			Name: "OK_Practice_15m",
			CandleReq: &pb.GetCandlesRequest{
				Mode:   practice,
				UserId: "",
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FifM,
				Mode:        practice,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: testResponseWithRequest,
		},
		{
			Name: "OK_Competition_5m",
			CandleReq: &pb.GetCandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FiveM,
				Mode:        competition,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
			CheckResp: testResponseWithRequest,
		},
		{
			Name: "OK_Competition_15m",
			CandleReq: &pb.GetCandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FifM,
				Mode:        competition,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
			CheckResp: testResponseWithRequest,
		},
		{
			Name: "No_Auth_Competition",
			CandleReq: &pb.GetCandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				Mode:   competition,
				UserId: "unauthorized user",
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
			CheckResp: func(t *testing.T, candleRes *pb.GetCandlesResponse, res *pb.GetCandlesResponse, req *pb.AnotherIntervalRequest, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "Fail_Validation_Practice",
			CandleReq: &pb.GetCandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: "Unsupported",
				Mode:        competition,
				UserId:      user,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
			CheckResp: func(t *testing.T, candleRes *pb.GetCandlesResponse, res *pb.GetCandlesResponse, req *pb.AnotherIntervalRequest, err error) {
				t.Log(err)
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := tc.SetUpAuth(t, s.tokenMaker)

			tc.CandleReq.Names = tn.Names
			candleRes, err := client.RequestCandles(ctx, tc.CandleReq)
			require.NoError(t, err)

			tn.append(candleRes.Name)
			if len(utilities.Splitnames(tn.Names)) > 10 {
				t.Log("The number of pairs has reached 10. End the test ")
				return
			}

			tc.Req.Identifier = candleRes.Identifier

			res, err := client.AnotherInterval(ctx, tc.Req)

			tc.CheckResp(t, candleRes, res, tc.Req, err)
		})
	}

}

func testIdentifier(t *testing.T, resIden, reqIden, requestedInterval string) {
	resInfo := new(utilities.IdentificationData)
	json.Unmarshal(utilities.DecryptByASE(resIden), resInfo)

	reqInfo := new(utilities.IdentificationData)
	json.Unmarshal(utilities.DecryptByASE(reqIden), reqInfo)

	require.Equal(t, reqInfo.Interval, db.OneH)
	require.Equal(t, resInfo.Interval, requestedInterval)

	require.Equal(t, resInfo.Name, reqInfo.Name)
	require.Equal(t, resInfo.RefTimestamp, reqInfo.RefTimestamp)
	require.Equal(t, resInfo.PriceFactor, reqInfo.PriceFactor)
	require.Equal(t, resInfo.VolumeFactor, reqInfo.VolumeFactor)
	require.Equal(t, resInfo.TimeFactor, reqInfo.TimeFactor)
}

func testResponseWithRequest(t *testing.T, candleRes, res *pb.GetCandlesResponse, req *pb.AnotherIntervalRequest, err error) {
	require.NoError(t, err)
	testIdentifier(t, res.Identifier, req.Identifier, req.ReqInterval)

	require.NotNil(t, candleRes)
	require.NotNil(t, candleRes.Onechart.PData)
	require.NotNil(t, candleRes.Onechart.VData)
	require.NotEmpty(t, candleRes.Identifier)
	require.NotNil(t, res)
	require.NotNil(t, res.Onechart.PData)
	require.NotNil(t, res.Onechart.VData)
	require.NotEmpty(t, res.Identifier)

	resClose := res.Onechart.PData[len(res.Onechart.PData)-1].Close
	candleClose := candleRes.Onechart.PData[len(candleRes.Onechart.PData)-1].Close

	require.Equal(t, res.Name, candleRes.Name)
	require.Equal(t, res.Onechart.PData[len(res.Onechart.PData)-1].Time, candleRes.Onechart.PData[len(candleRes.Onechart.PData)-1].Time)
	require.GreaterOrEqual(t, float64(0.01), math.Abs(resClose-candleClose)/resClose) // 각 단위의 캔들의 종가의 차이가 1% 이하여야 함
	require.Equal(t, res.EntryTime, candleRes.EntryTime)
	require.GreaterOrEqual(t, float64(0.01), math.Abs(res.EntryPrice-candleRes.EntryPrice)/res.EntryPrice)
}
