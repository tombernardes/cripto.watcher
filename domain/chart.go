package domain

type Chart struct {
	Ticker string
	Points []Point
}

type Point struct {
	Ticker     string
	OpenTime   int64
	CloseTime  int64
	OpenPrice  float64
	ClosePrice float64
	MinPrice   float64
	MaxPrice   float64
	Volume     float64
	Trades     int64
}
