package main

import (
	"bitmoi/backend/api"
	"bitmoi/backend/app"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
