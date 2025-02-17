package main

import (
	"flag"
	"tsunamiApi/internal/databases"
	"tsunamiApi/internal/earthquakes"
	"tsunamiApi/internal/middlewares"
	"tsunamiApi/internal/tdss"
	"tsunamiApi/internal/users"
	"tsunamiApi/internal/utilities"
	"tsunamiApi/internal/websocket"
	"tsunamiApi/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var addr = flag.String("addr", ":8080", "api address")

func init() {
	godotenv.Load(".env")
	databases.ConnectPGDB()
	databases.ConnectMRDB()
	databases.PGDB.AutoMigrate(&models.Earthquake{})
}

func main() {
	flag.Parse()
	router := gin.Default()
	ws := websocket.New()
	go ws.Run()
	utilities.CreateScheduler(ws)

	router.Use(middlewares.GinMiddleware("*"))
	router.POST("/login", users.Login)

	router.GET("/getbulletin", tdss.ServObservationPoints)
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(c.Writer, c.Request)
		ws.SetBroadcast(earthquakes.GetRecentEarthquakes())
	})

	router.Run(*addr)
}
