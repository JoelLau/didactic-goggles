package commands

import (
	"context"
	"didactic-goggles/internal/config"
	usecase "didactic-goggles/internal/usecases"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// stores command state
type IngestDBSStatementCommand struct {
	connPool *pgxpool.Pool
	slogger  *slog.Logger
}

func NewIngestDBSStatementCommand(connPool *pgxpool.Pool, logger *slog.Logger) *IngestDBSStatementCommand {
	return &IngestDBSStatementCommand{
		connPool: connPool,
		slogger:  logger,
	}
}

func (cmd *IngestDBSStatementCommand) logger() *slog.Logger {
	if cmd.slogger == nil {
		return config.NewSlogger(io.Discard, false)
	}

	return cmd.slogger
}

func (cmd *IngestDBSStatementCommand) Run(ctx context.Context, args []string) error {
	logger := cmd.logger()
	logger.InfoContext(ctx, "starting command run")

	filepath := args[0] // is this the best way to receive a filepath?
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file at at '%s': %+v", filepath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.ErrorContext(ctx, "error closing file", slog.Any("error", fmt.Errorf("error closing file: %+v", err)))
			return
		}
	}()

	uc := usecase.NewIngestDbsStatementUseCase(cmd.connPool, logger)
	if err := uc.Execute(ctx, file); err != nil {
		return fmt.Errorf("error executing ingest dbs statement use case: %+v", err)
	}

	return nil
}
