package card

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

	if req.CardExpiry != nil && req.GetCardExpiry() < 0 {
		return errors.New("card expiry cannot be negative")
	}

	return nil
}

func pairCard(c *pb.PairUserCardRequest) ([]*pb.UserCard, error) {
	var (
		hold []*pb.UserCard
	)

	userCard := &pb.UserCard{
		UserId:           c.UserId,
		CardId:           c.CardId,
		CardNickname:     c.CardNickname,
		CardStatus:       proto.Int64(int64(pb.UserCard_CARD_ACTIVE)),
		CardExpiry:       c.CardExpiry,
		AddedTimestamp:   proto.Int64(time.Now().Unix()),
		UpdatedTimestamp: proto.Int64(time.Now().Unix()),
	}

	if err := checkPairedCards(userCard); err != nil {
		return nil, err
	}

	if err := orm.DbInstance().Table(orm.UserCardTable).Create(&userCard).Error; err != nil {
		return nil, err
	}

	if err := orm.DbInstance().Raw("SELECT * FROM milestracker_db.user_card_table WHERE user_id = ?", c.GetUserId()).Scan(&hold).Error; err != nil {
		return nil, err
	}

	return hold, nil
}

func checkPairedCards(c *pb.UserCard) error {
	var (
		hold []*pb.UserCard
	)

	if err := orm.DbInstance().Raw("SELECT * FROM milestracker_db.user_card_table WHERE user_id = ? AND card_id = ? AND card_nickname = ?", c.GetUserId(), c.GetCardId(), c.GetCardNickname()).Scan(&hold).Error; err != nil {
		return err
	}

	if hold != nil {
		return errors.New("card already paired")
	}

	return nil
}
