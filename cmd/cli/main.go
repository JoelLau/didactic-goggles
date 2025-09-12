package main

import (
	"context"
	"log/slog"
	"os"
)

// CLA takes in exactly 1 argument - the path to file for ingestion
func main() {
	ctx := context.Background()
	logger := NewLogger()

	logger.InfoContext(ctx, "starting...",
		slog.Any("args", os.Args),
		slog.Any("envs", os.Environ()),
	)
}

func NewLogger() *slog.Logger {
	return slog.Default()
}
