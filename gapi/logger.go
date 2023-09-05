package gapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	resp interface{},
	err error,
) {
	startTime := time.Now()
	result, err := handler(ctx, req) // meneruskan ctx dan req pada handler
	duration := time.Since(startTime) // after getting result from the handler we can compute duration of the requerst

	statusCode := codes.Unknown
	// extract status dari error
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code() // get status code
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err) // jika ada error maka Log.Err() akan menimpa Log.Infor()
		// .Err(err) meambahkan error field jika ada error
	}

	// log.Info().
	logger.
		Str("protocol", "grpc").
		Str("method", info.FullMethod). // print request method
		Int("status_code", int(statusCode)). // print request status code
		Str("status_text", statusCode.String()). // print request status string
		Dur("duration", duration). // print request duration  to process
		Msg("menerima sebuah grpc request")

	return result, err
}
