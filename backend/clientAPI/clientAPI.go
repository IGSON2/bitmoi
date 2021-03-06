package clientAPI

import (
	db "bitmoi/backend/db/chartData"
	"bitmoi/backend/db/scoreData"
	"bitmoi/backend/utilities"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}), limiter.New(limiter.Config{
		Max:        30,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendString("too many request.")
		},
	}))

	app.Static("/", "./frontend/build")
	app.Static("/practice", "./frontend/build")
	app.Static("/competition", "./frontend/build")
	app.Static("/ranking", "./frontend/build")
	app.Static("/myscore", "./frontend/build")
	app.Static("/community", "./frontend/build")

	apiGroup := app.Group("/api")
	{
		apiGroup.Get("/", apihome)
		apiGroup.Get("/competition", competition)
		apiGroup.Get("/competition/:array", competition)
		apiGroup.Post("/competition", competition)
		apiGroup.Get("/practice", practice)
		apiGroup.Get("/practice/:array", practice)
		apiGroup.Post("/practice", practice)
		apiGroup.Get("/interval", sendInterval)
		apiGroup.Get("/myscore", myscore)
		apiGroup.Get("/ranking", ranking)
		apiGroup.Post("/ranking", ranking)
		apiGroup.Get("/moreinfo", moreinfo)
		apiGroup.Post("/moreinfo", moreinfo)
	}
	log.Panic(app.Listen(port))
}

func apihome(c *fiber.Ctx) error {
	data := []URLDescription{
		{
			URL:         "/api/competition",
			Description: "?????? ????????? ?????? ????????? ?????? ???????????? ???????????????.",
		},
		{
			URL:         "/api/practice",
			Description: "?????? ????????? ?????? ?????? ?????? ???????????? ???????????????.",
		},
		{
			URL:         "/api/interval/?interval=5m&mode=practice&identfier=",
			Description: "????????? ?????? ???????????? ??????(competition,practice)?????? ????????????(5m,15m,4h)??? ?????????(string)??? ?????? ?????? ???????????? ???????????????.",
		},
		{
			URL:         "/api/myscore",
			Description: "?????????????????? ????????? ???????????? ????????? ???????????? ???????????????.",
		},
		{
			URL:         "/api/ranking",
			Description: "?????? ????????? ?????? ????????? ??? ?????? 30?????? ????????? ???????????????.",
		},
		{
			URL:         "/api/moreinfo",
			Description: "????????? ????????? ???????????? ?????? ????????? ???????????????.",
		},
	}
	return c.JSON(data)
}

func competition(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		names := c.Params("array")
		return c.JSON(db.SendCharts(db.CompetitionMode, db.OneH, splitnames(names)))
	case "POST":
		var CompetitionOrder db.OrderStruct
		err := c.BodyParser(&CompetitionOrder)
		utilities.Errchk(err)
		compResult := db.SendCompResult(CompetitionOrder)
		scoreData.InsertStageScore(CompetitionOrder, compResult.ResultScore)
		return c.JSON(compResult)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}
func practice(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		names := c.Params("array")
		return c.JSON(db.SendCharts(db.PracticeMode, db.OneH, splitnames(names)))
	case "POST":
		var PracticeOrder db.OrderStruct
		err := c.BodyParser(&PracticeOrder)
		utilities.Errchk(err)
		return c.JSON(db.SendPracResult(PracticeOrder))
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func sendInterval(c *fiber.Ctx) error {
	i := new(IntervalQuery)
	if err := c.QueryParser(i); err != nil {
		return err
	}
	return c.JSON(db.SendOtherInterval(i.Identifier, i.ReqInterval, i.Mode))
}

func myscore(c *fiber.Ctx) error {
	q := new(UserQuery)
	if err := c.QueryParser(q); err != nil {
		return err
	}
	return c.JSON(scoreData.SelectStageScoreDB(q.User, q.Index))
}

func ranking(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		return c.JSON((scoreData.SelectTotalScoreDB()))
	case "POST":
		var t scoreData.TotalData
		err := c.BodyParser(&t)
		utilities.Errchk(err)
		scoreData.InsertTotalScore(t)
		return nil
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func moreinfo(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		q := new(UserQuery)
		if err := c.QueryParser(q); err != nil {
			return err
		}
		return c.JSON((scoreData.SendMoreInfo(q.User, q.Scoreid)))
	case "POST":
		var t PostedComment
		err := c.BodyParser(&t)
		utilities.Errchk(err)
		return scoreData.UpdateComment(t.Comment, t.User)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func splitnames(names string) []string {
	var splited []string
	if names != "" {
		withNilSlice := strings.Split(names, ",")
		for _, str := range withNilSlice {
			if str != "" && !strings.Contains(str, "STAGE") {
				splited = append(splited, str)
			}
		}
	}
	return splited
}
