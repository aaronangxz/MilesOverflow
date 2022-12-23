package main

import (
	"fmt"
	"math"
)

func main() {
	calculateMiles(mockTransaction())
}

func calculateMiles(t Transaction) {
	c := mockCard()[t.cardId-1]
	fmt.Println(c.cardName)
	fmt.Println("Trx amount:", t.currency, float64(t.amount)/100)
	if t.currency != "SGD" {
		calculateFCY(t)
	} else {
		r, m := calculateLocal(t, c)
		fmt.Println(c.rewardCurrency, ":", r, " Miles:", m)
	}
}

func isEligibleCategory(c Card, cat int64) bool {
	if cat == -1 {
		return false
	}
	eligibleCats := c.localBonusWhitelistCategories
	for _, eligibleCat := range eligibleCats {
		if eligibleCat == cat {
			return true
		}
	}
	return false
}

func isEligiblePaymentType(c Card, paymentType int64) bool {
	eligiblePaymentTypes := c.localBonusWhitelistPaymentTypes
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

func calculateBonusLocal(t Transaction, c Card) (float64, float64) {
	var (
		baseReward  float64
		bonusReward float64
		amount      float64
		miles       float64
	)

	baseReward, _ = calculateBaseLocal(t, c)
	amount = float64(t.amount) / 100 / c.amountBlock

	switch c.cardIssuer {
	case "UOB":
		bonusReward = math.Floor(amount) * float64(c.localBonusReward)
		break
	default:
		bonusReward = math.Floor(amount * float64(c.localBonusReward))
	}

	miles = (baseReward + bonusReward) * c.localBaseMiles * c.amountBlock

	return baseReward + bonusReward, math.Round(miles*100) / 100
}

func calculateBaseLocal(t Transaction, c Card) (float64, float64) {
	var (
		amount     float64
		baseReward float64
		baseMiles  float64
	)

	amount = float64(t.amount) / 100 / c.amountBlock

	switch c.rounding {
	case ROUNDING_ROUND_DOWN:
		baseReward = math.Floor(amount) * float64(c.localBaseReward)
		break
	case ROUNDING_ROUND:
		baseReward = math.Round(amount) * float64(c.localBaseReward)
		break
	}

	baseMiles = baseReward * c.localBaseMiles

	return baseReward, baseMiles
}

func calculateLocal(t Transaction, c Card) (float64, float64) {
	if isEligibleCategory(c, t.category) && isEligiblePaymentType(c, t.paymentType) && isWithinCap(c, float64(t.amount)/100) {
		return calculateBonusLocal(t, c)
	}
	return calculateBaseLocal(t, c)
}

func calculateFCY(t Transaction) {

}
