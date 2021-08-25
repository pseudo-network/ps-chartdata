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

func GetCryptoBarsByAddress(baseCurrency, quoteCurrency, sinceRFC3339, tillRFC3339 string, interval int) ([]model.Bar, error) {

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
		query ($baseCurrency: String!, $quoteCurrency: String!, $since: ISO8601DateTime, $till: ISO8601DateTime, $interval: Int) {
			ethereum(network: bsc) {
			dexTrades(
				options: {asc: "timeInterval.minute"}
				date: {since: $since, till: $till}
				exchangeAddress: {is: "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"}
				baseCurrency: {is: $baseCurrency},
				quoteCurrency: {is: $quoteCurrency},
				tradeAmountUsd: {gt: 10}
			) 
			{
				timeInterval {
					minute(count: $interval, format: "%Y-%m-%dT%H:%M:%SZ")  
				}
				volume: quoteAmount
				high: quotePrice(calculate: maximum)
				low: quotePrice(calculate: minimum)
				open: minimum(of: block, get: quote_price)
				close: maximum(of: block, get: quote_price) 
			}
		}
	}
	`

	// query := `
	// {
	// 	ethereum(network: bsc) {
	// 	  dexTrades(
	// 		options: {asc: "timeInterval.minute"}
	// 		date: {since: "2021-06-20T07:23:21.000Z", till: "2021-06-23T15:23:21.000Z"}
	// 		exchangeAddress: {is: "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"}
	// 		baseCurrency: {is: "0x2170ed0880ac9a755fd29b2688956bd959f933f8"},
	// 		quoteCurrency: {is: "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"},
	// 		tradeAmountUsd: {gt: 10}
	// 	  )
	// 	  {
	// 		timeInterval {
	// 		  minute(count: 15, format: "%Y-%m-%dT%H:%M:%SZ")
	// 		}
	// 		volume: quoteAmount
	// 		high: quotePrice(calculate: maximum)
	// 		low: quotePrice(calculate: minimum)
	// 		open: minimum(of: block, get: quote_price)
	// 		close: maximum(of: block, get: quote_price)
	// 	  }
	// 	}
	//   }`

	// query = strings.ReplaceAll(
	// 	query,
	// 	DATE_FROM,
	// 	sinceRFC3339,
	// )
	// query = strings.ReplaceAll(
	// 	query,
	// 	DATE_TILL,
	// 	tillRFC3339,
	// )
	// query = strings.ReplaceAll(
	// 	query,
	// 	CURRENCY_QUOTE,
	// 	quoteCurrency,
	// )
	// query = strings.ReplaceAll(
	// 	query,
	// 	CURRENCY_BASE,
	// 	baseCurrency,
	// )

	fmt.Println(query)

	vars := make(map[string]interface{})
	vars["baseCurrency"] = baseCurrency
	vars["quoteCurrency"] = quoteCurrency
	vars["since"] = sinceRFC3339
	vars["till"] = tillRFC3339
	vars["interval"] = interval

	resp, err := bitquery.Query(query, &vars)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]map[string][]bitquery.DexTrade)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	dexTrades := data["data"]["ethereum"]["dexTrades"]

	fmt.Println(dexTrades)

	var bars []model.Bar
	for _, t := range dexTrades {
		bar := model.Bar{}

		dateTime, err := time.Parse("2006-01-02T15:04:00Z", t.TimeInterval.Minute)
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

		bar.Median = t.Median

		bar.Volume = t.Volume

		if usdMultiplier != nil {
			bar.Open = bar.Open * *usdMultiplier
			bar.High = bar.High * *usdMultiplier
			bar.Low = bar.Low * *usdMultiplier
			bar.Close = bar.Close * *usdMultiplier
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

	resp, err := bitquery.Query(query, nil)
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

	resp, err := bitquery.Query(query, nil)
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

func GetCryptoDaySummaryByAddress(baseCurrency, quoteCurrency, sinceRFC3339 string) (*model.CryptoInfo, error) {

	query := `{
		ethereum(network: bsc) {
			daySummaries: dexTrades(
				options: {limit: 1, desc: "timeInterval.day"}
				date: {since:  "DATE_FROM"}
				exchangeName: {in: ["Pancake", "Pancake v2"]}
				any: [
					{
						baseCurrency: {is: "CURRENCY_BASE"}, 
						quoteCurrency: {is: "CURRENCY_QUOTE"}
					}
				]
			) {
				timeInterval {
					day(count: 1)
				}
				baseCurrency {
					name
					symbol
					address
				}
				quoteCurrency {
					name
					symbol
					address
				}
				quotePrice
				quoteAmount
				uniqueWallets: count(uniq: senders)
				tradeAmountUSD: tradeAmount(in:USD)
				tradeCount: count
				tradeVolume: baseAmount(calculate: sum)
				volumeValue: quoteAmount(calculate: sum)
				maxPrice: quotePrice(calculate: maximum)
				minPrice: quotePrice(calculate: minimum)
				openPrice: minimum(of: block, get: quote_price)
				closePrice: maximum(of: block, get: quote_price)
			}
		}
	}`

	/*
			overviews: transfers(date: {since: null, till: null}, amount: {gt: 0}) {
			minted: amount(
				calculate: sum
				sender: {is: "0x0000000000000000000000000000000000000000"}
			)
			burned: amount(
				calculate: sum
				receiver: {is: "0x000000000000000000000000000000000000dEaD"}
			)
			uniqueWallets: count(uniq: senders)
			currency(currency: {is: "CURRENCY_BASE"}) {
				symbol
				name
				tokenId
			}
		}
	*/

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
	query = strings.ReplaceAll(
		query,
		DATE_FROM,
		sinceRFC3339,
	)

	fmt.Println(query)

	resp, err := bitquery.Query(query, nil)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]bitquery.Summary)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	daySummary := data["data"]["ethereum"].DaySummaries[0]
	// overview := data["data"]["ethereum"].OverViews[0]

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

	openPrice, err := strconv.ParseFloat(daySummary.OpenPrice, 64)
	if err != nil {
		return nil, err
	}

	closePrice, err := strconv.ParseFloat(daySummary.ClosePrice, 64)
	if err != nil {
		return nil, err
	}

	cryptoInfo := model.NewCryptoInfo(
		daySummary.QuotePrice,
		daySummary.TradeVolume,
		// overview.Minted,
		// overview.Burned,
		daySummary.TradeCount,
		daySummary.TradeAmountUSD,
		daySummary.MaxPrice,
		daySummary.MinPrice,
		openPrice,
		closePrice,
		// overview.UniqueWalletsCount,
		usdMultiplier,
	)

	return cryptoInfo, nil
}
