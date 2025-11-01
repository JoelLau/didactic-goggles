package db

import (
	"context"
	dbgen "didactic-goggles/internal/db/gen"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbConnPool *pgxpool.Pool
var dbQueries *dbgen.Queries

// NOTE: unresolved when do i close the pool and connections if its a singleton?
func ConnectionPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	var err error

	if dbConnPool != nil {
		dbConnPool, err = pgxpool.New(ctx, dsn)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting connection pool: %+v", err)
	}

	return dbConnPool, err
}

func Queries(ctx context.Context, dsn string) (*dbgen.Queries, error) {
	if dbQueries != nil {
		p, err := ConnectionPool(ctx, dsn)
		if err != nil {
			return nil, fmt.Errorf("error getting connection pool for querier: %+v", err)
		}

		dbQueries = dbgen.New(p)
	}

	return dbQueries, nil
}
