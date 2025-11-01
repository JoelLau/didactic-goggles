package usecase

import (
	"context"
	"didactic-goggles/internal/config"
	dbgen "didactic-goggles/internal/db/gen"
	"didactic-goggles/internal/parsers"
	"fmt"
	"io"
	"log/slog"

	gocsv "github.com/JoelLau/go-csv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NOTE: debit  (usually left  hand side): money taken from account
// NOTE: credit (usually right hand side): money  put  into account

// takes a csv file, persist each row in debit table
type IngestDBSStatementUseCase struct {
	connPool *pgxpool.Pool
	slogger  *slog.Logger
}

func NewIngestDbsStatementUseCase(connPool *pgxpool.Pool, logger *slog.Logger) IngestDBSStatementUseCase {
	return IngestDBSStatementUseCase{
		connPool: connPool,
		slogger:  logger,
	}
}

func (uc *IngestDBSStatementUseCase) logger() *slog.Logger {
	if uc.slogger == nil {
		return config.NewSlogger(io.Discard, false)
	}

	return uc.slogger
}

func (uc *IngestDBSStatementUseCase) Execute(ctx context.Context, csv io.Reader) error {
	logger := uc.logger()
	logger.InfoContext(ctx, "Executing Ingest DBS Statement Use Case")

	b, err := io.ReadAll(csv)
	if err != nil {
		return fmt.Errorf("error reading bytes from reader: %+v", err)
	}

	strContents, err := gocsv.ReadAll([]byte(b))
	if err != nil {
		return fmt.Errorf("error reading bytes from reader: %+v", err)
	}
	logger.InfoContext(ctx, "parser", slog.Any("str", strContents))

	parser := parsers.DbsCreditCardParser{}
	statement, err := parser.Parse(strContents)
	if err != nil {
		return fmt.Errorf("error parsing dbs credit card statement: %+v", err)
	}

	tx, err := uc.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting db transaction: %+v", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			logger.ErrorContext(ctx, "error rolling back transaction", slog.Any("ctx", ctx), slog.Any("error", err))
		}
	}()

	queries := dbgen.New(uc.connPool).WithTx(tx)
	for _, lineItem := range statement.LineItems {
		_, err = queries.CreateCredit(ctx, dbgen.CreateCreditParams{
			Name:             lineItem.TransactionDescription,
			TransactedAt:     lineItem.TransactionDate.Time,
			AmountInMicrosgd: lineItem.CreditAmount.IntPart(),
		})

		if err != nil {
			return fmt.Errorf("error persisting credit: %+v", err)
		}
	}

	return nil
}
