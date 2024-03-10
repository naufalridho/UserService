package hashutil

import (
	"golang.org/x/crypto/bcrypt"
)

func HashSaltPassword(input string) (string, error) {
	hashedSalted, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hashedSalted), nil
}

func IsPasswordMatched(input, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input)) == nil
}
