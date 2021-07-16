package bitquery

// dex trades
type DexTrade struct {
	TimeInterval  TimeInterval `json:"timeInterval"`
	High          float64      `json:"high"`
	Low           float64      `json:"low"`
	Open          string       `json:"open"`
	Close         string       `json:"close"`
	BaseCurrency  Currency     `json:"baseCurrency"`
	QuoteCurrency Currency     `json:"quoteCurrency"`
	Date          Date         `json:"date"`
	Trades        int          `json:"trades"`
	TradeAmount   float64      `json:"tradeAmount"`
	UnixTimeMS    int64        `json:"unixTimeMS"`
}

type Date struct {
	Date string `json:"date"`
}

type Currency struct {
	Name string `json:"name"`
}

type TimeInterval struct {
	Minute string `json:"minute"`
}

type Crypto struct {
	Network Network `json:"network"`
	Subject Subject `json:"subject"`
}

type Network struct {
	Network string `json:"network"`
}

type Subject struct {
	Typename  string `json:"__typename"`
	Address   string `json:"address"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	TokenId   string `json:"tokenId"`
	TokenType string `json:"tokenType"`
}

// transactions
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
