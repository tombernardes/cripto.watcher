package crosscutting

import (
	"sort"
	"sync"
	"time"

	"github.com/tombernardes/cripto.watcher/configuration"
	"github.com/tombernardes/cripto.watcher/domain"
	"github.com/tombernardes/cripto.watcher/helpers"
	"github.com/tombernardes/cripto.watcher/repository"
)

var CriptoWallet *domain.Wallet
var Charts []domain.Chart
var Books []domain.Book
var Trades []domain.TimesAndTrades
var pointChan chan *domain.Point
var bookChan chan *domain.Book
var tradeChan chan *domain.Trade

func init() {
	pointChan = make(chan *domain.Point)
	bookChan = make(chan *domain.Book)
	tradeChan = make(chan *domain.Trade)
	updateInBackground()

	//go updateWallet()
}

func updateWallet() {
	wRepo := repository.WalletRepository{}
	for {
		CriptoWallet = wRepo.GetWallet()
		time.Sleep(1 * time.Millisecond)
	}
}

func updateInBackground() {
	cRepo := repository.ChartRepository{}
	tRepo := repository.TradeRepository{}
	bRepo := repository.BookRepository{}
	go receiveChartPoint()
	go receiveBook()
	go receiveTrade()
	for _, v := range configuration.Config.Tickers {
		chart := cRepo.GetChartHistory(v)
		Charts = append(Charts, *chart)
		go cRepo.StreamChart(v, pointChan)
		go tRepo.StreamTrades(v, tradeChan)
		go bRepo.StreamBook(v, bookChan)
	}

}

