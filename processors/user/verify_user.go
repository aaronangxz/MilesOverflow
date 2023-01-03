package user

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
)

func VerifyUser(userId int64) error {
	var (
		hold *pb.User
	)

	if err := orm.DbInstance().Raw(orm.Sql3(), userId).Scan(&hold).Error; err != nil {
		return err
	}

	if hold == nil {
		return errors.New("user does not exist")
	}

	return nil
}

func VerifyUserCard(userId int64, cardId int64) error {
	var (
		hold *pb.UserCard
	)

	if err := orm.DbInstance().Raw(orm.Sql4(), userId, cardId).Scan(&hold).Error; err != nil {
		return err
	}

	if hold == nil {
		return errors.New("card not paired to user")
	}

	return nil
}
