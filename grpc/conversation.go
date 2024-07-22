package grpc

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	pb "github.com/tel4vn/fins-microservices/gen/proto/conversation"
	"github.com/tel4vn/fins-microservices/middleware/auth"
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

	if len(request.GetStatus()) < 1 {
		log.Error("status is required")
		return nil, status.Errorf(codes.InvalidArgument, errors.New("status is required").Error())
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
