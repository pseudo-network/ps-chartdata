package handler

import (
	"net/http"
	"ps-chartdata/service"

	"github.com/labstack/echo"
)

// GET/wallets/:address/balances
func GetWalletBalancesByAddressHandler(c echo.Context) error {
	address := c.Param("address")

	balances, err := service.GetWalletBalancesByAddress(address)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, balances)
}
