package databases

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

type MainDB struct {
	DB *gorm.DB
}

func NewMainDB() *MainDB {
	godotenv.Load(".env")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("PG_DB_HOST"),
		os.Getenv("PG_DB_USER"),
		os.Getenv("PG_DB_PASS"),
		os.Getenv("PG_DB_NAME"),
		os.Getenv("PG_DB_PORT"),
		os.Getenv("PG_DB_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Earthquake{})

	return &MainDB{DB: db}
}

func (MainDB MainDB) FindAllEearthquakes() []byte {
	earthquakes := &[]models.Earthquake{}
	now := time.Now().UTC()
	t := now.AddDate(0, 0, -1).Truncate(24 * time.Hour)
	MainDB.DB.Model(&models.Earthquake{}).Where("time > ?", t).Order("time desc").Find(&earthquakes)
	b, e := json.Marshal(&earthquakes)

	if e != nil {
		return []byte{}
	}

	return b
}
