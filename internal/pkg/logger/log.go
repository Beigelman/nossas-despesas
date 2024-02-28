package logger

import (
	"log/slog"
	"os"
)

func New(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}

func NewProduction(logLevel ...string) *slog.Logger {
	attrMapper := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "level" {
			return slog.Attr{Key: "severity", Value: a.Value}
		}

		if a.Key == "msg" {
			return slog.Attr{Key: "message", Value: a.Value}
		}

		if a.Key == "time" {
			return slog.Attr{Key: "timestamp", Value: a.Value}
		}

		return a
	}

	return New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       getLogLevel(logLevel...),
		ReplaceAttr: attrMapper,
	}))
}

func NewDevelopment(logLevel ...string) *slog.Logger {
	return New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(logLevel...),
	}))
}

func getLogLevel(logLevel ...string) slog.Level {
	if len(logLevel) > 0 {
		return LogLevelMap(logLevel[0])
	}
	return slog.LevelInfo
}

func LogLevelMap(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
