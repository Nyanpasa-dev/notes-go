package services

import (
	"net/http"
	"simple-api/models"
	"simple-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authService struct {
	db *gorm.DB
}

type AuthService interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

func (s *authService) Login(c *gin.Context) {
	clientIp := c.ClientIP()

	userAgent := c.GetHeader("User-Agent")

	if userAgent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User-Agent header is required"})
		return
	}

	// Parse JSON
	var json struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if c.Bind(&json) == nil {
		user := models.User{Username: json.Username, Password: json.Password}

		if err := s.db.Select("username", "id").Where(&user).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неправильно нихуя"})
		}

		jwtToken, err := utils.CreateJWT(&user)
		refreshToken, err := utils.CreateRefreshToken(&user, clientIp)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": jwtToken, "user": user})
		}
	}

}

func (s *authService) RefreshToken(c *gin.Context) {
	clientIp := c.ClientIP()
	refreshToken := c.GetHeader("Authorization")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	userData, err := utils.ExtractUserFromToken(refreshToken)

	if userData.IpAdress != clientIp || userData.UserAgent != c.GetHeader("User-Agent") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Why are you stealing tokens?"})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user := models.User{ID: userData.ID}

	if err := s.db.First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if userData.IpAdress != clientIp || userData.UserAgent != c.GetHeader("User-Agent") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Why are you stealing tokens?"})
		return
	}

	jwtToken, err := utils.CreateJWT(&user{
		ID:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	})

}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db}
}
