package api

import (
	"bitmoi/backend/token"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
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
	s := newTestServer(t, newTestStore(t), nil)

	httpReq, _ := randomUserRequest(t)
	res, err := s.router.Test(httpReq)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.StatusCode)

	user := new(UserResponse)
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	json.Unmarshal(body, user)
	require.NotNil(t, user)

	s.router.Get("/auth", authMiddleware(s.tokenMaker), func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("Authorization passed")
	})

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker)
		checkResponse func(t *testing.T, recorder *http.Response)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, authorizationTypeBearer, user.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
		{
			name: "UnAuthorized",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, "not_supported_token", user.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "InvalidFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, "", user.UserID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "TokenExpired",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, &tokenMaker, authorizationTypeBearer, user.UserID, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *http.Response) {
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authPathReq, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)
			tc.setupAuth(t, authPathReq, *s.tokenMaker)
			authPathRes, err := s.router.Test(authPathReq)
			require.NoError(t, err)
			tc.checkResponse(t, authPathRes)
		})
	}
}
