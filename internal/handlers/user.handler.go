package handlers

import (
	"net/http"
	"simple-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var userRequest struct {
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
		Avatar          string `json:"avatar"`
		IsAdmin         bool   `json:"isAdmin"`
	}

	if ctx.ShouldBind(&userRequest) == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userService := services.NewUserService(h.db)

	err := userService.CreateUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}

func (h *UserHandler) GetUsers() {

}

func (h *UserHandler) GetUserById() {

}

func (h *UserHandler) UpdateUser() {

}

func (h *UserHandler) DeleteUser() {

}
