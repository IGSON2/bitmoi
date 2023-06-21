package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	finalstage  = 10
	competition = "competition"
	practice    = "practice"
)

var (
	errNotAuthenticated = errors.New("this service requires authentication in competition mode")
)

type Server struct {
	config     *utilities.Config
	store      db.Store
	router     *fiber.App
	tokenMaker *token.PasetoMaker
	pairs      []string
}

func NewServer(c *utilities.Config, s db.Store) (*Server, error) {
	tm, err := token.NewPasetoTokenMaker(c.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}
	server := &Server{
		config:     c,
		store:      s,
		tokenMaker: tm,
	}

	ps, err := server.store.GetAllParisInDB(context.Background())
	if err != nil {
		return nil, err
	}
	server.pairs = ps

	router := fiber.New(fiber.Config{})

	router.Use(allowOriginMiddleware, limiterMiddleware, server.createLoggerMiddleware())

	router.Get("/practice", server.practice)
	router.Post("/practice", server.practice)
	router.Get("/interval", server.sendInterval)
	router.Get("/rank", server.rank)
	router.Post("/rank", server.rank)
	router.Get("/moreinfo", server.moreinfo)
	router.Post("/user", server.createUser)
	router.Post("/user/login", server.loginUser)
	router.Post("/token/reissue_access", server.reissueAccessToken)

	authGroup := router.Group("/auth", authMiddleware(server.tokenMaker))
	authGroup.Get("/competition", server.competition)
	authGroup.Post("/competition", server.competition)
	authGroup.Get("/myscore/:user", server.myscore)

	server.router = router

	return server, nil
}

func (s *Server) Listen() error {
	return s.router.Listen(s.config.HTTPAddress)
}

func (s *Server) practice(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		r := new(CandlesRequest)
		err := c.QueryParser(r)
		if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
		}

		history := utilities.SplitPairnames(r.Names)
		if len(history) >= finalstage {
			return c.Status(fiber.StatusBadRequest).SendString("invalid current stage")
		}
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err := s.makeChartToRef(db.OneH, nextPair, practice, len(history), c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(oc)
	case "POST":
		var PracticeOrder ScoreRequest
		err := c.BodyParser(&PracticeOrder)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		errs := utilities.ValidateStruct(PracticeOrder)
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
		}

		err = validateOrderRequest(&PracticeOrder)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		r, err := s.createPracResult(&PracticeOrder, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(r)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func (s *Server) competition(c *fiber.Ctx) error {
	switch c.Method() {
	case "GET":
		r := new(CandlesRequest)
		err := c.QueryParser(r)
		if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
		}
		history := utilities.SplitPairnames(r.Names)

		if len(history) >= finalstage {
			return c.Status(fiber.StatusBadRequest).SendString("invalid current stage")
		}
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err := s.makeChartToRef(db.OneH, nextPair, competition, len(history), c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(oc)
	case "POST":
		var CompetitionOrder ScoreRequest
		err := c.BodyParser(&CompetitionOrder)
		if err != nil || CompetitionOrder.Mode != competition {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("%s, mode : %s", err, CompetitionOrder.Mode))
		}

		errs := utilities.ValidateStruct(CompetitionOrder)
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
		}

		err = validateOrderRequest(&CompetitionOrder)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		compResult, err := s.createCompResult(&CompetitionOrder, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		err = s.insertUserScore(&CompetitionOrder, compResult.Score, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(compResult)
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func (s *Server) sendInterval(c *fiber.Ctx) error {
	i := new(AnotherIntervalRequest)
	err := c.BodyParser(i)
	if errs := utilities.ValidateStruct(*i); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}
	if i.Mode == competition {
		if err := authMiddleware(s.tokenMaker)(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("%s %s", errNotAuthenticated, err))
		}
	}
	oc, err := s.sendAnotherInterval(i, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(oc)
}

func (s *Server) myscore(c *fiber.Ctx) error {
	u := c.Params("user")
	p := new(PageRequest)
	err := c.QueryParser(p)
	if errs := utilities.ValidateStruct(*p); err != nil || errs != nil {
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
		}
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	scores, err := s.getMyscores(u, p.Page, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
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
				return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
			}
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		ranks, err := s.getAllRanks(p.Page, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(ranks)
	case "POST":
		if err := authMiddleware(s.tokenMaker)(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("%s %s", errNotAuthenticated, err))
		}
		var r RankInsertRequest
		err := c.BodyParser(&r)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		errs := utilities.ValidateStruct(r)
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
		}
		s.insertScoreToRankBoard(&r, c)
		return nil
	default:
		return errors.New("Not allowed method : " + c.Method())
	}
}

func (s *Server) moreinfo(c *fiber.Ctx) error {
	q := new(MoreInfoRequest)
	if err := c.BodyParser(q); err != nil {
		return err
	}
	errs := utilities.ValidateStruct(q)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	}
	scores, err := s.sendMoreInfo(q.UserId, q.ScoreId, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}
