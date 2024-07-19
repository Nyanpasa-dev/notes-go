package services

import (
	"net/http"
	"simple-api/models"
	"simple-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userService struct {
	db *gorm.DB
}

type UserService interface {
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func (s *userService) CreateUser(c *gin.Context) {
	_, err := utils.CheckAdmin(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// Parse JSON
	var json struct {
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
		Avatar          string `json:"avatar"`
		Login           string `json:"login"`
		Email           string `json:"email"`
		IsAdmin         bool   `json:"isAdmin"`
	}

	if c.Bind(&json) == nil {
		if json.Password != json.ConfirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Passwords do not match",
			})
		}
		// Create user
		user := models.User{Username: json.Username, Password: json.Password, Avatar: json.Avatar}
		s.db.Create(&user)
		c.JSON(http.StatusCreated, user)
	}
}

func (s *userService) GetUsers(c *gin.Context) {
	var users []models.User

	s.db.Find(&users)

	c.JSON(http.StatusOK, users)
}

func (s *userService) GetUser(c *gin.Context) {
	var user models.User
	if err := s.db.First(&user, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func (s *userService) UpdateUser(c *gin.Context) {
	var user models.User
	if err := s.db.First(&user, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	} else {
		var json struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password"`
			Avatar   string `json:"avatar"`
		}

		if c.Bind(&json) == nil {
			user.Username = json.Username
			user.Password = json.Password
			user.Avatar = json.Avatar
			s.db.Save(&user)
			c.JSON(http.StatusOK, user)
		}
	}
}

func (s *userService) DeleteUser(c *gin.Context) {
	var user models.User
	if err := s.db.First(&user, c.Param("id")).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
	} else {
		s.db.Delete(&user)
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db}
}
