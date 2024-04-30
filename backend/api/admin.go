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

type AdminScoreResponse struct {
	Number        int64   `json:"number"`
	Nickname      string  `json:"nickname"`
	UserID        string  `json:"id"`
	BettingUsdp   float64 `json:"bettingusdp"`
	Position      string  `json:"position"`
	Leverage      int8    `json:"leverage"`
	Roe           float64 `json:"roe"`
	Pnl           float64 `json:"pnl"`
	MaxMinRoe     string  `json:"maxminroe"`
	SubmitTime    string  `json:"submittime"`
	EntryTime     string  `json:"entrytime"`
	ExitTime      string  `json:"exittime"`
	SettledAt     string  `json:"settledat"`
	AfterExitTime string  `json:"afterexittime"`
}

func (s *Server) GetScoresInfo(c *fiber.Ctx) error {
	mode := c.Query("mode", practice)
	switch mode {
	case practice:
		scores, err := s.store.GetAdminScores(c.Context(), db.GetAdminScoresParams{
			Limit:  1000,
			Offset: 0,
		})
		if err != nil {
			s.logger.Err(err).Msg("cannot get scores for admin")
			return c.Status(fiber.StatusInternalServerError).SendString("Cannot get scores")
		}
		var response []AdminScoreResponse
		for _, score := range scores {
			response = append(response, AdminScoreResponse{
				Number:        score.ID,
				Nickname:      score.Nickname,
				UserID:        score.UserID,
				BettingUsdp:   score.Entryprice * score.Quantity / float64(score.Leverage),
				Position:      score.Position,
				Leverage:      score.Leverage,
				Roe:           score.Roe,
				Pnl:           score.Pnl,
				EntryTime:     score.Entrytime,
				ExitTime:      score.Outtime,
				MaxMinRoe:     fmt.Sprintf("%.1f / %.1f", score.MaxRoe, score.MinRoe),
				SubmitTime:    score.CreatedAt.Format("06.01.02 15:04:05"),
				AfterExitTime: fmt.Sprintf("%.2f", float64(score.AfterOuttime/3600)),
			})
			if score.SettledAt.Valid {
				response[len(response)-1].SettledAt = score.SettledAt.Time.Format("06.01.02 15:04:05")
			}
		}
		return c.Status(fiber.StatusOK).JSON(response)
	default:
		return c.Status(fiber.StatusBadRequest).SendString("Invalid mode")
	}
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
