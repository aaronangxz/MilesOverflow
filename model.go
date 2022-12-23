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
	CATEGORY_AIR_TICKETS      = 1
	CATEGORY_DONATIONS        = 2
	CATEGORY_DINING           = 3
	CATEGORY_EDUCATION        = 4
	CATEGORY_FOOD_DELIVERY    = 5
	CATEGORY_RIDE_HAILING     = 6
	CATEGORY_GROCERIES        = 7
	CATEGORY_HOTELS           = 8
	CATEGORY_HOSPITALS        = 9
	CATEGORY_INSURANCE        = 10
	CATEGORY_ONLINE           = 11
	CATEGORY_SHOPPING         = 12
	CATEGORY_PUBLIC_TRANSPORT = 13
	CATEGORY_TELCO            = 14
	CATEGORY_UTILITIES        = 15

	PAYMENT_TYPE_CHIP        = 0
	PAYMENT_TYPE_STRIPE      = 1
	PAYMENT_TYPE_CONTACTLESS = 2
	PAYMENT_TYPE_ONLINE      = 3

	ROUNDING_ROUND_DOWN = 0
	ROUNDING_ROUND      = 1

	CAP_TYPE_NO_CAP          = 0
	CAP_TYPE_CALENDAR_MONTH  = 1
	CAP_TYPE_STATEMENT_MONTH = 2
)

type Card struct {
	cardId                          int64
	cardName                        string
	shortCardName                   string
	cardType                        int64
	cardImage                       string
	cardIssuer                      string
	localBaseReward                 int64
	localBaseMiles                  float64
	localBaseWhitelistCategories    []int64
	localBaseBlacklistCategories    []int64
	localBaseWhitelistPaymentTypes  []int64
	localBaseBlacklistPaymentTypes  []int64
	localBonusReward                int64
	localBonusMiles                 float64
	localBonusWhitelistCategories   []int64
	localBonusBlacklistCategories   []int64
	localBonusWhitelistPaymentTypes []int64
	localBonusBlacklistPaymentTypes []int64
	fcyBaseReward                   int64
	fcyBaseMiles                    float64
	fcyBaseWhitelistCategories      []int64
	fcyBaseBlacklistCategories      []int64
	fcyBonusReward                  int64
	fcyBonusMiles                   float64
	fcyBonusWhitelistCategories     []int64
	fcyBonusBlacklistCategories     []int64
	fcyBonusPaymentTypes            []int64
	rounding                        int64
	amountBlock                     float64
	rewardCurrency                  string
	capType                         int64
	cap                             int64
}
