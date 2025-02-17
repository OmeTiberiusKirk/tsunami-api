package databases

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PGDB *gorm.DB
var MRDB *gorm.DB

func ConnectPGDB() {
	godotenv.Load(".env")
	var err error

	dsn := fmt.Sprintf(
		`host=%s user=%s 
		password=%s dbname=%s port=%s
		sslmode=disable TimeZone=%s`,
		os.Getenv("PG_DB_HOST"),
		os.Getenv("PG_DB_USER"),
		os.Getenv("PG_DB_PASS"),
		os.Getenv("PG_DB_NAME"),
		os.Getenv("PG_DB_PORT"),
		os.Getenv("PG_DB_TIMEZONE"),
	)

	PGDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func ConnectMRDB() {
	var err error
	godotenv.Load(".env")
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MR_DB_USER"),
		os.Getenv("MR_DB_PASS"),
		os.Getenv("MR_DB_HOST"),
		os.Getenv("MR_DB_NAME"),
	)
	MRDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
