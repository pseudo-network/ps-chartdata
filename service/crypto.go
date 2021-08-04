package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ps-chartdata/bitquery"
	"ps-chartdata/model"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// GET/cryptos/:address/bars
func GetCryptoBarsHandler(c echo.Context) error {

	exchangeAddress := PANCAKESWAP_ADDRESS
	baseCurrency := c.Param("address")
	quoteCurrency := c.QueryParam("quote_currency")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	// resolution := c.QueryParam("resolution")

	fromDateUInt, err := strconv.ParseUint(string(from), 10, 64)
	fromDateString := time.Unix(int64(fromDateUInt), int64(0)).UTC().Format(time.RFC3339)

	tillDateUInt, err := strconv.ParseUint(string(to), 10, 64)
	tillDateString := time.Unix(int64(tillDateUInt), int64(0)).UTC().Format(time.RFC3339)

	// todo: these need to be refactored
	timeInterval := `minute(count: 15, format: "%Y-%m-%dT%H:%M:%SZ") `
	// if resolution == "5" {
	// 	timeInterval = "minute(count: 5)"
	// } else if resolution == "60" {
	// 	timeInterval = "minute(count: 60)"
	// } else if resolution == "240" {
	// 	timeInterval = "minute(count: 240)"
	// } else if resolution == "D" {
	// 	timeInterval = "minute(count: 1440)"
	// } else if resolution == "5D" {
	// 	timeInterval = "minute(count: 7200)"
	// } else if resolution == "1w" {
	// 	timeInterval = "minute(count: 10080)"
	// } else if resolution == "1m" {
	// 	timeInterval = "minute(count: 43200)"
	// }

	query := `
		{
			ethereum(network: bsc) {
				dexTrades(
					options: {asc: "timeInterval.second"}
					date: {since: "DATE_FROM", till: "DATE_TILL"}
					exchangeAddress: {is: "EXCHANGE_ADDRESS"}
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

	// query := `{
	// 	ethereum(network: bsc) {
	// 		dexTrades(
	// 			date: {since: "DATE_FROM" till:"DATE_TILL"}
	// 			exchangeAddress: {is: "EXCHANGE_ADDRESS"}
	// 			baseCurrency: {is: "CURRENCY_BASE"}
	// 			quoteCurrency: {is: "CURRENCY_QUOTE"}
	// 			)
	// 		{
	// 			timeInterval {
	// 				FORMATTED_INTERVAL
	// 			}
	// 			tradeAmount(in:USD)
	//       		trades:count
	// 			high: quotePrice(calculate: maximum)
	// 			low: quotePrice(calculate: minimum)
	// 			open: minimum(of: block, get: quote_price)
	// 			close: maximum(of: block, get: quote_price)
	// 			 baseCurrency {
	// 				symbol
	// 				name
	// 			}
	// 			quoteCurrency {
	// 				symbol
	// 				name
	// 			}
	// 			date {
	// 				date
	// 			}
	// 		}
	// 	}
	// }`

	query = strings.ReplaceAll(
		query,
		DATE_FROM,
		fromDateString,
	)
	query = strings.ReplaceAll(
		query,
		DATE_TILL,
		tillDateString,
	)
	query = strings.ReplaceAll(
		query,
		EXCHANGE_ADDRESS,
		exchangeAddress,
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
	query = strings.ReplaceAll(
		query,
		FORMATTED_INTERVAL,
		timeInterval,
	)

	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(query)
	fmt.Println()
	fmt.Println()
	fmt.Println()

	resp, err := bitquery.Query(query)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	data := make(map[string]map[string]map[string][]bitquery.DexTrade)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// todo:potentially revise this to be dynamic for other crypto networks
	dexTrades := data["data"]["ethereum"]["dexTrades"]

	var bars []Bar
	for _, t := range dexTrades {
		bar := Bar{}

		dateTime, err := time.Parse("2006-01-02T15:04:00Z", t.TimeInterval.Second)
		if err != nil {
			c.Logger().Error(err)
			continue
		}
		bar.Time = int64(dateTime.Unix()) * 1000

		open, err := strconv.ParseFloat(t.Open, 64)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		bar.Open = open

		bar.High = t.High

		bar.Low = t.Low

		close, err := strconv.ParseFloat(t.Close, 64)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		bar.Close = close

		bar.Volume = t.TradeAmount

		bars = append(bars, bar)
	}

	return c.JSON(http.StatusOK, bars)
}

// GET/cryptos
func GetCryptosHandler(c echo.Context) error {
	searchQuery := c.QueryParam("search_query")

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
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	data := make(map[string]map[string][]bitquery.Crypto)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	cryptos := []model.Crypto{}
	for _, c := range data["data"]["search"] {
		crypto := model.NewCrypto(c.Subject.Name, c.Subject.Address, c.Subject.Symbol, c.Subject.TokenType, c.Network.Network)
		cryptos = append(cryptos, *crypto)
	}

	return c.JSON(http.StatusOK, cryptos)
}

// GET/cryptos/:address/transactions
func GetCryptoTransactionsHandler(c echo.Context) error {
	baseCurrency := c.Param("address")

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
		baseCurrency,
	)

	resp, err := bitquery.Query(query)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	data := make(map[string]map[string]map[string][]bitquery.Transaction)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	dexTrades := data["data"]["ethereum"]["dexTrades"]

	return c.JSON(http.StatusOK, dexTrades)
}

// GET/cryptos/:address/info
func GetCryptoInfoByAddressHandler(c echo.Context) error {
	baseCurrency := c.Param("address")
	quoteCurrency := c.QueryParam("quote_currency")
	fromDateString := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	tillDateString := time.Now().Format(time.RFC3339)

	query := `{
		ethereum(network: bsc) {
			volume: dexTrades(
				date: {since: "DATE_FROM", till: "DATE_TILL"}
				baseCurrency: {is: "CURRENCY_BASE"}
			) {
				tradeAmount(in: USD)
			}
			begPrice: dexTrades(
				date: {is: "DATE_FROM"}
				baseCurrency: {is: "CURRENCY_BASE"}
				quoteCurrency: {is: "CURRENCY_QUOTE"}
			) {
				quotePrice
			}
			currentPrice: dexTrades(
				date: {is: "DATE_TILL"}
				baseCurrency: {is: "CURRENCY_BASE"}
				quoteCurrency: {is: "CURRENCY_QUOTE"}
			) {
				quotePrice
			}
		}
	}
	`
	query = strings.ReplaceAll(
		query,
		DATE_FROM,
		fromDateString,
	)
	query = strings.ReplaceAll(
		query,
		DATE_TILL,
		tillDateString,
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

	resp, err := bitquery.Query(query)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	data := make(map[string]map[string]map[string][]map[string]json.Number)
	err = json.Unmarshal(resp, &data)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	begPrice, _ := data["data"]["ethereum"]["begPrice"][0]["quotePrice"].Float64()
	curPrice, _ := data["data"]["ethereum"]["currentPrice"][0]["quotePrice"].Float64()
	volume, _ := data["data"]["ethereum"]["volume"][0]["tradeAmount"].Float64()

	cryptoInfo := model.NewCryptoInfo(begPrice, curPrice, volume)

	return c.JSON(http.StatusOK, cryptoInfo)
}
