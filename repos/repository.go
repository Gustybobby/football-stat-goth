package repos

import (
	"context"
	"football-stat-goth/queries"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	config *Config
	dsn    string
}

func (repo *Repository) Connect() (*queries.Queries, *pgx.Conn, context.Context, error) {
	ctx := context.Background()
	config, err := pgx.ParseConfig(repo.dsn)
	if err != nil {
		return nil, nil, nil, err
	}
	config.Tracer = NewMultiQueryTracer(NewLoggingQueryTracer(slog.Default()))
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return nil, nil, nil, err
	}
	queries := queries.New(conn)
	return queries, conn, ctx, nil
}
