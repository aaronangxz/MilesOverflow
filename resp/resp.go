package resp

import (
	pb "github.com/aaronangxz/RewardTracker/rewards_tracker.pb/rewards_tracker"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func AddCardResponseJSON(c echo.Context, id int64) error {
	return c.JSON(http.StatusOK, addCardResponse(id))
}

func addCardResponse(id int64) pb.AddCardResponse {
	return pb.AddCardResponse{
		CardId: proto.Int64(id),
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.AddCardRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully added card."),
		},
	}
}

func PairUserCardResponseJSON(c echo.Context, cards []*pb.UserCard) error {
	return c.JSON(http.StatusOK, pairUserCardResponse(cards))
}

func pairUserCardResponse(cards []*pb.UserCard) pb.PairUserCardResponse {
	return pb.PairUserCardResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.PairUserCardRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully paired card."),
		},
		UserCardsList: cards,
	}
}
