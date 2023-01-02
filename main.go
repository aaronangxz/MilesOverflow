package main

import (
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/processors/card"
	"github.com/aaronangxz/RewardTracker/processors/transaction"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	orm.ConnectMySQL()
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("add_card", card.AddCard)
	e.POST("pair_user_card", card.PairUserCard)
	e.POST("get_user_cards", card.GetUserCards)
	e.POST("calculate_transaction", transaction.CalculateTransaction)
	e.POST("add_transaction", transaction.AddTransaction)

	e.Logger.Fatal(e.Start(":1323"))
}
