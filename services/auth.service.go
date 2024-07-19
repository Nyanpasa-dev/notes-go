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
}

func (s *authService) Login(c *gin.Context) {
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

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": jwtToken, "user": user})
		}
	}

}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db}
}
