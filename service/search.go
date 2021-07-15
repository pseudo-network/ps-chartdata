package service

import (
	"encoding/json"
	"net/http"
	"ps-chartdata/bitquery"
	"ps-chartdata/model"
	"strings"

	"github.com/labstack/echo"
)

const (
	SEARCH_QUERY = "SEARCH_QUERY"
)

func GetCurrencyHandler(c echo.Context) error {

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
