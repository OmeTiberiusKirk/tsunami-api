package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tsunami/api/internal/database"
	"tsunami/api/internal/route"
	"tsunami/api/internal/scraper"
	"tsunami/api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	router := gin.New()
	db := database.Open()
	db.AutoMigrate(&models.Earthquake{})
	server := socketio.NewServer(nil)

	scraper.FeedTask(scraper.FeedTaskArgs{DB: db, Server: server})

	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.NewJob(
		gocron.CronJob(
			"*/5 * * * *",
			false,
		),
		gocron.NewTask(func() {
			scraper.FeedTask(scraper.FeedTaskArgs{
				DB:            db,
				Server:        server,
				WithBroadcast: true,
			})
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// start the scheduler
	s.Start()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected: ", s.ID())
		return nil
	})

	server.OnEvent("/", "provoke", func(s socketio.Conn, msg string) {
		events := &[]models.Earthquake{}
		db.Model(&models.Earthquake{}).Find(&events)
		b, _ := json.Marshal(&events)
		s.Emit("earthquakeData", string(b))
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	router.Use(route.GinMiddleware("http://localhost:3000"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("./assets"))

	if err := router.Run(":8081"); err != nil {
		log.Fatal("failed run app: ", err)
	}

}
