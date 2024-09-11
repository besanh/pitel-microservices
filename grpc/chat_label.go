package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/util"
	pb "github.com/tel4vn/pitel-microservices/gen/proto/chat_label"
	"github.com/tel4vn/pitel-microservices/middleware/auth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatLabel struct{}

func NewGRPCChatLabel() pb.ChatLabelServiceServer {
	return &GRPCChatLabel{}
}

func (g *GRPCChatLabel) InsertChatLabel(ctx context.Context, request *pb.PostChatLabelRequest) (result *pb.PostChatLabelResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatLabelRequest{}
	if err = util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err = payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatLabelService.InsertChatLabel(ctx, user, &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PostChatLabelResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return
}

func (g *GRPCChatLabel) GetChatLabels(ctx context.Context, request *pb.GetChatLabelsRequest) (result *pb.GetChatLabelsResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var labelStatus sql.NullBool
	if len(request.GetStatus()) > 0 {
		labelStatusTmp, _ := strconv.ParseBool(request.GetStatus())
		labelStatus.Valid = true
		labelStatus.Bool = labelStatusTmp
	}
	var isSearchExactly sql.NullBool
	if len(request.GetIsSearchExactly()) > 0 {
		isSearchExactlyTmp, _ := strconv.ParseBool(request.GetIsSearchExactly())
		isSearchExactly.Valid = true
		isSearchExactly.Bool = isSearchExactlyTmp
	}
	payload := model.ChatLabelFilter{
		AppId:           request.GetAppId(),
		OaId:            request.GetOaId(),
		LabelType:       request.GetLabelType(),
		LabelName:       request.GetLabelName(),
		LabelColor:      request.GetLabelColor(),
		LabelStatus:     labelStatus,
		ExternalLabelId: request.GetExternalLabelId(),
		IsSearchExactly: isSearchExactly,
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.ChatLabelService.GetChatLabels(ctx, user, payload, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatLabel, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			var tmp pb.ChatLabel
			tmp.CreatedAt = timestamppb.New(item.CreatedAt)
			tmp.UpdatedAt = timestamppb.New(item.UpdatedAt)

			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				result = &pb.GetChatLabelsResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, nil
			}
			resultData = append(resultData, &tmp)
		}
	}

	result = &pb.GetChatLabelsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}

func (g *GRPCChatLabel) GetChatLabelById(ctx context.Context, request *pb.GetChatLabelByIdRequest) (result *pb.GetChatLabelByIdResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatLabelService.GetChatLabelById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatLabel{}
	tmp.CreatedAt = timestamppb.New(data.CreatedAt)
	tmp.UpdatedAt = timestamppb.New(data.UpdatedAt)
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatLabelByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return
}

func (g *GRPCChatLabel) UpdateChatLabelById(ctx context.Context, request *pb.PutChatLabelRequest) (result *pb.PutChatLabelResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatLabelRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ChatLabelService.UpdateChatLabelById(ctx, user, request.GetId(), &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutChatLabelResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatLabel) DeleteChatLabelById(ctx context.Context, request *pb.DeleteChatLabelRequest) (result *pb.DeleteChatLabelResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if len(request.GetId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	err = service.ChatLabelService.DeleteChatLabelById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.DeleteChatLabelResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
