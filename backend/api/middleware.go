package api

import (
	"bitmoi/backend/token"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const (
	apiKeyHeader            = "X-API-Key"
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

var (
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

func authMiddleware(maker *token.PasetoMaker) fiber.Handler {
	abort := func(c *fiber.Ctx, err string) error {
		return c.Status(fiber.StatusUnauthorized).SendString(err)
	}

	return func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			return abort(c, "authorization header is not provided")
		}

		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			return abort(c, "invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return abort(c, "unsupported authorization type,"+authorizationType)
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			return abort(c, fmt.Sprintf("%v", err))
		}
		c.Locals(authorizationPayloadKey, payload)
		return c.Next()
	}
}
