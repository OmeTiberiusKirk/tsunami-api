package main

import (
	"net/http"
	"tsunami/api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.GinMiddleware("http://localhost:3000"))
	router.GET("/getbulletin", getAlbumByID)

	router.Run("localhost:8080")

}

func getAlbumByID(c *gin.Context) {
	// id := c.Param("id")

	c.String(http.StatusOK, "garon")

	// c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
