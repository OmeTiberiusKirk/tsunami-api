package main

import (
	"flag"
	"tsunami/api/internal/database"
	"tsunami/api/internal/middleware"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8081", "http service address")

func main() {
	flag.Parse()

	router := gin.Default()
	router.Use(middleware.GinMiddleware("http://localhost:3000"))

	dbHub := database.NewDB()
	hub := newHub()
	go hub.run()
	// create a scheduler
	createScheduler(dbHub, hub)

	router.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
		hub.broadcast <- dbHub.FindAllEearthquakes()
	})

	router.Run(*addr)

}
