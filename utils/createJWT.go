package utils

import (
	"simple-api/models"
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	IpAddress string `json:"ipAdress"`
	UserAgent string `json:"userAgent"`
	ID        uint   `json:"id"`
	jwt.StandardClaims
}

type RefreshParams struct {
	IpAddress string
	UserAgent string
	User      *models.User
}

type AccessParams struct {
	User *models.User
}

type JWTUtils interface {
	CreateJWT() (string, error)
	ExtractUserFromToken(tokenString string) (interface{}, error)
	VerifyToken(tokenString string) (interface{}, error)
}

func (u AccessParams) CreateJWT() (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute).Unix()

	claims := &AccessClaims{
		ID:       u.User.ID,
		Username: u.User.Username,
		IsAdmin:  u.User.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))

	return tokenString, err
}

func (u AccessParams) ExtractUserFromToken(tokenString string) (interface{}, error) {
	claims := &AccessClaims{}

	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if accessClaims, ok := parsed.Claims.(*AccessClaims); ok {
		return accessClaims, nil
	}

	return nil, err
}

func (u AccessParams) VerifyToken(tokenString string) (*AccessClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*AccessClaims)

	if !ok || !parsedToken.Valid {
		return nil, err
	}

	if claims.ID != u.User.ID {
		return nil, err
	}

	return claims, nil
}

func (r RefreshParams) CreateJWT() (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute).Unix()

	claims := &RefreshClaims{
		IpAddress: r.IpAddress,
		UserAgent: r.UserAgent,
		ID:        r.User.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))

	return tokenString, err
}

func (r RefreshParams) ExtractUserFromToken(tokenString string) (interface{}, error) {
	claims := &RefreshClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if refreshClaims, ok := parsed.Claims.(*RefreshClaims); ok {
		return refreshClaims, nil
	}

	return nil, err
}

func (r RefreshParams) VerifyToken(tokenString string) (*RefreshClaims, error) {
	refreshToken, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if !refreshToken.Valid {
		return nil, err
	}

	claims, ok := refreshToken.Claims.(*RefreshClaims)

	if !ok || !refreshToken.Valid {
		return nil, err
	}

	return claims, nil
}
