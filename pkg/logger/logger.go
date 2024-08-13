package logger

import (
	"context"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger
var once sync.Once

func Get() zerolog.Logger {

	once.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		logger = zerolog.New(os.Stdout).
			With().Timestamp().
			Logger()

	})

	return logger
}

func AttachLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}
