package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_connection_pipeline"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatConnectionPipeline struct{}

func NewGRPCChatConnectionPipeline() pb.ChatConnectionPipelineServiceServer {
	return &GRPCChatConnectionPipeline{}
}

func (g *GRPCChatConnectionPipeline) AttachConnectionQueueToApp(ctx context.Context, request *pb.ChatConnectionPipelineQueueRequest) (*pb.ChatConnectionPipelineQueueResponse, error) {
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

	id, err := service.ChatConnectionPipelineService.AttachConnectionQueueToApp(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.ChatConnectionPipelineQueueResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}
