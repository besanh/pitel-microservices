package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_script"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatScript struct{}

func NewGRPCChatScript() pb.ChatScriptServiceServer {
	return &GRPCChatScript{}
}

func (g *GRPCChatScript) InsertChatScript(ctx context.Context, request *pb.PostChatScriptRequest) (*pb.PostChatScriptResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatScriptRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatScriptService.InsertChatScript(ctx, user, payload, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostChatScriptResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatScript) GetChatScripts(ctx context.Context, request *pb.GetChatScriptsRequest) (*pb.GetChatScriptsResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatScriptFilter{
		TenantId:   user.TenantId,
		Channel:    request.GetChannel(),
		ScriptName: request.GetScriptName(),
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

	total, data, err := service.ChatScriptService.GetChatScripts(ctx, user, payload, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatScriptData, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			tmp, err := convertChatScriptToPbChatScript(item)
			if err != nil {
				log.Error(err)
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			resultData = append(resultData, tmp)
		}
	}

	result := &pb.GetChatScriptsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   limit,
		Offset:  offset,
	}
	return result, nil
}

func (g *GRPCChatScript) GetChatScriptById(ctx context.Context, request *pb.GetScriptByIdRequest) (*pb.GetScriptByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatScriptService.GetChatScriptById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp, err := convertChatScriptToPbChatScript(*data)
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetScriptByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatScript) UpdateChatScriptById(ctx context.Context, request *pb.PutChatScriptRequest) (*pb.PutChatScriptResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatScriptRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatScriptService.UpdateChatScriptById(ctx, user, request.GetId(), payload, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatScriptResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatScript) UpdateChatScriptStatusById(ctx context.Context, request *pb.PutChatScriptStatusRequest) (*pb.PutChatScriptResponse, error) {
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
	err := service.ChatScriptService.UpdateChatScriptStatusById(ctx, user, request.GetId(), scriptStatus)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatScriptResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatScript) DeleteChatScriptById(ctx context.Context, request *pb.DeleteChatScriptRequest) (*pb.DeleteChatScriptResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err := service.ChatScriptService.DeleteChatScriptById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteChatScriptResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
