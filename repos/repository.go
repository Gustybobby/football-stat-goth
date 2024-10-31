package repos

import (
	"context"
	"football-stat-goth/queries"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	config *Config
	dsn    string
}

func (repo *Repository) Connect() (*queries.Queries, *pgx.Conn, context.Context, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, repo.dsn)
	if err != nil {
		return nil, nil, nil, err
	}
	queries := queries.New(conn)
	return queries, conn, ctx, nil
}
