package resp

import (
	"github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func AddCardResponseJSON(c echo.Context, id int64) error {
	return c.JSON(http.StatusOK, addCardResponse(id))
}

func addCardResponse(id int64) rewards_tracker.AddCardResponse {
	return rewards_tracker.AddCardResponse{
		CardId: proto.Int64(id),
		ResponseMeta: &rewards_tracker.ResponseMeta{
			ErrorCode:    proto.Int64(int64(rewards_tracker.AddCardRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully added card."),
		},
	}
}
