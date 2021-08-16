package model

type Balance struct {
	Currency Currency `json:"currency"`
	Value    float64  `json:"value"`
}

type Currency struct {
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}
