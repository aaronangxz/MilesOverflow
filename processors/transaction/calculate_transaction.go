package transaction

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/processors/card"
	"github.com/aaronangxz/RewardTracker/processors/user"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/aaronangxz/RewardTracker/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"math"
)

func CalculateTransaction(c echo.Context) error {
	req := new(pb.CalculateTransactionRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyCalculateTransactionFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if err := user.VerifyUser(req.GetUserId()); err != nil {
		return resp.JSONResp(c, int64(pb.User_ERROR_USER_NOT_EXISTS), err.Error())
	}

	if trx, err := calculate(req); err != nil {
		return resp.JSONResp(c, int64(pb.CalculateTransactionRequest_ERROR_FAILED), err.Error())
	} else {
		return resp.GetCalculateTransactionResponseJSON(c, trx)
	}
}

func verifyCalculateTransactionFields(req *pb.CalculateTransactionRequest) error {
	if req.UserId == nil || req.GetUserId() < 0 {
		return errors.New("invalid user_id")
	}

	if req.TransactionDetails == nil {
		return errors.New("transaction details is required")
	}

	trx := req.GetTransactionDetails()

	if trx.Category == nil {
		return errors.New("category is required")
	}

	if _, ok := pb.CardCategory_name[int32(trx.GetCategory())]; !ok {
		return errors.New("invalid category")
	}

	if trx.PaymentType == nil {
		return errors.New("payment_type is required")
	}

	if _, ok := pb.CardPaymentType_name[int32(trx.GetPaymentType())]; !ok {
		return errors.New("invalid payment_type")
	}

	if trx.Amount == nil {
		return errors.New("amount is required")
	}

	if float64(trx.GetAmount())/100 < 0.1 {
		return errors.New("invalid amount")
	}

	if trx.Currency == nil {
		return errors.New("currency is required")
	}

	if trx.Time == nil {
		return errors.New("time is required")
	}

	if trx.CardId == nil {
		return errors.New("card_id is required")
	}

	if trx.GetCardId() < 0 {
		return errors.New("invalid card_id")
	}

	return nil
}

func calculate(c *pb.CalculateTransactionRequest) (*pb.CalculatedTransaction, error) {
	var (
		spending    *pb.CurrentSpending
		cardDetails *pb.CardDb
		err         error
	)

	if spending, err = user.GetCurrentSpendingByCard(c.GetUserId(), c.GetTransactionDetails().GetCardId()); err != nil {
		return nil, err
	}

	log.Infof("GetCurrentSpendingByCard: %v", spending.GetTotalSpending())

	if cardDetails, err = card.GetCardDetails(c.GetTransactionDetails().GetCardId()); err != nil {
		return nil, err
	}

	log.Infof("GetCardDetails: %v", cardDetails)

	if c.GetTransactionDetails().GetCurrency() != "SGD" {
		return calculateFCY(c.GetTransactionDetails(), cardDetails), nil
	} else {
		return calculateLocal(c.GetTransactionDetails(), cardDetails, spending), nil
	}
}

func calculateLocal(t *pb.Transaction, c *pb.CardDb, spending *pb.CurrentSpending) *pb.CalculatedTransaction {
	if isEligibleCategory(c, t.GetCategory()) && isEligiblePaymentType(c, t.GetPaymentType()) {
		return calculateBonusLocal(t, c, spending)
	}
	return calculateBaseLocal(t, c)
}

func calculateFCY(t *pb.Transaction, c *pb.CardDb) *pb.CalculatedTransaction {
	return calculateBaseFCY(t, c)
}

func isEligibleCategory(c *pb.CardDb, cat int64) bool {
	if cat == -1 {
		return false
	}

	var localBonusWhitelistCategories pb.Lists

	if err := proto.Unmarshal(c.GetLocalBonusWhitelistCategory(), &localBonusWhitelistCategories); err != nil {
		log.Error(err)
	}

	log.Info(localBonusWhitelistCategories.GetList())

	eligibleCats := localBonusWhitelistCategories.GetList()
	for _, eligibleCat := range eligibleCats {
		if eligibleCat == cat {
			log.Infof("isEligibleCategory: %v", eligibleCat)
			return true
		}
	}
	log.Info("isEligibleCategory: false")
	return false
}

func isEligiblePaymentType(c *pb.CardDb, paymentType int64) bool {
	var localBonusWhitelistPaymentTypes pb.Lists

	if err := proto.Unmarshal(c.GetLocalBonusPaymentTypes(), &localBonusWhitelistPaymentTypes); err != nil {
		log.Error(err)
	}
	log.Info(localBonusWhitelistPaymentTypes.GetList())

	eligiblePaymentTypes := localBonusWhitelistPaymentTypes.GetList()
	for _, eligiblePaymentType := range eligiblePaymentTypes {
		if eligiblePaymentType == paymentType {
			log.Infof("isEligiblePaymentType: %v", eligiblePaymentType)
			return true
		}
	}
	log.Info("isEligiblePaymentType: false")
	return false
}

func calculateBaseLocal(t *pb.Transaction, c *pb.CardDb) *pb.CalculatedTransaction {
	var (
		amount     float64
		baseReward float64
		baseMiles  float64
	)
	log.Info("start calculateBaseLocal")

	amount = float64(t.GetAmount()) / 100 / c.GetAmountBlock()
	log.Infof("amount: %v", amount)

	switch c.GetRounding() {
	case int64(pb.CardRounding_ROUND_DOWN):
		baseReward = math.Floor(amount) * float64(c.GetLocalBaseRewards())
		break
	case int64(pb.CardRounding_ROUND):
		baseReward = math.Round(amount) * float64(c.GetLocalBaseRewards())
		break
	}

	baseMiles = baseReward * c.GetLocalBaseMiles()

	log.Info(baseMiles, baseReward)

	return &pb.CalculatedTransaction{
		BaseMilesEarned:    proto.Float64(baseMiles),
		BonusMilesEarned:   proto.Float64(0),
		BaseRewardsEarned:  proto.Float64(baseReward),
		BonusRewardsEarned: proto.Float64(0),
	}
}

func calculateBonusLocal(t *pb.Transaction, c *pb.CardDb, current *pb.CurrentSpending) *pb.CalculatedTransaction {
	var (
		baseReward  float64
		bonusReward float64
		amount      float64
		miles       float64
	)

	log.Info("start calculateBonusLocal")
	//earns base regardless
	base := calculateBaseLocal(t, c)
	log.Infof("calculateBaseLocal: %v", base)

	baseReward = base.GetBaseRewardsEarned()

	//amount is inflated by 100
	//divided by card amount blocks
	amount = float64(t.GetAmount()) / 100 / c.GetAmountBlock()

	//calculate cap
	if willEarnBonus, amountToEarnBonus := processCap(c, amount, float64(current.GetTotalSpending()/100)); willEarnBonus {
		switch c.GetCardIssuer() {
		//UOB has $5 block policy
		case "UOB":
			bonusReward = math.Floor(amountToEarnBonus) * float64(c.GetLocalBonusRewards())
			break
		default:
			bonusReward = math.Floor(amountToEarnBonus * float64(c.GetLocalBonusRewards()))
		}
	}

	miles = bonusReward * c.GetLocalBaseMiles() * c.GetAmountBlock()

	return &pb.CalculatedTransaction{
		BaseMilesEarned:    proto.Float64(base.GetBaseMilesEarned()),
		BonusMilesEarned:   proto.Float64(math.Round(miles*100) / 100),
		BaseRewardsEarned:  proto.Float64(baseReward),
		BonusRewardsEarned: proto.Float64(bonusReward),
	}
}

func calculateBaseFCY(t *pb.Transaction, c *pb.CardDb) *pb.CalculatedTransaction {
	var (
		amount     float64
		baseReward float64
		baseMiles  float64
	)

	amount = float64(t.GetAmount()) / 100 / c.GetAmountBlock()

	switch c.GetRounding() {
	case int64(pb.CardRounding_ROUND_DOWN):
		baseReward = math.Floor(amount) * float64(c.GetFcyBaseRewards())
		break
	case int64(pb.CardRounding_ROUND):
		baseReward = math.Round(amount) * float64(c.GetFcyBaseRewards())
		break
	}

	baseMiles = baseReward * c.GetFcyBaseMiles()

	return &pb.CalculatedTransaction{
		BaseMilesEarned:    proto.Float64(baseMiles),
		BonusMilesEarned:   nil,
		BaseRewardsEarned:  proto.Float64(baseReward),
		BonusRewardsEarned: nil,
	}
}

func processCap(c *pb.CardDb, amount float64, current float64) (bool, float64) {
	//Fully exceeded cap
	if current >= c.GetCap() {
		log.Info("processCap: Fully exceeded cap")
		return false, 0
	}

	//Partially exceeded
	amountToEarnBonus := c.GetCap() - current

	finalAmount := utils.Min(amountToEarnBonus, amount)

	log.Infof("processCap: Not exceeded cap, %v", finalAmount)
	return true, finalAmount
}
