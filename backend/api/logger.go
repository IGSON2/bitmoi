package api

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
)

func (s *Server) createLoggerMiddleware() fiber.Handler {
	folderPath := filepath.Join(s.config.DataDir, "logs", "http")
	if err := os.MkdirAll(folderPath, 0700); err != nil {
		log.Panicf("cannot make http logs datadir: %s", err.Error())
	}

	logfilePath := filepath.Join(folderPath, time.Now().Format("06010215")+".log")

	f, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panicf("cannot create logger middleware: %s", err.Error())
	}

	w := io.MultiWriter(f, os.Stdout, os.Stderr)
	log.Printf("HTTP multi logger has created successfully. PATH=%s", logfilePath)

	lm := logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${time} ${status} - ${method} ${path} - ${latency}\n",
		Output:     w,
		TimeFormat: "2006-01-02T15:04:05",
	})

	logger := zerolog.New(zerolog.ConsoleWriter{Out: w}).With().Timestamp().Logger().Level(s.config.LogLevel)
	s.logger = &logger
	return lm
}
