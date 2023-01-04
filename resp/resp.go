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

func PairUserCardResponseJSON(c echo.Context, card *pb.UserCardWithInfo) error {
	return c.JSON(http.StatusOK, pairUserCardResponse(card))
}

func pairUserCardResponse(card *pb.UserCardWithInfo) pb.PairUserCardResponse {
	return pb.PairUserCardResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.PairUserCardRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully paired card."),
		},
		UserCardWithInfo: card,
	}
}

func GetUserCardsResponseJSON(c echo.Context, cards []*pb.UserCardWithInfo) error {
	return c.JSON(http.StatusOK, getUserCardsResponse(cards))
}

func getUserCardsResponse(cards []*pb.UserCardWithInfo) pb.GetUserCardsResponse {
	return pb.GetUserCardsResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.GetUserCardsRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully retrieved cards."),
		},
		UserCardsList: cards,
	}
}

func GetCalculateTransactionResponseJSON(c echo.Context, trx *pb.CalculatedTransaction) error {
	return c.JSON(http.StatusOK, getCalculateTransactionResponse(trx))
}

func getCalculateTransactionResponse(trx *pb.CalculatedTransaction) pb.CalculateTransactionResponse {
	return pb.CalculateTransactionResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.CalculateTransactionRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully calculated transaction."),
		},
		CalculatedTransaction: trx,
	}
}

func AddTransactionResponseJSON(c echo.Context) error {
	return c.JSON(http.StatusOK, addTransactionResponse())
}

func addTransactionResponse() pb.AddTransactionResponse {
	return pb.AddTransactionResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.AddTransactionRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully added transaction."),
		},
	}
}

func GetUserTransactionsResponseJSON(c echo.Context, trx []*pb.TransactionBasicWithCardInfo) error {
	return c.JSON(http.StatusOK, getUserTransactionsResponse(trx))
}

func getUserTransactionsResponse(trx []*pb.TransactionBasicWithCardInfo) pb.GetUserTransactionsResponse {
	return pb.GetUserTransactionsResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.GetUserTransactionsRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully retrieved transactions."),
		},
		TransactionList: trx,
	}
}

func GetUserTransactionByTrxIdResponseJSON(c echo.Context, trx *pb.TransactionDbWithCardInfo) error {
	return c.JSON(http.StatusOK, getUserTransactionByTrxIdResponse(trx))
}

func getUserTransactionByTrxIdResponse(trx *pb.TransactionDbWithCardInfo) pb.GetUserTransactionByTrxIdResponse {
	return pb.GetUserTransactionByTrxIdResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.GetUserTransactionByTrxIdRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully retrieved transaction."),
		},
		TransactionInfo: trx,
	}
}

func GetUserCardByCardIdResponseJSON(c echo.Context, cardInfo *pb.UserCardWithInfo, trx []*pb.TransactionBasic) error {
	return c.JSON(http.StatusOK, getUserCardByCardIdResponse(cardInfo, trx))
}

func getUserCardByCardIdResponse(cardInfo *pb.UserCardWithInfo, trx []*pb.TransactionBasic) pb.GetUserCardByUserCardIdResponse {
	return pb.GetUserCardByUserCardIdResponse{
		ResponseMeta: &pb.ResponseMeta{
			ErrorCode:    proto.Int64(int64(pb.GetUserTransactionByTrxIdRequest_ERROR_SUCCESS)),
			ErrorMessage: proto.String("successfully retrieved user card."),
		},
		UserCardInfo: cardInfo,
		Transactions: trx,
	}
}
