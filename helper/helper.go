package helper

import (
	"regexp"
)

func IsValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	return regexp.MustCompile(regex).MatchString(email)
}

func IsStrongPassword(password string) bool {
	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$`
	return regexp.MustCompile(regex).MatchString(password)
}
