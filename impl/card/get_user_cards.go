package card

import (
	"errors"
	"fmt"
	"github.com/aaronangxz/RewardTracker/impl/user"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"time"
)

func GetUserCards(c echo.Context) error {
	req := new(pb.GetUserCardsRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyGetUserCardFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if err := user.VerifyUser(req.GetUserId()); err != nil {
		return resp.JSONResp(c, int64(pb.User_ERROR_USER_NOT_EXISTS), err.Error())
	}

	if cards, err := getCards(req); err != nil {
		return resp.JSONResp(c, int64(pb.GetUserCardsRequest_ERROR_FAILED), err.Error())
	} else {
		return resp.GetUserCardsResponseJSON(c, cards)
	}
}

func verifyGetUserCardFields(req *pb.GetUserCardsRequest) error {
	if req.UserId == nil || req.GetUserId() < 0 {
		return errors.New("invalid user_id")
	}

	if _, ok := pb.GetUserCardsRequest_OrderByField_name[int32(req.GetOrderBy())]; !ok {
		return errors.New("invalid order by field")
	}

	if _, ok := pb.OrderBy_name[int32(req.GetDirection())]; !ok {
		return errors.New("invalid order by direction")
	}

	filters := req.GetFilter()
	for _, status := range filters.GetCardStatuses() {
		if _, ok := pb.UserCard_CardStatus_name[int32(status)]; !ok {
			return errors.New(fmt.Sprintf("contains invalid card status: %v", status))
		}
	}

	return nil
}

func getCards(c *pb.GetUserCardsRequest) ([]*pb.UserCardWithInfo, error) {
	var (
		hold                  []*pb.UserCardWithInfoDb
		cardStatusFilterQuery string
		cardExpiryFilterQuery string
		orderByDirection      string
		orderByQuery          string
	)

	query := orm.Sql9(c.GetUserId())
	if c.Filter != nil && c.Filter.CardStatuses != nil && len(c.GetFilter().GetCardStatuses()) > 0 {
		cardStatusFilter := ""
		statuses := c.Filter.GetCardStatuses()
		for i, status := range statuses {
			cardStatusFilter += fmt.Sprint(status)
			if len(statuses) > 1 && i != len(statuses)-1 {
				cardStatusFilter += ","
			}
		}
		cardStatusFilterQuery = fmt.Sprintf(" AND card_status IN (%v)", cardStatusFilter)
	}

	if c.Filter != nil && c.Filter.IsExpired != nil && c.Filter.GetIsExpired() == true {
		cardExpiryFilterQuery = fmt.Sprintf(" AND card_expiry < %v", time.Now().Unix())
	}

	if c.Direction != nil {
		switch c.GetDirection() {
		case int64(pb.OrderBy_ASC):
			orderByDirection = "ASC"
			break
		case int64(pb.OrderBy_DESC):
			orderByDirection = "DESC"
			break
		}
	}

	if c.OrderBy != nil {
		switch c.GetOrderBy() {
		case int64(pb.GetUserCardsRequest_USER_CARD_ADDED_TIME):
			orderByQuery = fmt.Sprintf(" ORDER BY added_timestamp %v", orderByDirection)
			break
		case int64(pb.GetUserCardsRequest_USER_CARD_EXPIRY):
			orderByQuery = fmt.Sprintf(" ORDER BY card_expiry %v", orderByDirection)
			break
		case int64(pb.GetUserCardsRequest_USER_CARD_NICKNAME):
			orderByQuery = fmt.Sprintf(" ORDER BY card_nickname %v", orderByDirection)
			break
		}
	} else {
		orderByQuery = " ORDER BY added_timestamp DESC"
	}

	finalQuery := query + cardStatusFilterQuery + cardExpiryFilterQuery + orderByQuery
	log.Info(finalQuery)
	if err := orm.DbInstance().Raw(finalQuery).Scan(&hold).Error; err != nil {
		return nil, err
	}

	var cards []*pb.UserCardWithInfo

	for i, h := range hold {
		cards = append(cards, new(pb.UserCardWithInfo))
		cards[i].UserCard = &pb.UserCard{
			Id:               h.Id,
			UserId:           h.UserId,
			CardId:           h.CardId,
			CardNickname:     h.CardNickname,
			CardStatus:       h.CardStatus,
			CardExpiry:       h.CardExpiry,
			AddedTimestamp:   h.AddedTimestamp,
			UpdatedTimestamp: h.UpdatedTimestamp,
		}
		cards[i].CardInfo = &pb.CardBasicInfo{
			CardName:      h.CardName,
			ShortCardName: h.ShortCardName,
			CardType:      h.CardType,
			CardImage:     h.CardImage,
			CardIssuer:    h.CardIssuer,
		}
	}
	return cards, nil
}
