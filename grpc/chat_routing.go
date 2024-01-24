package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_routing"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatRouting struct {
	pb.UnimplementedChatRoutingServer
}

func NewGRPCChatRouting() *GRPCChatRouting {
	return &GRPCChatRouting{}
}

func (s *GRPCChatRouting) PostChatRouting(ctx context.Context, req *pb.ChatRoutingBodyRequest) (result *pb.ChatRoutingResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.ChatRoutingRequest
	if err = util.ParseAnyToAny(req, &payload); err != nil {
		log.Error(err)
		result = &pb.ChatRoutingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}
	log.Println("payload -->", payload)

	if err = payload.Validate(); err != nil {
		log.Error(err)
		result = &pb.ChatRoutingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_VALIDATION_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	id, err := service.NewChatRouting().InsertChatRouting(ctx, authUser, &payload)
	if err != nil {
		log.Error(err)
		result = &pb.ChatRoutingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	result = &pb.ChatRoutingResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}

	return
}
