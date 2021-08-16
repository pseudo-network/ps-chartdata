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
	e.GET("/cryptos", handler.GetCryptosHandler)
	e.GET("/cryptos/:address/info", handler.GetCryptoInfoByAddressHandler)
	e.GET("/cryptos/:address/bars", handler.GetCryptoBarsHandler)
	e.GET("/cryptos/:address/transactions", handler.GetCryptoTransactionsHandler)
	e.GET("/wallets/:address/balances", handler.GetWalletBalancesByAddressHandler)
}
