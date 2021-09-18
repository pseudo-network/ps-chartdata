package bitquery

// dex trades
type DexTrade struct {
	TimeInterval TimeInterval `json:"timeInterval"`
	High         float64      `json:"high"`
	Low          float64      `json:"low"`
	Open         string       `json:"open"`
	Close        string       `json:"close"`
	Median       float64      `json:"median"`
	Date         Date         `json:"date"`
	Trades       int          `json:"trades"`
	Volume       float64      `json:"volume"`
	TradeAmount  float64      `json:"tradeAmount"`
	UnixTimeMS   int64        `json:"unixTimeMS"`
}

type Date struct {
	Date string `json:"date"`
}

type Currency struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}

type TimeInterval struct {
	Day    string `json:"day"`
	Second string `json:"second"`
	Minute string `json:"minute"`
}

type Token struct {
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
	Transaction  Hash         `json:"transaction"`
	TimeInterval TimeInterval `json:"timeInterval"`
	TradeAmount  float64      `json:"tradeAmount"`
	BuyCurrency  BuyCurrency  `json:"buyCurrency"`
	BuyAmount    float64      `json:"buyAmount"`
	SellCurrency SellCurrency `json:"sellCurrency"`
	SellAmount   float64      `json:"sellAmount"`
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

type Block struct {
	Height int `json:"height"`
}

type Info struct {
	Block       Block   `json:"block"`
	TradeAmount float64 `json:"tradeAmount"`
	QuotePrice  float64 `json:"quotePrice"`
}

type Summary struct {
	DaySummaries []DaySummary `json:"daySummaries"`
	// OverViews    []Overview   `json:"overviews"`
}

type DaySummary struct {
	TimeInterval   TimeInterval `json:"timeInterval"`
	BaseCurrency   Currency     `json:"baseCurrency"`
	QuoteCurrency  Currency     `json:"quoteCurrency"`
	QuotePrice     float64      `json:"quotePrice"`
	QuoteAmount    float64      `json:"quoteAmount"`
	TradeCount     float64      `json:"tradeCount"`
	TradeAmountUSD float64      `json:"tradeAmountUSD"`
	TradeVolume    float64      `json:"tradeVolume"`
	MaxPrice       float64      `json:"maxPrice"`
	MinPrice       float64      `json:"minPrice"`
	OpenPrice      string       `json:"openPrice"`
	ClosePrice     string       `json:"closePrice"`
}

type Overview struct {
	Minted             float64 `json:"minted"`
	Burned             float64 `json:"burned"`
	UniqueWalletsCount int     `json:"uniqueWallets"`
}
