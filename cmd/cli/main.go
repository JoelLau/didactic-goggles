package main

import (
	"context"
	dbgen "didactic-goggles/internal/db/gen"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/jackc/pgx/v5/pgtype"
)

const (
	ErrorCodeSuccess           = 0
	ErrorCodeFailedToConnectDB = 1
)

// CLA takes in exactly 1 argument - the path to file for ingestion

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f ../../internal/db/sqlc.yml
func main() {
	ctx := context.Background()
	logger := NewLogger()

	// set up sqlc generated libraries
	logger.InfoContext(ctx, "starting...",
		slog.Any("args", os.Args),
		slog.Any("envs", os.Environ()),
	)

	// TODO: refactor: singleton
	dsn := Dsn()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.ErrorContext(ctx, "failed to connect to database", slog.Any("error", err))

		os.Exit(ErrorCodeFailedToConnectDB)
		return
	}
	defer pool.Close()

	// TODO: refactor: singleton
	conn, err := pool.Acquire(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to connect to database", slog.Any("error", err))

		os.Exit(ErrorCodeFailedToConnectDB)
		return
	}

	// TODO: move this to app layer
	queries := dbgen.New(conn)
	categories, err := queries.ListCategories(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed to list categories", slog.Any("error", err))

		os.Exit(ErrorCodeFailedToConnectDB)
		return
	}

	logger.InfoContext(ctx, "fetched categories", slog.Any("categories", categories))
	os.Exit(ErrorCodeSuccess)

}

func Dsn() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	databaseName := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, password, host, port, databaseName,
	)
}

func NewLogger() *slog.Logger {
	return slog.Default()
}
