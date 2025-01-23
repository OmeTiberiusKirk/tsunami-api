package database

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"tsunami/api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBController struct {
	DB *gorm.DB
}

func Open() DBController {
	godotenv.Load(".env")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return DBController{DB: db}
}

func (dbC DBController) FindAllEearthquakes() []byte {
	earthquakes := &[]models.Earthquake{}
	now := time.Now().UTC()
	t := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
	dbC.DB.Model(&models.Earthquake{}).Where("time > ?", t).Order("time desc").Find(&earthquakes)
	b, e := json.Marshal(&earthquakes)

	if e != nil {
		return []byte{}
	}

	return b
}
