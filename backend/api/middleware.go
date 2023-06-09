package api

import (
	"bitmoi/backend/token"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	apiKeyHeader            = "X-API-Key"
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
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
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	})
)

// 기록을 갱신한 사용자에게만 접근 권한을 주고 이 권한이 오용되지 않아야함
// DB에 기록된 총 Score를 합산하여 Rank Board에 기록하는 방식이면 되지 않을까?
// 각 Stage에 대한 POST요청 또한 조작이 가능하지 않을까?
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
