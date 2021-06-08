package service

import (
	"net/http"
	"ps-chartdata/model"

	"github.com/labstack/echo"
)

// todo: refactor, this is all messed up
// GET/config handler
func SearchSymbolsHandler(c echo.Context) error {

	// todo: revise?
	search := model.Symbol{
		Name:                 "SafeMoon",
		ExchangeTraded:       "SafeMoon",
		ExchangeListed:       "SafeMoon",
		Timezone:             "America/New_York",
		PriceScale:           1000000000000000,
		Minmov:               1,
		Minmov2:              0,
		PointValue:           1,
		Session:              "24x7",
		HasIntraday:          true,
		HasNoVolume:          false,
		Description:          "Safe Moon Crypto",
		Type:                 "stock",
		SupportedResolutions: []string{"5", "240", "D", "5D", "1w", "1m"},
	}

	return c.JSON(http.StatusOK, search)
}
