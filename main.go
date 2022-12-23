package main

import (
	"fmt"
	"math"
)

func main() {
	currentRewards := float64(1200)
	calculateMiles(mockTransaction(), currentRewards)
}

func calculateMiles(t Transaction, current float64) {
	c := mockCard()[t.cardId-1]
	fmt.Println(c.cardName)
	fmt.Println("Trx amount:", t.currency, float64(t.amount)/100)
	if t.currency != "SGD" {
		calculateFCY(t, c)
	} else {
		r, m := calculateLocal(t, c, current)
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

func processCap(c Card, amount float64, current float64) (bool, float64) {
	//Fully exceeded cap
	if current >= c.cap {
		return false, 0
	}

	//Partially exceeded
	amountToEarnBonus := c.cap - current
	return true, min(amountToEarnBonus, amount)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func calculateBonusLocal(t Transaction, c Card, current float64) (float64, float64) {
	var (
		baseReward  float64
		bonusReward float64
		amount      float64
		miles       float64
	)

	//earns base regardless
	baseReward, _ = calculateBaseLocal(t, c)

	//amount is inflated by 100
	//divided by card amount blocks
	amount = float64(t.amount) / 100 / c.amountBlock

	//calculate cap
	if willEarnBonus, amountToEarnBonus := processCap(c, amount, current); willEarnBonus {
		switch c.cardIssuer {
		//UOB has $5 block policy
		case "UOB":
			bonusReward = math.Floor(amountToEarnBonus) * float64(c.localBonusReward)
			break
		default:
			bonusReward = math.Floor(amountToEarnBonus * float64(c.localBonusReward))
		}
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

func calculateLocal(t Transaction, c Card, current float64) (float64, float64) {
	if isEligibleCategory(c, t.category) && isEligiblePaymentType(c, t.paymentType) {
		return calculateBonusLocal(t, c, current)
	}
	return calculateBaseLocal(t, c)
}

func calculateBonusFCY(t Transaction, c Card, current float64) (float64, float64) {
	return 0, 0
}

func calculateBaseFCY(t Transaction, c Card) (float64, float64) {
	var (
		amount     float64
		baseReward float64
		baseMiles  float64
	)

	amount = float64(t.amount) / 100 / c.amountBlock

	switch c.rounding {
	case ROUNDING_ROUND_DOWN:
		baseReward = math.Floor(amount) * float64(c.fcyBaseReward)
		break
	case ROUNDING_ROUND:
		baseReward = math.Round(amount) * float64(c.fcyBaseReward)
		break
	}

	baseMiles = baseReward * c.fcyBaseMiles

	return baseReward, baseMiles
}

func calculateFCY(t Transaction, c Card) (float64, float64) {
	return calculateFCY(t, c)
}
