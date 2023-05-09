package clientAPI

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type URLDescription struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type IntervalQuery struct {
	ReqInterval string `query:"reqinterval"`
	Identifier  string `query:"identifier"`
	Mode        string `query:"mode"`
}

type UserQuery struct {
	User    string `query:"user"`
	Scoreid string `query:"scoreid"`
	Index   int    `query:"index"`
}

type PostedComment struct {
	User    string `json:"user"`
	Comment string `json:"comment"`
}

const port = ":80"

func Start() {
	app := fiber.New()

	app.Static("/", "./frontend/build")
	app.Static("/practice", "./frontend/build")
	app.Static("/competition", "./frontend/build")
	app.Static("/ranking", "./frontend/build")
	app.Static("/myscore", "./frontend/build")
	app.Static("/community", "./frontend/build")
	log.Panic(app.Listen(port))
}
