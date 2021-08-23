package manager

import (
	"time"

	"github.com/tombernardes/cripto.watcher/configuration"
	"github.com/tombernardes/cripto.watcher/crosscutting"
	"github.com/tombernardes/cripto.watcher/domain"
)

var books []domain.Book
var trades []domain.TimesAndTrades
var book domain.Book
var price float64
var openedOrders int64
var ManagedOrderHistory []ManagedOrder

type ManagedOrder struct {
	QtyBuyed      float64
	PriceBuyed    float64
	VolumeBuyed   float64
	QtySelled     float64
	PriceSelled   float64
	VolumeSelled  float64
	Result        float64
	ResultPercent float64
}

func init() {
	openedOrders = 0
	ManagedOrderHistory = []ManagedOrder{}

	go manageOrders()
}

func manageOrders() {
	money := float64(100)
	books = crosscutting.Books
	trades = crosscutting.Trades
	for {
		if price != 0 && len(books) > 0 {
			book = books[0]
			if len(book.BigPlayersBid) > 0 {
				if openedOrders < configuration.Config.OpenOrders && book.BigPlayersBid[0].Index < 5 {
					managedOrder := ManagedOrder{}
					managedOrder.QtyBuyed = money / price
					managedOrder.PriceBuyed = book.BigPlayersBid[0].Price
					managedOrder.VolumeBuyed = managedOrder.QtyBuyed * managedOrder.PriceBuyed
					ManagedOrderHistory = append(ManagedOrderHistory, managedOrder)
					for {
						book = books[0]
						price = trades[0].Trades[0].Price
						if len(book.BigPlayersBid) > 0 {
							if ManagedOrderHistory[len(ManagedOrderHistory)-1].PriceBuyed != book.BigPlayersBid[0].Price && book.BigPlayersBid[0].Index < 5 && ManagedOrderHistory[len(ManagedOrderHistory)-1].PriceBuyed < price && openedOrders < configuration.Config.OpenOrders {
								ManagedOrderHistory[len(ManagedOrderHistory)-1].PriceBuyed = book.BigPlayersBid[0].Price
								ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeBuyed = ManagedOrderHistory[len(ManagedOrderHistory)-1].PriceBuyed * ManagedOrderHistory[len(ManagedOrderHistory)-1].QtyBuyed
								openedOrders++
							} else if (openedOrders >= configuration.Config.OpenOrders &&
								(price*managedOrder.QtyBuyed) >= ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeBuyed+(ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeBuyed*0.01)) ||
								(openedOrders >= configuration.Config.OpenOrders &&
									(price*managedOrder.QtyBuyed) < ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeBuyed+(ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeBuyed*0.005)) {
								ManagedOrderHistory[len(ManagedOrderHistory)-1].QtySelled = managedOrder.QtyBuyed
								ManagedOrderHistory[len(ManagedOrderHistory)-1].PriceSelled = price
								ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeSelled = ManagedOrderHistory[len(ManagedOrderHistory)-1].QtySelled * ManagedOrderHistory[len(ManagedOrderHistory)-1].PriceSelled
								ManagedOrderHistory[len(ManagedOrderHistory)-1].Result = ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeSelled - ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeBuyed
								money += ManagedOrderHistory[len(ManagedOrderHistory)-1].Result
								ManagedOrderHistory[len(ManagedOrderHistory)-1].ResultPercent = (ManagedOrderHistory[len(ManagedOrderHistory)-1].Result / ManagedOrderHistory[len(ManagedOrderHistory)-1].VolumeSelled) * 100
								openedOrders--
								managedOrder = ManagedOrder{}
								break
							}
						}
						time.Sleep(1 * time.Millisecond)
					}
				}
			}
		}
		books = crosscutting.Books
		trades = crosscutting.Trades
		if len(trades) > 0 {
			price = trades[0].Trades[0].Price
		}
		time.Sleep(1 * time.Millisecond)
	}
}
