package api

import (
	"bitmoi/backend/utilities"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReissueAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ReissueAccessTokenResponse struct {
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
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
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
