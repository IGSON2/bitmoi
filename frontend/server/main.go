package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	port = ":80"
)

var (
	allowOriginMiddleware = cors.New(cors.Config{
		AllowOrigins: "*",
	})
	limiterMiddleware = limiter.New(limiter.Config{
		Max:        30,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString("too many request.")
		},
	})
	loggerMiddleware = logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${time} ${status} - ${method} ${path} - ${latency}\n",
		TimeFormat: "2006-01-02T15:04:05",
	})
)

func main() {
	app := fiber.New()
	app.Use(allowOriginMiddleware, limiterMiddleware, loggerMiddleware)

	app.Static("/.well-known/acme-challenge", "./acme-challenge")

	app.Static("/", "./build")
	app.Static("/competition", "./build")
	app.Static("/practice", "./build")
	app.Static("/community", "./build")
	app.Static("/ad_bidding", "./build")
	app.Static("/mypage", "./build")
	app.Static("/rank", "./build")
	app.Static("/login", "./build")
	app.Static("/signup", "./build")
	app.Static("/goto/:domain", "./build")
	app.Static("/freetoken", "./build")
	log.Fatalln(app.Listen(port))
}
