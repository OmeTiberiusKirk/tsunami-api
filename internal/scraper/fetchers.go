package scraper

import (
	"encoding/json"
	"tsunami/api/models"

	"github.com/go-zoox/fetch"
)

func FetchTmd(c chan []models.Earthquake) {
	resp, err := fetch.Get("https://earthquake.tmd.go.th/feed/rss_tmd.xml")
	if err != nil {
		c <- nil
		return
	}

	if items, err := ConvertXmlToJson[TmdFeedItemPropsType](resp.String()); err != nil {
		c <- nil
	} else {
		var events []models.Earthquake
		for _, v := range items {
			events = append(events, v.ModifyPropsTypeOfFeedItem())
		}
		c <- events
	}
}

func FetchGfz(c chan []models.Earthquake) {
	resp, err := fetch.Get("https://geofon.gfz-potsdam.de/eqinfo/list.php?fmt=rss")
	if err != nil {
		c <- nil
		return
	}

	if items, err := ConvertXmlToJson[GfzFeedItemPropsType](resp.String()); err != nil {
		c <- nil
	} else {
		var events []models.Earthquake
		for _, v := range items {
			events = append(events, v.ModifyPropsTypeOfFeedItem())
		}
		c <- events
	}
}

func FetchUsgs(c chan []models.Earthquake) {
	resp, err := fetch.Get("https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson")
	if err != nil {
		c <- nil
		return
	}

	var dat struct {
		Features []UsgsFeedItemPropsType
	}

	if err := json.Unmarshal(resp.Body, &dat); err != nil {
		c <- nil
		return
	}

	var events []models.Earthquake
	for _, v := range dat.Features {
		events = append(events, v.ModifyPropsTypeOfFeedItem())
	}

	c <- events
}
