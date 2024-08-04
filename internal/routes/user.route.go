package routes

import (
	"simple-api/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	var users = handlers.NewUserHandler(db)

	r.POST("/users", users.CreateUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:id", users.GetUserById)
	r.PUT("/users/:id", users.UpdateUser)
	r.DELETE("/users/:id", users.DeleteUser)
}
