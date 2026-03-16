package helpers

import (
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Zerolog default is json, which is good for production
	if os.Getenv("APP_ENV") == "development" {
		output := zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339}
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		// Production use JSON
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
