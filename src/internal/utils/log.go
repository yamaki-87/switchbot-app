package utils

import (
	"log/slog"
	"os"
)

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func FatalIfErr(err error, msg string, args ...any) {
	if err != nil {
		Fatal(msg+" Error:"+err.Error(), args...)
	}
}

func LogInit(level slog.Level) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}
