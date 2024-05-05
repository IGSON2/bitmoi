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
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	User                 UserResponse `json:"user"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// reissueAccessToken godoc
// @Summary		access token을 재발급 합니다.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		ReissueAccessTokenRequest	body		api.ReissueAccessTokenRequest	true	"refresh token"
// @Success		200		{object}	api.ReissueAccessTokenResponse
// @Router      /reissueAccess [post]
func (s *Server) reissueAccessToken(c *fiber.Ctx) error {
	r := new(ReissueAccessTokenRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(r.RefreshToken)
	if err != nil {
		s.logger.Error().Err(err).Msg("refresh token verification was failed")
		return c.Status(fiber.StatusUnauthorized).SendString("refresh token verification was failed")
	}

	session, err := s.store.GetSession(c.Context(), refreshPayload.SessionID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("session doesn't exist by token id: %s", err))
		}
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot find session: %s", err))
	}

	if session.IsBlocked {
		return c.Status(fiber.StatusUnauthorized).SendString("session is blocked")
	}

	if session.UserID != refreshPayload.UserID {
		return c.Status(fiber.StatusUnauthorized).SendString("user is incorrect")

	}

	if session.RefreshToken != r.RefreshToken {
		return c.Status(fiber.StatusUnauthorized).SendString("mismatched session token")
	}

	if time.Now().After(session.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).SendString("expired session")
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.UserID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot reissue access token: %s", err))
	}

	user, err := s.store.GetUser(c.Context(), refreshPayload.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("cannot find user")
	}
	userRes := convertUserResponse(user)

	rsp := ReissueAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		User:                 userRes,
	}

	return c.Status(fiber.StatusOK).JSON(rsp)
}

// verifyToken godoc
// @Summary		발급했던 access token을 검증합니다.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		VerifyTokenRequest	body		api.VerifyTokenRequest	true	"access token"
// @Success		200		{object}	api.UserResponse
// @Router      /verifyToken [post]
func (s *Server) verifyToken(c *fiber.Ctx) error {
	r := new(VerifyTokenRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	payload, err := s.tokenMaker.VerifyToken(r.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("token verification failed.")
	}

	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		s.logger.Error().Err(err).Any("user", user).Msg("cannot find user for verify token")
		return c.Status(fiber.StatusNotFound).SendString("cannot find user")
	}
	userRes := convertUserResponse(user)
	return c.Status(fiber.StatusOK).JSON(userRes)
}
