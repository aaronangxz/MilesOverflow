package resp

import (
	"github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func JSONResp(c echo.Context, code int64, err error) error {
	return c.JSON(http.StatusOK, response(code, err.Error()))
}

func response(code int64, msg string) rewards_tracker.GenericResponse {
	log.Error(msg)
	return rewards_tracker.GenericResponse{
		ResponseMeta: &rewards_tracker.ResponseMeta{
			ErrorCode:    proto.Int64(code),
			ErrorMessage: proto.String(msg),
		},
	}
}
