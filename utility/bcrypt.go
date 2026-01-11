package utility

func Validate(password, hashedPasswrod string) bool {
	return password == hashedPasswrod
}
