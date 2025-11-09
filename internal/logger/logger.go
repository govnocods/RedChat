package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}

	env := os.Getenv("APP_ENV")
	if env == "production" {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		Logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
}

func ParseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func WithError(err error) *slog.Logger {
	return Logger.With("error", err)
}

func WithFields(fields ...any) *slog.Logger {
	return Logger.With(fields...)
}
