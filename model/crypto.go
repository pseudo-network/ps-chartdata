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
	BeginningPrice    float64 `json:"beginning_price"`
	BeginningPriceUSD float64 `json:"beginning_price_usd"`
	CurrentPrice      float64 `json:"current_price"`
	CurrentPriceUSD   float64 `json:"current_price_usd"`
	VolumeUSD         float64 `json:"volume_usd"`
	PercentChange     float64 `json:"percent_change"`
}

func NewCryptoInfo(beginningPrice, currentPrice, volume float64, usdMultiplier *float64) *CryptoInfo {
	percentChange := ((beginningPrice - currentPrice) / beginningPrice) * 100

	info := &CryptoInfo{
		BeginningPrice: beginningPrice,
		CurrentPrice:   currentPrice,
		VolumeUSD:      volume,
		PercentChange:  percentChange,
	}

	if usdMultiplier != nil {
		info.BeginningPriceUSD = info.BeginningPrice * *usdMultiplier
		info.CurrentPriceUSD = info.CurrentPrice * *usdMultiplier
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
