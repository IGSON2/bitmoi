package api

import (
	"bitmoi/token"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/contrib/websocket"
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
			return abort(c, err.Error())
		}
		c.Locals(authorizationPayloadKey, payload)
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
	return nil
}

func adminAuthMiddleware(adminID string, maker *token.PasetoMaker) fiber.Handler {
	abort := func(c *fiber.Ctx, err string) error {
		return c.Status(fiber.StatusUnauthorized).SendString(err)
	}

	return func(c *fiber.Ctx) error {
		if err := checkAuthorization(c, maker); err != nil {
			return abort(c, err.Error())
		}

		payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
		if payload == nil || !ok {
			return abort(c, "admin token is required")
		}

		if payload.UserID != adminID {
			return abort(c, "not admin user")
		}

		return c.Next()
	}
}

func createNewOriginMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "https://bitmoi.co.kr, https://m.bitmoi.co.kr, http://localhost:3000, http://m.bitmoi.co.kr:444, http://m.bitmoi.co.kr:3000",
	})
}

func createNewLimitMiddleware(cnt int, logger *zerolog.Logger) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cnt,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn().Str("ip", c.Get("x-forwarded-for")).Str("path", c.Path()).Str("method", c.Method()).Str("ua", c.Get("user-agent")).Msg("too many request.")
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many request.")
		},
	})
}

func createWebsocketMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	}
}
