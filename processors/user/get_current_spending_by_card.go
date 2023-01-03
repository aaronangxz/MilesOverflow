package user

import (
	"github.com/aaronangxz/RewardTracker/orm"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/aaronangxz/RewardTracker/utils"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"time"
)

func GetCurrentSpendingByCard(userId, cardId int64) (*pb.CurrentSpending, error) {
	var (
		cardDetails *pb.CurrentSpending
	)

	start, end := utils.MonthStartEndDate(time.Now().Unix())

	if err := orm.DbInstance().
		Raw(orm.Sql5(), userId, cardId, start, end).Scan(&cardDetails).Error; err != nil {
		return nil, err
	}

	if cardDetails == nil {
		return &pb.CurrentSpending{
			TotalSpending:    proto.Int64(0),
			TransactionCount: proto.Int64(0),
		}, nil
	}
	log.Info(cardDetails)
	return cardDetails, nil
}
