package main

import (
	"github.com/aaronangxz/RewardTracker/processors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("add_card", processors.AddCard)

	e.Logger.Fatal(e.Start(":1323"))
}
