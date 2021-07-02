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
