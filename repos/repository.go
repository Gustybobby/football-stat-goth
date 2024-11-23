package repos

import (
	"context"
	"football-stat-goth/queries"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	Queries *queries.Queries
	Conn    *pgx.Conn
	Ctx     context.Context
	dsn     string
}

func DbConnect(dsn string) (*Repository, error) {
	ctx := context.Background()
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	if os.Getenv("DB_LOG") == "true" {
		config.Tracer = NewMultiQueryTracer(NewLoggingQueryTracer(slog.Default()))
	}

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	queries := queries.New(conn)
	return &Repository{Queries: queries, Conn: conn, Ctx: ctx, dsn: dsn}, nil
}
