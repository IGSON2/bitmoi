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

func TestMakeChart(t *testing.T) {

	store := newTestStore(t)
	s := newTestServer(t, store)
	client := newGRPCClient(t)

	testCases := []struct {
		Name      string
		req       *pb.GetCandlesRequest
		SetUpAuth func(t *testing.T, tm *token.PasetoMaker) context.Context
		CheckResp func(res *pb.GetCandlesResponse, pairs []string, err error)
	}{
		{
			Name: "OK_Practice",
			req: &pb.GetCandlesRequest{
				Names:  "",
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
			req: &pb.GetCandlesRequest{
				Names:  "",
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
			req: &pb.GetCandlesRequest{
				Names:  "",
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
			req: &pb.GetCandlesRequest{
				Names:  "",
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
			res, err := client.RequestCandles(ctx, tc.req)
			tc.CheckResp(res, s.pairs, err)
		})
	}

}
