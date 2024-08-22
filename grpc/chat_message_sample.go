package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_message_sample"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatMessageSample struct {
	pb.UnimplementedMessageSampleServiceServer
}

func NewGRPCChatMessageSample() *GRPCChatMessageSample {
	return &GRPCChatMessageSample{}
}

func (g *GRPCChatMessageSample) InsertMessageSample(ctx context.Context, request *pb.PostMessageSampleRequest) (*pb.PostMessageSampleResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatMsgSampleRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatMessageSampleService.InsertChatMsgSample(ctx, user, payload, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostMessageSampleResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatMessageSample) GetChatMessageSamples(ctx context.Context, request *pb.GetChatMessageSamplesRequest) (*pb.GetChatMessageSamplesResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatMsgSampleFilter{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	limit, offset := request.GetLimit(), request.GetOffset()

	total, data, err := service.ChatMessageSampleService.GetChatMsgSamples(ctx, user, payload, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := make([]*pb.ChatMessageSampleData, 0)
	for _, item := range *data {
		dataItem := &pb.ChatMessageSampleData{
			Id:            item.GetId(),
			CreatedAt:     timestamppb.New(item.CreatedAt),
			UpdatedAt:     timestamppb.New(item.UpdatedAt),
			TenantId:      item.TenantId,
			Keyword:       item.Keyword,
			Theme:         item.Theme,
			ConnectionId:  item.ConnectionId,
			ConnectionApp: nil,
			Channel:       item.Channel,
			Content:       item.Content,
			CreatedBy:     item.CreatedBy,
			UpdatedBy:     item.UpdatedBy,
			ImageUrl:      item.ImageUrl,
		}
		if item.ConnectionApp != nil {
			dataItem.ConnectionApp = &pb.ChatConnectionApp{
				Id:                item.ConnectionApp.Id,
				TenantId:          item.ConnectionApp.TenantId,
				ConnectionName:    item.ConnectionApp.ConnectionName,
				ConnectionType:    item.ConnectionApp.ConnectionType,
				ConnectionQueueId: item.ConnectionApp.ConnectionQueueId,
				ChatAppId:         item.ConnectionApp.ChatAppId,
				Status:            item.ConnectionApp.Status,
				CreatedAt:         timestamppb.New(item.ConnectionApp.CreatedAt),
				UpdatedAt:         timestamppb.New(item.ConnectionApp.UpdatedAt),
				OaInfo:            nil,
			}
			if err = util.ParseAnyToAny(item.ConnectionApp.OaInfo, &dataItem.ConnectionApp.OaInfo); err != nil {
				log.Error(err)
				result := &pb.GetChatMessageSamplesResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, err
			}
		}
		tmp = append(tmp, dataItem)
	}

	result := &pb.GetChatMessageSamplesResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
		Total:   int32(total),
		Limit:   limit,
		Offset:  offset,
	}
	return result, nil
}

func (g *GRPCChatMessageSample) GetMessageSampleById(ctx context.Context, request *pb.GetMessageSampleByIdRequest) (*pb.GetMessageSampleByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatMessageSampleService.GetChatMsgSampleById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatMessageSampleData{
		Id:           data.GetId(),
		CreatedAt:    timestamppb.New(data.CreatedAt),
		UpdatedAt:    timestamppb.New(data.UpdatedAt),
		TenantId:     data.TenantId,
		Keyword:      data.Keyword,
		Theme:        data.Theme,
		ConnectionId: data.ConnectionId,
		ConnectionApp: &pb.ChatConnectionApp{
			Id:        data.ConnectionApp.Id,
			Status:    data.ConnectionApp.Status,
			CreatedAt: timestamppb.New(data.ConnectionApp.CreatedAt),
			UpdatedAt: timestamppb.New(data.ConnectionApp.UpdatedAt),
		},
		Channel:   data.Channel,
		Content:   data.Content,
		CreatedBy: data.CreatedBy,
		UpdatedBy: data.UpdatedBy,
		ImageUrl:  data.ImageUrl,
	}
	if data.ConnectionApp != nil {
		if err = util.ParseAnyToAny(data.ConnectionApp.OaInfo, &tmp.ConnectionApp.OaInfo); err != nil {
			log.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	result := &pb.GetMessageSampleByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatMessageSample) UpdateMessageSampleById(ctx context.Context, request *pb.PutMessageSampleRequest) (*pb.PutMessageSampleResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatMsgSampleRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatMessageSampleService.UpdateChatMsgSampleById(ctx, user, request.GetId(), payload, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutMessageSampleResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatMessageSample) DeleteMessageSampleById(ctx context.Context, request *pb.DeleteMessageSampleRequest) (*pb.DeleteMessageSampleResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err := service.ChatMessageSampleService.DeleteChatMsgSampleById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteMessageSampleResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
