package handler

import (
	"net/http"
	"ps-chartdata/service"
	"time"

	"github.com/labstack/echo"
)

// GET/cryptos/:address/bars
func GetCryptoBarsHandler(c echo.Context) error {
	baseCurrency := c.Param("address")
	quoteCurrency := c.QueryParam("quote_currency")
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	fromRFC3339, err := unixStringToRFC3339String(from)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	toRFC3339, err := unixStringToRFC3339String(to)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	bars, err := service.GetCryptoBarsByAddress(baseCurrency, quoteCurrency, fromRFC3339, toRFC3339)

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
	fromRFC3339 := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)

	info, err := service.GetCryptoDaySummaryByAddress(baseCurrency, quoteCurrency, fromRFC3339)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, info)
}
