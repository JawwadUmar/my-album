package utility

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func ValidateEnteredPassword(password, hashedPasswrod string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswrod), []byte(password))

	return err == nil
}
