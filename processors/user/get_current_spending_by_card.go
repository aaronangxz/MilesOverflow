package user

import (
	"github.com/aaronangxz/RewardTracker/orm"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/aaronangxz/RewardTracker/utils"
	"google.golang.org/protobuf/proto"
	"time"
)

func GetCurrentSpendingByCard(userId, cardId int64) (*pb.CurrentSpending, error) {
	var (
		cardDetails *pb.CurrentSpending
	)

	start, end := utils.MonthStartEndDate(time.Now().Unix())

	if err := orm.DbInstance().
		Raw("SELECT SUM(amount_converted) AS total_spending, COUNT(*) FROM milestracker_db.expense_table "+
			"WHERE user_id = ? AND card_id = ? AND transaction_timestamp >= ? AND transaction_timestamp <= ? AND is_cancel = 0", userId, cardId, start, end).Scan(&cardDetails).Error; err != nil {
		return nil, err
	}

	if cardDetails == nil {
		return &pb.CurrentSpending{
			TotalSpending:    proto.Int64(0),
			TransactionCount: proto.Int64(0),
		}, nil
	}
	return cardDetails, nil
}
