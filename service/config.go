package service

import (
	"net/http"
	"ps-chartdata/model"

	"github.com/labstack/echo"
)

// GET/config handler
func GetConfigHandler(c echo.Context) error {

	// todo: revise?
	config := model.Config{
		SupportedResolutions:   []string{"5", "240", "D", "5D", "1w", "1m"},
		SupportsGroupRequest:   false,
		SupportsMarks:          false,
		SupportsSearch:         true,
		SupportsTimescaleMarks: true,
		SupportsTime:           false,
	}

	return c.JSON(http.StatusOK, config)
}
