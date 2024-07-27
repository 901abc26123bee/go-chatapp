package hashext

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// BcryptPassword hash and salt password
func BcryptPassword(s string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash string: %v", err)
	}
	return string(hashedPassword[:]), nil
}

// ComparePasswords compare the password with the hashedPassword in db, return nil if matched
func ComparePasswords(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

// Sha256 generates the hashed string by an input string.
func Sha256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return strings.ToUpper(fmt.Sprintf("%x", sum))
}
