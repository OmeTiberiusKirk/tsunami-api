package main

import (
	"flag"
	"tsunamiApi/internal/databases"
	"tsunamiApi/internal/middleware"
	"tsunamiApi/internal/tdss"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":80", "http service address")

func main() {
	flag.Parse()
	db := databases.ConnectMRDB()
	tdss := tdss.NewTdss(db)
	router := gin.Default()

	router.Use(middleware.GinMiddleware("http://localhost:3000"))
	router.GET("/getbulletin", tdss.Controllers().ServObservationPoints)

	router.Run(*addr)
}
