package gapi

import (
	"context"
	"net/http"
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
	result, err := handler(ctx, req)  // meneruskan ctx dan req pada handler
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
		Str("method", info.FullMethod).          // print request method
		Int("status_code", int(statusCode)).     // print request status code
		Str("status_text", statusCode.String()). // print request status string
		Dur("duration", duration).               // print request duration  to process
		Msg("menerima sebuah grpc request")

	return result, err
}

// kita akan membuat struct response recorder yang cara kerjanya seperti menjadi pengantar menuju default http.ResponseWriter, saat custom response recorder ini bekerja kita akan membuat salinan data dari status code dan status text terlebih dahulu dan selanjutnya baru data tadi diantarkan pada method method asli dari http.ResponseWriter
// simple nya ini akan merekam (record) data yang akan dikirimkan ke response writer
type ResponseRecorder struct {
	http.ResponseWriter // embed http.ResponseWriter disini akan membuat ResponseRecorder secara otomatis memiliki semua field dan method dari http.ResponseWriter
	Body       []byte
	StatusCode int
}

// override write pada responsewriter since this method will called by handler whenever it want to set  response body
func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)  // agar responsewriter tetap menulis body untuk response kita
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode) // agar responsewriter tetap menulis header untuk response kita
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK, // default saja
		}

		handler.ServeHTTP(rec, req) // meneruskan ctx dan req pada handler function to be processed // tidak seperti grpc interceptor function ini tidak mengembalikan result dan error object disini

		duration := time.Since(startTime) // after getting result from the handler we can compute duration of the requerst

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body) // pada server http gateway, error detail akan berada pada response body 
		}

		logger.
			Str("protocol", "http").
			Str("method", req.Method).   // it does'nt contain request path like FullMehtod in grpc
			Str("path", req.RequestURI). // we have to add one more string field to the log to show patch
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("menerima sebuah http request")
	})
}
