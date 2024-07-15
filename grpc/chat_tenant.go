package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_tenant"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatTenant struct {
	pb.UnimplementedChatTenantServiceServer
}

func NewGRPCChatTenant() *GRPCChatTenant {
	return &GRPCChatTenant{}
}

func (g *GRPCChatTenant) PostChatTenant(ctx context.Context, req *pb.PostChatTenantRequest) (result *pb.PostChatTenantResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if (user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	payload := model.ChatTenantRequest{}
	if err := util.ParseAnyToAny(req, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatTenantService.InsertChatTenant(ctx, &payload)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PostChatTenantResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}

	return
}
