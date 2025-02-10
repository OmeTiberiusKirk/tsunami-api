package main

import (
	"flag"
	"tsunamiApi/internal/databases"
	"tsunamiApi/internal/earthquake"
	"tsunamiApi/internal/middleware"
	"tsunamiApi/internal/utilities"
	"tsunamiApi/internal/websocket"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	router := gin.Default()
	router.Use(middleware.GinMiddleware("http://localhost:3000"))
	db := databases.ConnectPGDB()
	ws := websocket.New(db)
	go ws.Run()
	eq := earthquake.New(db)
	utilities.CreateScheduler(eq, ws)

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(c.Writer, c.Request)
		ws.SetBroadcast(eq.Services().GetRecentEarthquakes())
	})

	router.Run(*addr)
}
