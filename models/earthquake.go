package models

import (
	"time"

	"gorm.io/gorm"
)

type Earthquake struct {
	gorm.Model
	UID          string `gorm:"unique"`
	Latitude     float64
	Longitude    float64
	Magnitude    float64
	Depth        float64
	Time         time.Time
	Title        string
	Description  string
	FeedProvider string
}
