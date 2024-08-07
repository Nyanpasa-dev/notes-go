package handlers

import (
	"net/http"
	"simple-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		service: services.NewUserService(db),
	}
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

	err := h.service.CreateUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	user, err := h.service.GetUsers(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": user})
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	user, err := h.service.GetUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {

	err := h.service.UpdateUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	err := h.service.DeleteUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
