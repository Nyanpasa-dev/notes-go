package routes

import (
	services "simple-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
	var users = services.NewAuthService(db)

	r.POST("/login", users.Login)
	r.POST("/refresh", users.RefreshToken)
}
