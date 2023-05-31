package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
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

	serverLogger := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	server := &Server{
		config: c,
		store:  s,
		logger: &serverLogger,
	}

	ps, err := server.store.GetAllParisInDB(context.Background())
	if err != nil {
		return nil, err
	}
	server.pairs = ps

	router := fiber.New(fiber.Config{})

	router.Use(allowOriginMiddleware, limiterMiddleware, loggerMiddleware)

	router.Get("/practice", server.practice)
	router.Post("/practice", server.practice)
	router.Get("/competition", server.competition)
	router.Post("/competition", server.competition)
	// router.Get("/interval", sendInterval)
	router.Get("/myscore/:user", server.myscore)
	router.Get("/rank", server.rank)
	router.Post("/rank", server.rank)
	router.Get("/moreinfo", server.moreinfo)

	server.router = router

	return server, nil
}

func (s *Server) Listen() error {
	return s.router.Listen(s.config.Address)
}

func (s *Server) practice(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		r := new(ChartRequestQuery)
		err := c.QueryParser(r)
		if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
		}

		history := utilities.Splitnames(r.Names)
		if len(history) >= finalstage {
			return c.Status(fiber.StatusBadRequest).SendString("invalid current stage")
		}
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err := s.makeChartToRef(r.Interval, nextPair, PracticeMode, len(history), c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
		}
		return c.Status(fiber.StatusOK).JSON(oc)
	case "POST":
		var PracticeOrder OrderRequest
		// TODO : 유효한 주문인지 검사 필요 e.b 가격*수량 <= lev * bal
		err := c.BodyParser(&PracticeOrder)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(err))
		}

		errs := utilities.ValidateStruct(PracticeOrder)
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(errs.Error()))
		}

		r, err := s.createPracResult(&PracticeOrder, c)
		if err != nil {
			s.logger.Error().Err(err)
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
		}
		return c.Status(fiber.StatusOK).JSON(r)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func (s *Server) competition(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		r := new(ChartRequestQuery)
		err := c.QueryParser(r)
		if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
		}
		history := utilities.Splitnames(r.Names)

		if len(history) >= finalstage {
			return c.Status(fiber.StatusBadRequest).SendString("invalid current stage")
		}
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err := s.makeChartToRef(r.Interval, nextPair, CompetitionMode, len(history), c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
		}
		return c.Status(fiber.StatusOK).JSON(oc)
	case "POST":
		var CompetitionOrder OrderRequest
		// TODO : 유효한 주문인지 검사 필요 e.b 가격*수량 <= lev * bal
		// position에 따른 entryprice 대비 profit과 loss값 위치에 대한 검증
		// entryprice는 identifier에 삽입하여 전송해도 되지 않을까?
		err := c.BodyParser(&CompetitionOrder)
		if err != nil || CompetitionOrder.Mode != competition {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("%s, mode : %s", err, CompetitionOrder.Mode))
		}

		errs := utilities.ValidateStruct(CompetitionOrder)
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(errs.Error()))
		}

		compResult, err := s.createCompResult(&CompetitionOrder, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
		}
		err = s.insertUserScore(&CompetitionOrder, compResult.ResultScore, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
		}
		return c.Status(fiber.StatusOK).JSON(compResult)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

// func sendInterval(c *fiber.Ctx) error {
// 	i := new(IntervalQuery)
// 	if err := c.QueryParser(i); err != nil {
// 		return err
// 	}
// 	return c.JSON(db.SendOtherInterval(i.Identifier, i.ReqInterval, i.Mode))
// }

func (s *Server) myscore(c *fiber.Ctx) error {
	u := c.Params("user")
	p := new(PageRequest)
	err := c.QueryParser(p)
	if errs := utilities.ValidateStruct(*p); err != nil || errs != nil {
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(errs.Error()))
		}
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(err))
	}

	scores, err := s.getMyscores(u, p.Page, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}

func (s *Server) rank(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		p := new(PageRequest)
		err := c.QueryParser(p)
		if errs := utilities.ValidateStruct(*p); err != nil || errs != nil {
			if errs != nil {
				return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(errs.Error()))
			}
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(err))
		}

		ranks, err := s.getAllRanks(p.Page, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
		}
		return c.Status(fiber.StatusOK).JSON(ranks)
	case "POST":
		var r RankInsertRequest
		err := c.BodyParser(&r)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(err))
		}

		errs := utilities.ValidateStruct(r)
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(errs.Error()))
		}
		s.insertScoreToRankBoard(&r, c)
		return nil
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func (s *Server) moreinfo(c *fiber.Ctx) error {
	q := new(MoreInfoRequest)
	if err := c.QueryParser(q); err != nil {
		return err
	}
	errs := utilities.ValidateStruct(q)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprint(errs.Error()))
	}
	scores, err := s.sendMoreInfo(q.UserId, q.ScoreId, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprint(err))
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}
