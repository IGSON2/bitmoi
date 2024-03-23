package api

import "github.com/gofiber/fiber/v2"

func (s *Server) GetUsers(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) GetInvestInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) GetUsdpInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) GetTokenInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func (s *Server) GetReferralInfo(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
