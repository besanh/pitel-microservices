package grpc

import (
	"context"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/variables"
	pb "github.com/tel4vn/pitel-microservices/gen/proto/chat_user"
	"github.com/tel4vn/pitel-microservices/middleware/auth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"golang.org/x/exp/slices"
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

func (g *GRPCChatUser) UpdateChatUserStatusById(ctx context.Context, request *pb.PutChatUserStatusRequest) (result *pb.PutChatUserResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	if !slices.Contains(variables.USER_STATUSES, request.GetStatus()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid status: %v", request.GetStatus())
	}

	err = service.ChatUserService.UpdateChatUserStatusById(ctx, user, request.GetStatus())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutChatUserResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatUser) GetChatUserStatusById(ctx context.Context, request *pb.GetChatUserStatusRequest) (result *pb.GetChatUserResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	userStatus, err := service.ChatUserService.GetChatUserStatusById(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatUserResponse{
		Code:    "OK",
		Message: "ok",
		Status:  userStatus,
	}
	return
}
