package databases

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TdssDB struct {
	DB *gorm.DB
}

func NewTdssDB() *TdssDB {
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

	return &TdssDB{DB: db}
}

func (hub TdssDB) FindAllBulletins(c *gin.Context) {
	var (
		query struct {
			Latitude  float64 `form:"latitude"`
			Longitude float64 `form:"longitude"`
			Magnitude float64 `form:"magnitude"`
			Depth     float64 `form:"depth"`
		}
		result []map[string]interface{}
	)

	if c.ShouldBindQuery(&query) == nil {
		var depth int

		if query.Depth >= 0 && query.Depth <= 29.9 {
			depth = 10
		} else {
			depth = 30
		}

		hub.DB.Table("sim_result").Select(
			"id",
			"job_profile_id",
			"name",
			"magnitude",
			"depth",
			"decimal_lat",
			"decimal_long",
			fmt.Sprintf("ROUND(SQRT(POW(`decimal_lat`-(%v), 2) + POW(`decimal_long`-(%v), 2)), 3) AS `R`", query.Latitude, query.Longitude),
		).Where(
			"grp_id = ? AND magnitude >= ? AND depth >= ?",
			1, query.Magnitude, depth,
		).Order("depth, magnitude, R").Take(&result)
	}

	// var data = map[string]interface{}{
	// 	"eta_results": result,
	// }

	c.JSON(http.StatusOK, result)
}
