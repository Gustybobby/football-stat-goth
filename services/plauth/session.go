package plauth

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"football-stat-goth/queries"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ValidationResult struct {
	Session queries.Session
	User    queries.User
}

func ValidateSessionToken(token string, db *queries.Queries, ctx context.Context) (*ValidationResult, error) {
	encodedToken := encodeSessionToken(token)

	session, err := db.FindSessionByToken(ctx, encodedToken)
	if err != nil {
		return nil, err
	}

	user, err := db.FindUserByUsername(ctx, session.Username)
	if err != nil {
		return nil, err
	}

	if time.Now().UTC().Compare(session.ExpiresAt.Time) >= 0 {
		return nil, errors.New("session expired")
	}

	if time.Now().UTC().Compare(session.ExpiresAt.Time.AddDate(0, 0, -3)) >= 0 {
		var expiresAt = pgtype.Timestamp{Time: session.ExpiresAt.Time.AddDate(0, 0, 7), Valid: true}
		session, err := db.UpdateSessionExpiresAt(ctx, queries.UpdateSessionExpiresAtParams{
			Token:     encodedToken,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			return nil, err
		}
		return &ValidationResult{Session: session, User: user}, nil
	}

	return &ValidationResult{Session: session, User: user}, nil
}

func GenerateSessionToken() (string, error) {
	tokenBytes, err := generateRandomBytes(20)
	if err != nil {
		return "", err
	}

	token := strings.ToLower(base32.StdEncoding.EncodeToString(tokenBytes))

	return token, nil
}

func encodeSessionToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	encodedToken := strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))

	return encodedToken
}

func CreateSession(token string, username string, db *queries.Queries, ctx context.Context) (*queries.Session, error) {
	createdSession, err := db.CreateSession(ctx, queries.CreateSessionParams{
		Token:     encodeSessionToken(token),
		Username:  username,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &createdSession, nil
}
