package api

import (
	"bitmoi/backend/utilities"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReissueAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ReissueAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) reissueAccessToken(c *fiber.Ctx) error {
	r := new(ReissueAccessTokenRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(r.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("refresh token verification was failed: %s", err))
	}

	session, err := s.store.GetSession(c.Context(), refreshPayload.SessionID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("session doesn't exist by token id: %s", err))
		}
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot find session: %s", err))
	}

	if session.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprint("session is blocked"))
	}

	if session.UserID != refreshPayload.UserID {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprint("user is incorrect"))

	}

	if session.RefreshToken != r.RefreshToken {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprint("mismatched session token"))
	}

	if time.Now().After(session.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprint("expired session"))
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.UserID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot reissue access token: %s", err))
	}

	rsp := ReissueAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	return c.Status(fiber.StatusOK).JSON(rsp)
}
