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

	if err := orm.DbInstance().Raw("SELECT * FROM milestracker_db.user_table WHERE user_id = ?", userId).Scan(&hold).Error; err != nil {
		return err
	}

	if hold == nil {
		return errors.New("user does not exist")
	}

	return nil
}
