package transaction

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/processors/user"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
)

func GetUserTransactionByTrxId(c echo.Context) error {
	id := c.Param("id")
	req := new(pb.GetUserTransactionByTrxIdRequest)

	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyGetUserTransactionByTrxIdFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if err := user.VerifyUser(req.GetUserId()); err != nil {
		return resp.JSONResp(c, int64(pb.User_ERROR_USER_NOT_EXISTS), err.Error())
	}

	if trx, err := getUserTrxById(req, id); err != nil {
		if trx == nil {
			return resp.JSONResp(c, int64(pb.GetUserTransactionByTrxIdRequest_ERROR_TRX_NOT_FOUND), err.Error())
		}
		return resp.JSONResp(c, int64(pb.GetUserTransactionByTrxIdRequest_ERROR_FAILED), err.Error())
	} else {
		return resp.GetUserTransactionByTrxIdResponseJSON(c, trx)
	}
}

func verifyGetUserTransactionByTrxIdFields(req *pb.GetUserTransactionByTrxIdRequest) error {
	if req.UserId == nil || req.GetUserId() < 0 {
		return errors.New("invalid user_id")
	}
	return nil
}

func getUserTrxById(req *pb.GetUserTransactionByTrxIdRequest, trxId string) (*pb.TransactionDbWithCardInfo, error) {
	var (
		trxWithInfoDb *pb.TransactionDbWithCardInfoDb
	)

	if err := orm.DbInstance().Raw(orm.Sql7(), req.GetUserId(), trxId).Scan(&trxWithInfoDb).Error; err != nil {
		return nil, err
	}

	if trxWithInfoDb == nil {
		return nil, errors.New("transaction not found")
	}

	trxWithInfo := new(pb.TransactionDbWithCardInfo)
	trxWithInfo.Transaction = new(pb.TransactionDb)
	trxWithInfo.CardInfo = new(pb.CardBasicInfo)

	trxWithInfo.Transaction = &pb.TransactionDb{
		TrxId:                  trxWithInfoDb.TrxId,
		UserId:                 trxWithInfoDb.UserId,
		Description:            trxWithInfoDb.Description,
		Category:               trxWithInfoDb.Category,
		PaymentType:            trxWithInfoDb.PaymentType,
		Amount:                 trxWithInfoDb.Amount,
		AmountConverted:        trxWithInfoDb.AmountConverted,
		Currency:               trxWithInfoDb.Currency,
		TransactionTimestamp:   trxWithInfoDb.TransactionTimestamp,
		CreateTimestamp:        trxWithInfoDb.CreateTimestamp,
		UpdateTimestamp:        trxWithInfoDb.UpdateTimestamp,
		UserCardId:             trxWithInfoDb.UserCardId,
		IsCancel:               trxWithInfoDb.IsCancel,
		BaseMilesEarned:        trxWithInfoDb.BaseMilesEarned,
		BonusMilesEarned:       trxWithInfoDb.BonusMilesEarned,
		BaseRewardsEarned:      trxWithInfoDb.BaseRewardsEarned,
		BonusRewardsEarned:     trxWithInfoDb.BonusRewardsEarned,
		IsPromotion:            trxWithInfoDb.IsPromotion,
		PromotionId:            trxWithInfoDb.PromotionId,
		PromotionMilesEarned:   trxWithInfoDb.PromotionMilesEarned,
		PromotionRewardsEarned: trxWithInfoDb.PromotionRewardsEarned,
	}
	trxWithInfo.CardInfo = &pb.CardBasicInfo{
		CardName:      trxWithInfoDb.CardName,
		ShortCardName: trxWithInfoDb.ShortCardName,
		CardType:      trxWithInfoDb.CardType,
		CardImage:     trxWithInfoDb.CardImage,
		CardIssuer:    trxWithInfoDb.CardIssuer,
	}
	return trxWithInfo, nil
}
