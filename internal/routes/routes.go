package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	RegisterAuthRoutes(r, db)
	RegisterUserRoutes(r, db)
	RegisterNoteRoutes(r, db)
}
