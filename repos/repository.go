package repos

import (
	"context"
	"football-stat-goth/queries"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Queries *queries.Queries
	Pool    *pgxpool.Pool
	Ctx     context.Context
	dsn     string
}

func DbConnect(dsn string) (*Repository, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	if os.Getenv("DB_LOG") == "true" {
		config.ConnConfig.Tracer = NewMultiQueryTracer(NewLoggingQueryTracer(slog.Default()))
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	queries := queries.New(pool)
	return &Repository{Queries: queries, Pool: pool, Ctx: ctx, dsn: dsn}, nil
}
