package models

import (
	"errors"
	"time"

	"example.com/my-ablum/database"
	"example.com/my-ablum/utility"
)

// If a DB column allows NULL → use pointer in Go.
type User struct {
	UserId         int64   `json:"user_id"`
	Email          string  `json:"email" binding:"required"`
	Password       *string `json:"password" binding:"required"`
	GoogleId       *string
	FirstName      string    `json:"first_name" binding:"required"`
	LastName       string    `json:"last_name" binding:"required"`
	ProfilePic     *string   `json:"profile_pic"`
	AllowedStorage int64     `json:"allowed_storage"`
	CreatedAt      time.Time //YYYY-MM-DD HH:MM:SS.microseconds stored in DB
	UpdatedAt      time.Time
}

func (u *User) Save() error {
	query := `
				INSERT INTO users (email, password_hash, first_name, last_name, profile_pic, google_id)
				VALUES($1, $2, $3, $4, $5, $6)
				RETURNING id;
	`

	var hashedPassword *string
	var err error

	if u.Password != nil {

		hashedPassword, err = utility.HashPassword(*u.Password)

		if err != nil {
			return err
		}
	}

	//u.Profilepic is *string but thanks to go sqldatabase driver :)
	row := database.DB.QueryRow(query, u.Email, hashedPassword, u.FirstName, u.LastName, u.ProfilePic, u.GoogleId)

	err = row.Scan(&u.UserId)
	return err
}

func (u *User) ValidateCredential() error {
	// query := `SELECT id, password_hash FROM users where email = $1`

	query := `SELECT id, first_name, last_name, password_hash, profile_pic, created_at, updated_at, allowed_storage 
				FROM users 
				WHERE email = $1
			`

	row := database.DB.QueryRow(query, u.Email)

	var passwordHash string

	err := row.Scan(&u.UserId, &u.FirstName, &u.LastName, &passwordHash, &u.ProfilePic, &u.CreatedAt, &u.UpdatedAt, &u.AllowedStorage) //userId is updated here

	if err != nil {
		return err
	}

	isValid := utility.ValidateEnteredPassword(*u.Password, passwordHash)

	if !isValid {
		return errors.New("Invalid Credential")
	}

	u.Password = &passwordHash

	return nil
}

func (u *User) UpdateGoogleId() error {
	query := `
		UPDATE users
		SET google_id = $1
		WHERE email = $2
		RETURNING id;
	`
	row := database.DB.QueryRow(query, u.GoogleId, u.Email)
	err := row.Scan(&u.UserId)
	return err
}

func GetUserModelByEmail(email string) (User, error) {
	query := `SELECT id, email, first_name, last_name, password_hash, profile_pic, created_at, updated_at, google_id, allowed_storage
			  FROM users
			  WHERE email = $1
	`

	row := database.DB.QueryRow(query, email)

	var user User

	err := row.Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.ProfilePic,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.GoogleId,
		&user.AllowedStorage,
	)

	return user, err
}

func GetUserModelById(id int64) (User, error) {
	query := `SELECT id, email, first_name, last_name, password_hash, profile_pic, created_at, updated_at, google_id, allowed_storage
			  FROM users
			  WHERE id = $1
	`

	row := database.DB.QueryRow(query, id)

	var user User

	err := row.Scan(
		&user.UserId,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.ProfilePic,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.GoogleId,
		&user.AllowedStorage,
	)

	return user, err
}

func (u *User) Update() error {
	query := `UPDATE users
			SET first_name  = COALESCE($1, first_name),
				last_name   = COALESCE($2, last_name),
				profile_pic = COALESCE($3, profile_pic)
				WHERE id = $4`

	_, err := database.DB.Exec(query, u.FirstName, u.LastName, u.ProfilePic, u.UserId)

	return err
}

func GetUserStorage(userId int64) (int64, error) {
	var total int64
	var err error

	query := `
        SELECT COALESCE(SUM(file_size), 0)
        FROM files
        WHERE created_by = $1
    `
	row := database.DB.QueryRow(query, userId)
	err = row.Scan(&total)

	return total, err
}

func GetAllowedUserStorage(userId int64) (int64, error) {
	var allowedStorage int64
	var err error

	query := `
        SELECT COALESCE(allowed_storage, 0)
        FROM users
        WHERE id = $1
    `
	row := database.DB.QueryRow(query, userId)
	err = row.Scan(&allowedStorage)

	return allowedStorage, err
}
