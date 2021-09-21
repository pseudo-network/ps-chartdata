package handler

import (
	"fmt"
	"net/http"
	"ps-chartdata/model"
	"ps-chartdata/service"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

// GET/chains/:id/tokens/:address/bars
func GetBarsByTokenAddressHandler(c echo.Context) error {
	chainID := c.Param("id")
	baseTokenAddress := c.Param("address")
	since := c.QueryParam("since")
	till := c.QueryParam("till")
	interval := c.QueryParam("interval")
	limit := c.QueryParam("limit")

	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	chain, err := model.GetChainByID(chainIDInt)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	sinceRFC3339, err := unixStringToRFC3339String(since)
	if err != nil {
		sinceRFC3339 = ""
	}

	tillRFC3339, err := unixStringToRFC3339String(till)
	if err != nil {
		tillRFC3339 = ""
	}

	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	bars, err := service.GetBarsByTokenAddress(baseTokenAddress, sinceRFC3339, tillRFC3339, intervalInt, limitInt, *chain)

	return c.JSON(http.StatusOK, bars)
}

// GET/tokens
func GetTokensHandler(c echo.Context) error {
	chainID := c.Param("id")
	searchQuery := c.QueryParam("search_query")

	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	chain, err := model.GetChainByID(chainIDInt)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	tokens, err := service.GetTokens(searchQuery, *chain)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, tokens)
}

// GET/tokens/:address/transactions
func GetTransactionsByTokenAddressHandler(c echo.Context) error {
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

	dexTrades, err := service.GetTransactionsByTokenAddress(address, *chain)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, dexTrades)
}

// GET/tokens/:address/day-summary
func GetDaySummaryByTokenAddressHandler(c echo.Context) error {
	chainID := c.Param("id")
	baseTokenAddress := c.Param("address")
	sinceRFC3339 := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)

	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	chain, err := model.GetChainByID(chainIDInt)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	fmt.Println("chain", chain.Name)

	info, err := service.GetDaySummaryByTokenAddress(baseTokenAddress, sinceRFC3339, *chain)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, info)
}
