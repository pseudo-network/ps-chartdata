package handler

import (
	"net/http"
	"ps-chartdata/service"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// GET/cryptos/:address/bars
func GetCryptoBarsHandler(c echo.Context) error {
	baseCurrency := c.Param("address")
	quoteCurrency := c.QueryParam("quote_currency")
	since := c.QueryParam("since")
	till := c.QueryParam("till")
	interval := c.QueryParam("interval")

	sinceRFC3339, err := unixStringToRFC3339String(since)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	tillRFC3339, err := unixStringToRFC3339String(till)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	bars, err := service.GetCryptoBarsByAddress(baseCurrency, quoteCurrency, sinceRFC3339, tillRFC3339, intervalInt)

	return c.JSON(http.StatusOK, bars)
}

// GET/cryptos
func GetCryptosHandler(c echo.Context) error {
	searchQuery := c.QueryParam("search_query")

	cryptos, err := service.GetCryptos(searchQuery)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, cryptos)
}

// GET/cryptos/:address/transactions
func GetCryptoTransactionsHandler(c echo.Context) error {
	address := c.Param("address")

	dexTrades, err := service.GetCryptoTransactionsByAddress(address)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, dexTrades)
}

// GET/cryptos/:address/day-summary
func GetCryptoDaySummaryByAddressHandler(c echo.Context) error {
	baseCurrency := c.Param("address")
	quoteCurrency := c.QueryParam("quote_currency")
	sinceRFC3339 := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)

	info, err := service.GetCryptoDaySummaryByAddress(baseCurrency, quoteCurrency, sinceRFC3339)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, info)
}
