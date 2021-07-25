package domain

type Wallet struct {
	MakerCommission  float64   `json:"makerCommission"`
	TakerCommission  float64   `json:"takerCommission"`
	BuyerCommission  float64   `json:"buyerCommission"`
	SellerCommission float64   `json:"sellerCommission"`
	CanTrade         bool      `json:"canTrade"`
	CanWithdraw      bool      `json:"canWithdraw"`
	CanDeposit       bool      `json:"canDeposit"`
	UpdateTime       uint64    `json:"updateTime"`
	Balances         []Balance `json:"balances"`
}

// Balance define user balance of your account
type Balance struct {
	Asset  string  `json:"asset"`
	Free   float64 `json:"free"`
	Locked float64 `json:"locked"`
}
