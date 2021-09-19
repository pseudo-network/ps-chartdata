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

func GetTokenPrice(baseCurrency, quoteCurrency, chainName string) (*model.TokenInfo, error) {
	sinceRFC3339 := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	fmt.Println(sinceRFC3339)

	query := `
		query ($baseCurrency: String!, $quoteCurrency: String!, $since: ISO8601DateTime){
			ethereum(network: bsc) {
				currentPrice: dexTrades(
					options: {desc: ["block.height"], limit: 1}
					baseCurrency: {is: $baseCurrency}
					quoteCurrency: {is: $quoteCurrency}
					date: {after: $since}
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
		CHAIN_NAME,
		chainName,
	)

	vars := make(map[string]interface{})
	vars["baseCurrency"] = baseCurrency
	vars["quoteCurrency"] = quoteCurrency
	vars["since"] = sinceRFC3339

	resp, err := bitquery.Query(query, &vars)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string]map[string][]bitquery.Info)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	var curPrice *model.TokenInfo
	curPrice = &model.TokenInfo{
		CurrentPrice: data["data"]["ethereum"]["currentPrice"][0].QuotePrice,
	}

	fmt.Println(curPrice)

	return curPrice, nil
}

func GetBarsByTokenAddress(baseCurrency, quoteCurrency, sinceRFC3339, tillRFC3339 string, interval int, limit int, chainName string) ([]model.Bar, error) {
	// todo: cleanup
	var nativeCurrency string
	var usdCurrency string
	if chainName == "bsc" {
		nativeCurrency = WBNB_ADDRESS
		usdCurrency = BUSD_ADDRESS
	}
	if chainName == "ethereum" {
		nativeCurrency = WETH_ADDRESS
		usdCurrency = USDC_ADDRESS
	}

	info, err := GetTokenPrice(nativeCurrency, usdCurrency, chainName)
	if err != nil {
		return nil, err
	}
	usdMultiplier := info.CurrentPrice

	query := `
		query ($baseCurrency: String!, $quoteCurrency: String!, $since: ISO8601DateTime, $till: ISO8601DateTime, $interval: Int, $limit: Int) {
			ethereum(network: CHAIN_NAME) {
			dexTrades(
				options: {limit: $limit, desc: "timeInterval.minute"}
				date: {since: $since, till: $till}
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
	}`
	query = strings.ReplaceAll(
		query,
		CHAIN_NAME,
		chainName,
	)

	vars := make(map[string]interface{})
	vars["baseCurrency"] = baseCurrency
	vars["quoteCurrency"] = quoteCurrency

	if sinceRFC3339 == "" {
		vars["since"] = nil
	} else {
		vars["since"] = sinceRFC3339
	}

	if tillRFC3339 == "" {
		vars["till"] = nil
	} else {
		vars["till"] = tillRFC3339
	}

	vars["interval"] = interval
	vars["limit"] = limit

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

	// fmt.Println(dexTrades)

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

		bar.Open = bar.Open * usdMultiplier
		bar.High = bar.High * usdMultiplier
		bar.Low = bar.Low * usdMultiplier
		bar.Close = bar.Close * usdMultiplier

		bars = append(bars, bar)
	}

	fmt.Println(bars)

	return bars, nil
}

func GetTokens(searchQuery string, chainName string) ([]model.Token, error) {
	query := `
		query ($searchQuery: String!) {
			search(string: $searchQuery, network: CHAIN_NAME){
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
		CHAIN_NAME,
		chainName,
	)

	vars := make(map[string]interface{})
	vars["searchQuery"] = searchQuery

	resp, err := bitquery.Query(query, &vars)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string][]bitquery.Token)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	tokens := []model.Token{}
	for _, c := range data["data"]["search"] {
		token := model.NewToken(c.Subject.Name, c.Subject.Address, c.Subject.Symbol, c.Subject.TokenType, c.Network.Network)
		tokens = append(tokens, *token)
	}

	return tokens, nil
}

func GetTransactionsByTokenAddress(address, chainName string) ([]bitquery.Transaction, error) {
	query := `
		query ($baseCurrency: String!) {
			ethereum(network: CHAIN_NAME) {
			dexTrades(
				options: {limit: 100, desc: "timeInterval.second"}
				baseCurrency: {is: $baseCurrency}
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
		CHAIN_NAME,
		chainName,
	)

	vars := make(map[string]interface{})
	vars["baseCurrency"] = address

	resp, err := bitquery.Query(query, &vars)
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

func GetTokenDaySummaryByAddress(baseCurrency, quoteCurrency, sinceRFC3339, chainName string) (*model.TokenInfo, error) {
	query := `
		query ($since: ISO8601DateTime, $baseCurrency: String!, $quoteCurrency: String!) {
			ethereum(network: CHAIN_NAME) {
				daySummaries: dexTrades(
					options: {limit: 1, desc: "timeInterval.day"}
					date: {since:  $since}
					any: [
						{
							baseCurrency: {is: $baseCurrency}, 
							quoteCurrency: {is: $quoteCurrency}
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
		}
	`
	query = strings.ReplaceAll(
		query,
		CHAIN_NAME,
		chainName,
	)

	vars := make(map[string]interface{})
	vars["since"] = sinceRFC3339
	vars["baseCurrency"] = baseCurrency
	vars["quoteCurrency"] = quoteCurrency

	resp, err := bitquery.Query(query, &vars)
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

	info, err := GetTokenPrice(baseCurrency, quoteCurrency, chainName)
	if err != nil {
		return nil, err
	}
	usdMultiplier := info.CurrentPrice

	openPrice, err := strconv.ParseFloat(daySummary.OpenPrice, 64)
	if err != nil {
		return nil, err
	}

	closePrice, err := strconv.ParseFloat(daySummary.ClosePrice, 64)
	if err != nil {
		return nil, err
	}

	tokenInfo := model.NewTokenInfo(
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

	return tokenInfo, nil
}
