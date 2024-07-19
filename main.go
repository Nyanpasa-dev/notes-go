package main

import (
	"simple-api/config"
	"simple-api/middleware"
	"simple-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := config.GetDB()
	r.Use(middleware.ErrorHandlingMiddleware())
	routes.RegisterRoutes(r, db)

	r.Run(":8080")
}
