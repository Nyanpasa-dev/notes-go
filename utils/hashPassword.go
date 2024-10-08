package utils

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptPassword interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}

type BcryptHasher struct{}

func (BcryptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (BcryptHasher) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
