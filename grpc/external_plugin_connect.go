package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/external_plugin"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCExternalPluginConnect struct {
	pb.UnimplementedExternalPluginConnectServiceServer
}

func NewGRPCExternalPluginConnect() *GRPCExternalPluginConnect {
	return &GRPCExternalPluginConnect{}
}

func (g *GRPCExternalPluginConnect) PostCreateConnect(ctx context.Context, request *pb.ExternalPluginConnectBodyRequest) (result *pb.ExternalPluginConnectResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	var payload model.ExternalPluginConnect
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.ExternalPluginConnectResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}
	log.Println("payload -->", payload)
	if err := payload.Validate(); err != nil {
		result = &pb.ExternalPluginConnectResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	if err := service.NewExternalPluginConnect().InsertExternalPluginConnect(ctx, authUser, payload); err != nil {
		result = &pb.ExternalPluginConnectResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	result = &pb.ExternalPluginConnectResponse{
		Code:    "OK",
		Message: "success",
	}

	return
}
