package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
	pb "github.com/tel4vn/fins-microservices/gen/proto/conversation"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCConversation struct {
	pb.UnimplementedConversationServiceServer
}

func NewGRPCConversation() pb.ConversationServiceServer {
	return &GRPCConversation{}
}

func (g *GRPCConversation) UpdateStatusConversation(ctx context.Context, request *pb.PutConversationStatusRequest) (*pb.PutConversationStatusResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	payload := model.ConversationStatusRequest{
		AppId:          request.GetAppId(),
		ConversationId: request.GetConversationId(),
		Status:         request.GetStatus(),
	}
	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ConversationService.UpdateStatusConversation(ctx, user, request.GetAppId(), request.GetConversationId(), user.UserId, request.GetStatus())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutConversationStatusResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
