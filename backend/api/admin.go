package api

import (
	db "bitmoi/backend/db/sqlc"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AdminUserResponse struct {
	Number     int64     `json:"number"`
	Nickname   string    `json:"nickname"`
	UserID     string    `json:"id"`
	Usdp       float64   `json:"usdp"`
	Token      float64   `json:"token"`
	Attendance int       `json:"attendance"`
	PracScore  string    `json:"prac"`
	CompScore  string    `json:"comp"`
	Referral   int       `json:"referral"`
	RecomCode  string    `json:"recom"`
	SignUpDate time.Time `json:"signup"`
	LastAccess time.Time `json:"lastaccess"`
}

func (s *Server) GetUsers(c *fiber.Ctx) error {
	users, err := s.store.GetUsers(c.Context(), db.GetUsersParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot get users")
	}
	var response []AdminUserResponse
	for _, user := range users {
		response = append(response, AdminUserResponse{
			Number:     user.ID,
			Nickname:   user.Nickname,
			UserID:     user.UserID,
			Usdp:       user.PracBalance,
			Token:      user.WmoiBalance,
			Attendance: 0,
			PracScore:  "0/0/0",
			CompScore:  "0/0/0",
			Referral:   0,
			RecomCode:  user.RecommenderCode,
			SignUpDate: user.CreatedAt,
			LastAccess: user.LastAccessedAt.Time,
		})
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

type AdminInvestResponse struct {
	Number       int     `json:"number"`
	UserID       string  `json:"id"`
	BettingUsdp  float64 `json:"bettingusdp"`
	Position     string  `json:"position"`
	Leverage     int     `json:"leverage"`
	Roe          int     `json:"roe"`
	Pnl          int     `json:"pnl"`
	PositionTime string  `json:"positiontime"`
	OwnedTime    string  `json:"ownedtime"`
	MaxRoe       string  `json:"maxroe"`
	SubmitTime   string  `json:"submittime"`
	EntryTime    string  `json:"entrytime"`
	ExitTime     string  `json:"exittime"`
}

func (s *Server) GetInvestInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) GetUsdpInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) SetUsdpInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) GetTokenInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) GetReferralInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
