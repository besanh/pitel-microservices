package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/inbox_marketing"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCInboxMarketing struct {
	pb.UnsafeInboxMarketingServiceServer
}

func NewGRPCInboxMarketing() *GRPCInboxMarketing {
	return &GRPCInboxMarketing{}
}

func (g *GRPCInboxMarketing) SendInboxMarketing(ctx context.Context, request *pb.InboxMarketingRequestRequest) (result *pb.InboxMarketingResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.InboxMarketingRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.InboxMarketingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Status:  "failed",
			Message: err.Error(),
		}
		return result, err
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		result = &pb.InboxMarketingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Status:  "failed",
			Message: err.Error(),
		}
		return result, nil
	}

	res, err := service.NewInboxMarketing().SendInboxMarketing(ctx, authUser, payload)
	if err != nil {
		result = &pb.InboxMarketingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Status:  res.Status,
			Message: res.Message,
		}
		return result, nil
	}

	result = &pb.InboxMarketingResponse{
		Code:    res.Code,
		Status:  res.Status,
		Message: res.Message,
	}

	return result, nil
}
