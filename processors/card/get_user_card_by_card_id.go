package card

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/processors/user"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
)

func GetUserCardByUserCardId(c echo.Context) error {
	id := c.Param("id")
	req := new(pb.GetUserCardByUserCardIdRequest)

	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyGetUserCardByCardIdFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if err := user.VerifyUser(req.GetUserId()); err != nil {
		return resp.JSONResp(c, int64(pb.User_ERROR_USER_NOT_EXISTS), err.Error())
	}

	if cardInfo, trx, err := getUserCardByCardId(req, id); err != nil {
		if trx == nil {
			return resp.JSONResp(c, int64(pb.GetUserCardByUserCardIdRequest_ERROR_CARD_NOT_FOUND), err.Error())
		}
		return resp.JSONResp(c, int64(pb.GetUserCardByUserCardIdRequest_ERROR_FAILED), err.Error())
	} else {
		return resp.GetUserCardByCardIdResponseJSON(c, cardInfo, trx)
	}
}

func verifyGetUserCardByCardIdFields(req *pb.GetUserCardByUserCardIdRequest) error {
	if req.UserId == nil || req.GetUserId() < 0 {
		return errors.New("invalid user_id")
	}
	return nil
}

func getUserCardByCardId(req *pb.GetUserCardByUserCardIdRequest, userCardId string) (*pb.UserCardWithInfo, []*pb.TransactionBasic, error) {
	var (
		cardInfoDb *pb.UserCardWithInfoDb
		trx        []*pb.TransactionBasic
	)

	if err := orm.DbInstance().Raw(orm.Sql10(), req.GetUserId(), userCardId).Scan(&cardInfoDb).Error; err != nil {
		return nil, nil, err
	}

	cardInfo := new(pb.UserCardWithInfo)
	cardInfo.UserCard = new(pb.UserCard)
	cardInfo.CardInfo = new(pb.CardBasicInfo)

	cardInfo.UserCard = &pb.UserCard{
		Id:               cardInfoDb.Id,
		UserId:           cardInfoDb.UserId,
		CardId:           cardInfoDb.CardId,
		CardNickname:     cardInfoDb.CardNickname,
		CardStatus:       cardInfoDb.CardStatus,
		CardExpiry:       cardInfoDb.CardExpiry,
		AddedTimestamp:   cardInfoDb.AddedTimestamp,
		UpdatedTimestamp: cardInfoDb.UpdatedTimestamp,
	}
	cardInfo.CardInfo = &pb.CardBasicInfo{
		CardName:      cardInfoDb.CardName,
		ShortCardName: cardInfoDb.ShortCardName,
		CardType:      cardInfoDb.CardType,
		CardImage:     cardInfoDb.CardImage,
		CardIssuer:    cardInfoDb.CardIssuer,
	}

	if err := orm.DbInstance().Raw(orm.Sql11(), req.GetUserId(), userCardId).Scan(&trx).Error; err != nil {
		return nil, nil, err
	}

	if cardInfo == nil {
		return nil, nil, errors.New("card not found")
	}

	return cardInfo, trx, nil
}
