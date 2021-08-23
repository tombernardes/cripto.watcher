package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	c "github.com/tombernardes/cripto.watcher/configuration"
	"github.com/tombernardes/cripto.watcher/domain"
)

type ChartRepository struct {
}

func (cr *ChartRepository) getBinanceChartHistory(ticker string) []*binance.Kline {
	c.Config.Client.NewSetServerTimeService().Do(context.Background())
	res, err := c.Config.Client.NewKlinesService().Symbol(ticker).
		Interval(c.Config.ChartPeriod).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *ChartRepository) GetChartHistory(ticker string) *domain.Chart {
	bChartHistory := cr.getBinanceChartHistory(ticker)
	if bChartHistory == nil {
		fmt.Println("Error on getting chart history. Binance chart history is nil")
		return nil
	}
	history := domain.Chart{}
	history.Ticker = ticker
	for _, v := range bChartHistory {
		open := float64(0)
		close := float64(0)
		min := float64(0)
		max := float64(0)
		volume := float64(0)
		trades := v.TradeNum
		if val, err := strconv.ParseFloat(v.Open, 64); err == nil {
			open = val
		}
		if val, err := strconv.ParseFloat(v.Close, 64); err == nil {
			close = val
		}
		if val, err := strconv.ParseFloat(v.Low, 64); err == nil {
			min = val
		}
		if val, err := strconv.ParseFloat(v.High, 64); err == nil {
			max = val
		}
		if val, err := strconv.ParseFloat(v.Volume, 64); err == nil {
			volume = val
		}
		point := domain.Point{
			Ticker:     ticker,
			OpenTime:   v.OpenTime,
			CloseTime:  v.CloseTime,
			ClosePrice: close,
			OpenPrice:  open,
			MinPrice:   min,
			MaxPrice:   max,
			Volume:     volume,
			Trades:     trades,
		}
		history.Points = append(history.Points, point)
	}
	return &history
}

func (cr *ChartRepository) StreamChart(ticker string, lastPoint chan *domain.Point) {
	wsKlineHandler := func(v *binance.WsKlineEvent) {
		open := float64(0)
		close := float64(0)
		min := float64(0)
		max := float64(0)
		volume := float64(0)
		trades := v.Kline.TradeNum
		if val, err := strconv.ParseFloat(v.Kline.Open, 64); err == nil {
			open = val
		}
		if val, err := strconv.ParseFloat(v.Kline.Close, 64); err == nil {
			close = val
		}
		if val, err := strconv.ParseFloat(v.Kline.Low, 64); err == nil {
			min = val
		}
		if val, err := strconv.ParseFloat(v.Kline.High, 64); err == nil {
			max = val
		}
		if val, err := strconv.ParseFloat(v.Kline.Volume, 64); err == nil {
			volume = val
		}
		point := domain.Point{
			Ticker:     ticker,
			OpenTime:   v.Kline.StartTime,
			CloseTime:  v.Kline.EndTime,
			ClosePrice: close,
			OpenPrice:  open,
			MinPrice:   min,
			MaxPrice:   max,
			Volume:     volume,
			Trades:     trades,
		}
		lastPoint <- &point
	}
	errHandler := func(err error) {
		cr.StreamChart(ticker, lastPoint)
	}
	doneC, _, err := binance.WsKlineServe(ticker, c.Config.ChartPeriod, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC
}
