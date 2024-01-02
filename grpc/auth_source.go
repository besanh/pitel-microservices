package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/auth_source"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCAuthSoure struct {
	pb.UnimplementedAuthSourceServiceServer
}

func NewGRPCAuthSoure() *GRPCAuthSoure {
	return &GRPCAuthSoure{}
}

func (g *GRPCAuthSoure) PostAuthSource(ctx context.Context, req *pb.AuthSourceBodyRequest) (result *pb.AuthSourceResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.AuthSource
	if err := util.ParseAnyToAny(req, &payload); err != nil {
		result = &pb.AuthSourceResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}
	err = service.NewAuthSource().InsertAuthSource(ctx, authUser, payload)
	if err != nil {
		result = &pb.AuthSourceResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.AuthSourceResponse{
		Code:    "OK",
		Message: "ok",
	}

	return result, nil
}
