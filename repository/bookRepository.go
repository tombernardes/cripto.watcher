package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	c "github.com/tombernardes/cripto.watcher/configuration"
	"github.com/tombernardes/cripto.watcher/domain"
)

type BookRepository struct {
}

func (b *BookRepository) GetBookSnapshot(ticker string) *domain.Book {
	c.Config.Client.NewSetServerTimeService().Do(context.Background())
	res, err := c.Config.Client.NewDepthService().Symbol(ticker).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	book := domain.Book{}
	book.BestAskPrice = 0
	book.BestBidPrice = 0
	book.Ticker = ticker
	volume := float64(0)
	quantity := float64(0)
	price := float64(0)
	for i, b := range res.Bids {
		if val, err := strconv.ParseFloat(b.Price, 32); err == nil {
			price = val
		}
		if val, err := strconv.ParseFloat(b.Quantity, 32); err == nil {
			quantity = val
			book.BidSize += quantity
		}
		if quantity >= 0 {
			volume = price * quantity
			book.Bid = append(book.Bid, domain.Order{Quantity: quantity, Price: price, Volume: volume, Index: int64(i) + 1})
		}
	}
	volume = 0
	quantity = 0
	price = 0
	for i, a := range res.Asks {
		if val, err := strconv.ParseFloat(a.Price, 32); err == nil {
			price = val
		}
		if val, err := strconv.ParseFloat(a.Quantity, 32); err == nil {
			quantity = val
			book.AskSize += quantity
		}
		if quantity >= 0 {
			volume = price * quantity
			book.Ask = append(book.Ask, domain.Order{Quantity: quantity, Price: price, Volume: volume, Index: int64(i) + 1})
		}
	}
	if len(book.Ask) > 0 {
		book.BestAskPrice = book.Ask[0].Price
	}
	if len(book.Bid) > 0 {
		book.BestBidPrice = book.Bid[0].Price
	}
	book.Spread = book.BestAskPrice - book.BestBidPrice
	book.LastUpdateId = res.LastUpdateID
	return &book
}

func (b *BookRepository) StreamBook(ticker string, lastBook chan *domain.Book) {
	wsDepthHandler := func(v *binance.WsDepthEvent) {
		book := domain.Book{}
		book.BestAskPrice = 0
		book.BestBidPrice = 0
		book.Ticker = ticker
		book.LastUpdate = v.Time
		quantity := float64(0)
		price := float64(0)
		for i, b := range v.Bids {
			if val, err := strconv.ParseFloat(b.Price, 32); err == nil {
				price = val
			}
			if val, err := strconv.ParseFloat(b.Quantity, 32); err == nil {
				quantity = val
				book.BidSize += quantity
			}
			if quantity >= 0 {
				book.Bid = append(book.Bid, domain.Order{Quantity: quantity, Price: price, Volume: price * quantity, Index: int64(i) + 1})
			}
		}
		quantity = 0
		price = 0
		for i, a := range v.Asks {
			if val, err := strconv.ParseFloat(a.Price, 32); err == nil {
				price = val
			}
			if val, err := strconv.ParseFloat(a.Quantity, 32); err == nil {
				quantity = val
				book.AskSize += quantity
			}
			if quantity >= 0 {
				book.Ask = append(book.Ask, domain.Order{Quantity: quantity, Price: price, Volume: price * quantity, Index: int64(i) + 1})
			}
		}
		if len(book.Ask) > 0 {
			book.BestAskPrice = book.Ask[0].Price
		}
		if len(book.Bid) > 0 {
			book.BestBidPrice = book.Bid[0].Price
		}
		book.LastUpdateId = v.LastUpdateID
		book.Spread = book.BestAskPrice - book.BestBidPrice
		lastBook <- &book
	}
	errHandler := func(err error) {
		b.StreamBook(ticker, lastBook)
	}
	doneC, _, err := binance.WsDepthServe100Ms(ticker, wsDepthHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
