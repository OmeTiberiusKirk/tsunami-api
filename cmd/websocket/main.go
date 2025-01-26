package main

import (
	"flag"
	"tsunami/api/internal/databases"
	"tsunami/api/internal/middleware"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8081", "http service address")

func main() {
	flag.Parse()

	router := gin.Default()
	router.Use(middleware.GinMiddleware("http://localhost:3000"))

	MainDB := databases.NewMainDB()
	hub := newHub()
	go hub.run()
	// create a scheduler
	createScheduler(MainDB, hub)

	router.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
		hub.broadcast <- MainDB.FindAllEearthquakes()
	})

	router.Run(*addr)

}
