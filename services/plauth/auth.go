package plauth

import (
	"context"
	"football-stat-goth/queries"
	"net/http"
)

func Auth(r *http.Request, db *queries.Queries, ctx context.Context) (*queries.User, error) {
	token, err := GetSessionTokenFromCookie(r)
	if err != nil {
		return nil, err
	}

	session, err := ValidateSessionToken(token, db, ctx)
	if err != nil {
		return nil, err
	}

	return &session.User, nil
}
