package earthquake

import (
	"gorm.io/gorm"
)

type EarthquakeIntf interface {
	Services() ServicesIntf
}

type Earthquake struct {
	DB *gorm.DB
}

func New(db *gorm.DB) (tdss EarthquakeIntf) {
	tdss = &Earthquake{
		DB: db,
	}
	return
}
