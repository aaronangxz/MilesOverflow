package processors

import (
	"errors"
	"github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddCard(c echo.Context) error {
	req := new(rewards_tracker.AddCardRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, req.CardDetails)
	//return c.JSON(http.StatusOK, rewards_tracker.AddCardResponse{
	//	CardId: proto.Int64(0),
	//	ResponseMeta: &rewards_tracker.ResponseMeta{
	//		ErrorCode:    proto.Int64(0),
	//		ErrorMessage: proto.String("ok"),
	//	},
	//})
}

func verifyFields() error {
	return errors.New("")
}

func isCardExists() bool {
	return false
}
