package card

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
)

func isCardNameExists(cardName string) (bool, error) {
	var (
		hold []pb.CardDb
	)

	if err := orm.DbInstance().Raw(orm.Sql1(), cardName).Scan(&hold).Error; err != nil {
		return false, err
	}

	if hold != nil {
		return true, errors.New("card name already exists")
	}

	return false, nil
}

func isCardIdExists(cardId int64) (bool, error) {
	var (
		hold *pb.CardDb
	)

	if err := orm.DbInstance().Raw(orm.Sql2(), cardId).Scan(&hold).Error; err != nil {
		return false, err
	}

	if hold == nil {
		return false, errors.New("card id does not exist")
	}

	return true, nil
}
