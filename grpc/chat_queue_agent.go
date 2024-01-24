package grpc

import (
	"context"
	"log"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_queue_agent"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatQueueAgent struct {
	pb.UnimplementedQueueAgentServiceServer
}

func NewGRPCChatQueueAgent() *GRPCChatQueueAgent {
	return &GRPCChatQueueAgent{}
}

func (g *GRPCChatQueueAgent) InsertQueueAgent(ctx context.Context, req *pb.QueueAgentBodyRequest) (result *pb.QueueAgentResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.ChatQueueAgentRequest
	if err = util.ParseAnyToAny(req, &payload); err != nil {
		result = &pb.QueueAgentResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	log.Println("payload -->", payload)
	err = service.NewChatQueueAgent().InsertChatQueueAgent(ctx, authUser, payload)
	if err != nil {
		result = &pb.QueueAgentResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	result = &pb.QueueAgentResponse{
		Code:    "OK",
		Message: "ok",
	}

	return
}
