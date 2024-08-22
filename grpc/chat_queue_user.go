package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_queue_user"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatQueueUser struct{}

func NewGRPCChatQueueUser() pb.ChatQueueUserServiceServer {
	return &GRPCChatQueueUser{}
}

func (g *GRPCChatQueueUser) InsertChatQueueUser(ctx context.Context, request *pb.PostChatQueueUserRequest) (result *pb.PostChatQueueUserResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatQueueUserRequest{
		QueueId: request.QueueId,
		UserId:  request.UserId,
		Source:  user.Source,
	}

	if err = payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ChatQueueUserService.InsertChatQueueUser(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PostChatQueueUserResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatQueueUser) UpdateChatQueueUserById(ctx context.Context, request *pb.PutChatQueueUserRequest) (result *pb.PutChatQueueUserResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatQueueUserRequest{
		QueueId: request.QueueId,
		UserId:  request.UserId,
		Source:  user.Source,
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	data, err := service.ChatQueueUserService.UpdateChatQueueUserById(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutChatQueueUserResponse{
		Code:         "OK",
		Message:      "ok",
		TotalSuccess: int32(data.TotalSuccess),
		TotalFail:    int32(data.TotalFail),
		ListFail:     data.ListFail,
		ListSuccess:  data.ListSuccess,
	}
	return
}
