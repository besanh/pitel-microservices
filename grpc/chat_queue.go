package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_queue"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatQueue struct{}

func NewGRPCChatQueue() pb.ChatQueueServiceServer {
	return &GRPCChatQueue{}
}

func (g *GRPCChatQueue) InsertChatQueue(ctx context.Context, request *pb.PostChatQueueRequest) (*pb.PostChatQueueResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.AttachConnectionQueueToConnectionAppRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatQueueService.InsertChatQueueV2(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostChatQueueResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatQueue) GetChatQueues(ctx context.Context, request *pb.GetChatQueuesRequest) (*pb.GetChatQueuesResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.QueueFilter{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.ChatQueueService.GetChatQueues(ctx, user, payload, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatQueue, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			var tmp pb.ChatQueue
			tmp.CreatedAt = &timestamppb.Timestamp{
				Seconds: item.CreatedAt.Unix(),
			}
			tmp.UpdatedAt = &timestamppb.Timestamp{
				Seconds: item.UpdatedAt.Unix(),
			}

			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				result := &pb.GetChatQueuesResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, nil
			}
			for i := range item.ConnectionQueue {
				if tmp.ChatConnectionQueue[i] != nil && tmp.ChatConnectionQueue[i].ChatConnectionApp != nil {
					tmp.ChatConnectionQueue[i].ChatConnectionApp.CreatedAt = &timestamppb.Timestamp{
						Seconds: item.ConnectionQueue[i].ChatConnectionApp.CreatedAt.Unix(),
					}
					tmp.ChatConnectionQueue[i].ChatConnectionApp.UpdatedAt = &timestamppb.Timestamp{
						Seconds: item.ConnectionQueue[i].ChatConnectionApp.UpdatedAt.Unix(),
					}
				}
			}
			resultData = append(resultData, &tmp)
		}
	}

	result := &pb.GetChatQueuesResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return result, nil
}

func (g *GRPCChatQueue) GetChatQueueById(ctx context.Context, request *pb.GetChatQueueByIdRequest) (*pb.GetChatQueueByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatQueueService.GetChatQueueById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatQueue{}
	tmp.CreatedAt = &timestamppb.Timestamp{
		Seconds: data.CreatedAt.Unix(),
	}
	tmp.UpdatedAt = &timestamppb.Timestamp{
		Seconds: data.UpdatedAt.Unix(),
	}
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetChatQueueByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatQueue) UpdateChatQueueById(ctx context.Context, request *pb.PutChatQueueRequest) (*pb.PutChatQueueResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.AttachConnectionQueueToConnectionAppRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatQueueService.UpdateChatQueueByIdV2(ctx, user, request.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatQueueResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatQueue) DeleteChatQueueById(ctx context.Context, request *pb.DeleteChatQueueRequest) (*pb.DeleteChatQueueResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err := service.ChatQueueService.DeleteChatQueueByIdV2(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteChatQueueResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
