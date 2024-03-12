package api

import (
	db "bitmoi/backend/db/sqlc"
	btoken "bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	allowRpathes = []string{"practice", "mypage", "welcome"}
)

const (
	PlatformGoogle = "google"
	PlatformKakao  = "kakao"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	reqPathKey        = "req_url"
	platformKey       = "platform"
)

type GoogleOauthData struct {
	Email         string `json:"email"`
	ID            string `json:"id"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func NewGoogleOauthConfig(c *utilities.Config) *oauth2.Config {
	redirURL := fmt.Sprintf("http://localhost:%s/basic/login/google", strings.Split(c.HTTPAddress, ":")[1])
	if c.Environment == common.EnvProduction {
		redirURL = "https://api.bitmoi.co.kr/basic/login/google"
	}
	return &oauth2.Config{
		RedirectURL:  redirURL,
		ClientID:     c.GoogleOauthClientID,
		ClientSecret: c.GoogleOauthClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (s *Server) GoogleLogin(c *fiber.Ctx) error {
	rPath := c.Query("state")
	if rPath == "" || !strings.Contains(strings.Join(allowRpathes, ""), rPath) {
		s.logger.Warn().Str("platform", "google").Msgf("path is invalid. rPath: %s", rPath)
		if strings.HasPrefix(rPath, "v2") {
			p, err := s.tokenMaker.VerifyToken(rPath)
			if p != nil && err != nil {
				rPath = p.UserID
			}
		} else {
			rPath = "practice"
		}
	}

	code := c.Query("code")
	token, err := s.googleOauthCfg.Exchange(c.Context(), code)
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

	od := new(GoogleOauthData)
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
			Nickname:        fmt.Sprintf("Chartist%02d", idNum+1),
			Email:           od.Email,
			PhotoUrl:        sql.NullString{String: od.Picture, Valid: true},
			RecommenderCode: code,
		})
		if createErr != nil {
			s.logger.Error().Err(createErr).Str("user id", userId).Msg("cannot create user")
			return c.Status(fiber.StatusInternalServerError).SendString(createErr.Error())
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

type KakaoOauthData struct {
	Nicname string `json:"nickname"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
	Sub     string `json:"sub"` // user id
}

func (s *Server) KakaoLogin(c *fiber.Ctx) error {
	rPath := c.Query("state")
	if rPath == "" || !strings.Contains(strings.Join(allowRpathes, ""), rPath) {
		s.logger.Warn().Str("platform", "kakao").Msgf("path is invalid. rPath: %s", rPath)
		if strings.HasPrefix(rPath, "v2") {
			p, err := s.tokenMaker.VerifyToken(rPath)
			if p != nil && err != nil {
				rPath = p.UserID
			}
		} else {
			rPath = "practice"
		}
	}

	code := c.Query("code")

	redirURL := fmt.Sprintf("http://localhost:%s/basic/login/kakao", strings.Split(s.config.HTTPAddress, ":")[1])
	if s.config.Environment == common.EnvProduction {
		redirURL = "https://api.bitmoi.co.kr/basic/login/kakao"
	}

	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("client_id", s.config.KakaoOauthClientID)
	// v.Set("redirect_uri", fmt.Sprintf("%s/login/%s", s.config.OauthRedirectURL, rPath))
	v.Set("redirect_uri", redirURL)
	v.Set("code", code)

	req, err := http.NewRequest("POST", "https://kauth.kakao.com/oauth/token", strings.NewReader(v.Encode()))
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("charset", "utf-8")

	response, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	td := struct {
		IdToken string `json:"id_token"`
	}{}
	err = json.Unmarshal(contents, &td)
	if err != nil || td.IdToken == "" {
		return c.Status(fiber.StatusForbidden).SendString("invalid id token")
	}

	payload := strings.Split(td.IdToken, ".")[1]
	idBytes := utilities.Base64Decode(payload)
	if idBytes == nil {
		return c.Status(fiber.StatusForbidden).SendString("cannot decode id token")
	}

	od := new(KakaoOauthData)
	err = json.Unmarshal(idBytes, &od)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	userId := od.Email
	if strings.Contains(userId, "gmail") {
		s.logger.Debug().Str("user id", userId).Msg("kakao user id is gmail")
	}

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

		initNick := od.Nicname
		if initNick == "" {
			idNum, err := s.GetLastUserID(c.Context())
			if err != nil {
				s.logger.Error().Err(err).Str("user id", userId).Msg("cannot get last user id")
				return c.Status(fiber.StatusInternalServerError).SendString("cannot generate nickname")
			}

			initNick = fmt.Sprintf("Chartist%02d", idNum+1)
		}

		_, createErr := s.store.CreateUser(c.Context(), db.CreateUserParams{
			UserID:          od.Email,
			OauthUid:        sql.NullString{String: od.Sub, Valid: true},
			Nickname:        initNick,
			Email:           od.Email,
			PhotoUrl:        sql.NullString{String: od.Picture, Valid: true},
			RecommenderCode: code,
		})
		if createErr != nil {
			s.logger.Error().Err(createErr).Str("user id", userId).Msg("cannot create user")
			return c.Status(fiber.StatusInternalServerError).SendString(createErr.Error())
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
	platform := c.Query(platformKey)
	if platform == "" {
		return c.Status(fiber.StatusBadRequest).SendString("platform is required")
	}
	rpath := c.Params(reqPathKey)
	if !strings.Contains(strings.Join(allowRpathes, ""), rpath) {
		s.logger.Warn().Str("path", rpath).Msg("path is invalid.")
		rpath = "practice"
	}

	url := ""
	switch platform {
	case PlatformGoogle:
		url = s.googleOauthCfg.AuthCodeURL(rpath, oauth2.SetAuthURLParam("prompt", "select_account"))
	case PlatformKakao:
		redirURL := fmt.Sprintf("http://localhost:%s/basic/login/kakao", strings.Split(s.config.HTTPAddress, ":")[1])
		if s.config.Environment == common.EnvProduction {
			redirURL = "https://api.bitmoi.co.kr/basic/login/kakao"
		}
		url = fmt.Sprintf("https://kauth.kakao.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&state=%s", s.config.KakaoOauthClientID, redirURL, rpath)
	}
	s.logger.Info().Str("redirect url", url)
	return c.Redirect(url, fiber.StatusMovedPermanently)
}

func (s *Server) GetLastUserID(ctx context.Context) (int64, error) {
	lastID, err := s.store.GetLastUserID(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return lastID, nil
}
