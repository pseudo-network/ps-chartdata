package service

import (
	"encoding/json"
	"net/http"
	"ps-chartdata/bitquery"
	"strings"

	"github.com/labstack/echo"
)

func GetTransactionHandler(c echo.Context) error {

	baseCurrency := c.Param("address")

	query := `{
		ethereum(network: bsc) {
		  dexTrades(
			options: {limit: 100, desc: "timeInterval.second"}
			baseCurrency: {is: "0xb27adaffb9fea1801459a1a81b17218288c097cc"}
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
		BASE_CURRENCY,
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
