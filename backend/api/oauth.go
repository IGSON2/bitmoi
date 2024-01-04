package api

import (
	"bitmoi/backend/utilities"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConf *oauth2.Config

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func NewOauthConfig(c *utilities.Config) *oauth2.Config {
	port := strings.Split(c.HTTPAddress, ":")[1]
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://localhost:%s/auth/google/callback", port),
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

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(contents, &jsonMap)
	if err != nil {
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(jsonMap)
}

func (s *Server) GetLoginURL(c *fiber.Ctx) error {
	token, _, err := s.tokenMaker.CreateToken("state", s.config.AccessTokenDuration)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	url := s.oauthConfig.AuthCodeURL(token)
	return c.Status(fiber.StatusOK).SendString(url)
}
