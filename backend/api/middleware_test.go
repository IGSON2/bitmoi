package api

import (
	"bitmoi/backend/token"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func addAuthrization(
	t *testing.T,
	req *http.Request,
	tokenMaker *token.PasetoMaker,
	authType, userID string,
	duration time.Duration) {
	token, payload, err := tokenMaker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.Greater(t, token, "")
	require.NotNil(t, payload)

	bearerToken := fmt.Sprintf("%s %s", authType, token)
	req.Header.Set(authorizationHeaderKey, bearerToken)
}

func TestAuthMiddleware(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	s := newTestServer(t, newTestStore(t), nil)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker)
		checkResponse func(t *testing.T, recorder *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, authorizationTypeBearer, masterID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
		{
			name: "UnAuthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, "not_supported_token", masterID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "InvalidFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, "", masterID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "TokenExpired",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, authorizationTypeBearer, masterID, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := http.DefaultClient
			myscoreReq, err := http.NewRequest(http.MethodGet, apiAddress+"/competition", nil)
			require.NoError(t, err)
			tc.setupAuth(t, myscoreReq, *s.tokenMaker)
			myscoreRes, err := client.Do(myscoreReq)
			require.NoError(t, err)
			tc.checkResponse(t, myscoreRes)
		})
	}
}
