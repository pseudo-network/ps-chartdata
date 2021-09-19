package model

type Token struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Symbol    string `json:"symbol"`
	TokenType string `json:"token_type"`
	Network   string `json:"network"`
}

func NewToken(name, address, symbol, tokenType, network string) *Token {
	return &Token{
		Name:      name,
		Address:   address,
		Symbol:    symbol,
		TokenType: tokenType,
		Network:   network,
	}
}

type TokenInfo struct {
	CurrentPrice       float64 `json:"current_price"`
	CurrentPriceUSD    float64 `json:"current_price_usd"`
	MintedCount        float64 `json:"minted_count"`
	BurnedCount        float64 `json:"burned_count"`
	UniqueWalletsCount int     `json:"unique_wallets_count"`
	TradeVolume        float64 `json:"trade_volume"`
	TradeCount         float64 `json:"trade_count"`
	TradeAmountUSD     float64 `json:"trade_amount_usd"`
	MaxPrice           float64 `json:"max_price"`
	MinPrice           float64 `json:"min_price"`
	OpenPrice          float64 `json:"open_price"`
	ClosePrice         float64 `json:"close_price"`
}

func NewTokenInfo(currentPrice,
	volume,
	// minted,
	// burned,
	tradeCount,
	tradeAmountUSD,
	maxPrice,
	minPrice,
	openPrice,
	closePrice float64,
	// uniqueWalletsCount int,
	usdMultiplier float64) *TokenInfo {

	info := &TokenInfo{
		CurrentPrice: currentPrice,
		TradeVolume:  volume,
		// MintedCount:        minted,
		// BurnedCount:        burned,
		TradeCount:     tradeCount,
		TradeAmountUSD: tradeAmountUSD,
		MaxPrice:       maxPrice,
		MinPrice:       minPrice,
		OpenPrice:      openPrice,
		ClosePrice:     closePrice,
		// UniqueWalletsCount: uniqueWalletsCount,
	}

	// todo: cleanup
	info.CurrentPriceUSD = info.CurrentPrice * usdMultiplier
	info.MaxPrice = info.MaxPrice * usdMultiplier
	info.MinPrice = info.MinPrice * usdMultiplier
	info.OpenPrice = info.OpenPrice * usdMultiplier
	info.ClosePrice = info.ClosePrice * usdMultiplier

	return info
}

type Bar struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	Time   int64   `json:"time"`
	Median float64 `json:"median"`
}
