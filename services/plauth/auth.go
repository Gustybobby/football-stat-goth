package plauth

import (
	"encoding/json"
	"football-stat-goth/queries"
	"football-stat-goth/repos"
	"log/slog"
	"net/http"
	"os"
)

type AuthContextKey string

const (
	ContextUserKey AuthContextKey = "pl.user"
)

func Auth(w http.ResponseWriter, r *http.Request, repo *repos.Repository) (*queries.FindUserByUsernameRow, error) {
	token, err := GetSessionTokenFromCookie(r)
	if err != nil {
		slog.Warn("cookie session token not found: " + err.Error())
		return nil, nil
	}

	session, err := ValidateSessionToken(token, repo.Queries, repo.Ctx)
	if err != nil {
		slog.Error("validate session token error: " + err.Error())
		DeleteSessionTokenCookie(w, os.Getenv("ENV") == "production")
		return nil, nil
	}

	stringifiedSession, err := json.MarshalIndent(session, "-", "  ")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	slog.Info("plauth session: " + string(stringifiedSession))

	return &session.User, nil
}

func GetContextUser(r *http.Request) *queries.FindUserByUsernameRow {
	return r.Context().Value(ContextUserKey).(*queries.FindUserByUsernameRow)
}
