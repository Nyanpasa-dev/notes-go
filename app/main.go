package main

import (
	"fmt"
	"simple-api/config"
	"simple-api/internal/routes"
	"simple-api/utils/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	r := gin.Default()

	db := config.GetDB()
	r.Use(middleware.ErrorHandlingMiddleware())
	routes.RegisterRoutes(r, db)

	port := config.AppConfig.Application.Port
	r.Run(fmt.Sprintf(":%d", port))
}
