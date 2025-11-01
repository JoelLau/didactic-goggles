package config

import (
	"io"
	"log/slog"
)

func NewSlogger(r io.Writer, debug bool) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			r,
			&slog.HandlerOptions{
				AddSource: debug,
				Level:     slog.LevelDebug,
			},
		),
	)
}
