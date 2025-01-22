package scraper

type TmdFeedItemPropsType struct {
	Lat         string
	Long        string
	Magnitude   string
	Depth       string
	Time        string
	Title       string
	Description string
	Link        string
}

type GfzFeedItemPropsType struct {
	Title       string
	Description string
	Guid        struct {
		Content string `json:"#content"`
	}
}

type UsgsFeedItemPropsType struct {
	ID         string
	Properties struct {
		Mag   float64
		Title string
		Time  int64
		Place string
	}
	Geometry struct {
		Coordinates [3]float64
	}
}
