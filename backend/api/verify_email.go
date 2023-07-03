package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type VerifyEmailResponse struct {
	IsVerified bool `json:"is_verified"`
}

func (s *Server) VerifyEmail(c *fiber.Ctx) error {
	r := new(VerifyEmailRequest)
	err := c.QueryParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	txResult, err := s.store.VerifyEmailTx(c.Context(), db.VerifyEmailTxParams{
		EmailId:    r.EmailId,
		SecretCode: r.SecretCode,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to verify email")
	}
	c.Status(fiber.StatusOK).SendFile("./welcome.html", false)
	return c.Status(fiber.StatusOK).JSON(&VerifyEmailResponse{IsVerified: txResult.User.IsEmailVerified})
}
