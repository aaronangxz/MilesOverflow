package main

type Transaction struct {
	description string
	category    int64
	paymentType int64
	amount      int64
	currency    string
	time        int64
	cardId      int64
}

const (
	CATEGORY_GENERAL_SPENDING = 0
	CATEGORY_ONLINE           = 1
	CATEGORY_DINING           = 2
	CATEGORY_SHOPPING         = 3
	CATEGORY_TRAVEL           = 4

	PAYMENT_TYPE_CHIP        = 0
	PAYMENT_TYPE_STRIPE      = 1
	PAYMENT_TYPE_CONTACTLESS = 2
	PAYMENT_TYPE_ONLINE      = 3
)

type Card struct {
	cardId                 int64
	cardName               string
	shortCardName          string
	cardType               int64
	cardImage              string
	cardIssuer             string
	localBaseReward        int64
	localBaseMiles         float64
	localBonusCategories   []int64
	localBonusReward       int64
	localBonusMiles        float64
	localBonusPaymentTypes []int64
	fcyBaseReward          int64
	fcyBaseMiles           float64
	fcyBonusCategories     []int64
	fcyBonusReward         int64
	fcyBonusMiles          float64
	fcyBonusPaymentTypes   []int64
	rounding               int64
	amountBlock            int64
	rewardCurrency         string
	capType                int64
	cap                    int64
}
