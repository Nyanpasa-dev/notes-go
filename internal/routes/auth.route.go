package routes

import (
	"simple-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
	var users = handlers.NewAuthHandler(db)

	r.POST("/login", users.Authenticate)
	r.POST("/refresh", users.RefreshToken)
}
