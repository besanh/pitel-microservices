package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_queue"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatQueue struct {
	pb.UnimplementedChatQueueServiceServer
}

func NewGRPCChatQueue() *GRPCChatQueue {
	return &GRPCChatQueue{}
}

func (g *GRPCChatQueue) InsertChatQueue(ctx context.Context, req *pb.ChatQueueBodyRequest) (result *pb.ChatQueueResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.ChatQueueRequest
	if err = util.ParseAnyToAny(req, &payload); err != nil {
		result = &pb.ChatQueueResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	id, err := service.NewChatQueue().InsertChatQueue(ctx, authUser, payload)
	if err != nil {
		result = &pb.ChatQueueResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	result = &pb.ChatQueueResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}

	return
}
