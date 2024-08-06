package grpc

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/message"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCMessage struct{}

func NewGRPCMessage() pb.MessageServiceServer {
	return &GRPCMessage{}
}

func (g *GRPCMessage) GetMessagesWithScrollAPI(ctx context.Context, request *pb.GetMessagesScrollRequest) (result *pb.GetMessagesScrollResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	filter := model.MessageFilter{
		ConversationId: request.GetConversationId(),
		ExternalUserId: request.GetExternalUserId(),
	}

	limit := util.ParseLimit(request.GetLimit())

	total, data, respScrollId, err := service.MessageService.GetMessagesWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.Message, 0)
	for _, item := range data {
		tmp, err := convertMessageToPbMessage(*item)
		if err != nil {
			log.Error(err)
			result = &pb.GetMessagesScrollResponse{
				Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
				Message: err.Error(),
			}
			return result, err
		}
		resultData = append(resultData, tmp)
	}

	result = &pb.GetMessagesScrollResponse{
		Code:     "OK",
		Message:  "ok",
		Data:     resultData,
		Total:    int32(total),
		Limit:    int32(limit),
		ScrollId: respScrollId,
	}
	return
}

func (g *GRPCMessage) GetMessages(ctx context.Context, request *pb.GetMessagesRequest) (result *pb.GetMessagesResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	filter := model.MessageFilter{
		ConversationId: request.GetConversationId(),
		ExternalUserId: request.GetExternalUserId(),
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.MessageService.GetMessages(ctx, user, filter, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.Message, 0)
	if data != nil {
		for _, item := range *data {
			tmp, err := convertMessageToPbMessage(item)
			if err != nil {
				log.Error(err)
				result = &pb.GetMessagesResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, err
			}
			resultData = append(resultData, tmp)
		}

	}
	result = &pb.GetMessagesResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}

func (g *GRPCMessage) SendMessage(ctx context.Context, request *pb.PostMessageRequest) (result *pb.PostMessageResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.MessageRequest{}
	if err = util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if request.EventName == "form" {
		err = errors.New("not supported form event")
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err = payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	data, err := service.MessageService.SendMessageToOTT(ctx, user, payload, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp, err := convertMessageToPbMessage(data)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	result = &pb.PostMessageResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return
}

func (g *GRPCMessage) MarkReadMessages(ctx context.Context, request *pb.MarkReadMessagesRequest) (result *pb.MarkReadMessagesResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.MessageMarkRead{}
	if err = util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err = payload.ValidateMarkRead(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	data, err := service.MessageService.MarkReadMessages(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.MarkReadMessagesResponse{
		Code:         "OK",
		Message:      "ok",
		TotalSuccess: int32(data.TotalSuccess),
		TotalFail:    int32(data.TotalFail),
		ListFail:     data.ListFail,
		ListSuccess:  data.ListSuccess,
	}
	return
}

func (g *GRPCMessage) ShareInfo(ctx context.Context, request *pb.ShareInfoRequest) (result *pb.ShareInfoResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ShareInfo{}
	if err = util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	data, err := service.MessageService.ShareInfo(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.ShareInfoResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data.Data,
	}
	return
}

func (g *GRPCMessage) GetMessageMediasWithScrollAPI(ctx context.Context, request *pb.GetMessageMediasScrollRequest) (result *pb.GetMessageMediasScrollResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	filter := model.MessageFilter{
		ConversationId: request.GetConversationId(),
		AttachmentType: request.GetAttachmentType(),
	}

	limit := util.ParseLimit(request.GetLimit())

	total, data, respScrollId, err := service.MessageService.GetMessageMediasWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.MessageAttachmentsDetails, 0)
	for _, item := range data {
		tmp, err := convertMessageAttachmentToPbMessageAttachment(*item)
		if err != nil {
			log.Error(err)
			result = &pb.GetMessageMediasScrollResponse{
				Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
				Message: err.Error(),
			}
			return result, err
		}
		resultData = append(resultData, tmp)
	}

	result = &pb.GetMessageMediasScrollResponse{
		Code:     "OK",
		Message:  "ok",
		Data:     resultData,
		Total:    int32(total),
		Limit:    int32(limit),
		ScrollId: respScrollId,
	}
	return
}
