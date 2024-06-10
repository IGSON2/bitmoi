package gapi

import (
	db "bitmoi/db/sqlc"
	"bitmoi/gapi/pb"
	"bitmoi/token"
	"bitmoi/utilities"
	"bitmoi/utilities/common"
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
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
	if c.Environment == common.EnvProduction {
		setMultiLogger(c)
	}

	server := &Server{
		config:     c,
		store:      s,
		tokenMaker: tm,
	}

	ps, err := server.store.GetAllPairsInDB1H(context.Background())
	if err != nil {
		return nil, err
	}
	server.pairs = ps

	return server, nil
}

func (s *Server) ListenGRPC(errCh chan error) {
	grpcInterceptor := grpc.UnaryInterceptor(GrpcLogger)
	grpcServer := grpc.NewServer(grpcInterceptor)

	pb.RegisterBitmoiServer(grpcServer, s)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", s.config.GRPCAddress)
	if err != nil {
		log.Panic().Err(fmt.Errorf("cannot create gRPC listener: %w", err))
	}
	log.Info().Msgf("Start gRPC server at %s", s.config.GRPCAddress)
	errCh <- grpcServer.Serve(listener)
}

func (s *Server) ListenGRPCGateWay(errCh chan error) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterBitmoiHandlerServer(ctx, grpcMux, s)
	if err != nil {
		log.Panic().Err(fmt.Errorf("cannot register handler server: %w", err))
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", s.config.GRPCHTTPAddress)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	errCh <- http.Serve(listener, GatewayLogger(mux))
}

func (s *Server) RequestCandles(c context.Context, r *pb.CandlesRequest) (*pb.CandlesResponse, error) {
	next, prevStage, err := validateGetCandlesRequest(r, s.pairs)
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
		p, err := s.authorizeUser(c)

		if p == nil || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "cannot generate payload: err = %v, payload nil = %v", err, p == nil)
		}

		if p.UserID != r.UserId {
			log.Error().Msgf("payload: %s, request: %s", p.UserID, r.UserId)
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized user: %v", r.UserId)
		}
		oc, err := s.makeChartToRef(db.OneH, next, competition, prevStage, c)
		if err != nil {
			return nil, err
		}
		return convertGetCandlesRes(oc), nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "mode must be specified")
	}
}

func (s *Server) PostScore(c context.Context, r *pb.ScoreRequest) (*pb.ScoreResponse, error) {
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
		p, err := s.authorizeUser(c)

		if p == nil || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "cannot generate payload: err = %v, payload nil = %v", err, p == nil)
		}

		if p.UserID != r.UserId {
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized user: %s", err)
		}
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

func (s *Server) AnotherInterval(c context.Context, r *pb.AnotherIntervalRequest) (*pb.CandlesResponse, error) {
	if err := validateAnotherIntervalRequest(r); err != nil {
		return nil, err
	}
	if r.Mode == competition {
		p, err := s.authorizeUser(c)
		if p == nil || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "cannot generate payload: err = %v, payload nil = %v", err, p == nil)
		}
		if p.UserID != r.UserId {
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized user: err = %v", err)
		}
	}
	oc, err := s.sendAnotherInterval(r, c)
	if err != nil {
		return nil, err
	}
	return convertGetCandlesRes(oc), nil
}
