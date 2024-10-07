package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_notify_message"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatNotifyMessage struct{}

func NewGRPCChatNotifyMessage() pb.NotifyMessageServiceServer {
	return &GRPCChatNotifyMessage{}
}

func (g *GRPCChatNotifyMessage) InsertNotifyMessage(ctx context.Context, request *pb.PostNotifyMessageRequest) (*pb.PostNotifyMessageResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatNotifyMessageRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatNotifyMessageService.InsertChatNotifyMessage(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostNotifyMessageResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatNotifyMessage) GetNotifyMessages(ctx context.Context, request *pb.GetNotifyMessagesRequest) (*pb.GetNotifyMessagesResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatNotifyMessageFilter{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	limit, offset := request.GetLimit(), request.GetOffset()

	total, data, err := service.ChatNotifyMessageService.GetChatNotifyMessages(ctx, user, payload, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatNotifyMessageData, 0)
	if err = util.ParseAnyToAny(data, &resultData); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetNotifyMessagesResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   limit,
		Offset:  offset,
	}
	return result, nil
}

func (g *GRPCChatNotifyMessage) GetNotifyMessageById(ctx context.Context, request *pb.GetNotifyMessageByIdRequest) (*pb.GetNotifyMessageByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatNotifyMessageService.GetChatNotifyMessageById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatNotifyMessageData{}
	tmp.CreatedAt = timestamppb.New(data.CreatedAt)
	tmp.UpdatedAt = timestamppb.New(data.UpdatedAt)
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetNotifyMessageByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatNotifyMessage) UpdateNotifyMessageById(ctx context.Context, request *pb.PutNotifyMessageRequest) (*pb.PutNotifyMessageResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatNotifyMessageRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatNotifyMessageService.UpdateChatNotifyMessageById(ctx, user, request.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutNotifyMessageResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatNotifyMessage) DeleteNotifyMessageById(ctx context.Context, request *pb.DeleteNotifyMessageRequest) (*pb.DeleteNotifyMessageResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err := service.ChatNotifyMessageService.DeleteChatNotifyMessageById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteNotifyMessageResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
