package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"

	"github.com/gofiber/fiber/v2"
)

const (
	wmoiRows = 15
)

func (s *Server) getWmoiMintingHist(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	if page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("page must be greater than 0")
	}
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get authorization payload")
	}

	histories, err := s.store.GetWmoiMintingHist(c.Context(), db.GetWmoiMintingHistParams{
		ToUser: payload.UserID,
		Title:  db.RecommendationTitle,
		Limit:  wmoiRows,
		Offset: (int32(page) - 1) * wmoiRows,
	})

	if err != nil {
		s.logger.Error().Err(err).Msgf("cannot get recommendation rewards. user_id: %s", payload.UserID)
		return c.Status(fiber.StatusInternalServerError).SendString("cannot get recommendation rewards")
	}

	return c.Status(fiber.StatusOK).JSON(histories)
}
