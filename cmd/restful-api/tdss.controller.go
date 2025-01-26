package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TdssController struct {
	services TdssServices
}

func NewTdssController() TdssController {
	services := NewTdssServices()
	return TdssController{services: services}
}

func (c TdssController) FindObservationPoints(ct *gin.Context) {
	var (
		query struct {
			Latitude  float64 `form:"latitude"`
			Longitude float64 `form:"longitude"`
			Magnitude float64 `form:"magnitude"`
			Depth     float64 `form:"depth"`
		}
	)
	ct.Bind(&query)

	simId := c.services.GetSimResult(query.Magnitude, query.Depth, query.Latitude, query.Longitude)
	result := c.services.FindObservationPoints(simId)

	ct.JSON(http.StatusOK, map[string]interface{}{"eta_results": result})
}
