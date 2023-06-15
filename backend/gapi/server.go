package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	finalstage  = 10
	competition = "competition"
	practice    = "practice"
)

type Server struct {
	pb.UnimplementedBitmoiServer
	config     *utilities.Config
	store      db.Store
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

	return server, nil
}

func (s *Server) GetCandles(c context.Context, r *pb.GetCandlesRequest) (*pb.GetCandlesResponse, error) {
	err, next, prevStage := validateAndGetNextPair(r, s.pairs)
	if err != nil {
		return nil, err
	}
	switch r.Mode {
	case practice:
		oc, err := s.makeChartToRef(db.OneH, next, practice, prevStage, c)
		if err != nil {
			return nil, err
		}
		return convertGetCandlesRes(oc), nil
	case competition:
		oc, err := s.makeChartToRef(db.OneH, next, competition, prevStage, c)
		if err != nil {
			return nil, err
		}
		return convertGetCandlesRes(oc), nil
	default:
		return nil, fmt.Errorf("error: mode must be specified")
	}
}

func (s *Server) PostScore(c context.Context, r *pb.OrderRequest) (*pb.OrderResponse, error) {
	if err := validateOrderRequest(r); err != nil {
		return nil, err
	}
	switch r.Mode {
	case practice:
		pracResult, err := s.createPracResult(r, c)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err)
		}
		return pracResult, nil
	case competition:
		compResult, err := s.createCompResult(r, c)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err)
		}
		err = s.insertUserScore(r, compResult.Score, c)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err)
		}
		return compResult, nil
	default:
		return nil, status.Errorf(codes.Internal, "error: mode must be specified")
	}
}
