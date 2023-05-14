package clientAPI

import (
	"github.com/gofiber/fiber/v2"
)

func Start() {
	app := fiber.New()

	app.Static("/", "./frontend/build")
	app.Static("/practice", "./frontend/build")
	app.Static("/competition", "./frontend/build")
	app.Static("/ranking", "./frontend/build")
	app.Static("/myscore", "./frontend/build")
	app.Static("/community", "./frontend/build")
}
