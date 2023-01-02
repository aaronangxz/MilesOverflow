package transaction

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/processors/user"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

func AddTransaction(c echo.Context) error {
	req := new(pb.AddTransactionRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyAddTransactionFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if err := user.VerifyUser(req.GetUserId()); err != nil {
		return resp.JSONResp(c, int64(pb.User_ERROR_USER_NOT_EXISTS), err.Error())
	}

	calReq := &pb.CalculateTransactionRequest{
		RequestMeta:        req.RequestMeta,
		UserId:             req.UserId,
		TransactionDetails: req.TransactionDetails,
	}
	if trx, err := calculate(calReq); err != nil {
		return resp.JSONResp(c, int64(pb.AddTransactionRequest_ERROR_FAILED), err.Error())
	} else {
		if addErr := add(req, trx); addErr != nil {
			return resp.JSONResp(c, int64(pb.AddTransactionRequest_ERROR_FAILED), addErr.Error())
		}
		return resp.GetAddTransactionResponseJSON(c)
	}
}

func verifyAddTransactionFields(req *pb.AddTransactionRequest) error {
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

func add(req *pb.AddTransactionRequest, t *pb.CalculatedTransaction) error {
	//TODO currency conversion
	expense := pb.TransactionDb{
		UserId:                 req.UserId,
		Description:            req.TransactionDetails.Description,
		Category:               req.TransactionDetails.Category,
		PaymentType:            req.TransactionDetails.PaymentType,
		Amount:                 req.TransactionDetails.Amount,
		AmountConverted:        nil,
		Currency:               req.TransactionDetails.Currency,
		TransactionTimestamp:   req.TransactionDetails.Time,
		CreateTimestamp:        proto.Int64(time.Now().Unix()),
		UpdateTimestamp:        proto.Int64(time.Now().Unix()),
		CardId:                 req.TransactionDetails.CardId,
		IsCancel:               proto.Int64(0),
		BaseMilesEarned:        proto.Int64(int64(t.GetBaseMilesEarned() * float64(100))),
		BonusMilesEarned:       proto.Int64(int64(t.GetBonusMilesEarned() * float64(100))),
		BaseRewardsEarned:      proto.Int64(int64(t.GetBaseRewardsEarned() * float64(100))),
		BonusRewardsEarned:     proto.Int64(int64(t.GetBonusRewardsEarned() * float64(100))),
		IsPromotion:            nil,
		PromotionId:            nil,
		PromotionMilesEarned:   nil,
		PromotionRewardsEarned: nil,
	}
	if err := orm.DbInstance().Table(orm.ExpenseTable).Create(&expense).Error; err != nil {
		return err
	}
	return nil
}
