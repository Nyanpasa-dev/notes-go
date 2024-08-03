package routes

import (
	services "simple-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	var users = services.NewUserService(db)

	r.POST("/users", users.CreateUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:id", users.GetUser)
	r.PUT("/users/:id", users.UpdateUser)
	r.DELETE("/users/:id", users.DeleteUser)
}