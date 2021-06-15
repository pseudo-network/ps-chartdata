package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ps-chartdata/model"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

const (
	FROM_TIME          = "FROM_TIME"
	TO_TIME            = "TO_TIME"
	BASE_CURRENCY      = "BASE_CURRENCY"
	QUOTE_CURRENCY     = "QUOTE_CURRENCY"
	FORMATTED_INTERVAL = "FORMATTED_INTERVAL"
)

// GET/config handler
func GetHistoryHandler(c echo.Context) error {

	tickerName := c.QueryParam("symbol")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	resolution := c.QueryParam("resolution")
	countback := c.QueryParam("countback")

	fromTimeUInt, err := strconv.ParseUint(string(from), 10, 64)
	fromTimeString := time.Unix(int64(fromTimeUInt), int64(0)).UTC().Format(time.RFC3339)

	toTimeUInt, err := strconv.ParseUint(string(to), 10, 64)
	toTimeString := time.Unix(int64(toTimeUInt), int64(0)).UTC().Format(time.RFC3339)

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

	fmt.Println(fromTimeString, toTimeString)

	fmt.Println(tickerName, from, to, resolution, countback)

	bitQueryURL := "https://graphql.bitquery.io"

	query := `{
		ethereum(network: bsc) {
			dexTrades(
				date: {since: "FROM_TIME" till:"TO_TIME"}
				exchangeAddress: {is: "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"} 
				baseCurrency: {is: "0x8076c74c5e3f5852037f31ff0093eeb8c8add8d3"}
				quoteCurrency: {is: "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"}
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
					name
				}
				quoteCurrency {
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
		FORMATTED_INTERVAL,
		timeInterval,
	)

	reqBody, err := json.Marshal(map[string]string{
		"query": query,
	})

	req, err := http.NewRequest("POST", bitQueryURL, bytes.NewBuffer(reqBody))
	if err != nil {
		c.Logger().Error(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "BQYug1u2azt1EzuPggXfnhdhzFObRW0g")
	fmt.Println("req.Body")
	fmt.Println(req.Body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	data := make(map[string]map[string]map[string][]model.DexTrade)
	err = json.Unmarshal(body, &data)

	// todo: map data to history struct
	history := model.History{}

	// fmt.Println(data)
	if len(data["data"]["ethereum"]["dexTrades"]) == 0 {
		n := model.NoHistory{
			StatusCode: "no_data",
			NextTime:   "",
		}
		return c.JSON(http.StatusOK, n)
	}
	// map the data...
	for _, trade := range data["data"]["ethereum"]["dexTrades"] {

		// parse and map from year/month/date to unix time
		t, err := time.Parse("2006-01-02 15:04:00", trade.TimeInterval.Minute)
		if err != nil {
			c.Logger().Error(err.Error())
			break
		}
		//1584403200 //1623110400
		//next time 1522108800
		/*fmt.Println(trade.TimeInterval.Minute)
		fmt.Println(t)
		fmt.Println(t.Unix())*/

		history.BarTime = append(history.BarTime, t.Unix())

		openPrice, err := strconv.ParseFloat(trade.Open, 64)
		if err != nil {
			c.Logger().Error(err.Error())
		}
		history.OpeningPrice = append(history.OpeningPrice, openPrice)

		closePrice, err := strconv.ParseFloat(trade.Close, 64)
		if err != nil {
			c.Logger().Error(err.Error())
		}
		history.ClosingPrice = append(history.ClosingPrice, closePrice)
		history.HighPrice = append(history.HighPrice, trade.High)

		history.LowPrice = append(history.LowPrice, trade.Low)

		history.Volume = append(history.Volume, trade.TradeAmount)
	}

	history.StatusCode = model.STATUS_OK

	return c.JSON(http.StatusOK, history)
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
