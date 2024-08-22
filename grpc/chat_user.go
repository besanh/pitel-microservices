package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_user"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatUser struct {
	pb.UnimplementedChatUserServiceServer
}

func NewGRPCChatUser() *GRPCChatUser {
	return &GRPCChatUser{}
}

func (g *GRPCChatUser) PostChatUser(ctx context.Context, request *pb.PostChatUserRequest) (result *pb.PostChatUserResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if (user.GetLevel() != "admin" || user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		log.Error("user level ", user.GetLevel(), " is not admin or superadmin")
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	payload := model.ChatUserRequest{
		Username: request.GetUsername(),
		Password: request.GetPassword(),
		Email:    request.GetEmail(),
		Level:    request.GetLevel(),
		Fullname: request.GetFullName(),
		RoleId:   request.GetRoleId(),
		Status:   request.GetStatus(),
	}
	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatUserService.InsertChatUser(ctx, &payload)
	if err != nil {
		result = &pb.PostChatUserResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.PostChatUserResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}

	return
}
