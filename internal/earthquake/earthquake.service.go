package earthquake

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
	"tsunamiApi/models"

	"github.com/paulmach/orb"
	"gorm.io/gorm"
)

type ServicesIntf interface {
	GetRecentEarthquakes() []byte
	GetGeometryOfAndaman() orb.MultiPoint
	GetDB() *gorm.DB
}

func (eq *Earthquake) Services() (ser ServicesIntf) {
	ser = eq
	return
}

func (eq *Earthquake) GetRecentEarthquakes() []byte {
	earthquakes := &[]models.Earthquake{}
	isDev, _ := strconv.ParseBool(os.Getenv("DEV"))
	if isDev {
		eq.DB.Model(&models.Earthquake{}).Order("time desc").Find(&earthquakes)
	} else {
		now := time.Now().UTC()
		t := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
		eq.DB.Model(&models.Earthquake{}).Where("time > ?", t).Order("time desc").Find(&earthquakes)
	}

	b, e := json.Marshal(&earthquakes)
	if e != nil {
		return []byte{}
	}

	return b
}

func (eq *Earthquake) GetGeometryOfAndaman() orb.MultiPoint {
	b, err := os.ReadFile("data/tsunami.geojson")
	if err != nil {
		panic(err)
	}

	j := map[string]orb.MultiPoint{}
	if err := json.Unmarshal(b, &j); err != nil {
		panic(err)
	}

	return j["coordinates"]
}

func (eq *Earthquake) GetDB() *gorm.DB {
	return eq.DB
}
