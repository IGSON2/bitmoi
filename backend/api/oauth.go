package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
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

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type OauthData struct {
	Email         string `json:"email"`
	ID            string `json:"id"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func NewOauthConfig(c *utilities.Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "https://api.bitmoi.co.kr/login/callback",
		ClientID:     c.OauthClientID,
		ClientSecret: c.OauthClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (s *Server) CallBackLogin(c *fiber.Ctx) error {
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

	userID := strings.Split(od.Email, "@")[0]

	user, err := s.store.GetUserByEmail(c.Context(), od.Email)

	if user.UserID == "" || err == sql.ErrNoRows {
		_, err = s.store.CreateUser(c.Context(), db.CreateUserParams{
			UserID:   userID,
			OauthUid: sql.NullString{String: od.ID, Valid: true},
			Email:    od.Email,
			PhotoUrl: sql.NullString{String: od.Picture, Valid: true},
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	} else {
		userID = user.UserID
	}

	accessToken, _, err := s.tokenMaker.CreateToken(userID, s.config.AccessTokenDuration)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		userID,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	_, err = s.store.CreateSession(c.Context(), db.CreateSessionParams{
		SessionID:    refreshPayload.SessionID.String(),
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    string(c.Request().Header.UserAgent()),
		ClientIp:     c.IP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	redirectURL := fmt.Sprintf("%s/welcome?accessToken=%s&refreshToken=%s", s.config.OauthRedirectURL, accessToken, refreshToken)

	return c.Redirect(redirectURL, fiber.StatusMovedPermanently)
}

func (s *Server) GetLoginURL(c *fiber.Ctx) error {
	token, _, err := s.tokenMaker.CreateToken("state", s.config.AccessTokenDuration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	url := s.oauthConfig.AuthCodeURL(token)
	return c.Redirect(url, fiber.StatusMovedPermanently)
}
