package helpers

import (
	"math"
	"strconv"
	"time"

	"github.com/tombernardes/cripto.watcher/domain"
)

var loc *time.Location

func init() {
	loc, _ = time.LoadLocation("America/Sao_Paulo")
}

func RoundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

func FloatFromString(value string) float64 {
	if s, err := strconv.ParseFloat(value, 64); err == nil {
		return s
	} else {
		return 0
	}
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ParseDateTime(str string) time.Time {
	format := "2006-01-02T15:04:05.999"
	t, _ := time.ParseInLocation(format, str, loc)
	return t
}

func FormatDateTime(t time.Time) string {
	format := "2006-01-02T15:04:05.999"
	return t.Local().In(loc).Format(format)
}

func ContainsAndIndexString(s []string, e string) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, 0
}

func ContainsAndIndexBook(s []domain.Book, ticker string) (bool, int) {
	for i, a := range s {
		if a.Ticker == ticker {
			return true, i
		}
	}
	return false, 0
}

func ContainsAndIndexTrade(s []domain.TimesAndTrades, ticker string) (bool, int) {
	for i, a := range s {
		if a.Ticker == ticker {
			return true, i
		}
	}
	return false, 0
}
