package grpc

import (
	"context"
	"net/http"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/profile"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCProfile struct{}

func NewGRPCProfile() pb.ProfileServiceServer {
	return &GRPCProfile{}
}

func (g *GRPCProfile) GetUpdateProfile(ctx context.Context, request *pb.GetUpdateProfileRequest) (result *pb.GetUpdateProfileResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ProfileRequest{
		AppId:          request.AppId,
		OaId:           request.OaId,
		UserId:         request.UserId,
		ProfileType:    request.ProfileType,
		ConversationId: request.ConversationId,
	}

	if err = payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	code, responseData := service.ProfileService.GetUpdateProfileByUserId(ctx, user, payload)
	if code != http.StatusOK {
		return nil, status.Errorf(codes.Internal, response.ERR_GET_FAILED)
	}
	var data model.Conversation
	if err = util.ParseAnyToAny(responseData, &data); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.Conversation{
		TenantId:               data.TenantId,
		ConversationId:         data.ConversationId,
		ConversationType:       data.ConversationType,
		AppId:                  data.AppId,
		OaId:                   data.OaId,
		OaName:                 data.OaName,
		OaAvatar:               data.OaAvatar,
		ShareInfo:              nil,
		ExternalUserId:         data.ExternalUserId,
		Username:               data.Username,
		Avatar:                 data.Avatar,
		Major:                  data.Major,
		Following:              data.Following,
		Label:                  make([]*pb.ChatLabel, 0),
		IsDone:                 data.IsDone,
		IsDoneAt:               timestamppb.New(data.IsDoneAt),
		IsDoneBy:               data.IsDoneBy,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
		ExternalConversationId: data.ExternalConversationId,
	}
	if data.ShareInfo != nil {
		if err = util.ParseAnyToAny(data.ShareInfo, &tmp.ShareInfo); err != nil {
			log.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	if err = util.ParseStringToAny(string(data.Labels), &tmp.Label); err != nil {
		log.Error(err)
		result = &pb.GetUpdateProfileResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	result = &pb.GetUpdateProfileResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return
}
