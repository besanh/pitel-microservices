package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_app"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatApp struct {
	pb.UnimplementedChatAppServiceServer
}

func NewGRPCChatApp() *GRPCChatApp {
	return &GRPCChatApp{}
}

func (g *GRPCChatApp) PostChatApp(ctx context.Context, req *pb.PostChatAppRequest) (result *pb.PostChatAppResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	if (user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	payload := model.ChatAppRequest{}
	if err := util.ParseAnyToAny(req, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatAppService.InsertChatApp(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PostChatAppResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return
}
