package services

import (
	"errors"
	"simple-api/internal/repositories"
	"simple-api/models"
	"simple-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userService struct {
	userRepo repositories.UserRepository
}
type UserService interface {
	CreateUser(c *gin.Context) error
	GetUsers(c *gin.Context) ([]models.User, error)
	GetUser(c *gin.Context) (*models.User, error)
	UpdateUser(c *gin.Context) error
	DeleteUser(c *gin.Context) error
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		userRepo: repositories.NewUserRepository(db),
	}
}

func (s *userService) CreateUser(c *gin.Context) error {
	var err error
	hasher := utils.BcryptHasher{}

	var json struct {
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
		Avatar          string `json:"avatar"`
		IsAdmin         bool   `json:"isAdmin"`
	}

	if json.Password != json.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	json.Password, err = hasher.HashPassword(json.Password)

	if err != nil {
		return errors.New("failed to hash password")
	}

	user := models.User{Username: json.Username, Password: json.Password, Avatar: json.Avatar}
	s.userRepo.Create(&user)

	return nil

}

func (s *userService) GetUsers(c *gin.Context) ([]models.User, error) {
	var users []models.User
	var err error

	users, err = s.userRepo.GetAll()

	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	return users, nil
}

func (s *userService) GetUser(c *gin.Context) (*models.User, error) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(uint(id))
	if err != nil {
		return nil, errors.New("failed to fetch user")
	}

	return user, nil
}

func (s *userService) UpdateUser(c *gin.Context) error {
	var user *models.User
	var err error

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.New("invalid user ID")
	}

	if user, err = s.userRepo.GetByID(uint(id)); err != nil {
		return errors.New("user not found")
	}

	var json struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password"`
		Avatar   string `json:"avatar"`
	}

	if err := c.Bind(&json); err != nil {
		return errors.New("invalid request")
	}

	user.Username = json.Username
	user.Password = json.Password
	user.Avatar = json.Avatar

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update user")
	}
	return nil
}

func (s *userService) DeleteUser(c *gin.Context) error {
	var err error

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.New("invalid user ID")
	}

	if _, err = s.userRepo.GetByID(uint(id)); err != nil {
		return errors.New("user not found")
	}

	if err := s.userRepo.Delete(uint(id)); err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}
