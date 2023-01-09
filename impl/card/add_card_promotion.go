package card

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/orm"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

func AddCardPromotion(c echo.Context) error {
	req := new(pb.AddCardPromotionRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if code, err := verifyAddCardPromotionFields(req); err != nil {
		if code == nil {
			return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
		}
		return resp.JSONResp(c, int64(code.(pb.AddCardPromotionRequest_ErrorCode)), err.Error())
	}
	return nil
}

func verifyAddCardPromotionFields(req *pb.AddCardPromotionRequest) (interface{}, error) {
	if req.CardPromotion == nil {
		return nil, errors.New("card_promotion cannot be empty")
	}

	if req.CardPromotion.PromotionName == nil {
		return nil, errors.New("card_promotion_name cannot be empty")
	}

	if len(req.GetCardPromotion().GetPromotionName()) > 50 {
		return nil, errors.New("card_promotion_name must not exceed 50 characters")
	}

	if req.CardPromotion.PromotionDescription == nil {
		return nil, errors.New("card_promotion_description cannot be empty")
	}

	if len(req.GetCardPromotion().GetPromotionDescription()) > 500 {
		return nil, errors.New("card_promotion_description must not exceed 500 characters")
	}

	if req.GetCardPromotion().EligibleCardIds == nil {
		return nil, errors.New("eligible_card_ids cannot be empty")
	}

	if req.GetCardPromotion().GetEligibleCardIds().EligibleCards == nil {
		return nil, errors.New("eligible_cards cannot be empty")
	}

	if req.GetCardPromotion().GetEligibleCardIds().IneligibleCards == nil {
		return nil, errors.New("ineligible_cards cannot be empty")
	}

	if req.GetCardPromotion().PromotionType == nil {
		return nil, errors.New("promotion_type cannot be empty")
	}

	if req.GetCardPromotion().GetPromotionType() < 0 {
		return nil, errors.New("invalid promotion_type")
	}

	if req.GetCardPromotion().StartTime == nil {
		return nil, errors.New("start_time cannot be empty")
	}

	if req.GetCardPromotion().EndTime == nil {
		return nil, errors.New("end_time cannot be empty")
	}

	if req.GetCardPromotion().GetEndTime() != 0 && req.GetCardPromotion().GetStartTime() <= req.GetCardPromotion().GetEndTime() {
		return nil, errors.New("start_time cannot be later than end_time")
	}

	if req.GetCardPromotion().PromotionConditions != nil && req.GetCardPromotion().GetPromotionConditions().CardRules != nil {
		rules := req.GetCardPromotion().GetPromotionConditions().GetCardRules()

		if rules.WhitelistCategories != nil && rules.WhitelistCategories.List != nil {
			for _, wl := range rules.WhitelistCategories.GetList() {
				if _, ok := pb.CardCategory_name[int32(wl)]; !ok {
					return pb.AddCardPromotionRequest_ERROR_CATEGORY_NOT_EXIST, errors.New("invalid whitelist_category")
				}
			}
		}

		if rules.BlacklistCategories != nil && rules.BlacklistCategories.List != nil {
			for _, bl := range rules.BlacklistCategories.GetList() {
				if _, ok := pb.CardCategory_name[int32(bl)]; !ok {
					return pb.AddCardPromotionRequest_ERROR_CATEGORY_NOT_EXIST, errors.New("invalid blacklist_category")
				}
			}
		}

		if rules.WhitelistPaymentTypes != nil && rules.WhitelistPaymentTypes.List != nil {
			for _, wlp := range rules.WhitelistPaymentTypes.GetList() {
				if _, ok := pb.CardPaymentType_name[int32(wlp)]; !ok {
					return pb.AddCardPromotionRequest_ERROR_PAYMENT_TYPE_NOT_EXIST, errors.New("invalid whitelist_payment_types")
				}
			}
		}

		if rules.BlacklistPaymentTypes != nil && rules.BlacklistPaymentTypes.List != nil {
			for _, blp := range rules.BlacklistPaymentTypes.GetList() {
				if _, ok := pb.CardPaymentType_name[int32(blp)]; !ok {
					return pb.AddCardPromotionRequest_ERROR_PAYMENT_TYPE_NOT_EXIST, errors.New("invalid blacklist_payment_types")
				}
			}
		}
	}

	if req.GetCardPromotion().PromotionConditions != nil && req.GetCardPromotion().GetPromotionConditions().MinAmount != nil && req.GetCardPromotion().GetPromotionConditions().MaxAmount != nil {
		if req.GetCardPromotion().GetPromotionConditions().GetMaxAmount() < req.GetCardPromotion().GetPromotionConditions().GetMinAmount() {
			return nil, errors.New("min_amount cannot be larger than max_amount")
		}
	}

	if req.GetCardPromotion().PromotionRewards != nil {
		if req.GetCardPromotion().GetPromotionRewards().LocalBonusRewards == nil {
			return nil, errors.New("local_bonus_rewards cannot be empty")
		}

		if req.GetCardPromotion().GetPromotionRewards().LocalBonusMiles == nil {
			return nil, errors.New("local_bonus_miles cannot be empty")
		}

		if req.GetCardPromotion().GetPromotionRewards().FcyBonusRewards == nil {
			return nil, errors.New("fcy_bonus_rewards cannot be empty")
		}

		if req.GetCardPromotion().GetPromotionRewards().FcyBonusMiles == nil {
			return nil, errors.New("fcy_bonus_miles cannot be empty")
		}
	}
	return nil, nil
}

func addCardPromotion(req *pb.AddCardPromotionRequest) error {
	if req.GetCardPromotion().PromotionConditions == nil {
		if req.GetCardPromotion().GetPromotionConditions().MinAmount == nil {
			req.CardPromotion.PromotionConditions.MinAmount = proto.Float64(0)
		}
		if req.GetCardPromotion().GetPromotionConditions().MaxAmount == nil {
			req.CardPromotion.PromotionConditions.MaxAmount = proto.Float64(0)
		}
		if req.GetCardPromotion().GetPromotionConditions().IsSGDOnly == nil {
			req.CardPromotion.PromotionConditions.IsSGDOnly = proto.Bool(true)
		}
		if req.GetCardPromotion().GetPromotionConditions().IsRecurringOnly == nil {
			req.CardPromotion.PromotionConditions.IsRecurringOnly = proto.Bool(false)
		}
	}

	if err := orm.DbInstance().Table(orm.CardPromotionTable).Create(&req).Error; err != nil {
		return err
	}

	//return card_promotion_id
	return nil
}
