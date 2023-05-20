package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(lvl zerolog.Level) zerolog.Logger {
	logger := zerolog.New(os.Stdout)
	return logger.Level(lvl)
}
