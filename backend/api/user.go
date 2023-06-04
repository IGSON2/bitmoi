package api

import (
	"bitmoi/backend/utilities"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) loginUser(c *fiber.Ctx) error {
	loginReq := new(LoginUserRequest)
	err := c.BodyParser(&loginReq)
	if errs := utilities.ValidateStruct(loginReq); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	user, err := s.store.GetUser(c.Context(), loginReq.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if err := utilities.CheckPassword(loginReq.Password, user.HashedPassword); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("password is not correct err : %v", err))
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.UserID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.UserID,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

}
