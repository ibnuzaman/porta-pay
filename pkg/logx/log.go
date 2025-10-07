package logx

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(service string, env string) zerolog.Logger {
	l := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", service).
		Str("env", env).
		Logger()

	zerolog.TimeFieldFormat = time.RFC3339Nano

	return l
}
