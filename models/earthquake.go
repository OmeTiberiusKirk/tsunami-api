package models

import (
	"time"

	"gorm.io/gorm"
)

type Earthquake struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UID          string         `gorm:"unique" json:"uid"`
	Latitude     float64        `json:"latitude"`
	Longitude    float64        `json:"longitude"`
	Magnitude    float64        `json:"magnitude"`
	Depth        float64        `json:"depth"`
	Time         time.Time      `json:"local_time"`
	Title        string         `json:"location"`
	Description  string         `json:"description"`
	FeedProvider string         `json:"feed_from"`
	CreatedAt    time.Time      `json:"created_date"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
}
