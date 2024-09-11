package grpc

import (
	"context"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/util"
	pb "github.com/tel4vn/pitel-microservices/gen/proto/assign_conversation"
	"github.com/tel4vn/pitel-microservices/middleware/auth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCAssignConversation struct{}

func NewGRPCAssignConversation() pb.AssignConversationServiceServer {
	return &GRPCAssignConversation{}
}

func (g *GRPCAssignConversation) InsertUserInQueue(ctx context.Context, request *pb.PostUserInQueueRequest) (*pb.PostUserInQueueResponse, error) {
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

	err := service.AssignConversationService.AllocateConversation(ctx, user, &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostUserInQueueResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCAssignConversation) GetUserAssigned(ctx context.Context, request *pb.GetUserAssignedRequest) (*pb.GetUserAssignedResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	conversationId := request.GetId()
	statusTmp := request.GetStatus()
	data, err := service.AssignConversationService.GetUserAssigned(ctx, user, conversationId, statusTmp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	result := &pb.GetUserAssignedResponse{
		Code:    "OK",
		Message: "ok",
	}
	if data == nil {
		result.DataDetail = &pb.GetUserAssignedResponse_NoData{NoData: &emptypb.Empty{}}
	} else {
		var tmp pb.AllocateUserData
		err = util.ParseAnyToAny(data, &tmp)
		if err != nil {
			log.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		tmp.CreatedAt = timestamppb.New(data.CreatedAt)
		tmp.UpdatedAt = timestamppb.New(data.UpdatedAt)
		result.DataDetail = &pb.GetUserAssignedResponse_Data{Data: &tmp}
	}
	return result, nil
}

func (g *GRPCAssignConversation) GetUserInQueue(ctx context.Context, request *pb.GetUserInQueueRequest) (result *pb.GetUserInQueueResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		err = status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
		return
	}

	filter := model.UserInQueueFilter{
		AppId:            request.GetAppId(),
		OaId:             request.GetOaId(),
		ConversationId:   request.GetConversationId(),
		ConversationType: request.GetConversationType(),
		Status:           request.GetStatus(),
	}
	data, err := service.AssignConversationService.GetUserInQueue(ctx, user, filter)
	if err != nil {
		return
	}
	var tmp []*pb.ChatQueueUserData
	err = util.ParseAnyToAny(data, &tmp)
	if err != nil {
		log.Error(err)
		return
	}

	result = &pb.GetUserInQueueResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return
}
