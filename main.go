package main

import (
	"flag"
	"tsunamiApi/internal/databases"
	"tsunamiApi/internal/earthquake"
	"tsunamiApi/internal/middleware"
	"tsunamiApi/internal/tdss"
	"tsunamiApi/internal/utilities"
	"tsunamiApi/internal/websocket"
	"tsunamiApi/models"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8080", "api address")

func init() {
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

	router.Use(middleware.GinMiddleware("*"))
	router.GET("/getbulletin", tdss.ServObservationPoints)
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(c.Writer, c.Request)
		ws.SetBroadcast(earthquake.GetRecentEarthquakes())
	})

	router.Run(*addr)
}
