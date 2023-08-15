package main

import (
	_ "bitmoi/frontend/server/docs"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const (
	port = ":443"
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

	redirectMiddleware = func(c *fiber.Ctx) error {
		zlog.Info().Msgf(`[%s]:%s Redirected from "%s://%s%s" to https`, c.IP(), c.Port(), c.Protocol(), c.Hostname(), c.OriginalURL())
		return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
	}
)

func init() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	noTLSApp := fiber.New()
	noTLSApp.Use(redirectMiddleware, loggerMiddleware)
	go func() {
		log.Fatalln(noTLSApp.Listen(":80"))
	}()
	app := fiber.New()

	app.Use(allowOriginMiddleware, limiterMiddleware, loggerMiddleware)
	app.Get("/swagger/*", swagger.HandlerDefault)

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
	log.Fatalln(app.ListenTLS(port, "./server.crt", "./server.key"))
	// log.Fatalln(app.Listen(port))
}
