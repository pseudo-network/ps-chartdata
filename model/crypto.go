package model

type Crypto struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Symbol    string `json:"symbol"`
	TokenType string `json:"token_type"`
	Network   string `json:"network"`
}

func NewCrypto(name, address, symbol, tokenType, network string) *Crypto {
	return &Crypto{
		Name:      name,
		Address:   address,
		Symbol:    symbol,
		TokenType: tokenType,
		Network:   network,
	}
}

type CryptoInfo struct {
	CurrentPrice       float64 `json:"current_price"`
	CurrentPriceUSD    float64 `json:"current_price_usd"`
	VolumeUSD          float64 `json:"volume_usd"`
	Minted             float64 `json:"minted"`
	Burned             float64 `json:"burned"`
	UniqueWalletsCount int     `json:"uniqueWallets"`
	TradeCount         float64 `json:"tradeCount"`
	MaxPrice           float64 `json:"maxPrice"`
	MinPrice           float64 `json:"minPrice"`
	OpenPrice          float64 `json:"openPrice"`
	ClosePrice         float64 `json:"closePrice"`
}

func NewCryptoInfo(currentPrice,
	volume,
	minted,
	burned,
	tradeCount,
	maxPrice,
	minPrice,
	openPrice,
	closePrice float64,
	uniqueWalletsCount int,
	usdMultiplier *float64) *CryptoInfo {

	info := &CryptoInfo{
		CurrentPrice:       currentPrice,
		VolumeUSD:          volume,
		Minted:             minted,
		Burned:             burned,
		TradeCount:         tradeCount,
		MaxPrice:           maxPrice,
		MinPrice:           minPrice,
		UniqueWalletsCount: uniqueWalletsCount,
	}

	if usdMultiplier != nil {
		info.CurrentPriceUSD = info.CurrentPrice * *usdMultiplier
		info.VolumeUSD = info.VolumeUSD / *usdMultiplier
		info.MaxPrice = info.MaxPrice * *usdMultiplier
		info.MinPrice = info.MinPrice * *usdMultiplier
		info.OpenPrice = info.OpenPrice * *usdMultiplier
		info.ClosePrice = info.ClosePrice * *usdMultiplier
	}

	return info
}

type Bar struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	Time   int64   `json:"time"`
}
