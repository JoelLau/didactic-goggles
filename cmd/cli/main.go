package main

import (
	"context"
	"didactic-goggles/internal/commands"
	"didactic-goggles/internal/config"
	"didactic-goggles/internal/db"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// wrap main command to facilitate end-to-end testing later

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f ../../internal/db/sqlc.yml
func main() {
	app := App{
		Environment:  config.Env(),
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "panic: %+v\n", r)
		}
	}()

	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		slog.Error("error while running main", slog.Any("error", err))
		return
	}
}

type App struct {
	Environment map[string]string

	InputStream  io.Reader
	OutputStream io.Writer
	ErrorStream  io.Writer
}

func (app *App) Run(ctx context.Context, args []string) error {
	logger := config.NewSlogger(app.ErrorStream, true)
	logger.InfoContext(ctx, "running main command", slog.Any("args", args))

	file, err := os.Open(args[0]) // TODO: find a better way to resolve file name
	if err != nil {
		return fmt.Errorf("error opening file at path '%s': %+v", args[0], err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.ErrorContext(ctx, "error closing file", slog.String("file path", args[0]), slog.Any("error", err))
			return
		}
	}()

	connectionPool, err := db.ConnectionPool(ctx, app.Dsn())
	if err != nil {
		return fmt.Errorf("error getting connection pool: %v", err)
	}

	cmd := commands.NewIngestDBSStatementCommand(connectionPool, logger)
	if err := cmd.Run(ctx, args); err != nil {
		return fmt.Errorf("error running ingest dbs command: %+v", err)
	}

	return nil
}

func (app *App) Dsn() string {
	host := app.Environment["DB_HOST"]
	port := app.Environment["DB_PORT"]
	user := app.Environment["DB_USER"]
	password := app.Environment["DB_PASS"]
	databaseName := app.Environment["DB_NAME"]

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, password, host, port, databaseName,
	)
}
