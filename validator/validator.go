package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString // matchstring tanpa () berarti sekarang variable isvalidusername menjadi function
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

// general function to check if a string has an appropriate length or not
func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(name string) error {
	if err := ValidateString(name, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(name) {
		return fmt.Errorf("must contain only lowercase letters, digits or underscore")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidFullname(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 4, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}

	// function ini mengembalikan parse email address dan error, kita akan memanfaatkan return errornya untuk check validitas email
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not valid email address")
	}

	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("it must be a positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}
