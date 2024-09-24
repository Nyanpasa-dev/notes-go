package main

import (
	"fmt"
	"simple-api/config"
	"simple-api/internal/routes"
	"simple-api/utils/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	if recover() != nil {
		fmt.Println("Error loading config")
	}

	r := gin.Default()

	db := config.GetDB()
	r.Use(middleware.ErrorHandlingMiddleware())
	routes.RegisterRoutes(r, db)

	port := config.AppConfig.Application.Port
	err := r.Run(fmt.Sprintf(":%d", port))

	if err != nil {
		panic(err)
	}
}

func init() {
	config.LoadConfig()
	config.RunDB()
}
