package processors

import (
	"encoding/json"
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func AddCard(c echo.Context) error {
	req := new(rewards_tracker.AddCardRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	if isExists, err := isCardExists(req.GetCardDetails()); err != nil {
		return c.JSON(http.StatusOK, rewards_tracker.AddCardResponse{
			ResponseMeta: &rewards_tracker.ResponseMeta{
				ErrorCode:    proto.Int64(0),
				ErrorMessage: proto.String(err.Error()),
			},
		})
	} else if isExists {
		return c.JSON(http.StatusOK, rewards_tracker.AddCardResponse{
			ResponseMeta: &rewards_tracker.ResponseMeta{
				ErrorCode:    proto.Int64(0),
				ErrorMessage: proto.String("card exists"),
			},
		})
	}

	if idx, err := createCard(req.GetCardDetails()); err != nil {
		return c.JSON(http.StatusOK, rewards_tracker.AddCardResponse{
			ResponseMeta: &rewards_tracker.ResponseMeta{
				ErrorCode:    proto.Int64(0),
				ErrorMessage: proto.String(err.Error()),
			},
		})
	} else {
		return c.JSON(http.StatusOK, rewards_tracker.AddCardResponse{
			CardId: proto.Int64(idx),
			ResponseMeta: &rewards_tracker.ResponseMeta{
				ErrorCode:    proto.Int64(0),
				ErrorMessage: proto.String("ok"),
			},
		})
	}
}

func verifyFields(c *rewards_tracker.Card) error {

	return errors.New("")
}

func isCardExists(c *rewards_tracker.Card) (bool, error) {
	var (
		hold []rewards_tracker.CardDb
	)

	if err := orm.DbInstance().Raw("SELECT * FROM milestracker_db.card_table WHERE card_name = ?", c.GetCardName()).Scan(&hold).Error; err != nil {
		return false, err
	}
	return hold != nil, nil
}

func createCard(c *rewards_tracker.Card) (int64, error) {
	LocalBaseWhitelistCategory, _ := json.Marshal(c.LocalBaseCardRules.WhitelistCategories)
	LocalBaseBlacklistCategory, _ := json.Marshal(c.LocalBaseCardRules.BlacklistCategories)
	LocalBonusWhitelistCategory, _ := json.Marshal(c.LocalBonusCardRules.WhitelistCategories)
	LocalBonusBlacklistCategory, _ := json.Marshal(c.LocalBonusCardRules.BlacklistCategories)
	LocalBonusPaymentTypes, _ := json.Marshal(c.LocalBonusCardRules.WhitelistPaymentTypes)

	FcyBaseWhitelistCategory, _ := json.Marshal(c.FcyBaseCardRules.WhitelistCategories)
	FcyBaseBlacklistCategory, _ := json.Marshal(c.FcyBaseCardRules.BlacklistCategories)
	FcyBonusWhitelistCategory, _ := json.Marshal(c.FcyBonusCardRules.WhitelistCategories)
	FcyBonusBlacklistCategory, _ := json.Marshal(c.FcyBonusCardRules.BlacklistCategories)
	FcyBonusPaymentTypes, _ := json.Marshal(c.FcyBonusCardRules.WhitelistPaymentTypes)

	cc := rewards_tracker.CardDb{
		CardId:                      nil,
		CardName:                    c.CardName,
		ShortCardName:               c.ShortCardName,
		CardType:                    c.CardType,
		CardImage:                   c.CardImage,
		CardIssuer:                  c.CardIssuer,
		LocalBaseRewards:            c.LocalBaseRewards,
		LocalBaseMiles:              c.LocalBaseMiles,
		LocalBaseWhitelistCategory:  LocalBaseWhitelistCategory,
		LocalBaseBlacklistCategory:  LocalBaseBlacklistCategory,
		LocalBonusRewards:           c.LocalBonusRewards,
		LocalBonusMiles:             c.LocalBonusMiles,
		LocalBonusWhitelistCategory: LocalBonusWhitelistCategory,
		LocalBonusBlacklistCategory: LocalBonusBlacklistCategory,
		LocalBonusPaymentTypes:      LocalBonusPaymentTypes,
		FcyBaseRewards:              c.FcyBaseRewards,
		FcyBaseMiles:                c.FcyBaseMiles,
		FcyBaseWhitelistCategory:    FcyBaseWhitelistCategory,
		FcyBaseBlacklistCategory:    FcyBaseBlacklistCategory,
		FcyBonusRewards:             c.FcyBonusRewards,
		FcyBonusMiles:               c.FcyBonusMiles,
		FcyBonusWhitelistCategory:   FcyBonusWhitelistCategory,
		FcyBonusBlacklistCategory:   FcyBonusBlacklistCategory,
		FcyBonusPaymentTypes:        FcyBonusPaymentTypes,
		Rounding:                    c.Rounding,
		AmountBlock:                 c.AmountBlock,
		RewardCurrency:              c.RewardCurrency,
		CapType:                     c.CapType,
		Cap:                         c.Cap,
	}

	if err := orm.DbInstance().Table("milestracker_db.card_table").Create(&cc).Error; err != nil {
		return -1, err
	}

	if err := orm.DbInstance().Raw("SELECT * FROM milestracker_db.card_table WHERE card_name = ?", c.GetCardName()).Scan(&cc).Error; err != nil {
		return -1, err
	}

	return cc.GetCardId(), nil
}
