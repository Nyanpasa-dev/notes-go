package main

import (
	"simple-api/config"
	"simple-api/internal/routes"
	"simple-api/utils/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := config.GetDB()
	r.Use(middleware.ErrorHandlingMiddleware())
	routes.RegisterRoutes(r, db)

	r.Run(":8080")
}
