package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"ps-chartdata/model"

	"github.com/labstack/echo"
)

//const ()

func GetTransactionHandler(c echo.Context) error {

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
		}`

	reqBodyTransaction, err := json.Marshal(map[string]string{
		"query": query,
	})

	reqTransaction, err := http.NewRequest("POST", BITQUERY_URL, bytes.NewBuffer(reqBodyTransaction))
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	reqTransaction.Header.Set("Content-Type", "application/json")
	reqTransaction.Header.Set("X-API-KEY", "BQYug1u2azt1EzuPggXfnhdhzFObRW0g")

	respTransaction, err := http.DefaultClient.Do(reqTransaction)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	defer respTransaction.Body.Close()

	body, err := ioutil.ReadAll(respTransaction.Body)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	data := make(map[string]map[string]map[string][]model.Transaction)
	err = json.Unmarshal(body, &data)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// todo:potentially revise this to be dynamic for other crypto networks
	dexTrades := data["data"]["ethereum"]["dexTrades"]

	return c.JSON(http.StatusOK, dexTrades)
}
