package api

import (
	"bitmoi/backend/token"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog"
)

const (
	apiKeyHeader            = "X-API-Key"
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
	authorizedUserKey       = "authorization_user"
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

		if len(fields) < 2 { // Authorization 값은 "Bearer <token>" 형태로 전달되어야 함
			return abort(c, "invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer { // Auth type은 표준인 Bearer만 허용
			return abort(c, "unsupported authorization type,"+authorizationType)
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			return abort(c, fmt.Sprintf("%v", err))
		}
		c.Locals(authorizationPayloadKey, payload)
		c.Locals(authorizedUserKey, payload.UserID)
		return c.Next()
	}
}

func checkAuthorization(c *fiber.Ctx, maker *token.PasetoMaker) error {

	authorizationHeader := c.Get(authorizationHeaderKey)

	if len(authorizationHeader) == 0 {
		return fmt.Errorf("authorization header is not provided")
	}

	fields := strings.Fields(authorizationHeader)

	if len(fields) < 2 {
		return fmt.Errorf("invalid authorization header format")

	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return fmt.Errorf("unsupported authorization type, %s", authorizationType)
	}

	accessToken := fields[1]
	payload, err := maker.VerifyToken(accessToken)
	if err != nil {
		return err
	}
	c.Locals(authorizationPayloadKey, payload)
	c.Locals(authorizedUserKey, payload.UserID)
	return nil
}

func createNewOriginMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "https://bitmoi.co.kr, https://m.bitmoi.co.kr, http://localhost:3000",
	})
}

func createNewLimitMiddleware(cnt int, logger zerolog.Logger) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cnt,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn().Str("ip", c.Get("x-forwarded-for")).Str("path", c.Path()).Str("method", c.Method()).Str("ua", c.Get("user-agent")).Str("authorization", c.Get(authorizationHeaderKey)).Msg("too many request.")
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many request.")
		},
	})
}
