package service

import (
	"encoding/json"
	"ps-chartdata/bitquery"
	"ps-chartdata/model"
	"strings"
)

func GetWalletBalancesByAddress(address string) ([]model.Balance, error) {
	query := `{
		ethereum(network: bsc) {
			address(address: {is: "WALLET_ADDRESS"}) {
				balances {
					currency {
						symbol
						address
					}
					value
				}
			}
		}
	}	
	`
	query = strings.ReplaceAll(
		query,
		WALLET_ADDRESS,
		address,
	)

	resp, err := bitquery.Query(query, nil)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]map[string][]map[string][]model.Balance)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	balances, _ := data["data"]["ethereum"]["address"][0]["balances"]

	return balances, nil
}
