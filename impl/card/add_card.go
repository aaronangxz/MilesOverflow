package card

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

func AddCard(c echo.Context) error {
	req := new(pb.AddCardRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyAddCardFields(req.GetCardDetails()); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}

	if isExists, err := isCardNameExists(req.GetCardDetails().GetCardName()); err != nil {
		if isExists {
			return resp.JSONResp(c, int64(pb.AddCardRequest_ERROR_CARD_EXISTS), err.Error())
		}
		return resp.JSONResp(c, int64(pb.AddCardRequest_ERROR_FAILED), err.Error())
	}

	if idx, err := createCard(req.GetCardDetails()); err != nil {
		return resp.JSONResp(c, int64(pb.AddCardRequest_ERROR_FAILED), err.Error())
	} else {
		return resp.AddCardResponseJSON(c, idx)
	}
}

func verifyAddCardFields(c *pb.Card) error {
	if len(c.GetCardName()) > 50 {
		return errors.New("card name must not exceed 50 characters")
	}

	if len(c.GetShortCardName()) > len(c.GetCardName()) {
		return errors.New("card short name must not be longer than card name")
	}

	if _, ok := pb.CardType_name[int32(c.GetCardType())]; !ok {
		return errors.New("invalid card type")
	}

	if len(c.GetCardIssuer()) > 50 {
		return errors.New("card issuer must not exceed 50 characters")
	}

	if c.GetLocalBaseRewards() < 0 || c.GetLocalBonusRewards() < 0 || c.GetFcyBaseRewards() < 0 || c.GetFcyBonusRewards() < 0 ||
		c.GetLocalBaseMiles() < 0 || c.GetLocalBonusMiles() < 0 || c.GetFcyBaseMiles() < 0 || c.GetFcyBonusMiles() < 0 {
		return errors.New("miles and rewards cannot be negative")
	}

	if _, ok := pb.CardRounding_name[int32(c.GetRounding())]; !ok {
		return errors.New("invalid rounding type")
	}

	if _, ok := pb.CardPaymentType_name[int32(c.GetCapType())]; !ok {
		return errors.New("invalid cap type")
	}
	return nil
}

func fillCardToCardDb(c *pb.Card) *pb.CardDb {
	LocalBaseWhitelistCategory, _ := proto.Marshal(c.LocalBaseCardRules.WhitelistCategories)
	LocalBaseBlacklistCategory, _ := proto.Marshal(c.LocalBaseCardRules.BlacklistCategories)
	LocalBonusWhitelistCategory, _ := proto.Marshal(c.LocalBonusCardRules.WhitelistCategories)
	LocalBonusBlacklistCategory, _ := proto.Marshal(c.LocalBonusCardRules.BlacklistCategories)
	LocalBonusWhitelistPaymentTypes, _ := proto.Marshal(c.LocalBonusCardRules.WhitelistPaymentTypes)
	LocalBonusBlacklistPaymentTypes, _ := proto.Marshal(c.LocalBonusCardRules.BlacklistPaymentTypes)

	FcyBaseWhitelistCategory, _ := proto.Marshal(c.FcyBaseCardRules.WhitelistCategories)
	FcyBaseBlacklistCategory, _ := proto.Marshal(c.FcyBaseCardRules.BlacklistCategories)
	FcyBonusWhitelistCategory, _ := proto.Marshal(c.FcyBonusCardRules.WhitelistCategories)
	FcyBonusBlacklistCategory, _ := proto.Marshal(c.FcyBonusCardRules.BlacklistCategories)
	FcyBonusWhitelistPaymentTypes, _ := proto.Marshal(c.FcyBonusCardRules.WhitelistPaymentTypes)
	FcyBonusBlacklistPaymentTypes, _ := proto.Marshal(c.FcyBonusCardRules.BlacklistPaymentTypes)

	cc := &pb.CardDb{
		CardId:                          nil,
		CardName:                        c.CardName,
		ShortCardName:                   c.ShortCardName,
		CardType:                        c.CardType,
		CardImage:                       c.CardImage,
		CardIssuer:                      c.CardIssuer,
		LocalBaseRewards:                c.LocalBaseRewards,
		LocalBaseMiles:                  c.LocalBaseMiles,
		LocalBaseWhitelistCategory:      LocalBaseWhitelistCategory,
		LocalBaseBlacklistCategory:      LocalBaseBlacklistCategory,
		LocalBonusRewards:               c.LocalBonusRewards,
		LocalBonusMiles:                 c.LocalBonusMiles,
		LocalBonusWhitelistCategory:     LocalBonusWhitelistCategory,
		LocalBonusBlacklistCategory:     LocalBonusBlacklistCategory,
		LocalBonusWhitelistPaymentTypes: LocalBonusWhitelistPaymentTypes,
		LocalBonusBlacklistPaymentTypes: LocalBonusBlacklistPaymentTypes,
		FcyBaseRewards:                  c.FcyBaseRewards,
		FcyBaseMiles:                    c.FcyBaseMiles,
		FcyBaseWhitelistCategory:        FcyBaseWhitelistCategory,
		FcyBaseBlacklistCategory:        FcyBaseBlacklistCategory,
		FcyBonusRewards:                 c.FcyBonusRewards,
		FcyBonusMiles:                   c.FcyBonusMiles,
		FcyBonusWhitelistCategory:       FcyBonusWhitelistCategory,
		FcyBonusBlacklistCategory:       FcyBonusBlacklistCategory,
		FcyBonusWhitelistPaymentTypes:   FcyBonusWhitelistPaymentTypes,
		FcyBonusBlacklistPaymentTypes:   FcyBonusBlacklistPaymentTypes,
		Rounding:                        c.Rounding,
		AmountBlock:                     c.AmountBlock,
		RewardCurrency:                  c.RewardCurrency,
		CapType:                         c.CapType,
		Cap:                             c.Cap,
	}
	return cc
}

func createCard(c *pb.Card) (int64, error) {
	cDb := fillCardToCardDb(c)
	if err := orm.DbInstance().Table(orm.CardTable).Create(&cDb).Error; err != nil {
		return -1, err
	}

	if err := orm.DbInstance().Raw(orm.Sql1(), c.GetCardName()).Scan(&cDb).Error; err != nil {
		return -1, err
	}
	return cDb.GetCardId(), nil
}
