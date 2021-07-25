package repository

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/tombernardes/cripto.watcher/domain"
)

type TradeRepository struct {
}

func (t *TradeRepository) StreamTrades(ticker string, lastPoint chan *domain.Point) {
	wsAggTradeHandler := func(v *binance.WsAggTradeEvent) {
		close := float64(0)
		quantity := float64(0)
		if val, err := strconv.ParseFloat(v.Price, 64); err == nil {
			close = val
		}
		if val, err := strconv.ParseFloat(v.Quantity, 64); err == nil {
			quantity = val
		}
		point := domain.Point{
			Ticker:     ticker,
			ClosePrice: close,
			Volume:     close * quantity,
			Trades:     v.LastBreakdownTradeID - v.FirstBreakdownTradeID,
		}
		lastPoint <- &point
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := binance.WsAggTradeServe(ticker, wsAggTradeHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
