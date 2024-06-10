package api

import (
	db "bitmoi/db/sqlc"
	"bitmoi/token"
	"bitmoi/utilities"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) rewardRecommender(c *fiber.Ctx) error {
	r := new(CreateRecommendHistoryRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get authorization payload")
	}

	recmNick, err := s.store.RewardRecommenderTx(c.Context(), db.RewardRecommenderTxParams{
		NewMember:       payload.UserID,
		RecommenderCode: r.Code,
	})

	if err != nil {
		s.logger.Error().Err(err).Msg("cannot reward recommender")
		if err == db.ErrRecommenderNotFound {
			return c.Status(fiber.StatusBadRequest).SendString("cannot find recommender.")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("cannot reward recommender")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"recommender": recmNick})
}
