package main

import (
	"ps-chartdata/config"
	"ps-chartdata/service"
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
	e.GET("/currencies", service.GetCurrencyHandler)
	e.GET("/currencies/:address/history", service.GetCurrencyHistoryHandler)
	e.GET("/currencies/:address/transactions", service.GetTransactionHandler)
}
