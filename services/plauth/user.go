package plauth

import (
	"context"
	"football-stat-goth/queries"
)

func CreateUser(username string, password string, firstName string, lastName string, db *queries.Queries, ctx context.Context) (*queries.User, error) {
	passwordHash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := db.CreateUser(ctx, queries.CreateUserParams{Username: username, PasswordHash: passwordHash, Firstname: firstName, Lastname: lastName})
	if err != nil {
		return nil, err
	}

	return &user, nil
}
