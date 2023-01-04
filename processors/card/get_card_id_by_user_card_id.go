package card

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
)

func GetCardIdByUserCardId(userCardId int64) (int64, error) {
	var (
		hold *pb.CardDb
	)

	if err := orm.DbInstance().Raw(orm.Sql12(), userCardId).Scan(&hold).Error; err != nil {
		return 0, err
	}

	if hold == nil {
		return 0, errors.New("card not found")
	}

	return hold.GetCardId(), nil
}
