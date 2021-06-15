package model

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
