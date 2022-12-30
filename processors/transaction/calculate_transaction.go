package transaction

import (
	"encoding/json"
	"errors"
	"github.com/aaronangxz/RewardTracker/processors/card"
	"github.com/aaronangxz/RewardTracker/processors/user"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/aaronangxz/RewardTracker/utils"
	"github.com/labstack/echo/v4"
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
		calculated  *pb.CalculatedTransaction
		spending    *pb.CurrentSpending
		cardDetails *pb.CardDb
		err         error
	)

	if spending, err = user.GetCurrentSpendingByCard(c.GetUserId(), c.GetTransactionDetails().GetCardId()); err != nil {
		return nil, err
	}

	if cardDetails, err = card.GetCardDetails(c.GetTransactionDetails().GetCardId()); err != nil {
		return nil, err
	}

	if c.GetTransactionDetails().GetCurrency() != "SGD" {
		calculateFCY(c.GetTransactionDetails(), cardDetails)
	} else {
		return calculateLocal(c.GetTransactionDetails(), cardDetails, spending), nil
	}

	return calculated, nil
}

func calculateLocal(t *pb.Transaction, c *pb.CardDb, spending *pb.CurrentSpending) *pb.CalculatedTransaction {
	if isEligibleCategory(c, t.GetCategory()) && isEligiblePaymentType(c, t.GetPaymentType()) {
		return calculateBonusLocal(t, c, spending)
	}
	return calculateBaseLocal(t, c)
}

func calculateFCY(t *pb.Transaction, c *pb.CardDb) (float64, float64) {
	return calculateBaseFCY(t, c)
}

func isEligibleCategory(c *pb.CardDb, cat int64) bool {
	if cat == -1 {
		return false
	}

	var localBonusWhitelistCategories *pb.CardRules

	if err := json.Unmarshal(c.GetLocalBonusWhitelistCategory(), &localBonusWhitelistCategories); err != nil {

	}

	eligibleCats := localBonusWhitelistCategories.GetWhitelistCategories()
	for _, eligibleCat := range eligibleCats {
		if eligibleCat == cat {
			return true
		}
	}
	return false
}

func isEligiblePaymentType(c *pb.CardDb, paymentType int64) bool {
	var localBonusWhitelistPaymentTypes *pb.CardRules

	if err := json.Unmarshal(c.GetLocalBonusPaymentTypes(), &localBonusWhitelistPaymentTypes); err != nil {

	}

	eligiblePaymentTypes := localBonusWhitelistPaymentTypes.GetWhitelistPaymentTypes()
	for _, eligiblePaymentType := range eligiblePaymentTypes {
		if eligiblePaymentType == paymentType {
			return true
		}
	}
	return false
}

func calculateBaseLocal(t *pb.Transaction, c *pb.CardDb) *pb.CalculatedTransaction {
	var (
		amount     float64
		baseReward float64
		baseMiles  float64
	)

	amount = float64(t.GetAmount()) / 100 / c.GetAmountBlock()

	switch c.GetRounding() {
	case int64(pb.CardRounding_ROUND_DOWN):
		baseReward = math.Floor(amount) * float64(c.GetLocalBaseRewards())
		break
	case int64(pb.CardRounding_ROUND):
		baseReward = math.Round(amount) * float64(c.GetLocalBonusRewards())
		break
	}

	baseMiles = baseReward * c.GetLocalBaseMiles()

	return &pb.CalculatedTransaction{
		BaseMilesEarned:    proto.Int64(int64(baseMiles * 100)),
		BonusMilesEarned:   proto.Int64(0),
		BaseRewardsEarned:  proto.Int64(int64(baseReward * 100)),
		BonusRewardsEarned: proto.Int64(0),
	}
}

func calculateBonusLocal(t *pb.Transaction, c *pb.CardDb, current *pb.CurrentSpending) *pb.CalculatedTransaction {
	var (
		baseReward  float64
		bonusReward float64
		amount      float64
		miles       float64
	)

	//earns base regardless
	base := calculateBaseLocal(t, c)
	baseReward = float64(base.GetBaseRewardsEarned() / 100)

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
		BaseMilesEarned:    proto.Int64(base.GetBaseMilesEarned() * 100),
		BonusMilesEarned:   proto.Int64(int64(math.Round(miles*100) / 100)),
		BaseRewardsEarned:  proto.Int64(int64(baseReward * 100)),
		BonusRewardsEarned: proto.Int64(int64(bonusReward * 100)),
	}
}

func calculateBaseFCY(t *pb.Transaction, c *pb.CardDb) (float64, float64) {
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

	return baseReward, baseMiles
}

func processCap(c *pb.CardDb, amount float64, current float64) (bool, float64) {
	//Fully exceeded cap
	if current >= c.GetCap() {
		return false, 0
	}

	//Partially exceeded
	amountToEarnBonus := c.GetCap() - current
	return true, utils.Min(amountToEarnBonus, amount)
}
