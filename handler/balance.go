package handler

import (
	"net/http"
	"ps-chartdata/model"
	"ps-chartdata/service"
	"strconv"

	"github.com/labstack/echo"
)

// GET/wallets/:address/balances
func GetBalancesByAddressHandler(c echo.Context) error {
	chainID := c.Param("id")
	address := c.Param("address")

	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	chain, err := model.GetChainByID(chainIDInt)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	balances, err := service.GetBalancesByAddress(address, *chain)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, balances)
}
