package transaction

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/processors/user"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/aaronangxz/RewardTracker/utils"
	"github.com/labstack/echo/v4"
	"time"
)

func GetUserTransactions(c echo.Context) error {
	req := new(pb.GetUserTransactionsRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}
	if err := verifyGetUserTransactionsFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if err := user.VerifyUser(req.GetUserId()); err != nil {
		return resp.JSONResp(c, int64(pb.User_ERROR_USER_NOT_EXISTS), err.Error())
	}

	if trx, err := getUserTrx(req); err != nil {
		return resp.JSONResp(c, int64(pb.GetUserTransactionsRequest_ERROR_FAILED), err.Error())
	} else {
		return resp.GetUserTransactionsResponseJSON(c, trx)
	}
}

func verifyGetUserTransactionsFields(req *pb.GetUserTransactionsRequest) error {
	if req.UserId == nil || req.GetUserId() < 0 {
		return errors.New("invalid user_id")
	}

	return nil
}

func getUserTrx(req *pb.GetUserTransactionsRequest) ([]*pb.TransactionBasicWithCardInfo, error) {
	var (
		trx []*pb.TransactionBasicWithCardInfoDb
	)

	start, end := utils.MonthStartEndDate(time.Now().Unix())
	if err := orm.DbInstance().Raw(orm.Sql6(), req.GetUserId(), start, end).Scan(&trx).Error; err != nil {
		return nil, err
	}

	var trxWithInfo []*pb.TransactionBasicWithCardInfo

	for i, t := range trx {
		trxWithInfo = append(trxWithInfo, new(pb.TransactionBasicWithCardInfo))
		trxWithInfo[i].Transaction = &pb.TransactionBasic{
			TrxId:                t.TrxId,
			Description:          t.Description,
			Category:             t.Category,
			Amount:               t.Amount,
			Currency:             t.Currency,
			TransactionTimestamp: t.TransactionTimestamp,
			UserCardId:           t.UserCardId,
			TotalMilesEarned:     t.TotalMilesEarned,
		}
		trxWithInfo[i].CardInfo = &pb.CardBasicInfo{
			CardName:      t.CardName,
			ShortCardName: t.ShortCardName,
			CardType:      t.CardType,
			CardImage:     t.CardImage,
			CardIssuer:    t.CardIssuer,
		}
	}
	return trxWithInfo, nil
}
