package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	globalBinder = &echo.DefaultBinder{}
)

func Bind(c echo.Context, s interface{}) error {
	if err := BindBody(c, s); err != nil {
		return err
	}
	if err := BindPath(c, s); err != nil {
		return err
	}
	if c.Request().Method == http.MethodGet || c.Request().Method == http.MethodDelete {
		if err := BindQuery(c, s); err != nil {
			return err
		}
	}
	return nil
}

func BindBody(c echo.Context, s interface{}) error {
	return globalBinder.BindBody(c, s)
}

func BindQuery(c echo.Context, s interface{}) error {
	return globalBinder.BindQueryParams(c, s)
}

func BindPath(c echo.Context, s interface{}) error {
	return globalBinder.BindPathParams(c, s)
}

func BindHeader(c echo.Context, s interface{}) error {
	return globalBinder.BindHeaders(c, s)
}
