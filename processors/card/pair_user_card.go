package card

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/processors/user"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/aaronangxz/RewardTracker/utils"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"time"
)

func PairUserCard(c echo.Context) error {
	req := new(pb.PairUserCardRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyPairUserCardFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if err := user.VerifyUser(req.GetUserId()); err != nil {
		return resp.JSONResp(c, int64(pb.User_ERROR_USER_NOT_EXISTS), err.Error())
	}

	if isExists, err := isCardIdExists(req.GetCardId()); err != nil {
		if !isExists {
			return resp.JSONResp(c, int64(pb.PairUserCardRequest_ERROR_CARD_NOT_EXISTS), err.Error())
		}
		return resp.JSONResp(c, int64(pb.PairUserCardRequest_ERROR_FAILED), err.Error())
	}

	if cards, err := pairCard(req); err != nil {
		return resp.JSONResp(c, int64(pb.PairUserCardRequest_ERROR_FAILED), err.Error())
	} else {
		return resp.PairUserCardResponseJSON(c, cards)
	}
}

func verifyPairUserCardFields(req *pb.PairUserCardRequest) error {
	if req.UserId == nil || req.GetUserId() < 0 {
		return errors.New("invalid user_id")
	}

	if req.CardId == nil || req.GetCardId() < 0 {
		return errors.New("invalid card_id")
	}

	if req.CardNickname != nil && len(req.GetCardNickname()) > 50 {
		return errors.New("card nickname must not exceed 50 characters")
	}

	return nil
}

func pairCard(c *pb.PairUserCardRequest) ([]*pb.UserCard, error) {
	var (
		hold []*pb.UserCard
	)

	timeStamp, err := time.Parse("1/2006", c.GetCardExpiry())
	if err != nil {
		return nil, err
	}
	_, end := utils.MonthStartEndDate(timeStamp.Unix())

	userCard := &pb.UserCard{
		UserId:           c.UserId,
		CardId:           c.CardId,
		CardNickname:     c.CardNickname,
		CardStatus:       proto.Int64(int64(pb.UserCard_CARD_ACTIVE)),
		CardExpiry:       proto.Int64(end),
		AddedTimestamp:   proto.Int64(time.Now().Unix()),
		UpdatedTimestamp: proto.Int64(time.Now().Unix()),
	}

	if err := checkPairedCards(userCard); err != nil {
		return nil, err
	}

	if err := orm.DbInstance().Table(orm.UserCardTable).Create(&userCard).Error; err != nil {
		return nil, err
	}

	if err := orm.DbInstance().Raw("SELECT * FROM milestracker_db.user_card_table WHERE user_id = ? ORDER BY added_timestamp DESC", c.GetUserId()).Scan(&hold).Error; err != nil {
		return nil, err
	}

	return hold, nil
}

func checkPairedCards(c *pb.UserCard) error {
	var (
		hold []*pb.UserCard
	)

	if err := orm.DbInstance().Raw("SELECT * FROM milestracker_db.user_card_table WHERE user_id = ? AND card_id = ? AND card_nickname = ? AND card_expiry = ?", c.GetUserId(), c.GetCardId(), c.GetCardNickname(), c.GetCardExpiry()).Scan(&hold).Error; err != nil {
		return err
	}

	if hold != nil {
		return errors.New("card already paired")
	}

	return nil
}
