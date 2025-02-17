package utilities

import (
	"log"
	"os"
	"strconv"
	"tsunamiApi/internal/databases"
	"tsunamiApi/internal/earthquakes"
	"tsunamiApi/internal/websocket"
	"tsunamiApi/models"

	"github.com/go-co-op/gocron/v2"
	"github.com/paulmach/orb"
	"gorm.io/gorm/clause"
)

func CreateScheduler(
	ws *websocket.Hub,
) {
	Task()

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
			Task()
			ws.SetBroadcast(earthquakes.GetRecentEarthquakes())
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}

func Task() {
	isDev, _ := strconv.ParseBool(os.Getenv("DEV"))
	columns := []string{
		"latitude",
		"longitude",
		"magnitude",
		"depth",
		"time",
		"title",
		"description",
		"feed_provider",
		"updated_at",
	}
	tmd := make(chan []models.Earthquake)
	gfz := make(chan []models.Earthquake)
	usgs := make(chan []models.Earthquake)

	if isDev {
		go ReadTmd(tmd)
		go ReadGfz(gfz)
		go ReadUsgs(usgs)
	} else {
		go FetchTmd(tmd)
		go FetchGfz(gfz)
		go FetchUsgs(usgs)
	}

	rs := append(<-tmd, <-gfz...)
	rs = append(rs, <-usgs...)

	close(tmd)
	close(gfz)
	close(usgs)

	coordinates := earthquakes.GetGeometryOfAndaman()
	rs = Filter(rs, func(item models.Earthquake, idx int) bool {
		p := orb.Point{item.Longitude, item.Latitude}
		bound := coordinates.Bound()
		return bound.Contains(p)
	})

	if len(rs) != 0 {
		databases.PGDB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "uid"}},
			DoUpdates: clause.AssignmentColumns(columns),
		}).Create(rs)
	}

}
