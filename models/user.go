package models

import (
	"errors"
	"time"

	"example.com/my-ablum/database"
	"example.com/my-ablum/utility"
)

// If a DB column allows NULL → use pointer in Go.
type User struct {
	UserId     int64
	Email      string    `json:"email" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	FirstName  string    `json:"first_name" binding:"required"`
	LastName   string    `json:"last_name" binding:"required"`
	ProfilePic *string   `json:"profile_pic"`
	CreatedAt  time.Time //YYYY-MM-DD HH:MM:SS.microseconds stored in DB
	UpdatedAt  time.Time
}

func (u *User) Save() error {
	query := `
				INSERT INTO users (email, password_hash, first_name, last_name, profile_pic)
				VALUES($1, $2, $3, $4, $5)
				RETURNING id;
	`

	hashedPassword, err := utility.HashPassword(u.Password)

	if err != nil {
		return err
	}

	row := database.DB.QueryRow(query, u.Email, hashedPassword, u.FirstName, u.LastName, u.ProfilePic)

	err = row.Scan(&u.UserId)
	return err
}

func (u *User) ValidateCredential() error {
	query := `SELECT id, password_hash FROM users where email = $1`

	row := database.DB.QueryRow(query, u.Email)

	var passwordHash string

	err := row.Scan(&u.UserId, &passwordHash)

	if err != nil {
		return err
	}

	isValid := utility.ValidateEnteredPassword(u.Password, passwordHash)

	if !isValid {
		return errors.New("Invalid Credential")
	}

	return nil
}
