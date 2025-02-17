package earthquake

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
	"tsunamiApi/internal/databases"
	"tsunamiApi/models"

	"github.com/joho/godotenv"
	"github.com/paulmach/orb"
)

func GetRecentEarthquakes() []byte {
	godotenv.Load(".env")
	isDev, _ := strconv.ParseBool(os.Getenv("DEV"))
	earthquakes := &[]models.Earthquake{}

	if isDev {
		databases.PGDB.Model(&models.Earthquake{}).Order("time desc").Find(&earthquakes)
	} else {
		now := time.Now().UTC()
		t := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
		databases.PGDB.Model(&models.Earthquake{}).Where("time > ?", t).Order("time desc").Find(&earthquakes)
	}

	b, e := json.Marshal(&earthquakes)
	if e != nil {
		return []byte{}
	}

	return b
}

func GetGeometryOfAndaman() orb.MultiPoint {
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
