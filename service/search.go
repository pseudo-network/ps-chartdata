package service

import (
	"net/http"
	"ps-chartdata/model"

	"github.com/labstack/echo"
)

// todo: refactor, this is all messed up
// GET/config handler
func SearchHandler(c echo.Context) error {

	// todo: revise?
	search := model.Search{
		Symbol:      "test",
		FullName:    "test",
		Description: "test",
		Exchange:    "test",
		Ticker:      "te",
		Type:        "bitcoin",
	}

	return c.JSON(http.StatusOK, search)
}
