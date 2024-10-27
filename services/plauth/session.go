package plauth

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"football-stat-goth/models"
	"football-stat-goth/repos"
	"strings"
	"time"
)

type ValidationResult struct {
	Session models.Session
	User    models.User
}

func ValidateSessionToken(token string, repo *repos.Repository) (*ValidationResult, error) {
	encodedToken := encodeSessionToken(token)

	var session = models.Session{Token: encodedToken}
	if results := repo.DB.First(&session); results.Error != nil {
		return nil, results.Error
	}

	var user = models.User{Username: session.Username}
	if results := repo.DB.First(&user); results.Error != nil {
		return nil, results.Error
	}
	fmt.Printf("%#v\n%#v", session, user)

	if time.Now().Compare(session.ExpiresAt) >= 0 {
		return nil, errors.New("session expired")
	}

	if time.Now().Compare(session.ExpiresAt.AddDate(0, 0, -3)) >= 0 {
		session.ExpiresAt = session.ExpiresAt.AddDate(0, 0, 7)
		repo.DB.Save(&session)
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

func CreateSession(token string, username string, repo *repos.Repository) (*models.Session, error) {
	session := &models.Session{
		Token:     encodeSessionToken(token),
		Username:  username,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	}

	results := repo.DB.Create(session)
	if results.Error != nil {
		return nil, results.Error
	}

	return session, nil
}
