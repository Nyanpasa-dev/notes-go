package routes

import (
	services "simple-api/internal/services"
	"simple-api/utils/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterNoteRoutes(r *gin.Engine, db *gorm.DB) {
	var notes = services.NewNoteService(db)

	r.Use(middleware.CheckAuthMiddleware())
	r.POST("/notes", notes.CreateNote)
	r.GET("/notes", notes.GetNotes)
	r.GET("/notes/:id", notes.GetNote)
	r.PUT("/notes/:id", notes.UpdateNote)
	r.DELETE("/notes/:id", notes.DeleteNote)
}
