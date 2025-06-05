package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	LogLevel string
	Pretty bool
}

func InitLogger(config Config) {
	zerolog.TimeFieldFormat = time.RFC3339
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	var output io.Writer = os.Stdout
	if config.Pretty {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	}
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}

func Debug(message string, fields ...map[string]interface{}) {
	event := log.Debug().Str("message", message)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Send()
}

func Info(message string, fields ...map[string]interface{}) {
	event := log.Info().Str("message", message)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Send()
}

func Warn(message string, fields ...map[string]interface{}) {
	event := log.Warn().Str("message", message)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Send()
}

func Error(message string, err error, fields ...map[string]interface{}) {
	event := log.Error().Str("message", message)
	if err != nil {
		event = event.Err(err)
	}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Send()
}

func Fatal(message string, err error, fields ...map[string]interface{}) {
	event := log.Fatal().Str("message", message)
	if err != nil {
		event = event.Err(err)
	}
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Send()
}
