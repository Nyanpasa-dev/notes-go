package utils

import (
	"simple-api/models"

	"github.com/golang-jwt/jwt"
)

type claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func CreateJWT(user *models.User) (string, error) {

	claims := &claims{
		ID:       user.ID,
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("secret"))

	return tokenString, err
}
