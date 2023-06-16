package main

import (
	"bitmoi/backend/api"
	"bitmoi/backend/app"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var bApp = app.NewApp()

func init() {
	bApp.Commands = []*cli.Command{
		app.StoreCommand,
	}
	bApp.Action = bitmoi
}

func main() {
	if err := bApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bitmoi(ctx *cli.Context) error {
	config := utilities.GetConfig("./")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		return fmt.Errorf("cannot connect db %w", err)
	}
	dbStore := db.NewStore(conn)
	go runHttpServer(config, dbStore)
	go runGateWayServer(config, dbStore)
	runGrpcServer(config, dbStore)
	return nil
}

func runGrpcServer(config *utilities.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create gRPC server")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBitmoiServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create gRPC listener")
	}

	log.Info().Msgf("Start gRPC server at %s", listener.Addr().String())
	log.Panic().Err(grpcServer.Serve(listener)).Msg("cannot start gRPC server")
}

func runHttpServer(config *utilities.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create server")
	}

	log.Panic().Err(server.Listen()).Msg("cannot start http server")
}

func runGateWayServer(config *utilities.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Panic().Err(err).Msg("cannot create gateway server")
	}

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

	err = pb.RegisterBitmoiHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Panic().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", "0.0.0.0:7001")
	if err != nil {
		log.Panic().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	log.Panic().Err(http.Serve(listener, mux)).Msg("cannot start HTTP gate werver")
}
