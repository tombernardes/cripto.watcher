package domain

type Order struct {
	Index    int64
	Quantity float64
	Price    float64
	Volume   float64
}

type Book struct {
	LastUpdateId  int64
	FirstUpdateId int64
	Ticker        string
	LastUpdate    int64
	BidMeanOrder  Order
	BidSize       float64
	BestBidPrice  float64
	AskMeanOrder  Order
	AskSize       float64
	BestAskPrice  float64
	Spread        float64
	Bid           []Order
	Ask           []Order
	BigPlayersBid []Order
	BigPlayersAsk []Order
}
