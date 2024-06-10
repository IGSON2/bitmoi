package api

import (
	db "bitmoi/db/sqlc"
	"bitmoi/token"
	"bitmoi/utilities"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	rankRows = 17
)

// getRank godoc
// @Summary      랭크에 등재된 사용자들을 불러옵니다.
// @Tags         rank
// @Param 		 mode path string true "practice, competition"
// @Param 		 category path string false "순위를 집계할 종목"
// @Param 		 start query int false "시작일자 (YY-MM-DD)"
// @Param 		 end query int false "종료일자 (YY-MM-DD)"
// @Param 		 page query int false "페이지"
// @param 		 offset query int false "오프셋"
// @Produce      json
// @Success      200  {array}  db.RankingBoard
// @Router       /rank/{category} [get]
func (s *Server) getRank(c *fiber.Ctx) error {
	now := time.Now()

	r := new(GetRankRequest)
	r.Mode = c.Query("mode", practice)
	r.Category = c.Query("category", "pnl")
	r.Start = c.Query("start", now.AddDate(0, 0, -int(now.Weekday())+1).Format("06-01-02"))
	r.End = c.Query("end", now.AddDate(0, 0, -int(now.Weekday())+8).Format("06-01-02"))
	r.Page = int32(c.QueryInt("page", 1))

	errs := utilities.ValidateStruct(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	}

	start, err := time.Parse("06-01-02", r.Start)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	end, err := time.Parse("06-01-02", r.End)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	switch r.Mode {
	case practice:
		switch r.Category {
		case "pnl":
			pnlRanks, err := s.store.GetUserPracRankByPNL(c.Context(), db.GetUserPracRankByPNLParams{
				CreatedAt:   start,
				CreatedAt_2: end.Add(24 * time.Hour), // 24시간을 더해줘야 해당일까지의 데이터를 가져올 수 있음
				Limit:       rankRows,
				Offset:      (r.Page - 1) * rankRows,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			return c.Status(fiber.StatusOK).JSON(pnlRanks)
		case "roe":
			roeRanks, err := s.store.GetUserPracRankByROE(c.Context(), db.GetUserPracRankByROEParams{
				CreatedAt:   start,
				CreatedAt_2: end,
				Limit:       rankRows,
				Offset:      (r.Page - 1) * rankRows,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			return c.Status(fiber.StatusOK).JSON(roeRanks)
		}
	case competition:
		return c.Status(fiber.StatusBadRequest).SendString("competition mode is not supported")
	}

	return c.Status(fiber.StatusBadRequest).SendString("invalid mode or category")
}

// postRank godoc
// @Summary      사용자를 랭크에 등재합니다.
// @Tags         rank
// @Param 		 rankInsertRequest body api.RankInsertRequest true "랭크 등재 요청에 대한 정보"
// @param 		 Authorization header string true "Authorization"
// @Produce      json
// @Success      200
// @Router       /rank [post]
func (s *Server) postRank(c *fiber.Ctx) error {
	payload := c.Locals(authorizationPayloadKey).(*token.Payload)
	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get user by token payload. err: %w", err).Error())
	}

	var r RankInsertRequest
	err = c.BodyParser(&r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	errs := utilities.ValidateStruct(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	}

	err = s.insertCompScoreToRankBoard(&r, &user, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}
