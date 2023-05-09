package api

import (
	"bitmoi/backend/db/scoreData"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Server struct {
	config utilities.Config
	store  db.Store
	router *fiber.App
}

func NewServer(c utilities.Config, s db.Store) (*Server, error) {
	server := &Server{
		config: c,
		store:  s,
	}

	router := fiber.New()

	router.Use(cors.New(cors.Config{
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

	router.Get("/competition", competition)
	router.Get("/competition/:array", competition)
	router.Post("/competition", competition)
	router.Get("/practice", practice)
	router.Get("/practice/:array", practice)
	router.Post("/practice", practice)
	router.Get("/interval", sendInterval)
	router.Get("/myscore", myscore)
	router.Get("/ranking", ranking)
	router.Post("/ranking", ranking)
	router.Get("/moreinfo", moreinfo)
	router.Post("/moreinfo", moreinfo)

	server.router = router

	return server, nil
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
