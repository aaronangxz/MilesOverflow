package main

import (
	"github.com/aaronangxz/RewardTracker/impl/card"
	"github.com/aaronangxz/RewardTracker/impl/transaction"
	"github.com/aaronangxz/RewardTracker/orm"
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
	//card/list - GetCards
	//card/delete/:id - DeleteCard
	//card/update/:id - UpdateCard

	//card/promotion/add - AddCardPromotion
	//card/promotion/list - GetCardPromotions
	//card/promotion/delete/:id - DeleteCardPromotion
	//card/promotion/update/:id - UpdateCardPromotion

	//FE
	e.POST("transaction/calculate", transaction.CalculateTransaction) //return card info for FE
	e.POST("transaction/add", transaction.AddTransaction)

	e.POST("user/card/pair", card.PairUserCard)
	//user/card/unpair - UnpairUserCard
	e.POST("user/card/list", card.GetUserCards)           //return user_card_id for FE to call GetUserCardByUserCardId
	e.POST("user/card/:id", card.GetUserCardByUserCardId) //return total monthly transaction amount on this card

	e.POST("user/transaction/list", transaction.GetUserTransactions)
	e.POST("user/transaction/:id", transaction.GetUserTransactionByTrxId) //return total monthly transaction amount on this card
	//user/transaction/delete/:id - DeleteUserTransaction

	e.Logger.Fatal(e.Start(":1323"))
}
