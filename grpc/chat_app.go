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

func (g *GRPCChatApp) GetChatApps(ctx context.Context, req *pb.GetChatAppRequest) (result *pb.GetChatAppResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatAppFilter{}
	if err := util.ParseAnyToAny(req, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	limit, offset := req.GetLimit(), req.GetOffset()

	total, chatApps, err := service.ChatAppService.GetChatApp(ctx, user, payload, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	data := make([]*pb.GetChatAppData, 0)
	if err = util.ParseAnyToAny(chatApps, &data); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatAppResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data,
		Total:   int32(total),
	}
	return
}

func (g *GRPCChatApp) GetChatAppById(ctx context.Context, req *pb.GetChatAppByIdRequest) (result *pb.GetChatAppByIdResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	chatApp, err := service.ChatAppService.GetChatAppById(ctx, user, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	data := &pb.GetChatAppData{}
	if err = util.ParseAnyToAny(chatApp, &data); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatAppByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data,
	}
	return
}

func (g *GRPCChatApp) UpdateChatAppById(ctx context.Context, req *pb.UpdateChatAppByIdRequest) (result *pb.UpdateChatAppByIdResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatAppRequest{}
	if err = util.ParseAnyToAny(req.Data, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err = payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ChatAppService.UpdateChatAppById(ctx, user, req.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.UpdateChatAppByIdResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatApp) DeleteChatAppById(ctx context.Context, req *pb.DeleteChatAppByIdRequest) (result *pb.DeleteChatAppByIdResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err = service.ChatAppService.DeleteChatAppById(ctx, user, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.DeleteChatAppByIdResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
