package api

import (
	db "bitmoi/backend/db/sqlc"
	btoken "bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	reqPathKey        = "req_url"
)

type OauthData struct {
	Email         string `json:"email"`
	ID            string `json:"id"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func NewOauthConfig(c *utilities.Config) *oauth2.Config {
	redirURL := fmt.Sprintf("http://localhost:%s", strings.Split(c.HTTPAddress, ":")[1])
	if c.Environment == common.EnvProduction {
		redirURL = "https://api.bitmoi.co.kr"
	}
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/login/callback", redirURL),
		ClientID:     c.OauthClientID,
		ClientSecret: c.OauthClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (s *Server) CallBackLogin(c *fiber.Ctx) error {
	state := c.Query("state", "practice")
	rPayload, err := s.tokenMaker.VerifyToken(state)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("state mismatched.")
	}
	rPath := rPayload.UserID

	code := c.Query("code")
	token, err := s.oauthConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	od := new(OauthData)
	err = json.Unmarshal(contents, &od)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	userId := od.Email

	user, err := s.store.GetUserByEmail(c.Context(), od.Email)

	if user.UserID == "" || err != nil {
		if err != sql.ErrNoRows {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		code, err := btoken.GenerateRecCode()
		if err != nil {
			s.logger.Error().Err(err).Str("user id", userId).Msg("cannot generate recommender code")
			return c.Status(fiber.StatusInternalServerError).SendString("cannot generate recommender code")
		}

		idNum, err := s.GetLastUserID(c.Context())
		if err != nil {
			s.logger.Error().Err(err).Str("user id", userId).Msg("cannot get last user id")
			return c.Status(fiber.StatusInternalServerError).SendString("cannot generate nickname")
		}

		_, createErr := s.store.CreateUser(c.Context(), db.CreateUserParams{
			UserID:          od.Email,
			OauthUid:        sql.NullString{String: od.ID, Valid: true},
			Nickname:        sql.NullString{String: fmt.Sprintf("Chartist %d", idNum+1), Valid: true},
			Email:           od.Email,
			PhotoUrl:        sql.NullString{String: od.Picture, Valid: true},
			RecommenderCode: sql.NullString{String: code, Valid: true},
		})
		if createErr != nil {
			s.logger.Error().Err(err).Str("user id", userId).Msg("cannot create user")
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		rPath = "welcome"

	} else {
		userId = user.UserID
	}

	accessToken, _, err := s.tokenMaker.CreateToken(userId, s.config.AccessTokenDuration)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		userId,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	_, err = s.store.CreateSession(c.Context(), db.CreateSessionParams{
		SessionID:    refreshPayload.SessionID.String(),
		UserID:       userId,
		RefreshToken: refreshToken,
		UserAgent:    string(c.Request().Header.UserAgent()),
		ClientIp:     c.IP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	rewardStr := fmt.Sprintf("%d", db.AttendanceReward)

	redirectURL := fmt.Sprintf("%s/auth?accessToken=%s&refreshToken=%s&path=%s&attendanceReward=", s.config.OauthRedirectURL, accessToken, refreshToken, rPath)

	err = s.checkAttendance(c.Context(), userId)
	if err != nil {
		s.logger.Warn().Err(err).Str("user id", userId).Msg("cannot check attendance")
		rewardStr = ""
	}

	return c.Redirect(redirectURL+rewardStr, fiber.StatusMovedPermanently)
}

func (s *Server) GetLoginURL(c *fiber.Ctx) error {
	rPath := c.Params(reqPathKey, "practice")

	token, _, err := s.tokenMaker.CreateToken(rPath, s.config.AccessTokenDuration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	url := s.oauthConfig.AuthCodeURL(token)
	return c.Redirect(url, fiber.StatusMovedPermanently)
}
