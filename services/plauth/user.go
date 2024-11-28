package plauth

import (
	"context"
	"errors"
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

func FindPasswordHash(username string, db *queries.Queries, ctx context.Context) (string, error) {
	passwordHash, err := db.FindPasswordHashByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func UpdatePassword(username string, currentPassword string, newPassword string, db *queries.Queries, ctx context.Context) error {
	userPasswordHash, err := FindPasswordHash(username, db, ctx)
	if err != nil {
		return err
	}

	valid, err := VerifyPassword(currentPassword, userPasswordHash)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid password")
	}

	newPasswordHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	db.UpdatePasswordByUsername(ctx, queries.UpdatePasswordByUsernameParams{
		Username:     username,
		PasswordHash: newPasswordHash,
	})

	return nil
}