func receiveChartPoint() {
	for {
		select {
		case point := <-pointChan:
			for i := range Charts {
				if Charts[i].Ticker == point.Ticker && Charts[i].Points[len(Charts[i].Points)-1].OpenTime == point.OpenTime {
					Charts[i].Points = removePoint(Charts[i].Points, len(Charts[i].Points)-1)
					Charts[i].Points = append(Charts[i].Points, *point)
				} else if Charts[i].Ticker == point.Ticker && point.OpenTime != 0 {
					Charts[i].Points = append(Charts[i].Points, *point)
					Charts[i].Points = removePoint(Charts[i].Points, 0)
				} else if point.OpenTime == 0 {
					Charts[i].Points[len(Charts[i].Points)-1].ClosePrice = point.ClosePrice
					if point.ClosePrice < Charts[i].Points[len(Charts[i].Points)-1].MinPrice {
						Charts[i].Points[len(Charts[i].Points)-1].MinPrice = point.ClosePrice
					}
					if point.ClosePrice > Charts[i].Points[len(Charts[i].Points)-1].MaxPrice {
						Charts[i].Points[len(Charts[i].Points)-1].MaxPrice = point.ClosePrice
					}
				}
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func receiveBook() {
	var wg sync.WaitGroup
	for {
		select {
		case book := <-bookChan:
			has, i := helpers.ContainsAndIndexBook(Books, book.Ticker)
			if has {
				if book.LastUpdateId >= Books[i].LastUpdateId+1 && book.FirstUpdateId <= Books[i].LastUpdateId+1 {
					wg.Add(2)
					go func() {
						defer wg.Done()
						for j := range book.Bid {
							hasPriceLevel, level := containsAndIndexPriceLevel(Books[i].Bid, book.Bid[j].Price)
							if hasPriceLevel {
								if book.Bid[j].Quantity > 0 {
									Books[i].Bid[level].Quantity = book.Bid[j].Quantity
									Books[i].Bid[level].Volume = book.Bid[j].Quantity * Books[i].Bid[level].Price
								} else {
									Books[i].Bid = removePriceLevelFromBook(Books[i].Bid, level)
								}
							} else {
								if book.Bid[j].Quantity > 0 {
									Books[i].Bid = append(Books[i].Bid, book.Bid[j])
								}
							}
						}
						for j := 0; j < len(Books[i].Bid); j++ {
							if Books[i].Bid[j].Quantity == 0 {
								Books[i].Bid = removePriceLevelFromBook(Books[i].Bid, j)
								j--
							}
						}
					}()
					go func() {
						defer wg.Done()
						for j := range book.Ask {
							hasPriceLevel, level := containsAndIndexPriceLevel(Books[i].Ask, book.Ask[j].Price)
							if hasPriceLevel {
								if book.Ask[j].Quantity > 0 {
									Books[i].Ask[level].Quantity = book.Ask[j].Quantity
									Books[i].Ask[level].Volume = book.Ask[j].Quantity * Books[i].Ask[level].Price
								} else {
									Books[i].Ask = removePriceLevelFromBook(Books[i].Ask, level)
								}
							} else {
								if book.Ask[j].Quantity > 0 {
									Books[i].Ask = append(Books[i].Ask, book.Ask[j])
								}
							}
						}
						for j := 0; j < len(Books[i].Ask); j++ {
							if Books[i].Ask[j].Quantity == 0 {
								Books[i].Ask = removePriceLevelFromBook(Books[i].Ask, j)
								j--
							}
						}
					}()
					wg.Wait()
					Books[i].AskSize = getBookSize(Books[i].Ask)
					if len(Books[i].Ask) > 0 {
						Books[i].BestAskPrice = Books[i].Ask[0].Price
					}
					Books[i].BidSize = getBookSize(Books[i].Bid)
					if len(Books[i].Bid) > 0 {
						Books[i].BestBidPrice = Books[i].Bid[0].Price
					}
					Books[i].Spread = Books[i].BestAskPrice - Books[i].BestBidPrice
					Books[i].LastUpdateId = book.LastUpdateId
					Books[i].BidMeanOrder = getMeanOrder(Books[i].Bid)
					Books[i].AskMeanOrder = getMeanOrder(Books[i].Ask)
					Books[i].Ask = sortAskBook(Books[i].Ask, true)
					Books[i].Bid = sortBidBook(Books[i].Bid, true)
					Books[i].BigPlayersAsk = getBigPlayers(Books[i].Ask)
					Books[i].BigPlayersBid = getBigPlayers(Books[i].Bid)
					Books[i].BigPlayersAsk = sortAskBook(Books[i].BigPlayersAsk, false)
					Books[i].BigPlayersBid = sortBidBook(Books[i].BigPlayersBid, false)
				}
			} else {
				Books = append(Books, *book)
			}

		}
		time.Sleep(1 * time.Millisecond)
	}
}

func receiveTrade() {
	for {
		select {
		case trade := <-tradeChan:
			has, i := helpers.ContainsAndIndexTrade(Trades, trade.Ticker)
			if has {

				Trades[i].Trades = append([]domain.Trade{*trade}, Trades[i].Trades...)
				if len(Trades[i].Trades) > 5000 {
					Trades[i].Trades = Trades[i].Trades[:len(Trades[i].Trades)-1]
				}
				hasPrice, p := containsAndIndexVolumeAtPriceLevel(Trades[i].VAP, trade.Price)
				if hasPrice {
					if trade.IsMaker {
						bAnds := Trades[i].VAP[p]
						bAnds.Time = trade.Time
						bAnds.Sellers += trade.Quantity
						Trades[i].VAP[p] = bAnds
						Trades[i].VAP[p].Trades++
						Trades[i].VAP[p].TradesSellers++
					} else {
						bAnds := Trades[i].VAP[p]
						bAnds.Time = trade.Time
						bAnds.Buyers += trade.Quantity
						Trades[i].VAP[p] = bAnds
						Trades[i].VAP[p].Trades++
						Trades[i].VAP[p].TradesBuyers++
					}
				} else {
					if trade.IsMaker {
						bAnds := domain.BuyersAndSellers{}
						bAnds.Time = trade.Time
						bAnds.Price = trade.Price
						bAnds.Sellers = trade.Quantity
						bAnds.Trades = 1
						bAnds.TradesSellers = 1
						Trades[i].VAP = append(Trades[i].VAP, bAnds)
					} else {
						bAnds := domain.BuyersAndSellers{}
						bAnds.Time = trade.Time
						bAnds.Price = trade.Price
						bAnds.Buyers = trade.Quantity
						bAnds.Trades = 1
						bAnds.TradesBuyers = 1
						Trades[i].VAP = append(Trades[i].VAP, bAnds)
					}
				}
				Trades[i].ClearOldVAP()
				Trades[i].SortVAP()
			} else {
				Trades = append(Trades, domain.TimesAndTrades{Ticker: trade.Ticker, Trades: []domain.Trade{*trade}, VAP: []domain.BuyersAndSellers{}})
				hasPrice, p := containsAndIndexVolumeAtPriceLevel(Trades[i].VAP, trade.Price)
				if hasPrice {
					if trade.IsMaker {
						bAnds := Trades[i].VAP[p]
						bAnds.Time = trade.Time
						bAnds.Sellers += trade.Quantity
						Trades[i].VAP[p] = bAnds
						Trades[i].VAP[p].Trades++
						Trades[i].VAP[p].TradesSellers++
					} else {
						bAnds := Trades[i].VAP[p]
						bAnds.Time = trade.Time
						bAnds.Buyers += trade.Quantity
						Trades[i].VAP[p] = bAnds
						Trades[i].VAP[p].Trades++
						Trades[i].VAP[p].TradesBuyers++
					}
					Trades[i].ClearOldVAP()
					Trades[i].SortVAP()
				} else {
					if trade.IsMaker {
						bAnds := domain.BuyersAndSellers{}
						bAnds.Time = trade.Time
						bAnds.Price = trade.Price
						bAnds.Sellers = trade.Quantity
						bAnds.Trades = 1
						bAnds.TradesSellers = 1
						Trades[i].VAP = append(Trades[i].VAP, bAnds)
					} else {
						bAnds := domain.BuyersAndSellers{}
						bAnds.Time = trade.Time
						bAnds.Price = trade.Price
						bAnds.Buyers = trade.Quantity
						bAnds.Trades = 1
						bAnds.TradesBuyers = 1
						Trades[i].VAP = append(Trades[i].VAP, bAnds)
					}
				}
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func removePoint(slice []domain.Point, index int) []domain.Point {
	return append(slice[:index], slice[index+1:]...)
}

func sortBidBook(book []domain.Order, index bool) []domain.Order {
	sort.SliceStable(book, func(i, j int) bool {
		return book[i].Price > book[j].Price
	})
	if index {
		for i := range book {
			book[i].Index = int64(i) + 1
		}
	}
	return book
}

func sortAskBook(book []domain.Order, index bool) []domain.Order {
	sort.SliceStable(book, func(i, j int) bool {
		return book[i].Price < book[j].Price
	})
	if index {
		for i := range book {
			book[i].Index = int64(i) + 1
		}
	}
	return book
}

func getBookSize(book []domain.Order) float64 {
	size := float64(0)
	for i := range book {
		size += book[i].Quantity
	}
	return size
}

func removePriceLevelFromBook(slice []domain.Order, index int) []domain.Order {
	return append(slice[:index], slice[index+1:]...)
}

func containsAndIndexPriceLevel(s []domain.Order, priceLevel float64) (bool, int) {
	for i, a := range s {
		if a.Price == priceLevel {
			return true, i
		}
	}
	return false, 0
}

func containsAndIndexVolumeAtPriceLevel(s []domain.BuyersAndSellers, priceLevel float64) (bool, int) {
	for i, a := range s {
		if a.Price == priceLevel {
			return true, i
		}
	}
	return false, 0
}

func getMeanOrder(s []domain.Order) domain.Order {
	meanOrder := domain.Order{}
	book := s
	sort.SliceStable(book, func(i, j int) bool {
		return book[i].Quantity < book[j].Quantity
	})

	qty := float64(0)
	volume := float64(0)
	for _, o := range book {
		qty += o.Quantity
		volume += o.Volume
	}
	meanOrder.Price = volume / qty
	meanOrder.Quantity = volume / meanOrder.Price
	meanOrder.Volume = meanOrder.Price * meanOrder.Quantity
	return meanOrder
}

func getBigPlayers(s []domain.Order) []domain.Order {
	bigPlayers := []domain.Order{}
	for _, o := range s {
		if o.Quantity >= 1 {
			bigPlayers = append(bigPlayers, o)
		}
	}
	return bigPlayers
}
