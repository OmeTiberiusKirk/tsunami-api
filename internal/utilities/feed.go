package utilities

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
	"tsunamiApi/models"

	xj "github.com/basgys/goxml2json"
)

func ConvertXmlToJson[Item comparable](d string) ([]Item, error) {
	type FeedPropsType struct {
		Rss struct {
			Channel struct{ Item []Item }
		}
	}

	d = regexp.MustCompile(`<!\[CDATA\[|<br>|]]>`).ReplaceAllString(d, " ")
	reader := strings.NewReader(d)
	data, err := xj.Convert(reader)
	if err != nil {
		return nil, err
	}

	res := FeedPropsType{}
	if err := json.Unmarshal(data.Bytes(), &res); err != nil {
		return nil, err
	}
	return res.Rss.Channel.Item, nil
}

func (p TmdFeedItemPropsType) ModifyPropTypesOfFeedItem() models.Earthquake {
	uid := regexp.MustCompile("earthquake=").Split(p.Link, 2)[1]
	lat, _ := strconv.ParseFloat(p.Lat, 64)
	long, _ := strconv.ParseFloat(p.Long, 64)
	mag, _ := strconv.ParseFloat(p.Magnitude, 64)
	depth, _ := strconv.ParseFloat(p.Depth, 64)
	t := regexp.MustCompile(`\s+UTC`).ReplaceAllString(p.Time, "Z")
	t = regexp.MustCompile(`\s+`).ReplaceAllString(t, "T")
	timeOfEvent, _ := time.Parse(time.RFC3339, t)

	return models.Earthquake{
		UID:          uid,
		Latitude:     lat,
		Longitude:    long,
		Magnitude:    mag,
		Depth:        depth,
		Time:         timeOfEvent,
		Title:        p.Title,
		Description:  p.Description,
		FeedProvider: "tmd",
		UpdatedAt:    time.Now(),
	}
}

func (p GfzFeedItemPropsType) ModifyPropTypesOfFeedItem() models.Earthquake {
	reg := regexp.MustCompile(`\s+`)
	title := reg.Split(p.Title, -1)
	desc := reg.Split(p.Description, -1)
	lat, _ := strconv.ParseFloat(desc[2], 64)
	long, _ := strconv.ParseFloat(desc[3], 64)
	mag, _ := strconv.ParseFloat(title[1][:len(title[1])-1], 64)
	depth, _ := strconv.ParseFloat(desc[4], 64)
	timeOfEvent, _ := time.Parse(time.RFC3339, desc[0]+"T"+desc[1]+"Z")
	return models.Earthquake{
		UID:          p.Guid.Content,
		Latitude:     lat,
		Longitude:    long,
		Magnitude:    mag,
		Depth:        depth,
		Time:         timeOfEvent,
		Title:        p.Title,
		Description:  p.Description,
		FeedProvider: "gfz",
		UpdatedAt:    time.Now(),
	}
}

func (p UsgsFeedItemPropsType) ModifyPropTypesOfFeedItem() models.Earthquake {
	lat := p.Geometry.Coordinates[1]
	long := p.Geometry.Coordinates[0]
	title := p.Properties.Title
	mag := p.Properties.Mag
	depth := p.Geometry.Coordinates[2]
	timeOfEvent := time.UnixMilli(p.Properties.Time).UTC()
	desc := p.Properties.Place
	return models.Earthquake{
		UID:          p.ID,
		Latitude:     lat,
		Longitude:    long,
		Magnitude:    mag,
		Depth:        depth,
		Time:         timeOfEvent,
		Title:        title,
		Description:  desc,
		FeedProvider: "usgs",
		UpdatedAt:    time.Now(),
	}
}

func FilterEarthquakesByArea(eq *[]models.Earthquake) {
	minLat, maxLat, minLong, maxLong := -5, 25, 87, 123
	var tmp []models.Earthquake
	for _, v := range *eq {
		if v.Latitude > float64(minLat) &&
			v.Latitude < float64(maxLat) &&
			v.Longitude > float64(minLong) &&
			v.Longitude < float64(maxLong) {
			tmp = append(tmp, v)
		}
	}
	*eq = tmp
}
