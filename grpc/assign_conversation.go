package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/assign_conversation"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCAssignConversation struct {
	pb.UnimplementedAssignConversationServiceServer
}

func NewGRPCAssignConversation() *GRPCAssignConversation {
	return &GRPCAssignConversation{}
}

func (g *GRPCChatApp) InsertUserInQueue(ctx context.Context, request *pb.PostUserInQueueRequest) (*pb.PostUserInQueueResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.AssignConversation{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	code, data := service.AssignConversationService.AllocateConversation(ctx, user, &payload)
	if code < 200 || code >= 300 {
		return nil, status.Errorf(codes.Internal, "assign conversation failed, code: %d", code)
	}
	tmp, err := util.ToStructPb(data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostUserInQueueResponse{
		Data:    tmp,
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatApp) GetUserAssigned(ctx context.Context, request *pb.GetUserAssignedRequest) (*pb.GetUserAssignedResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	conversationId := request.GetId()
	statusTmp := request.GetStatus()
	code, data := service.AssignConversationService.GetUserAssigned(ctx, user, conversationId, statusTmp)
	if code < 200 || code >= 300 {
		return nil, status.Errorf(codes.Internal, "assign conversation failed, code: %d", code)
	}
	tmp, err := util.ToStructPb(data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetUserAssignedResponse{
		Data:    tmp,
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatApp) GetUserInQueue(ctx context.Context, request *pb.GetUserInQueueRequest) (*pb.GetUserInQueueResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	filter := model.UserInQueueFilter{
		AppId:            request.GetAppId(),
		OaId:             request.GetOaId(),
		ConversationId:   request.GetConversationId(),
		ConversationType: request.GetConversationType(),
		Status:           request.GetStatus(),
	}
	code, data := service.AssignConversationService.GetUserInQueue(ctx, user, filter)
	if code < 200 || code >= 300 {
		return nil, status.Errorf(codes.Internal, "assign conversation failed, code: %d", code)
	}
	tmp, err := util.ToStructPb(data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetUserInQueueResponse{
		Data:    tmp,
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
