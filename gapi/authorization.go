package gapi

import (
	"bitmoi/token"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func (s *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	v := md.Get(authorizationHeaderKey)
	if len(v) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := v[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationTypeBearer {
		return nil, fmt.Errorf("unsupported authrization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := s.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}
