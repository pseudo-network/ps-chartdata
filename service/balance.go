package service

import (
	"encoding/json"
	"ps-chartdata/bitquery"
	"ps-chartdata/model"
	"strings"
)

func GetBalancesByAddress(address string, chain model.Chain) ([]model.Balance, error) {
	query := `
	query ($address: String!) {
		ethereum(network: bsc) {
			address(address: {is: $address}) {
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
		CHAIN_NAME,
		chain.Name,
	)

	vars := make(map[string]interface{})
	vars["address"] = address

	resp, err := bitquery.Query(query, &vars)
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
