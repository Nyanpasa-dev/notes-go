package services

import (
	"errors"
	"fmt"
	"net/http"
	"simple-api/internal/repositories"
	"simple-api/models"
	"simple-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authService struct {
	db *gorm.DB
}

type AuthService interface {
	Authenticate(ctx *gin.Context) (*models.User, error)
	RefreshToken(ctx *gin.Context) (string, error)
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db}
}

func (s *authService) Authenticate(ctx *gin.Context) (*models.User, error) {
	hasher := utils.BcryptHasher{}
	clientIp := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	if userAgent == "" {
		return nil, errors.New("User-Agent header is required")
	}

	var json struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	newUser := &models.User{Username: json.Username}

	userRepo := repositories.NewUserRepository(s.db)

	var user *models.User
	var err error

	user, err = userRepo.GetUserByUserName(newUser.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if hasher.ComparePassword(user.Password, json.Password) != nil {
		return nil, errors.New("invalid password")
	}

	refreshParams := utils.RefreshParams{
		IpAddress: clientIp,
		UserAgent: userAgent,
		User:      user,
	}

	var refreshUtils = refreshParams
	refreshToken, err := refreshUtils.CreateJWT()

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	accessParams := utils.AccessParams{
		User: user,
	}

	var accessUtils = accessParams
	accessToken, err := accessUtils.CreateJWT()

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	ectryptedRefreshToken, err := utils.Encrypt(refreshToken)
	fmt.Println("before crypt", refreshToken)
	fmt.Println("after crypt", ectryptedRefreshToken)

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	ctx.SetCookie(
		"refresh_token",       // cookie name
		ectryptedRefreshToken, // cookie value
		3600*24*7,             // max age in seconds (1 week)
		"/",                   // path
		"",                    // domain (leave empty for default)
		true,                  // secure flag
		true,                  // HttpOnly flag
	)

	ctx.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "user": user})

	s.db.Model(user).Update("ip_address", clientIp)
	s.db.Model(user).Update("user_agent", userAgent)

	return user, nil
}

func (s *authService) RefreshToken(c *gin.Context) (string, error) {
	clientIp := c.ClientIP()
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		return "", errors.New("refresh token not found")
	}

	decryptedRefreshToken, err := utils.Decrypt(refreshToken)

	fmt.Println("refresh", decryptedRefreshToken)

	if err != nil {
		return "", errors.New("access error")
	}

	refreshParams := utils.RefreshParams{
		IpAddress: clientIp,
		UserAgent: c.GetHeader("User-Agent"),
	}

	var refreshUtils = refreshParams
	userData, err := refreshUtils.ExtractUserFromToken(decryptedRefreshToken)

	if err != nil {
		return "", errors.New("access error")
	}

	refreshClaims := userData.(*utils.RefreshClaims)

	userRepo := repositories.NewUserRepository(s.db)

	var user *models.User

	user, err = userRepo.GetByID(refreshClaims.ID)

	if err != nil {
		return "", errors.New("user not found")
	}

	if refreshClaims.IpAddress != clientIp || refreshClaims.UserAgent != c.GetHeader("User-Agent") {
		return "", errors.New("invalid token")
	}

	if _, err = refreshUtils.VerifyToken(decryptedRefreshToken); err != nil {
		return "", errors.New("invalid token")
	}

	jwtToken, err := utils.AccessParams{
		User: user,
	}.CreateJWT()

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return jwtToken, nil
}
