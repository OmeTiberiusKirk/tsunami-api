package tdss

import (
	"fmt"
	"strings"
	"tsunamiApi/internal/databases"
)

func GetSimResult(
	mag float64,
	depth float64,
	lat float64,
	long float64,
) int {
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

	databases.MRDB.Table("sim_result").Select(
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

func FindObservationPoints(simId int) []map[string]interface{} {
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

	databases.MRDB.Raw(strings.Join(sql, " ")).Scan(&result)
	return result
}
