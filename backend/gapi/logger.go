package gapi

import (
	"bitmoi/backend/utilities"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func setMultiLogger(config *utilities.Config) error {
	folderPath := filepath.Join(config.DataDir, "logs", "grpc")
	if err := os.MkdirAll(folderPath, 0700); err != nil {
		return fmt.Errorf("cannot make grpc logs datadir: %w", err)
	}

	logfilePath := filepath.Join(folderPath, time.Now().Format("06010215")+".log")

	f, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Multilogger는 log level이 nil이다.
	fileLogger := log.Output(zerolog.ConsoleWriter{Out: f,
		FormatLevel: func(i interface{}) string {
			if i == nil {
				return strings.ToUpper(fmt.Sprintf("| %-6s|", "INFO"))
			}
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		},
		TimeFormat: "2006-01-02T15:04:05"})
	if err != nil {
		return err
	}

	log.Logger = log.Output(zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05"}, fileLogger))
	log.Info().Msgf("gRPC multi logger has set successfully. PATH=%s", logfilePath)
	return nil
}

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")

	return result, err
}

type ResponseRecoder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseRecoder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecoder) Write(b []byte) (int, error) {
	r.Body = b
	return r.ResponseWriter.Write(b)
}

func GatewayLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recoder := ResponseRecoder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(&recoder, r)
		duration := time.Since(startTime)

		logger := log.Info()
		if recoder.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("Body", recoder.Body)
		}

		logger.Str("protocol", "grpc HTTP").
			Str("method", r.Method).
			Str("path", r.RequestURI).
			Int("status_code", recoder.StatusCode).
			Str("status_text", http.StatusText(recoder.StatusCode)).
			Dur("duration", duration).
			Msg("received a grpc HTTP request")

	})
}
