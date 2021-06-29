package model

type Transaction struct {
	Transaction  Hash               `json:"transaction"`
	TimeInterval TimeIntervalSecond `json:"timeInterval"`
	TradeAmount  float64            `json:"tradeAmount"`
	BuyCurrency  BuyCurrency        `json:"buyCurrency"`
	BuyAmount    float64            `json:"buyAmount"`
	SellCurrency SellCurrency       `json:"sellCurrency"`
	SellAmount   float64            `json:"sellAmount"`
}

type BuyCurrency struct {
	Symbol string `json:"symbol"`
}

type SellCurrency struct {
	Symbol string `json:"symbol"`
}

type Hash struct {
	Hash string `json:"hash"`
}

type TimeIntervalSecond struct {
	Second string `json:"second"`
}
