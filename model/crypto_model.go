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
	BeginningPrice float64 `json:"beginning_price"`
	CurrentPrice   float64 `json:"current_price"`
	Volume         float64 `json:"volume"`
	PercentChange  float64 `json:"percent_change"`
}

func NewCryptoInfo(beginningPrice, currentPrice, volume float64) *CryptoInfo {

	percentChange := ((beginningPrice - currentPrice) / beginningPrice) * 100

	return &CryptoInfo{
		BeginningPrice: beginningPrice,
		CurrentPrice:   currentPrice,
		Volume:         volume,
		PercentChange:  percentChange,
	}
}
