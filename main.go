package main

import (
	"fmt"
	"time"
)

func main() {
	calculateMiles(mockTransaction())
}

func calculateMiles(t Transaction) {
	fmt.Println("Trx amount:", t.currency, float64(t.amount)/100)
	if t.currency != "SGD" {
		calculateFCY(t)
	} else {
		r, m := calculateLocal(t)
		fmt.Println("Rewards:", r, " Miles:", m)
	}
}

func isEligibleCategory(c Card, cat int64) bool {
	eligibleCats := c.localBonusCategories
	for _, eligibleCat := range eligibleCats {
		if eligibleCat == cat {
			return true
		}
	}
	return false
}

func isEligiblePaymentType(c Card, paymentType int64) bool {
	eligiblePaymentTypes := c.localBonusPaymentTypes
	for _, eligiblePaymentType := range eligiblePaymentTypes {
		if eligiblePaymentType == paymentType {
			return true
		}
	}
	return false
}

func isWithinCap(c Card, amount float64) bool {
	return true
}

func processLocalCategory(c Card, cat int64) bool {
	if !isEligibleCategory(c, cat) {
		return false
	}
	return true
}

func processLocalPaymentType(c Card, paymentType int64) bool {
	if !isEligiblePaymentType(c, paymentType) {
		return false
	}
	return true
}

func processCap(c Card, amount int64) bool {
	if !isWithinCap(c, float64(amount)/100) {
		return false
	}
	return true
}

func calculateBonusLocal(t Transaction, c Card) (int64, float64) {
	amount := float64(t.amount) / 100
	baseReward := amount * float64(c.localBaseReward)
	baseMiles := amount * c.localBaseMiles

	bonusReward := amount * float64(c.localBonusReward)
	bonusMiles := amount * c.localBonusMiles

	return int64(baseReward + bonusReward), baseMiles + bonusMiles
}

func calculateBaseLocal(t Transaction, c Card) (int64, float64) {
	amount := float64(t.amount) / 100
	baseReward := int64(amount) * c.localBaseReward
	baseMiles := amount * c.localBaseMiles

	return baseReward, baseMiles
}

func calculateLocal(t Transaction) (int64, float64) {
	c := mockCard()
	if t.category != -1 {
		if processLocalCategory(c, t.category) {
			if processLocalPaymentType(c, t.paymentType) {
				if processCap(c, t.amount) {
					return calculateBonusLocal(t, c)
				} else {
				}
			}
		}
	}
	return calculateBaseLocal(t, c)
}

func calculateFCY(t Transaction) {

}

func mockCard() Card {
	return Card{
		cardId:                 1,
		cardName:               "HSBC Revolution",
		shortCardName:          "HRV",
		cardType:               0,
		cardImage:              "",
		cardIssuer:             "HSBC",
		localBaseReward:        1,
		localBaseMiles:         0.4,
		localBonusCategories:   []int64{CATEGORY_ONLINE, CATEGORY_SHOPPING, CATEGORY_DINING},
		localBonusReward:       9,
		localBonusMiles:        3.6,
		localBonusPaymentTypes: []int64{PAYMENT_TYPE_CONTACTLESS},
		fcyBaseReward:          1,
		fcyBaseMiles:           0.4,
		fcyBonusCategories:     nil,
		fcyBonusReward:         0,
		fcyBonusMiles:          0,
		fcyBonusPaymentTypes:   nil,
		rounding:               0,
		amountBlock:            1,
		rewardCurrency:         "HSBC Reward Points",
		capType:                0,
		cap:                    0,
	}
}

func mockTransaction() Transaction {
	return Transaction{
		description: "Mock Trx",
		category:    CATEGORY_ONLINE,
		paymentType: PAYMENT_TYPE_CONTACTLESS,
		amount:      999,
		currency:    "SGD",
		time:        time.Now().Unix(),
		cardId:      1,
	}
}
