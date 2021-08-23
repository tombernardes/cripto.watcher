package main

import (
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/leekchan/accounting"
	"github.com/rivo/tview"
	"github.com/tombernardes/cripto.watcher/crosscutting"
	"github.com/tombernardes/cripto.watcher/domain"
	"github.com/tombernardes/cripto.watcher/manager"
)

var candleCharts []domain.Chart
var books []domain.Book
var timesAndTrades []domain.TimesAndTrades
var app *tview.Application

var orders []manager.ManagedOrder

func init() {
	go updateInBackground()
}

func updateInBackground() {
	candleCharts = crosscutting.Charts
	books = crosscutting.Books
	timesAndTrades = crosscutting.Trades
	//orders = manager.ManagedOrderHistory
	for {
		candleCharts = crosscutting.Charts
		books = crosscutting.Books
		timesAndTrades = crosscutting.Trades
		//orders = manager.ManagedOrderHistory
		time.Sleep(1 * time.Millisecond)
	}
}

func draw() {
	for {
		app.Draw()
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	app = tview.NewApplication()

	//Header
	Header := tview.NewTable().
		SetBorders(true)
	go updateHeader(books, timesAndTrades, Header)
	Header.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			Header.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		Header.GetCell(row, column).SetTextColor(tcell.ColorRed)
		Header.SetSelectable(false, false)
	})

	//Orders
	Orders := tview.NewTable().
		SetBorders(true)
	go updateOrders(Orders)
	Orders.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			Orders.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		Orders.GetCell(row, column).SetTextColor(tcell.ColorRed)
		Orders.SetSelectable(false, false)
	})

	//Book
	Book := tview.NewTable().
		SetBorders(true)
	go updateBook(books, Book)
	Book.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			Book.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		Book.GetCell(row, column).SetTextColor(tcell.ColorRed)
		Book.SetSelectable(false, false)
	})

	//VAP
	VAP := tview.NewTable().
		SetBorders(true)
	go updateVAP(timesAndTrades, VAP)
	VAP.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			Book.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		VAP.GetCell(row, column).SetTextColor(tcell.ColorRed)
		VAP.SetSelectable(false, false)
	})
	//TimesAndTrades
	TimesAndTrades := tview.NewTable().
		SetBorders(true)
	go updateTimesAndTrades(timesAndTrades, TimesAndTrades)
	TimesAndTrades.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			TimesAndTrades.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		TimesAndTrades.GetCell(row, column).SetTextColor(tcell.ColorRed)
		TimesAndTrades.SetSelectable(false, false)
	})

	grid := tview.NewGrid().
		SetRows(5, 0).
		SetColumns(0, 90).
		SetBorders(true)

	grid.AddItem(Header, 0, 0, 1, 2, 0, 100, false).
		AddItem(Book, 1, 0, 1, 1, 0, 100, true).
		AddItem(TimesAndTrades, 1, 1, 1, 1, 0, 100, false)

	go draw()
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func updateOrders(table *tview.Table) {
	p := accounting.Accounting{Symbol: "US$ ", Precision: 4, Thousand: ".", Decimal: ","}
	q := accounting.Accounting{Symbol: "", Precision: 8, Thousand: ".", Decimal: ","}
	ordersTexts := strings.Split("Qty_Compra Preco_Compra Total_Compra Qty_Venda Preco_Venda Total_Venda Resultado_Percent Resultado", " ")
	for {
		if len(orders) > 0 {
			for i, col := range ordersTexts {
				table.SetCell(0, i,
					tview.NewTableCell(col).
						SetTextColor(tcell.ColorYellow).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}
			for i, o := range orders {
				table.SetCell(i+1, 0,
					tview.NewTableCell(q.FormatMoney(o.QtyBuyed)).
						SetMaxWidth(0).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))

				table.SetCell(i+1, 1,
					tview.NewTableCell(p.FormatMoney(o.PriceBuyed)).
						SetMaxWidth(0).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 2,
					tview.NewTableCell(p.FormatMoney(o.VolumeBuyed)).
						SetMaxWidth(0).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 3,
					tview.NewTableCell(q.FormatMoney(o.QtySelled)).
						SetMaxWidth(0).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))

				table.SetCell(i+1, 4,
					tview.NewTableCell(p.FormatMoney(o.PriceSelled)).
						SetMaxWidth(0).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 5,
					tview.NewTableCell(p.FormatMoney(o.VolumeSelled)).
						SetMaxWidth(0).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				if o.Result > 0 {
					table.SetCell(i+1, 6,
						tview.NewTableCell(q.FormatMoney(o.ResultPercent)+"%").
							SetMaxWidth(0).
							SetTextColor(tcell.ColorGreen).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
					table.SetCell(i+1, 7,
						tview.NewTableCell(p.FormatMoney(o.Result)).
							SetMaxWidth(0).
							SetTextColor(tcell.ColorGreen).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				} else {
					table.SetCell(i+1, 6,
						tview.NewTableCell(q.FormatMoney(o.ResultPercent)+"%").
							SetMaxWidth(0).
							SetTextColor(tcell.ColorRed).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
					table.SetCell(i+1, 7,
						tview.NewTableCell(p.FormatMoney(o.Result)).
							SetMaxWidth(0).
							SetTextColor(tcell.ColorRed).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				}
			}

		}
		//app.Draw()
		table.ScrollToEnd()
		orders = manager.ManagedOrderHistory
		time.Sleep(time.Duration(300) * time.Millisecond)
		table.Clear()
	}
}

func updateHeader(books []domain.Book, timesAndTrades []domain.TimesAndTrades, table *tview.Table) {
	p := accounting.Accounting{Symbol: "US$ ", Precision: 4, Thousand: ".", Decimal: ","}
	q := accounting.Accounting{Symbol: "", Precision: 8, Thousand: ".", Decimal: ","}
	lastPrice := float64(0)
	lastColor := int32(0)
	headerTexts := strings.Split("Data/Hora Cripto Tamanho_Book_Compra Melhor_Preco_Compra Spread Melhor_Preco_Venda Tamanho_Book_Venda Preco_Ultimo_Negocio", " ")
	for {
		if len(books) > 0 && len(timesAndTrades) > 0 {
			for i, col := range headerTexts {
				table.SetCell(0, i,
					tview.NewTableCell(col).
						SetTextColor(tcell.ColorYellow).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}
			when := " " + FormatDateTimeNoMili(time.Now()) + " "
			table.SetCell(1, 0,
				tview.NewTableCell(when).
					SetMaxWidth(0).
					SetTextColor(tcell.ColorWhite).
					SetSelectable(true).
					SetAlign(tview.AlignCenter))
			table.SetCell(1, 1,
				tview.NewTableCell(" "+books[0].Ticker+" ").
					SetMaxWidth(0).
					SetTextColor(tcell.ColorWhite).
					SetSelectable(true).
					SetAlign(tview.AlignCenter))

			table.SetCell(1, 2,
				tview.NewTableCell(q.FormatMoney(books[0].BidSize)).
					SetMaxWidth(5).
					SetTextColor(tcell.ColorWhite).
					SetSelectable(true).
					SetAlign(tview.AlignCenter))

			table.SetCell(1, 3,
				tview.NewTableCell(p.FormatMoney(books[0].BestBidPrice)).
					SetMaxWidth(5).
					SetTextColor(tcell.ColorGreen).
					SetSelectable(true).
					SetAlign(tview.AlignCenter))

			if books[0].BidSize > 0 && books[0].AskSize > 0 {
				table.SetCell(1, 4,
					tview.NewTableCell(" "+p.FormatMoney(books[0].Spread)+" ").
						SetMaxWidth(10).
						SetTextColor(tcell.ColorYellow).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			} else {
				table.SetCell(1, 4,
					tview.NewTableCell(" "+p.FormatMoney(0)+" ").
						SetMaxWidth(10).
						SetTextColor(tcell.ColorYellow).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}

			table.SetCell(1, 5,
				tview.NewTableCell(p.FormatMoney(books[0].BestAskPrice)).
					SetMaxWidth(5).
					SetMaxWidth(5).
					SetTextColor(tcell.ColorRed).
					SetSelectable(true).
					SetAlign(tview.AlignCenter))

			table.SetCell(1, 6,
				tview.NewTableCell(q.FormatMoney(books[0].AskSize)).
					SetMaxWidth(5).
					SetTextColor(tcell.ColorWhite).
					SetSelectable(true).
					SetAlign(tview.AlignCenter))

			if timesAndTrades[0].Trades[0].Price > lastPrice {
				table.SetCell(1, 7,
					tview.NewTableCell(p.FormatMoney(timesAndTrades[0].Trades[0].Price)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorGreen).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				lastColor = tcell.ColorGreen.Hex()
			} else if timesAndTrades[0].Trades[0].Price < lastPrice {
				table.SetCell(1, 7,
					tview.NewTableCell(p.FormatMoney(timesAndTrades[0].Trades[0].Price)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorRed).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				lastColor = tcell.ColorRed.Hex()
			} else {
				table.SetCell(1, 7,
					tview.NewTableCell(p.FormatMoney(timesAndTrades[0].Trades[0].Price)).
						SetMaxWidth(5).
						SetTextColor(tcell.NewHexColor(lastColor)).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))

			}
			lastPrice = timesAndTrades[0].Trades[0].Price
		}
		//app.Draw()
		books = crosscutting.Books
		timesAndTrades = crosscutting.Trades
		time.Sleep(time.Duration(300) * time.Millisecond)
		table.Clear()
	}
}

func updateBook(books []domain.Book, table *tview.Table) {
	id := accounting.Accounting{Symbol: "", Precision: 0, Thousand: ".", Decimal: ","}
	p := accounting.Accounting{Symbol: "US$ ", Precision: 4, Thousand: ".", Decimal: ","}
	q := accounting.Accounting{Symbol: "", Precision: 8, Thousand: ".", Decimal: ","}
	bookTexts := strings.Split("Pos Volume_Total_Negociado Quantidade_Negociada Preco_Negociado Preco_Negociado Quantidade_Negociada Volume_Total_Negociado Pos", " ")
	for {
		if len(books) > 0 {
			for i, col := range bookTexts {
				table.SetCell(0, i,
					tview.NewTableCell(col).
						SetTextColor(tcell.ColorYellow).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}
			for i, col := range books[0].BigPlayersBid {
				table.SetCell(i+1, 0,
					tview.NewTableCell(id.FormatMoneyInt(int(col.Index))).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 1,
					tview.NewTableCell(p.FormatMoney(col.Volume)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				if col.Quantity >= 1 {
					table.SetCell(i+1, 2,
						tview.NewTableCell(q.FormatMoney(col.Quantity)).
							SetMaxWidth(5).
							SetTextColor(tcell.ColorYellow).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				} else {
					table.SetCell(i+1, 2,
						tview.NewTableCell(q.FormatMoney(col.Quantity)).
							SetMaxWidth(5).
							SetTextColor(tcell.ColorDarkGrey).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				}
				table.SetCell(i+1, 3,
					tview.NewTableCell(p.FormatMoney(col.Price)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorGreen).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}
			for i, col := range books[0].BigPlayersAsk {
				table.SetCell(i+1, 4,
					tview.NewTableCell(p.FormatMoney(col.Price)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorRed).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				if col.Quantity >= 1 {
					table.SetCell(i+1, 5,
						tview.NewTableCell(q.FormatMoney(col.Quantity)).
							SetMaxWidth(5).
							SetTextColor(tcell.ColorYellow).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				} else {
					table.SetCell(i+1, 5,
						tview.NewTableCell(q.FormatMoney(col.Quantity)).
							SetTextColor(tcell.ColorDarkGrey).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				}
				table.SetCell(i+1, 6,
					tview.NewTableCell(p.FormatMoney(col.Volume)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				table.SetCell(i+1, 7,
					tview.NewTableCell(id.FormatMoneyInt(int(col.Index))).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}
		}
		//app.Draw()
		table.ScrollToBeginning()
		time.Sleep(time.Duration(300) * time.Millisecond)
		books = crosscutting.Books
		table.Clear()
	}
}

func updateTimesAndTrades(timesAndTrades []domain.TimesAndTrades, TimesAndTrades *tview.Table) {
	p := accounting.Accounting{Symbol: "US$ ", Precision: 4, Thousand: ".", Decimal: ","}
	q := accounting.Accounting{Symbol: "", Precision: 8, Thousand: ".", Decimal: ","}
	timesAndTradesTexts := strings.Split("Data_e_Hora Quantidade Preco_Negociado Volume_Total_Negociado Agr", " ")
	for {
		if len(timesAndTrades) > 0 {
			for i, col := range timesAndTradesTexts {
				TimesAndTrades.SetCell(0, i,
					tview.NewTableCell(col).
						SetTextColor(tcell.ColorYellow).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}
			for i, col := range timesAndTrades[0].Trades {
				when := " " + FormatDateTime(time.Unix(0, col.Time*1000*1000)) + " "
				TimesAndTrades.SetCell(i+1, 0,
					tview.NewTableCell(when).
						SetMaxWidth(0).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				TimesAndTrades.SetCell(i+1, 1,
					tview.NewTableCell(" "+q.FormatMoney(col.Quantity)+" ").
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				TimesAndTrades.SetCell(i+1, 2,
					tview.NewTableCell(p.FormatMoney(col.Price)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))

				TimesAndTrades.SetCell(i+1, 3,
					tview.NewTableCell(p.FormatMoney(col.Volume)).
						SetMaxWidth(5).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))

				if col.IsMaker {
					TimesAndTrades.SetCell(i+1, 4,
						tview.NewTableCell(" Seller ").
							SetTextColor(tcell.ColorRed).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				} else {
					TimesAndTrades.SetCell(i+1, 4,
						tview.NewTableCell(" Buyer ").
							SetTextColor(tcell.ColorGreen).
							SetSelectable(true).
							SetAlign(tview.AlignCenter))
				}
			}
		}
		//app.Draw()
		TimesAndTrades.ScrollToBeginning()
		time.Sleep(time.Duration(300) * time.Millisecond)
		timesAndTrades = crosscutting.Trades
		TimesAndTrades.Clear()
	}
}

func updateVAP(timesAndTrades []domain.TimesAndTrades, VAP *tview.Table) {
	p := accounting.Accounting{Symbol: "US$ ", Precision: 4, Thousand: ".", Decimal: ","}
	q := accounting.Accounting{Symbol: "", Precision: 8, Thousand: ".", Decimal: ","}
	timesAndTradesTexts := strings.Split("Preco_Negociado Buyers Sellers", " ")
	for {
		if len(timesAndTrades) > 0 {
			for i, col := range timesAndTradesTexts {
				VAP.SetCell(0, i,
					tview.NewTableCell(col).
						SetTextColor(tcell.ColorYellow).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}
			prices := getVAPPrices()
			for i := range prices {

				VAP.SetCell(i+1, 0,
					tview.NewTableCell(p.FormatMoney(timesAndTrades[0].VAP[i].Price)).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				VAP.SetCell(i+1, 1,
					tview.NewTableCell(q.FormatMoney(timesAndTrades[0].VAP[i].Buyers)).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
				VAP.SetCell(i+1, 2,
					tview.NewTableCell(q.FormatMoney(timesAndTrades[0].VAP[i].Sellers)).
						SetTextColor(tcell.ColorWhite).
						SetSelectable(true).
						SetAlign(tview.AlignCenter))
			}

		}
		//app.Draw()
		VAP.ScrollToBeginning()
		time.Sleep(time.Duration(300) * time.Millisecond)
		timesAndTrades = crosscutting.Trades
		VAP.Clear()
	}
}

func getVAPPrices() []float64 {
	ret := []float64{}
	for _, price := range timesAndTrades[0].VAP {
		ret = append(ret, price.Price)
	}
	// return sortVAP(ret)
	return ret
}

func sortVAPBuy(vap []float64) []float64 {
	sort.SliceStable(vap, func(i, j int) bool {
		return vap[i] < vap[j]
	})
	return vap
}

func sortVAPSell(vap []float64) []float64 {
	sort.SliceStable(vap, func(i, j int) bool {
		return vap[i] < vap[j]
	})
	return vap
}

func FormatDateTime(t time.Time) string {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	format := "02/01/2006 15:04:05.999"
	return t.Local().In(loc).Format(format)
}

func FormatDateTimeNoMili(t time.Time) string {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	format := "02/01/2006 15:04:05"
	return t.Local().In(loc).Format(format)
}
