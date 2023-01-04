package main

import (
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/processors/card"
	"github.com/aaronangxz/RewardTracker/processors/transaction"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
	orm.ConnectMySQL()
	e := echo.New()

	//Admin
	e.POST("card/add", card.AddCard)
	//card/delete/:id - DeleteCard
	//card/update/:id - UpdateCard

	//FE
	e.POST("transaction/calculate", transaction.CalculateTransaction)
	e.POST("transaction/add", transaction.AddTransaction)

	e.POST("user/card/pair", card.PairUserCard)
	//user/card/unpair - UnpairUserCard
	e.POST("user/card", card.GetUserCards)
	e.POST("user/card/:id", card.GetUserCardByUserCardId)

	//WIP return card info
	e.POST("user/transaction", transaction.GetUserTransactions)
	e.POST("user/transaction/:id", transaction.GetUserTransactionByTrxId)
	//user/transaction/delete/:id

	e.Logger.Fatal(e.Start(":1323"))
}
