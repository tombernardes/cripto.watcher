package crosscutting

import (
	"time"

	"github.com/tombernardes/cripto.watcher/configuration"
	"github.com/tombernardes/cripto.watcher/domain"
	"github.com/tombernardes/cripto.watcher/repository"
)

var CriptoWallet *domain.Wallet
var Charts []domain.Chart
var pointChan chan *domain.Point

func init() {
	pointChan = make(chan *domain.Point)
	go receiveChartPoint()
	go getCharts()
	go updateWallet()
}

func updateWallet() {
	wRepo := repository.WalletRepository{}
	for {
		CriptoWallet = wRepo.GetWallet()
		time.Sleep(1 * time.Millisecond)
	}
}

func getCharts() {
	cRepo := repository.ChartRepository{}
	tRepo := repository.TradeRepository{}
	for _, v := range configuration.Config.Tickers {
		chart := cRepo.GetChartHistory(v)
		Charts = append(Charts, *chart)
		go cRepo.StreamChart(v, pointChan)
		go tRepo.StreamTrades(v, pointChan)
	}

	for {
		time.Sleep(1 * time.Millisecond)
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
					Charts[i].Points[len(Charts[i].Points)-1].Volume += point.Volume
					Charts[i].Points[len(Charts[i].Points)-1].Trades += point.Trades
				}
				/*
					raw, _ := json.Marshal(Charts[i].Points[len(Charts[i].Points)-1])
					println(string(raw))*/
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func removePoint(slice []domain.Point, index int) []domain.Point {
	return append(slice[:index], slice[index+1:]...)
}
