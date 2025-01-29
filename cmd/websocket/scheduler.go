package main

import (
	"log"
	"tsunami-api/internal/databases"
	"tsunami-api/models"

	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm/clause"
)

func CreateScheduler(MainDB *databases.MainDB, hub *Hub) {
	Task(MainDB)

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
			Task(MainDB)
			hub.broadcast <- MainDB.FindAllEearthquakes()
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}

func Task(MainDB *databases.MainDB) {
	columns := []string{
		"latitude",
		"longitude",
		"magnitude",
		"depth",
		"time",
		"title",
		"description",
		"feed_provider",
	}
	tmd := make(chan []models.Earthquake)
	gfz := make(chan []models.Earthquake)
	usgs := make(chan []models.Earthquake)

	go FetchTmd(tmd)
	go FetchGfz(gfz)
	go FetchUsgs(usgs)
	rs := append(<-tmd, <-gfz...)
	rs = append(rs, <-usgs...)

	close(tmd)
	close(gfz)
	close(usgs)

	FilterEarthquakesByArea(&rs)

	if len(rs) != 0 {
		MainDB.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uid"}},
			DoUpdates: clause.AssignmentColumns(columns),
		}).Create(rs)
	}

}
