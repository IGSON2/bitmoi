package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestMakeChart(t *testing.T) {
	store := newTestStore(t)
	s := newTestServer(t, store)

	httpReq, _ := randomUserRequest(t)

	httpRes, err := s.router.Test(httpReq)
	require.NoError(t, err)
	require.NotNil(t, httpRes)
	require.Equal(t, http.StatusOK, httpRes.StatusCode)

	user := new(UserResponse)
	body, err := ioutil.ReadAll(httpRes.Body)
	require.NoError(t, err)

	json.Unmarshal(body, user)
	require.NotNil(t, user)

	testCases := []struct {
		Name      string
		Path      string
		Method    string
		SetUpAuth func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker)
		Body      interface{}
		CheckResp func(resp *http.Response, pairs []string)
	}{
		{
			Name:   "GetPracticeChart",
			Path:   "/practice",
			Method: http.MethodGet,
			Body: ChartRequestQuery{
				Names: "",
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {},
			CheckResp: func(resp *http.Response, pairs []string) {
				oc := new(OnePairChart)

				require.Equal(t, fiber.StatusOK, resp.StatusCode)
				b, err := ioutil.ReadAll(resp.Body)

				require.NotNil(t, b)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(b, oc))
				require.Greater(t, oc.Identifier, "")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(oc.Identifier), info))

				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, oc.BtcRatio, float64(0))
				require.Greater(t, oc.EntryTime, "")

				require.NotNil(t, oc.OneChart.PData)
				require.NotNil(t, oc.OneChart.VData)

				require.Contains(t, pairs, info.Name)
				require.Equal(t, info.Interval, db.OneH)
				require.Equal(t, info.PriceFactor, float64(0))
				require.Equal(t, info.VolumeFactor, float64(0))
				require.Equal(t, info.TimeFactor, int64(0))
			},
		},
		{
			Name:   "GetCompetitonChart",
			Path:   "/competition",
			Method: http.MethodGet,
			Body: ChartRequestQuery{
				Names: "",
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, s.tokenMaker, authorizationTypeBearer, user.UserID, time.Minute)
			},
			CheckResp: func(resp *http.Response, pairs []string) {
				oc := new(OnePairChart)

				require.Equal(t, fiber.StatusOK, resp.StatusCode)
				b, err := ioutil.ReadAll(resp.Body)

				require.NotNil(t, b)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(b, oc))
				require.Greater(t, oc.Identifier, "")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(oc.Identifier), info))

				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, oc.BtcRatio, float64(0))
				require.Greater(t, oc.EntryTime, "")

				require.NotNil(t, oc.OneChart.PData)
				require.NotNil(t, oc.OneChart.VData)

				require.Contains(t, pairs, info.Name)
				require.Equal(t, info.Interval, db.OneH)
				require.Greater(t, info.PriceFactor, float64(0))
				require.Greater(t, info.VolumeFactor, float64(0))
				require.Greater(t, info.TimeFactor, int64(0))
			},
		},
		{
			Name:   "CompetitionUnAuthorized",
			Path:   "/competition",
			Method: http.MethodGet,
			Body: ChartRequestQuery{
				Names: "",
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
			},
			CheckResp: func(resp *http.Response, pairs []string) {
				require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			b, err := json.Marshal(tc.Body)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.Method, tc.Path, bytes.NewReader(b))
			require.NoError(t, err)

			tc.SetUpAuth(t, req, *s.tokenMaker)

			res, err := s.router.Test(req)
			require.NoError(t, err)

			tc.CheckResp(res, s.pairs)
		})
	}

}
