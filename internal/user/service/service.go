package userservice

import (
	"connect/internal/db"
	errormessage "connect/internal/pkg/error_message"
	usermodel "connect/internal/user/model"
	"context"
	"database/sql"
	"errors"
	"log"
)

func CreateTableIfNotExists() {
	_, err := db.Pool().Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(100) NOT NULL UNIQUE,
			username VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(128) NOT NULL
		);`)

	if err != nil {
		log.Fatalf("Failed to create table users %v\n", err)
	}
}

func InsertUser(user usermodel.User) error {
	_, err := db.Pool().Exec(context.Background(), `INSERT INTO users ( email, username, password ) VALUES ( $1, $2, $3 );`, user.Email, user.Username, user.Password)

	if err != nil {
		return errors.New(errormessage.InternalServerError)
	}

	return nil
}

func UserExists(email string) (bool, error) {
	var exists bool

	userExistsErr := db.Pool().QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);", email).Scan(&exists)

	if userExistsErr != nil {
		return exists, errors.New(errormessage.InternalServerError)
	}

	return exists, nil
}

func GetMinimalUserById(id int) (*usermodel.MinimalUser, error) {
	user := new(usermodel.MinimalUser)

	userErr := db.Pool().QueryRow(context.Background(), "SELECT id, username FROM users WHERE id = $1;", id).Scan(&user.Id, &user.Username)

	if userErr != nil {
		if errors.Is(userErr, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.New(errormessage.InternalServerError)
	}

	return user, nil
}

func GetMinimalUserByUsername(username string) (*usermodel.MinimalUser, error) {
	user := new(usermodel.MinimalUser)

	userErr := db.Pool().QueryRow(context.Background(), "SELECT id, username FROM users WHERE username = $1;", username).Scan(&user.Id, &user.Username)

	if userErr != nil {
		if errors.Is(userErr, sql.ErrNoRows) {
			return nil, nil
		}

		return user, errors.New(errormessage.InternalServerError)
	}

	return user, nil
}
