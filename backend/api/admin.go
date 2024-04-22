package api

import (
	db "bitmoi/backend/db/sqlc"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AdminUserResponse struct {
	Number     int64   `json:"number"`
	Nickname   string  `json:"nickname"`
	UserID     string  `json:"id"`
	Usdp       float64 `json:"usdp"`
	Token      float64 `json:"token"`
	Attendance int64   `json:"attendance"`
	PracScore  string  `json:"prac"`
	CompScore  string  `json:"comp"`
	Referral   int64   `json:"referral"`
	RecomCode  string  `json:"recom"`
	SignUpDate string  `json:"signup"`
	LastAccess string  `json:"lastaccess"`
}

func (s *Server) GetUsers(c *fiber.Ctx) error {
	users, err := s.store.GetAdminUsers(c.Context(), db.GetAdminUsersParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		s.logger.Err(err).Msg("cannot get users for admin")
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
			Attendance: user.Attendance,
			PracScore:  fmt.Sprintf("%d/%d/%d", user.PracWin+user.PracLose, user.PracWin, user.PracLose),
			CompScore:  fmt.Sprintf("%d/%d/%d", user.CompWin+user.CompLose, user.CompWin, user.CompLose),
			Referral:   user.Referral,
			RecomCode:  user.RecommenderCode,
			SignUpDate: user.CreatedAt.Format("06.01.02 15:04:05"),
			LastAccess: user.LastAccessedAt.Time.Format("06.01.02 15:04:05"),
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

type AdminUsdpResponse struct {
}

func (s *Server) GetUsdpInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

type AdminSetUsdpRequest struct{}

func (s *Server) SetUsdpInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

type AdminTokenResponse struct{}

func (s *Server) GetTokenInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

type AdminReferralResponse struct{}

func (s *Server) GetReferralInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
