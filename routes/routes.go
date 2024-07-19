package routes

import (
	"simple-api/middleware"
	"simple-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	var notes = services.NewNoteService(db)
	var users = services.NewUserService(db)
	var auth = services.NewAuthService(db)

	r.POST("/login", auth.Login)
	r.Use(middleware.CheckAuthMiddleware())
	r.POST("/notes", notes.CreateNote)
	r.GET("/notes", notes.GetNotes)
	r.GET("/notes/:id", notes.GetNote)
	r.PUT("/notes/:id", notes.UpdateNote)

	r.POST("/users", users.CreateUser)
	r.GET("/users", users.GetUsers)
	r.GET("/users/:id", users.GetUser)
	r.PUT("/users/:id", users.UpdateUser)
	r.DELETE(("users/:id"), users.DeleteUser)

}
