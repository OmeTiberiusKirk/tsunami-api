package tdss

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServObservationPoints(ctx *gin.Context) {
	var (
		query struct {
			Latitude  float64 `form:"latitude"`
			Longitude float64 `form:"longitude"`
			Magnitude float64 `form:"magnitude"`
			Depth     float64 `form:"depth"`
		}
	)
	ctx.Bind(&query)

	result := map[string]interface{}{
		"results": map[string]interface{}{
			"eta_results": []interface{}{},
		}}
	simId := GetSimResult(
		query.Magnitude,
		query.Depth,
		query.Latitude,
		query.Longitude,
	)
	if simId != 0 {
		result["results"] = map[string]interface{}{
			"eta_results": FindObservationPoints(simId),
		}
	}

	ctx.JSON(http.StatusOK, result)
}
