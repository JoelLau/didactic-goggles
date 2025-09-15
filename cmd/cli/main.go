package main

import (
	"context"
	"log/slog"
	"os"
)

// CLA takes in exactly 1 argument - the path to file for ingestion

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f ../../internal/db/sqlc.yml
func main() {
	ctx := context.Background()
	logger := NewLogger()

	// set up pgx connection pool
	// set up sqlc generated libraries
	logger.InfoContext(ctx, "starting...",
		slog.Any("args", os.Args),
		slog.Any("envs", os.Environ()),
	)
}

func NewLogger() *slog.Logger {
	return slog.Default()
}
