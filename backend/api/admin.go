package api

import (
	db "bitmoi/backend/db/sqlc"
	"fmt"
	"math"
	"strings"
	"time"

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
	LastAccess string  `json:"last_access"`
}

var asiaSeoul, _ = time.LoadLocation("Asia/Seoul")

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
			Usdp:       math.Round(100*user.PracBalance) / 100,
			Token:      user.WmoiBalance,
			Attendance: user.Attendance,
			PracScore:  fmt.Sprintf("%d/%d/%d", user.PracWin+user.PracLose, user.PracWin, user.PracLose),
			CompScore:  fmt.Sprintf("%d/%d/%d", user.CompWin+user.CompLose, user.CompWin, user.CompLose),
			Referral:   user.Referral,
			RecomCode:  user.RecommenderCode,
			SignUpDate: user.CreatedAt.Format("06.01.02 15:04:05"),
			// Mysql DB의 time_zone이 Asia/Seoul이여도, Default로 생성되는 값에만 적용되고 Update로 변환되는 argument에는 적용되지 않는다.
			// LastAccess 필드는 Default로 생성되지 않고, update되기 때문에 client에게 제공할 때 별도의 time zone 변환이 필요하다.
			LastAccess: user.LastAccessedAt.Time.In(asiaSeoul).Format("06.01.02 15:04:05"),
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
				BettingUsdp:   math.Round(1000*score.Entryprice*score.Quantity/float64(score.Leverage)) / 1000,
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

type AdminUsdpInfoResponse struct {
	Number    int64   `json:"number"`
	Nickname  string  `json:"nickname"`
	UserID    string  `json:"id"`
	Amount    float64 `json:"amount"`
	Title     string  `json:"title"`
	Method    string  `json:"method"`
	Giver     string  `json:"giver"`
	CreatedAt string  `json:"created_at"`
}

func (s *Server) GetPracUsdpInfo(c *fiber.Ctx) error {
	infos, err := s.store.GetAdminUsdpInfo(c.Context(), db.GetAdminUsdpInfoParams{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).SendString("Cannot get usdp info")
	}

	var response []AdminUsdpInfoResponse
	for _, info := range infos {
		response = append(response, AdminUsdpInfoResponse{
			Number:    info.ID,
			Nickname:  info.Nickname,
			UserID:    info.ToUser,
			Amount:    info.Amount,
			Title:     info.Title,
			Method:    info.Method,
			Giver:     info.Giver,
			CreatedAt: info.CreatedAt.Format("06.01.02 15:04:05"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

type AdminSetUsdpRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Title  string  `json:"title"`
}

func (s *Server) SetPracUsdpInfo(c *fiber.Ctx) error {
	req := new(AdminSetUsdpRequest)
	err := c.BodyParser(req)
	if err != nil {
		s.logger.Err(err).Msg("cannot parse admin set usdp info request")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	err = s.store.AppendPracBalanceTx(c.Context(), db.AppendPracBalanceTxParams{
		UserID: req.UserID,
		Amount: req.Amount,
		Title:  req.Title,
		Giver:  "관리자",
		Method: "수동",
	})

	if err != nil {
		s.logger.Err(err).Msg("cannot set usdp info")
		if strings.Contains(err.Error(), "chk_prac_bal") {
			return c.Status(fiber.StatusBadRequest).SendString("Cannot set usdp under 0")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot set usdp info")
	}

	return c.SendStatus(fiber.StatusOK)
}

type AdminTokenResponse struct{}

func (s *Server) GetTokenInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

type AdminReferralResponse struct{}

func (s *Server) GetReferralInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
