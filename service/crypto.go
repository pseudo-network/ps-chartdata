package service

import (
	"encoding/json"
	"fmt"
	"ps-chartdata/bitquery"
	"ps-chartdata/model"
	"strconv"
	"strings"
	"time"
)

func GetCryptoBarsByAddress(baseCurrency, quoteCurrency, fromRFC3339, toRFC3339 string) ([]model.Bar, error) {

	var usdMultiplier *float64
	if quoteCurrency == WBNB_ADDRESS {
		bnbUSD, err := GetBNBInfo()
		if err != nil {
			return nil, err
		}
		usdMultiplier = &bnbUSD.CurrentPrice
	} else {
		usdMultiplier = nil
	}

	query := `
		{
			ethereum(network: bsc) {
				dexTrades(
					options: {asc: "timeInterval.second"}
					exchangeName: {in: ["Pancake", "Pancake v2"]}
					date: {since: "DATE_FROM", till: "DATE_TILL"}
					baseCurrency: {is: "CURRENCY_BASE"},
					quoteCurrency: {is: "CURRENCY_QUOTE"},
					tradeAmountUsd: {gt: 10}
				)
				{
					timeInterval {
						second(count: 60, format: "%Y-%m-%dT%H:%M:%SZ")
					}
					tradeAmount(in:USD)
	       		trades:count
					volume: quoteAmount
					high: quotePrice(calculate: maximum)
					low: quotePrice(calculate: minimum)
					open: minimum(of: block, get: quote_price)
					close: maximum(of: block, get: quote_price)
				}
			}
		}
	`

	query = strings.ReplaceAll(
		query,
		DATE_FROM,
		fromRFC3339,
	)
	query = strings.ReplaceAll(
		query,
		DATE_TILL,
		toRFC3339,
	)
	query = strings.ReplaceAll(
		query,
		CURRENCY_QUOTE,
		quoteCurrency,
	)
	query = strings.ReplaceAll(
		query,
		CURRENCY_BASE,
		baseCurrency,
	)

	resp, err := bitquery.Query(query)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]map[string][]bitquery.DexTrade)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	dexTrades := data["data"]["ethereum"]["dexTrades"]

	var bars []model.Bar
	for _, t := range dexTrades {
		bar := model.Bar{}

		dateTime, err := time.Parse("2006-01-02T15:04:00Z", t.TimeInterval.Second)
		if err != nil {
			continue
		}
		bar.Time = int64(dateTime.Unix()) * 1000

		open, err := strconv.ParseFloat(t.Open, 64)
		if err != nil {
			return nil, err
		}
		bar.Open = open

		bar.High = t.High

		bar.Low = t.Low

		close, err := strconv.ParseFloat(t.Close, 64)
		if err != nil {
			return nil, err
		}
		bar.Close = close

		bar.Volume = t.TradeAmount

		if usdMultiplier != nil {
			bar.Open = bar.Open * *usdMultiplier
			bar.High = bar.High * *usdMultiplier
			bar.Low = bar.Low * *usdMultiplier
			bar.Close = bar.Close * *usdMultiplier
			bar.Volume = bar.Volume * *usdMultiplier
		}

		bars = append(bars, bar)
	}

	return bars, nil
}

func GetCryptos(searchQuery string) ([]model.Crypto, error) {
	query := `
		query {
			search(string: "SEARCH_QUERY", network:bsc){
				subject{
					__typename
					... on Address {
						address
						annotation
					}
					... on Currency {
						symbol
						name
						address
						tokenId
						tokenType
					}
					... on SmartContract {
						address
						annotation
						contractType
						protocol
					}
					... on TransactionHash {
						hash
					}
				}
			}
		}
	`
	query = strings.ReplaceAll(
		query,
		SEARCH_QUERY,
		searchQuery,
	)

	resp, err := bitquery.Query(query)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string][]bitquery.Crypto)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	cryptos := []model.Crypto{}
	for _, c := range data["data"]["search"] {
		crypto := model.NewCrypto(c.Subject.Name, c.Subject.Address, c.Subject.Symbol, c.Subject.TokenType, c.Network.Network)
		cryptos = append(cryptos, *crypto)
	}

	return cryptos, nil
}

func GetCryptoTransactionsByAddress(address string) ([]bitquery.Transaction, error) {
	query := `{
		ethereum(network: bsc) {
		  dexTrades(
			options: {limit: 100, desc: "timeInterval.second"}
			baseCurrency: {is: "CURRENCY_BASE"}
		  ) {
				transaction {
					hash
				}
				timeInterval {
					second
				}
				buyAmount
				buyCurrency {
					symbol
					address
				}
				sellAmount
				sellCurrency {
					symbol
					address
				}
				tradeAmount(in: USD)
				}
		  }
		}
	`
	query = strings.ReplaceAll(
		query,
		CURRENCY_BASE,
		address,
	)

	resp, err := bitquery.Query(query)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]map[string][]bitquery.Transaction)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	dexTrades := data["data"]["ethereum"]["dexTrades"]

	return dexTrades, nil
}

func GetCryptoInfoByAddress(baseCurrency, quoteCurrency, fromRFC3339, toRFC3339 string) (*model.CryptoInfo, error) {

	query := `{
		ethereum(network: bsc) {
			volume: dexTrades(
				date: {since: "DATE_FROM", till: "DATE_TILL"}
				baseCurrency: {is: "CURRENCY_BASE"}
			) {
      	tradeAmount(in: USD)
			}
			begPrice: dexTrades(
				date: {in: "DATE_FROM"}
				baseCurrency: {is: "CURRENCY_BASE"}
				quoteCurrency: {is: "CURRENCY_QUOTE"}
			) {
				quotePrice
			}
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
		DATE_FROM,
		fromRFC3339,
	)
	query = strings.ReplaceAll(
		query,
		DATE_TILL,
		toRFC3339,
	)
	query = strings.ReplaceAll(
		query,
		CURRENCY_BASE,
		baseCurrency,
	)
	query = strings.ReplaceAll(
		query,
		CURRENCY_QUOTE,
		quoteCurrency,
	)

	fmt.Println(query)

	resp, err := bitquery.Query(query)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]map[string][]bitquery.Info)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	begPrice := data["data"]["ethereum"]["begPrice"][0].QuotePrice
	curPrice := data["data"]["ethereum"]["currentPrice"][0].QuotePrice
	volume := data["data"]["ethereum"]["volume"][0].TradeAmount

	var usdMultiplier *float64
	if quoteCurrency == WBNB_ADDRESS {
		bnbUSD, err := GetBNBInfo()
		if err != nil {
			return nil, err
		}
		usdMultiplier = &bnbUSD.CurrentPrice
	} else {
		usdMultiplier = nil
	}

	fmt.Println(usdMultiplier)

	cryptoInfo := model.NewCryptoInfo(begPrice, curPrice, volume, usdMultiplier)

	return cryptoInfo, nil
}
