package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
	c "github.com/tombernardes/cripto.watcher/configuration"
	w "github.com/tombernardes/cripto.watcher/domain"
)

type WalletRepository struct {
}

func (wr *WalletRepository) getBinanceWallet() *binance.Account {
	c.Config.Client.NewSetServerTimeService().Do(context.Background())
	res, err := c.Config.Client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (wr *WalletRepository) GetWallet() *w.Wallet {
	bWallet := wr.getBinanceWallet()
	if bWallet == nil {
		return wr.GetWallet()
	}
	wallet := w.Wallet{}
	wallet.UpdateTime = bWallet.UpdateTime
	wallet.BuyerCommission = float64(bWallet.TakerCommission) / 100
	wallet.MakerCommission = float64(bWallet.MakerCommission) / 100
	wallet.SellerCommission = float64(bWallet.SellerCommission) / 100
	wallet.TakerCommission = float64(bWallet.TakerCommission) / 100
	wallet.CanTrade = bWallet.CanTrade
	wallet.CanDeposit = bWallet.CanDeposit
	wallet.CanWithdraw = bWallet.CanWithdraw
	for _, v := range bWallet.Balances {
		free := float64(0)
		locked := float64(0)
		if f, err := strconv.ParseFloat(v.Free, 64); err == nil {
			free = f
		}
		if l, err := strconv.ParseFloat(v.Locked, 64); err == nil {
			locked = l
		}
		balance := w.Balance{
			Asset:  v.Asset,
			Free:   free,
			Locked: locked,
		}
		if free > 0 || locked > 0 {
			wallet.Balances = append(wallet.Balances, balance)
		}
	}
	return &wallet
}
