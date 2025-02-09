package main

import (
	"tsunamiApi/internal/databases"
	"tsunamiApi/internal/middleware"
	"tsunamiApi/internal/tdss"

	"github.com/gin-gonic/gin"
)

func main() {
	db := databases.ConnectMRDB()
	tdss := tdss.NewTdss(db)
	router := gin.Default()

	router.Use(middleware.GinMiddleware("http://localhost:3000"))
	router.GET("/getbulletin", tdss.Controllers().ServObservationPoints)

	router.Run("localhost:8080")
}
