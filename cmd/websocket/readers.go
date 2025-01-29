package main

import (
	"encoding/json"
	"os"
	"tsunami-api/models"
)

func ReadTmd(c chan []models.Earthquake) {
	b, err := os.ReadFile("data/tmd.xml")
	if err != nil {
		panic(err)
	}

	if items, err := ConvertXmlToJson[TmdFeedItemPropsType](string(b)); err != nil {
		c <- nil
	} else {
		var events []models.Earthquake
		for _, v := range items {
			events = append(events, v.ModifyPropsTypeOfFeedItem())
		}
		c <- events
	}
}

func ReadGfz(c chan []models.Earthquake) {
	b, err := os.ReadFile("data/gfz.xml")
	if err != nil {
		panic(err)
	}

	if items, err := ConvertXmlToJson[GfzFeedItemPropsType](string(b)); err != nil {
		c <- nil
	} else {
		var events []models.Earthquake
		for _, v := range items {
			events = append(events, v.ModifyPropsTypeOfFeedItem())
		}
		c <- events
	}
}

func ReadUsgs(c chan []models.Earthquake) {
	res, err := os.ReadFile("data/usgs.json")
	if err != nil {
		panic(err)
	}

	var dat struct {
		Features []UsgsFeedItemPropsType
	}
	if err := json.Unmarshal(res, &dat); err != nil {
		c <- nil
	}

	var events []models.Earthquake
	for _, v := range dat.Features {
		events = append(events, v.ModifyPropsTypeOfFeedItem())
	}

	c <- events
}
