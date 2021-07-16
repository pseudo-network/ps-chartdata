package service

import (
	"encoding/json"
	"net/http"
	"ps-chartdata/bitquery"
	"ps-chartdata/model"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// GET/config handler
func GetCryptoBarsHandler(c echo.Context) error {

	exchangeAddress := PANCAKESWAP_ADDRESS
	baseCurrency := c.Param("address")
	quoteCurrency := c.QueryParam("quote_currency")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	resolution := c.QueryParam("resolution")

	fromTimeUInt, err := strconv.ParseUint(string(from), 10, 64)
	fromTimeString := time.Unix(int64(fromTimeUInt), int64(0)).UTC().Format(time.RFC3339)

	toTimeUInt, err := strconv.ParseUint(string(to), 10, 64)
	toTimeString := time.Unix(int64(toTimeUInt), int64(0)).UTC().Format(time.RFC3339)

	// todo: these need to be refactored
	timeInterval := "minute(count: 1)"
	if resolution == "5" {
		timeInterval = "minute(count: 5)"
	} else if resolution == "60" {
		timeInterval = "minute(count: 60)"
	} else if resolution == "240" {
		timeInterval = "minute(count: 240)"
	} else if resolution == "D" {
		timeInterval = "minute(count: 1440)"
	} else if resolution == "5D" {
		timeInterval = "minute(count: 7200)"
	} else if resolution == "1w" {
		timeInterval = "minute(count: 10080)"
	} else if resolution == "1m" {
		timeInterval = "minute(count: 43200)"
	}

	query := `{
		ethereum(network: bsc) {
			dexTrades(
				date: {since: "FROM_TIME" till:"TO_TIME"}
				exchangeAddress: {is: "EXCHANGE_ADDRESS"} 
				baseCurrency: {is: "BASE_CURRENCY"}
				quoteCurrency: {is: "QUOTE_CURRENCY"}
				)
			{
				timeInterval {
					FORMATTED_INTERVAL
				}
				tradeAmount(in:USD)
        		trades:count
				high: quotePrice(calculate: maximum)
				low: quotePrice(calculate: minimum)
				open: minimum(of: block, get: quote_price)
				close: maximum(of: block, get: quote_price)
				 baseCurrency {
					symbol
					name
				}
				quoteCurrency {
					symbol
					name
				}
				date {
					date
				}
			}
		}
	}`

	query = strings.ReplaceAll(
		query,
		FROM_TIME,
		fromTimeString,
	)
	query = strings.ReplaceAll(
		query,
		TO_TIME,
		toTimeString,
	)
	query = strings.ReplaceAll(
		query,
		EXCHANGE_ADDRESS,
		exchangeAddress,
	)
	query = strings.ReplaceAll(
		query,
		QUOTE_CURRENCY,
		quoteCurrency,
	)
	query = strings.ReplaceAll(
		query,
		BASE_CURRENCY,
		baseCurrency,
	)
	query = strings.ReplaceAll(
		query,
		FORMATTED_INTERVAL,
		timeInterval,
	)

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

	// map time interval to unix time ms for tradingview
	for i, t := range dexTrades {
		dateTime, err := time.Parse("2006-01-02 15:04:00", t.TimeInterval.Minute)
		if err != nil {
			c.Logger().Error(err)
			continue
		}
		t.UnixTimeMS = int64(dateTime.Unix()) * 1000
		dexTrades[i] = t
	}

	return c.JSON(http.StatusOK, dexTrades)
}

func GetCryptoHandler(c echo.Context) error {

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
