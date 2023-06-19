package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnotherInterval(t *testing.T) {

	store := newTestStore(t)
	s := newTestServer(t, store)
	client := newGRPCClient(t)

	// 1h 요청 -> 5m, 15m 해당 구간만큼 데이터 있어야함 -> 1h요청 identifier 해석 -> another interval요청 -> 결과 비교
	// identifier := utilities.EncrtpByASE(utilities.IdentificationData{})

	testCases := []struct {
		Name      string
		Req       *pb.AnotherIntervalRequest
		SetUpAuth func(t *testing.T, tm *token.PasetoMaker) context.Context
		CheckResp func(res *pb.GetCandlesResponse, pairs []string, err error)
	}{
		{
			Name: "OK_Practice",
			Req: &pb.AnotherIntervalRequest{
				Mode:   practice,
				UserId: user,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: func(res *pb.GetCandlesResponse, pairs []string, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.Identifier, "")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(res.Identifier), info))

				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, res.BtcRatio, float64(0))
				require.Greater(t, res.EntryTime, "")

				require.NotNil(t, res.Onechart.PData)
				require.NotNil(t, res.Onechart.VData)

				require.Contains(t, pairs, res.Name)
				require.Equal(t, info.Interval, db.OneH)
				require.Equal(t, info.PriceFactor, float64(0))
				require.Equal(t, info.VolumeFactor, float64(0))
				require.Equal(t, info.TimeFactor, int64(0))
			},
		},
		{
			Name: "OK_Competition",
			Req: &pb.AnotherIntervalRequest{
				Mode:   competition,
				UserId: user,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
			CheckResp: func(res *pb.GetCandlesResponse, pairs []string, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.Identifier, "")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(res.Identifier), info))

				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, res.BtcRatio, float64(0))
				require.Greater(t, res.EntryTime, "")

				require.NotNil(t, res.Onechart.PData)
				require.NotNil(t, res.Onechart.VData)

				require.Contains(t, res.Name, "STAGE")
				require.Equal(t, info.Interval, db.OneH)
				require.Greater(t, info.PriceFactor, float64(0))
				require.Greater(t, info.VolumeFactor, float64(0))
				require.Greater(t, info.TimeFactor, int64(0))
			},
		},
		{
			Name: "No_Auth_Competition",
			Req: &pb.AnotherIntervalRequest{
				Mode:   competition,
				UserId: "unauthorized user",
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: func(res *pb.GetCandlesResponse, pairs []string, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "Fail_Validation_Practice",
			Req: &pb.AnotherIntervalRequest{
				Mode:   "Unsupported",
				UserId: user,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: func(res *pb.GetCandlesResponse, pairs []string, err error) {
				t.Log(err)
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := tc.SetUpAuth(t, s.tokenMaker)
			res, err := client.AnotherInterval(ctx, tc.Req)
			tc.CheckResp(res, s.pairs, err)
		})
	}

}
