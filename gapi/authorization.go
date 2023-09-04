package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/tech_school/simple_bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	// biasanya accessToken akan dikirimkan oleh client didalam metadata
	md, ok := metadata.FromIncomingContext(ctx) // get metadata stored inside the context
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0] // Bearer abcsdf : <authorization-type> <authorization-data> : type accessToken

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("format Authorization header tidak valid")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return nil, fmt.Errorf("bentuk authorization ini tidak disupport oleh server %s", authorizationType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token : %s", err)
	}

	return payload, nil

}
