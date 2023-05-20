package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLogger(t *testing.T) {
	// zerolog.TimeFieldFormat = zerolog.TimestampFunc().Format("2006-01-02 15:04:05")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Print("JustMessage")
}

func TestMessage(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	logger.Info().Msg("Hello")
}
