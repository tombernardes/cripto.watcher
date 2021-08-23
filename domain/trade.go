package domain

import (
	"sort"
	"time"
)

type Trade struct {
	Ticker   string
	Time     int64
	Quantity float64
	Price    float64
	Volume   float64
	IsMaker  bool
	WhoTake  string
}

type TimesAndTrades struct {
	Ticker     string
	Trades     []Trade
	VAP        []BuyersAndSellers
	VAPBuyers  []BuyersAndSellers
	VAPSellers []BuyersAndSellers
}

type BuyersAndSellers struct {
	Time          int64
	Price         float64
	Trades        int64
	TradesSellers int64
	TradesBuyers  int64
	Buyers        float64
	Sellers       float64
}

func (tt *TimesAndTrades) ClearOldVAP() {
	for i := 0; i < len(tt.VAP); i++ {
		now := time.Now().Add(time.Duration(-120)*time.Minute).UnixNano() / int64(time.Millisecond)
		if tt.VAP[i].Time < now {
			tt.removeVAP(i)
			i--
		}
	}
}

func (tt *TimesAndTrades) SortVAP() {
	sort.SliceStable(tt.VAP, func(i, j int) bool {
		return tt.VAP[i].Price > tt.VAP[j].Price
	})
}

func (tt *TimesAndTrades) GetBuyersVAP() {
	sort.SliceStable(tt.VAP, func(i, j int) bool {
		return tt.VAP[i].Buyers > tt.VAP[j].Buyers
	})
	tt.VAPBuyers = []BuyersAndSellers{}
	for i := 0; i < len(tt.VAP); i++ {
		tt.VAPBuyers = append(tt.VAPBuyers, tt.VAP[i])
	}

}

func (tt *TimesAndTrades) GetSellersVAP() {
	sort.SliceStable(tt.VAP, func(i, j int) bool {
		return tt.VAP[i].Sellers > tt.VAP[j].Sellers
	})
	tt.VAPSellers = []BuyersAndSellers{}
	for i := 0; i < len(tt.VAP); i++ {
		tt.VAPSellers = append(tt.VAPSellers, tt.VAP[i])
	}
}

func (tt *TimesAndTrades) removeVAP(index int) {
	tt.VAP = append(tt.VAP[:index], tt.VAP[index+1:]...)
}
