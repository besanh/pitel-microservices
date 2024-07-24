package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_connection_queue"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatConnectionQueue struct{}

func NewGRPCChatConnectionQueue() pb.ChatConnectionQueueServiceServer {
	return &GRPCChatConnectionQueue{}
}

func (g *GRPCChatConnectionQueue) GetChatConnectionQueueById(ctx context.Context, request *pb.GetChatConnectionQueueByIdRequest) (result *pb.GetChatConnectionQueueByIdResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if len(request.GetId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, response.ERR_GET_FAILED)
	}

	data, err := service.ChatConnectionQueueService.GetChatConnectionQueueById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if data == nil {
		result = &pb.GetChatConnectionQueueByIdResponse{
			Code:    "OK",
			Message: "ok",
			Data:    nil,
		}
		return
	}
	tmp := &pb.ConnectionQueue{
		Id:                data.Id,
		CreatedAt:         timestamppb.New(data.CreatedAt),
		UpdatedAt:         timestamppb.New(data.UpdatedAt),
		TenantId:          data.TenantId,
		ConnectionId:      data.ConnectionId,
		ChatConnectionApp: 0,
		QueueId:           data.QueueId,
		ChatQueue:         0,
	}

	result = &pb.GetChatConnectionQueueByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return
}
