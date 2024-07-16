package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_auto_script"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatAutoScript struct{}

func NewGRPCChatAutoScript() pb.ChatAutoScriptServiceServer {
	return &GRPCChatAutoScript{}
}

func (g *GRPCChatAutoScript) InsertChatAutoScript(ctx context.Context, request *pb.PostChatAutoScriptRequest) (*pb.PostChatAutoScriptResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatAutoScriptRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatAutoScriptService.InsertChatAutoScript(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostChatAutoScriptResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatAutoScript) GetChatAutoScripts(ctx context.Context, request *pb.GetChatAutoScriptsRequest) (*pb.GetChatAutoScriptsResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatAutoScriptFilter{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	limit, offset := request.GetLimit(), request.GetOffset()
	statusTmp := request.GetStatus()
	var scriptStatus sql.NullBool
	if len(statusTmp) > 0 {
		tmp, _ := strconv.ParseBool(statusTmp)
		scriptStatus.Valid = true
		scriptStatus.Bool = tmp
	}
	payload.Status = scriptStatus

	total, data, err := service.ChatAutoScriptService.GetChatAutoScripts(ctx, user, payload, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatAutoScriptData, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			var tmp pb.ChatAutoScriptData
			tmp.CreatedAt = &timestamppb.Timestamp{
				Seconds: item.CreatedAt.Unix(),
			}
			tmp.UpdatedAt = &timestamppb.Timestamp{
				Seconds: item.UpdatedAt.Unix(),
			}
			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				result := &pb.GetChatAutoScriptsResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, nil
			}
			resultData = append(resultData, &tmp)
		}
	}

	result := &pb.GetChatAutoScriptsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   limit,
		Offset:  offset,
	}
	return result, nil
}

func (g *GRPCChatAutoScript) GetChatAutoScriptById(ctx context.Context, request *pb.GetAutoScriptByIdRequest) (*pb.GetAutoScriptByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatAutoScriptService.GetChatAutoScriptById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatAutoScriptData{}
	tmp.CreatedAt = &timestamppb.Timestamp{
		Seconds: data.CreatedAt.Unix(),
	}
	tmp.UpdatedAt = &timestamppb.Timestamp{
		Seconds: data.UpdatedAt.Unix(),
	}
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetAutoScriptByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatAutoScript) UpdateChatAutoScriptById(ctx context.Context, request *pb.PutChatAutoScriptRequest) (*pb.PutChatAutoScriptResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatAutoScriptRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatAutoScriptService.UpdateChatAutoScriptById(ctx, user, request.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatAutoScriptResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatAutoScript) UpdateChatAutoScriptStatusById(ctx context.Context, request *pb.PutChatAutoScriptStatusRequest) (*pb.PutChatAutoScriptResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	statusTmp := request.GetStatus()
	var scriptStatus sql.NullBool
	if len(statusTmp) > 0 {
		tmp, _ := strconv.ParseBool(statusTmp)
		scriptStatus.Valid = true
		scriptStatus.Bool = tmp
	}
	err := service.ChatAutoScriptService.UpdateChatAutoScriptStatusById(ctx, user, request.GetId(), scriptStatus)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatAutoScriptResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatAutoScript) DeleteChatAutoScriptById(ctx context.Context, request *pb.DeleteChatAutoScriptRequest) (*pb.DeleteChatAutoScriptResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err := service.ChatAutoScriptService.DeleteChatAutoScriptById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteChatAutoScriptResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
