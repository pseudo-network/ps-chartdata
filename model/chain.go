package model

import (
	"errors"
)

// should come from db
var Chains = map[string]Chain{
	"ethereum": {
		Id:           1,
		Name:         "Ethereum",
		BitqueryName: "ethereum",
		NativeToken: Token{
			Name:    "Wrapped Ether",
			Symbol:  "WETH",
			Address: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		},
		USDToken: Token{
			Name:    "USD Coin",
			Symbol:  "USDC",
			Address: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		},
	},
	"bsc": {
		Id:           2,
		Name:         "Binance Smart Chain",
		BitqueryName: "bsc",
		NativeToken: Token{
			Name:    "Wrapped BNB",
			Symbol:  "WBNB",
			Address: "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c",
		},
		USDToken: Token{
			Name:    "Binance USD",
			Symbol:  "BUSD",
			Address: "0xe9e7cea3dedca5984780bafc599bd69add087d56",
		},
	},
	"cardano": {
		Id:           3,
		Name:         "Cardano",
		BitqueryName: "cardano",
		NativeToken: Token{
			Name:    "Cardano",
			Symbol:  "ADA",
			Address: "...",
		},
		// USDToken: Token{
		// 	Name:    "Binane USD",
		// 	Symbol:  "BUSD",
		// 	Address: "0xe9e7cea3dedca5984780bafc599bd69add087d56",
		// },
	},
}

type Chain struct {
	Id           int    `json:"id,omitempty" bson:"id"` // todo: should be primitive.ObjectId
	Name         string `json:"name" bson:"name"`
	BitqueryName string `json:"bitquery_name" bson:"bitquery_name"`
	NativeToken  Token  `json:"native_token" bson:"native_token"`
	USDToken     Token  `json:"usd_token" bson:"usd_token"`
	// todo...
}

// placeholder
func GetChainByID(chainID int) (*Chain, error) {
	for k, chain := range Chains {
		if Chains[k].Id == chainID {
			return &chain, nil
		}
	}

	return nil, errors.New("chain not found")
}

// placeholder
func GetChainByName(chainName string) (*Chain, error) {
	for k, chain := range Chains {
		if Chains[k].Name == chainName {
			return &chain, nil
		}
	}

	return nil, errors.New("chain not found")
}
