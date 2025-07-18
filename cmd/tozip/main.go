package main

import (
	"smthtozip/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"
	router := gin.New()
	routes.Routes(router)
	router.Run(":" + port)
}
