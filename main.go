package main

import (
	"github.com/aaronangxz/RewardTracker/impl/card"
	"github.com/aaronangxz/RewardTracker/impl/transaction"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	//Admin
	e.POST("api/v1/card/add", card.AddCard)
	//api/v1/card/list - GetCards
	//api/v1/card/delete/:id - DeleteCard
	//api/v1/card/update/:id - UpdateCard

	//api/v1/card/promotion/add - AddCardPromotion
	//api/v1/card/promotion/list - GetCardPromotions
	//api/v1/card/promotion/delete/:id - DeleteCardPromotion
	//api/v1/card/promotion/update/:id - UpdateCardPromotion

	//FE
	e.POST("api/v1/transaction/calculate", transaction.CalculateTransaction) //return card info for FE
	e.POST("api/v1/transaction/add", transaction.AddTransaction)

	e.POST("api/v1/user/card/pair", card.PairUserCard)
	//api/v1/user/card/unpair - UnpairUserCard
	e.POST("api/v1/user/card/list", card.GetUserCards)           //return user_card_id for FE to call GetUserCardByUserCardId
	e.POST("api/v1/user/card/:id", card.GetUserCardByUserCardId) //return total monthly transaction amount on this card

	e.POST("api/v1/user/transaction/list", transaction.GetUserTransactions)
	e.POST("api/v1/user/transaction/:id", transaction.GetUserTransactionByTrxId) //return total monthly transaction amount on this card
	//api/v1/user/transaction/delete/:id - DeleteUserTransaction

	e.Logger.Fatal(e.Start(":4000"))
}
