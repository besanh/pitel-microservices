package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_policy_setting"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatPolicySetting struct{}

func NewGRPCChatPolicySetting() pb.ChatPolicySettingServiceServer {
	return &GRPCChatPolicySetting{}
}

func (g *GRPCChatPolicySetting) InsertChatPolicySetting(ctx context.Context, request *pb.PostChatPolicySettingRequest) (*pb.PostChatPolicySettingResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatPolicyConfigRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatPolicySettingService.InsertChatPolicySetting(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostChatPolicySettingResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatPolicySetting) GetChatPolicySettings(ctx context.Context, request *pb.GetChatPolicySettingsRequest) (*pb.GetChatPolicySettingsResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatPolicyFilter{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.ChatPolicySettingService.GetChatPolicySettings(ctx, user, payload, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatPolicySetting, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			var tmp pb.ChatPolicySetting
			tmp.CreatedAt = timestamppb.New(item.CreatedAt)
			tmp.UpdatedAt = timestamppb.New(item.UpdatedAt)

			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				result := &pb.GetChatPolicySettingsResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, nil
			}
			resultData = append(resultData, &tmp)
		}
	}

	result := &pb.GetChatPolicySettingsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return result, nil
}

func (g *GRPCChatPolicySetting) GetChatPolicySettingById(ctx context.Context, request *pb.GetChatPolicySettingByIdRequest) (*pb.GetChatPolicySettingByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatPolicySettingService.GetChatPolicySettingById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatPolicySetting{}
	tmp.CreatedAt = timestamppb.New(data.CreatedAt)
	tmp.UpdatedAt = timestamppb.New(data.UpdatedAt)
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetChatPolicySettingByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatPolicySetting) UpdateChatPolicySettingById(ctx context.Context, request *pb.PutChatPolicySettingRequest) (*pb.PutChatPolicySettingResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatPolicyConfigRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatPolicySettingService.UpdateChatPolicySettingById(ctx, user, request.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatPolicySettingResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatPolicySetting) DeleteChatPolicySettingById(ctx context.Context, request *pb.DeleteChatPolicySettingRequest) (*pb.DeleteChatPolicySettingResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if len(request.GetId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, response.ERR_DELETE_FAILED)
	}

	err := service.ChatPolicySettingService.DeleteChatPolicySettingById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteChatPolicySettingResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
