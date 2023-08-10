package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
)

func main() {
	app := fiber.New()
	app.Use(allowOriginMiddleware, limiterMiddleware)
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
