package repository

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/tombernardes/cripto.watcher/domain"
)

type TradeRepository struct {
}

func (t *TradeRepository) StreamTrades(ticker string, trades chan *domain.Trade) {
	wsTradeHandler := func(v *binance.WsTradeEvent) {
		close := float64(0)
		quantity := float64(0)
		if val, err := strconv.ParseFloat(v.Price, 64); err == nil {
			close = val
		}
		if val, err := strconv.ParseFloat(v.Quantity, 64); err == nil {
			quantity = val
		}

		trade := domain.Trade{
			Ticker:   ticker,
			Time:     v.Time,
			Quantity: quantity,
			Price:    close,
			Volume:   close * quantity,
			IsMaker:  v.IsBuyerMaker,
		}
		if v.IsBuyerMaker {
			trade.WhoTake = "Seller"
		} else {
			trade.WhoTake = "Buyer"
		}
		trades <- &trade
	}
	errHandler := func(err error) {
		t.StreamTrades(ticker, trades)
	}
	doneC, _, err := binance.WsTradeServe(ticker, wsTradeHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
