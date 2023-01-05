package card

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/resp"
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
)

func AddCardPromotion(c echo.Context) error {
	req := new(pb.AddCardPromotionRequest)
	if err := c.Bind(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_JSON_BIND), err.Error())
	}

	if err := verifyAddCardPromotionFields(req); err != nil {
		return resp.JSONResp(c, int64(pb.ErrorCode_ERROR_PARAMS), err.Error())
	}
	return nil
}

func verifyAddCardPromotionFields(req *pb.AddCardPromotionRequest) error {
	if req.CardPromotion == nil {
		return errors.New("card_promotion cannot be empty")
	}

	if req.CardPromotion.PromotionName == nil {
		return errors.New("card_promotion_name cannot be empty")
	}

	if len(req.GetCardPromotion().GetPromotionName()) > 50 {
		return errors.New("card_promotion_name must not exceed 50 characters")
	}

	if req.CardPromotion.PromotionDescription == nil {
		return errors.New("card_promotion_description cannot be empty")
	}

	if len(req.GetCardPromotion().GetPromotionDescription()) > 500 {
		return errors.New("card_promotion_description must not exceed 500 characters")
	}

	if req.GetCardPromotion().EligibleCardIds == nil {
		return errors.New("eligible_card_ids cannot be empty")
	}

	if req.GetCardPromotion().GetEligibleCardIds().EligibleCards == nil {
		return errors.New("eligible_cards cannot be empty")
	}

	if req.GetCardPromotion().GetEligibleCardIds().IneligibleCards == nil {
		return errors.New("ineligible_cards cannot be empty")
	}

	if req.GetCardPromotion().PromotionType == nil {
		return errors.New("promotion_type cannot be empty")
	}

	if req.GetCardPromotion().GetPromotionType() < 0 {
		return errors.New("invalid promotion_type")
	}

	if req.GetCardPromotion().StartTime == nil {
		return errors.New("start_time cannot be empty")
	}

	if req.GetCardPromotion().EndTime == nil {
		return errors.New("end_time cannot be empty")
	}

	if req.GetCardPromotion().GetEndTime() != 0 && req.GetCardPromotion().GetStartTime() <= req.GetCardPromotion().GetEndTime() {
		return errors.New("start_time cannot be later than end_time")

	}

	//TODO verify nested fields
	return nil
}
