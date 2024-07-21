package utils

import (
	"simple-api/models"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

func CreateJWT(user *models.User) (string, error) {

	claims := &Claims{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("secret"))

	return tokenString, err
}

func ExtractUserFromToken(user *models.User, tokenString string) (Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return *claims, err
	}

	user.ID = claims.ID
	user.Username = claims.Username
	user.IsAdmin = claims.IsAdmin

	return *claims, nil

}