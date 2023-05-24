package api

import "github.com/gofiber/fiber/v2"

const (
	apiKeyHeader = "X-API-Key"
)

// 기록을 갱신한 사용자에게만 접근 권한을 주고 이 권한이 오용되지 않아야함
// DB에 기록된 총 Score를 합산하여 Rank Board에 기록하는 방식이면 되지 않을까?
// 각 Stage에 대한 POST요청 또한 조작이 가능하지 않을까?
func ApiAuthMiddleWare(c *fiber.Ctx) error {
	// apikey := c.Get(apiKeyHeader)
	return nil
}
