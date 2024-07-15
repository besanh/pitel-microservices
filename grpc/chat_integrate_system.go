package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_integrate_system"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatIntegrateSystem struct {
	pb.UnimplementedChatIntegrateSystemServer
}

func NewGRPCChatIntegrateSystem() *GRPCChatIntegrateSystem {
	return &GRPCChatIntegrateSystem{}
}

func (g *GRPCChatIntegrateSystem) GetChatIntegrateSystems(ctx context.Context, request *pb.GetChatIntegrateSystemRequest) (result *pb.GetChatIntegrateSystemResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	statusTmp := request.GetStatus()
	var statusIntegrate sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		statusIntegrate.Valid = true
		statusIntegrate.Bool = statusTmp
	}
	filter := model.ChatIntegrateSystemFilter{
		SystemName: request.GetSystemName(),
		VendorName: request.GetVendorName(),
		Status:     statusIntegrate,
		SystemId:   request.GetSystemId(),
	}

	total, chatIntegrateSystems, err := service.ChatIntegrateSystemService.GetChatIntegrateSystems(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.GetChatIntegrateSystemResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}
	data := make([]*pb.ChatIntegrateSystemData, 0)
	if err = util.ParseAnyToAny(chatIntegrateSystems, &data); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatIntegrateSystemResponse{
		Code:    "OK",
		Message: "ok",
		Total:   int32(total),
		Data:    data,
	}

	return
}

func (g *GRPCChatIntegrateSystem) PostChatIntegrateSystem(ctx context.Context, req *pb.PostChatIntegrateSystemRequest) (result *pb.PostChatIntegrateSystemResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatIntegrateSystemRequest{
		SystemName:    req.GetSystemName(),
		VendorId:      req.GetVendorId(),
		Status:        req.GetStatus(),
		AuthType:      req.GetAuthType(),
		Username:      req.GetUsername(),
		Password:      req.GetPassword(),
		Token:         req.GetToken(),
		WebsocketUrl:  req.GetWebsocketUrl(),
		ApiUrl:        req.GetApiUrl(),
		ApiGetUserUrl: req.GetApiGetUserUrl(),
		ChatApps:      req.GetChatApps(),
	}

	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, systemId, err := service.ChatIntegrateSystemService.InsertChatIntegrateSystem(ctx, authUser, &payload)
	if err != nil {
		result = &pb.PostChatIntegrateSystemResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.PostChatIntegrateSystemResponse{
		Code:     "OK",
		Message:  "ok",
		Id:       id,
		SystemId: systemId,
	}
	return
}

func (g *GRPCChatIntegrateSystem) GetChatIntegrateSystemById(ctx context.Context, req *pb.GetChatIntegrateSystemByIdRequest) (result *pb.GetChatIntegrateSystemByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	integrateSystem, err := service.ChatIntegrateSystemService.GetChatIntegrateSystemById(ctx, authUser, req.GetId())
	if err != nil {
		log.Error(err)
		result = &pb.GetChatIntegrateSystemByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	data := &pb.ChatIntegrateSystemData{}
	if err = util.ParseAnyToAny(integrateSystem, data); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatIntegrateSystemByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data,
	}
	return
}

func (g *GRPCChatIntegrateSystem) UpdateChatIntegrateSystemById(ctx context.Context, req *pb.PutChatIntegrateSystemRequest) (result *pb.PostChatIntegrateSystemResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatIntegrateSystemRequest{
		SystemName:    req.GetSystemName(),
		VendorId:      req.GetVendorId(),
		Status:        req.GetStatus(),
		AuthType:      req.GetAuthType(),
		Username:      req.GetUsername(),
		Password:      req.GetPassword(),
		Token:         req.GetToken(),
		WebsocketUrl:  req.GetWebsocketUrl(),
		ApiUrl:        req.GetApiUrl(),
		ApiGetUserUrl: req.GetApiGetUserUrl(),
	}

	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ChatIntegrateSystemService.UpdateChatIntegrateSystemById(ctx, authUser, req.GetId(), &payload)
	if err != nil {
		log.Error(err)
		result = &pb.PostChatIntegrateSystemResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_PUT_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.PostChatIntegrateSystemResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatIntegrateSystem) DeleteChatIntegrateSystemById(ctx context.Context, req *pb.GetChatIntegrateSystemByIdRequest) (result *pb.PostChatIntegrateSystemResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err = service.ChatIntegrateSystemService.DeleteChatIntegrateSystemById(ctx, authUser, req.GetId())
	if err != nil {
		log.Error(err)
		result = &pb.PostChatIntegrateSystemResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DELETE_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.PostChatIntegrateSystemResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatIntegrateSystem) GetChatIntegrateSystemById(ctx context.Context, request *pb.GetChatIntegrateSystemByIdRequest) (result *pb.GetChatIntegrateSystemByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	system, err := service.ChatIntegrateSystemService.GetChatIntegrateSystemById(ctx, authUser, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	data, err := util.ToStructPb(system)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatIntegrateSystemByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data,
	}
	return
}
func (g *GRPCChatIntegrateSystem) UpdateChatIntegrateSystemById(ctx context.Context, request *pb.UpdateChatIntegrateSystemByIdRequest) (result *pb.UpdateChatIntegrateSystemByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatIntegrateSystemRequest{
		SystemName:    request.GetData().GetSystemName(),
		VendorId:      request.GetData().GetVendorId(),
		Status:        request.GetData().GetStatus(),
		AuthType:      request.GetData().GetAuthType(),
		Username:      request.GetData().GetUsername(),
		Password:      request.GetData().GetPassword(),
		Token:         request.GetData().GetToken(),
		WebsocketUrl:  request.GetData().GetWebsocketUrl(),
		ApiUrl:        request.GetData().GetApiUrl(),
		ApiGetUserUrl: request.GetData().GetApiGetUserUrl(),
		ChatApps:      request.GetData().GetChatApps(),
	}

	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ChatIntegrateSystemService.UpdateChatIntegrateSystem(ctx, authUser, request.GetId(), &payload)
	if err != nil {
		result = &pb.UpdateChatIntegrateSystemByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_PUT_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.UpdateChatIntegrateSystemByIdResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
func (g *GRPCChatIntegrateSystem) DeleteChatIntegrateSystemById(ctx context.Context, request *pb.DeleteChatIntegrateSystemByIdRequest) (result *pb.DeleteChatIntegrateSystemByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err = service.ChatIntegrateSystemService.DeleteChatIntegrateSystem(ctx, authUser, request.GetId())
	if err != nil {
		result = &pb.DeleteChatIntegrateSystemByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DELETE_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.DeleteChatIntegrateSystemByIdResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
