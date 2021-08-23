package domain

type IndicatorsBase struct {
	Ticker                   string
	AverageTrueRange         []float64
	ExponentialMovingAverage []float64
	AverageVolumn            []float64
	Trend                    []float64
}
