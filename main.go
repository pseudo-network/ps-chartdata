package main

import (
	"ps-chartdata/config"
	"ps-chartdata/handler"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	config.InitConf()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	initPublicRoutes(e)

	go e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.PORT)))
}

func initPublicRoutes(e *echo.Echo) {
	e.GET("/chains/:id/tokens", handler.GetTokensHandler)
	e.GET("/chains/:id/tokens/:address/day-summary", handler.GetDaySummaryByTokenAddressHandler)
	e.GET("/chains/:id/tokens/:address/bars", handler.GetBarsByTokenAddressHandler)
	e.GET("/chains/:id/tokens/:address/transactions", handler.GetTransactionsByTokenAddressHandler)
	e.GET("/chains/:id/addresses/:address/balances", handler.GetBalancesByAddressHandler)
}
