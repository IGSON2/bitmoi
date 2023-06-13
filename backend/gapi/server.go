package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"fmt"
)

type Server struct {
	pb.UnimplementedBitmoiServer
	config     *utilities.Config
	store      db.Store
	tokenMaker *token.PasetoMaker
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

	return server, nil
}
