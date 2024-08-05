package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/conversation_media"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCConversationMedia struct{}

func NewGRPCConversationMedia() pb.ConversationMediaServiceServer {
	return &GRPCConversationMedia{}
}

func (g *GRPCConversationMedia) GetConversationMedias(ctx context.Context, request *pb.GetConversationMediasRequest) (result *pb.GetConversationMediasResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	filter := model.ConversationMediaFilter{
		TenantId:               request.GetTenantId(),
		ConversationId:         request.GetConversationId(),
		ExternalConversationId: request.GetExternalConversationId(),
		ConversationType:       request.GetConversationType(),
		MediaType:              request.GetMediaType(),
		MediaName:              request.GetMediaName(),
		SendTimestamp:          request.GetSendTimestamp(),
	}
	if err = filter.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.ConversationMediaService.GetConversationMedias(ctx, user, filter, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ConversationMedia, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			tmp := &pb.ConversationMedia{
				Id:                     item.Id,
				CreatedAt:              timestamppb.New(item.CreatedAt),
				UpdatedAt:              timestamppb.New(item.UpdatedAt),
				TenantId:               item.TenantId,
				ConversationId:         item.ConversationId,
				ExternalConversationId: item.ExternalConversationId,
				ConversationType:       item.ConversationType,
				MessageId:              item.MessageId,
				MediaType:              item.MediaType,
				MediaHeader:            item.MediaHeader,
				MediaUrl:               item.MediaUrl,
				MediaSize:              item.MediaSize,
				SendTimestamp:          timestamppb.New(item.SendTimestamp),
			}

			resultData = append(resultData, tmp)
		}
	}

	result = &pb.GetConversationMediasResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}
