package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TdssServices struct {
	DB *gorm.DB
}

func NewTdssServices() TdssServices {
	godotenv.Load(".env")
	dsn := fmt.Sprintf(
		"%v:%v@tcp(127.0.0.1:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MR_DB_USER"),
		os.Getenv("MR_DB_PASS"),
		os.Getenv("MR_DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return TdssServices{DB: db}
}

func (s TdssServices) GetSimResult(mag float64, depth float64, lat float64, long float64) int {
	var (
		result []struct {
			Id int
		}
		d int
	)

	if depth >= 0 && depth <= 29.9 {
		d = 10
	} else {
		d = 30
	}

	s.DB.Table("sim_result").Select(
		"id",
		"job_profile_id",
		"name",
		"magnitude",
		"depth",
		"decimal_lat",
		"decimal_long",
		fmt.Sprintf("ROUND(SQRT(POW(`decimal_lat`-(%v), 2) + POW(`decimal_long`-(%v), 2)), 3) AS `R`", lat, long),
	).Where(
		"grp_id = ? AND magnitude >= ? AND depth >= ?",
		1, mag, d,
	).Order("depth, magnitude, R").Take(&result)

	if len(result) != 0 {
		fmt.Println("result")
		fmt.Println(result[0].Id)
		return result[0].Id
	}

	return 0
}

func (s TdssServices) FindObservationPoints(simId int) []map[string]interface{} {
	var result []map[string]interface{}
	sql := []string{
		"SELECT",
		"`observe_point`.`observ_point_id`,",
		"`observe_point`.`province_t`,",
		"`observe_point`.`name_t`,",
		"`observe_point`.`lat_1`,",
		"`observe_point`.`lat_2`,",
		"`observe_point`.`lat_3`,",
		"`observe_point`.`long_1`,",
		"`observe_point`.`long_2`,",
		"`observe_point`.`long_3`,",
		"`observe_point`.`decimal_lat`,",
		"`observe_point`.`decimal_long`,",
		"`sim_point_val`.`values`,",
		"`sim_point_val`.`type`,",
		"`sim_point_val`.`region_no`",
		"FROM `observe_point`, `sim_point_val`",
		"WHERE",
		fmt.Sprintf("`sim_point_val`.`sim_result_id` = %d", simId),
		"AND `observe_point`.`observ_point_id` = `sim_point_val`.`id_point`",
		"AND `sim_point_val`.`type` = 'ETA'",
	}

	s.DB.Raw(strings.Join(sql, " ")).Scan(&result)
	return result
}
