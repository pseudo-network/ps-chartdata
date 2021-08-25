package service

import (
	"encoding/json"
	"ps-chartdata/bitquery"
	"ps-chartdata/model"
	"strings"
	"time"
)

func GetBNBInfo() (*model.CryptoInfo, error) {
	sinceRFC3339 := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)

	query := `{
		ethereum(network: bsc) {
			currentPrice: dexTrades(
				options: {desc: ["block.height"], limit: 1}
        exchangeName: {in: ["Pancake", "Pancake v2"]}
        baseCurrency: {is: "CURRENCY_BASE"}
				quoteCurrency: {is: "CURRENCY_QUOTE"}
				date: {after: "DATE_FROM"}
			) {
        block {
          height
        }
				quotePrice
			}
		}
	}
	`

	query = strings.ReplaceAll(
		query,
		CURRENCY_BASE,
		WBNB_ADDRESS,
	)
	query = strings.ReplaceAll(
		query,
		CURRENCY_QUOTE,
		BUSD_ADDRESS,
	)
	query = strings.ReplaceAll(
		query,
		DATE_FROM,
		sinceRFC3339,
	)

	resp, err := bitquery.Query(query, nil)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]map[string][]bitquery.Info)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	var curPrice *model.CryptoInfo
	curPrice = &model.CryptoInfo{
		CurrentPrice: data["data"]["ethereum"]["currentPrice"][0].QuotePrice,
	}

	return curPrice, nil
}
