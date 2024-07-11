package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_role"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatRole struct {
	pb.UnimplementedChatRoleServiceServer
}

func NewChatRole() *GRPCChatRole {
	return &GRPCChatRole{}
}

func (g *GRPCChatRole) PostChatRole(ctx context.Context, request *pb.PostChatRoleRequest) (result *pb.PostChatRoleResponse, err error) {
	_, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	setting := model.ChatRoleSetting{
		AssignConversation:   request.GetSetting().GetAssignConversation(),
		ReassignConversation: request.GetSetting().GetReassignConversation(),
		MakeDone:             request.GetSetting().GetMakeDone(),
		AddLabel:             request.GetSetting().GetAddLabel(),
		RemoveLabel:          request.GetSetting().GetRemoveLabel(),
		Major:                request.GetSetting().GetMajor(),
		Following:            request.GetSetting().GetFollowing(),
		SubmitForm:           request.GetSetting().GetSubmitForm(),
	}
	payload := model.ChatRoleRequest{
		RoleName: request.GetRoleName(),
		Status:   request.GetStatus(),
		Setting:  setting,
	}

	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatRoleService.InsertChatRole(ctx, &payload)
	if err != nil {
		result = &pb.PostChatRoleResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.PostChatRoleResponse{
		Code:    "OK",
		Message: "OK",
		Id:      id,
	}

	return
}
