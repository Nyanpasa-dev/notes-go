package utils

import (
	"simple-api/models"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"isAdmin"`
	IpAdress  string `json:"ipAdress"`
	UserAgent string `json:"userAgent"`
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

func ExtractUserFromToken(tokenString string) (Claims, error) {
	claims := Claims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return *claims, err
	}

	claims.ID = claims.ID
	claims.Username = claims.Username
	claims.IsAdmin = claims.IsAdmin
	claims.IpAdress = claims.IpAdress
	claims.UserAgent = claims.UserAgent

	return *claims, nil

}
