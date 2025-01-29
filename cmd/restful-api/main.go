package main

import (
	"tsunami-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	controller := NewTdssController()

	router := gin.Default()
	router.Use(middleware.GinMiddleware("http://localhost:3000"))
	router.GET("/getbulletin", controller.FindObservationPoints)

	router.Run("localhost:8080")

}
