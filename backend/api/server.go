package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/futureclient"
	"bitmoi/backend/utilities"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog"
)

const (
	finalstage = 10
)

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

type Server struct {
	config *utilities.Config
	store  db.Store
	router *fiber.App
	logger *zerolog.Logger
	pairs  []string
}

func NewServer(c *utilities.Config, s db.Store) (*Server, error) {
	fc, err := futureclient.NewFutureClient(c)
	if err != nil {
		return nil, fmt.Errorf("cannot create futureclient during creat server, err : %w", err)
	}
	if err := fc.GetAllPairs(); err != nil {
		return nil, fmt.Errorf("cannot get all future pair during creat server, err : %w", err)
	}
	serverLogger := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	server := &Server{
		config: c,
		store:  s,
		logger: &serverLogger,
		pairs:  fc.Pairs,
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
			return c.Status(fiber.StatusTooManyRequests).SendString("too many request.")
		},
	}))

	router.Get("/practice", server.practice)
	router.Post("/practice", server.practice)
	router.Get("/competition/:array", server.competition)
	router.Post("/competition", server.competition)
	router.Get("/interval", sendInterval)
	router.Get("/myscore", myscore)
	router.Get("/ranking", ranking)
	router.Post("/ranking", ranking)
	router.Get("/moreinfo", moreinfo)
	router.Post("/moreinfo", moreinfo)
	router.Get("/test", server.test)

	server.router = router

	return server, nil
}

func (s *Server) Listen() error {
	return s.router.Listen(s.config.Address)
}

func (s *Server) practice(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		history := utilities.Splitnames(c.Query("names", ""))
		if len(history) >= finalstage {
			//TODO : Handle this
		}
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err := s.makeChartToRef(c.Query("interval", db.FourH), nextPair, PracticeMode, len(history))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		return c.Status(fiber.StatusOK).JSON(oc)
	case "POST":
		var PracticeOrder OrderStruct
		err := c.BodyParser(&PracticeOrder)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}
		r, err := s.createPracResult(&PracticeOrder, nil)
		if err != nil {
			s.logger.Error().Err(err)
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		return c.Status(fiber.StatusOK).JSON(r)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func (s *Server) competition(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		history := utilities.Splitnames(c.Query("names", ""))
		if len(history) >= finalstage {
			//TODO : Handle this
		}
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err := s.makeChartToRef(c.Query("interval", db.FourH), nextPair, CompetitionMode, len(history))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		return c.Status(fiber.StatusOK).JSON(oc)
	case "POST":
		var CompetitionOrder OrderStruct
		err := c.BodyParser(&CompetitionOrder)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}
		compResult, err := s.createCompResult(&CompetitionOrder)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		db.InsertStageScore(CompetitionOrder, compResult.ResultScore)
		return c.Status(fiber.StatusOK).JSON(compResult)
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
	return c.JSON(db.SelectStageScoreDB(q.User, q.Index))
}

func ranking(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		return c.JSON((db.SelectTotalScoreDB()))
	case "POST":
		var t db.TotalData
		err := c.BodyParser(&t)
		utilities.Errchk(err)
		db.InsertTotalScore(t)
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
		return c.JSON((db.SendMoreInfo(q.User, q.Scoreid)))
	case "POST":
		var t PostedComment
		err := c.BodyParser(&t)
		utilities.Errchk(err)
		return db.UpdateComment(t.Comment, t.User)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func (s *Server) test(c *fiber.Ctx) error {
	interval := c.Query("interval")
	name := c.Query("name")
	candles, refTimestamp, err := s.selectCandlesToRef(interval, name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	oc, err := makeChart(candles, CompetitionMode, refTimestamp)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	return c.Status(fiber.StatusOK).JSON(oc)
}
