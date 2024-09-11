package grpc

import (
	"context"
	"strconv"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	pb "github.com/tel4vn/pitel-microservices/gen/proto/chat_manage_queue"
	"github.com/tel4vn/pitel-microservices/middleware/auth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatManageQueue struct{}

func NewGRPCChatManageQueue() pb.ChatManageQueueServiceServer {
	return &GRPCChatManageQueue{}
}

func (g *GRPCChatManageQueue) InsertChatManageQueue(ctx context.Context, request *pb.PostChatManageQueueRequest) (result *pb.PostChatManageQueueResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	isNew, _ := strconv.ParseBool(request.GetIsNew())
	payload := model.ChatManageQueueUserRequest{
		ConnectionId: request.ConnectionId,
		QueueId:      request.QueueId,
		UserId:       request.UserId,
		IsNew:        isNew,
	}

	if err = payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ManageQueueService.PostManageQueue(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PostChatManageQueueResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return
}

func (g *GRPCChatManageQueue) UpdateChatManageQueueById(ctx context.Context, request *pb.PutChatManageQueueRequest) (result *pb.PutChatManageQueueResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	if len(request.GetId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	isNew, _ := strconv.ParseBool(request.GetIsNew())
	payload := model.ChatManageQueueUserRequest{
		ConnectionId: request.ConnectionId,
		QueueId:      request.QueueId,
		UserId:       request.UserId,
		IsNew:        isNew,
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ManageQueueService.UpdateManageQueueById(ctx, user, request.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutChatManageQueueResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatManageQueue) DeleteChatManageQueueById(ctx context.Context, request *pb.DeleteChatManageQueueRequest) (result *pb.DeleteChatManageQueueResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if len(request.GetId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	err = service.ManageQueueService.DeleteManageQueueById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.DeleteChatManageQueueResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
