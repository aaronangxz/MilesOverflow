package card

import (
	"github.com/aaronangxz/RewardTracker/orm"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
)

func GetCardDetails(cardId int64) (*pb.CardDb, error) {
	var (
		cardDetails *pb.CardDb
	)

	if err := orm.DbInstance().Raw(orm.Sql2(), cardId).Scan(&cardDetails).Error; err != nil {
		return nil, err
	}

	return cardDetails, nil
}
