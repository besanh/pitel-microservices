package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_connection_app"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatConnectionApp struct{}

func NewGRPCChatConnectionApp() pb.ChatConnectionAppServiceServer {
	return &GRPCChatConnectionApp{}
}

func (g *GRPCChatConnectionApp) InsertChatConnectionApp(ctx context.Context, request *pb.PostChatConnectionAppRequest) (*pb.PostChatConnectionAppResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatConnectionAppRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatConnectionAppService.InsertChatConnectionApp(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostChatConnectionAppResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatConnectionApp) GetChatConnectionApps(ctx context.Context, request *pb.GetChatConnectionAppsRequest) (*pb.GetChatConnectionAppsResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatConnectionAppFilter{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	limit, offset := request.GetLimit(), request.GetOffset()

	total, data, err := service.ChatConnectionAppService.GetChatConnectionApp(ctx, user, payload, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatConnectionAppView, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			var tmp pb.ChatConnectionAppView
			tmp.CreatedAt = &timestamppb.Timestamp{
				Seconds: item.CreatedAt.Unix(),
			}
			tmp.UpdatedAt = &timestamppb.Timestamp{
				Seconds: item.UpdatedAt.Unix(),
			}
			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				result := &pb.GetChatConnectionAppsResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, nil
			}
			resultData = append(resultData, &tmp)
		}
	}

	result := &pb.GetChatConnectionAppsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   limit,
		Offset:  offset,
	}
	return result, nil
}

func (g *GRPCChatConnectionApp) GetChatConnectionAppById(ctx context.Context, request *pb.GetChatConnectionAppByIdRequest) (*pb.GetChatConnectionAppByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatConnectionAppService.GetChatConnectionAppById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatConnectionApp{}
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

	result := &pb.GetChatConnectionAppByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatConnectionApp) UpdateChatConnectionAppById(ctx context.Context, request *pb.PutChatConnectionAppRequest) (*pb.PutChatConnectionAppResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatConnectionAppRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	payload.Id = request.GetConnectionId()

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatConnectionAppService.UpdateChatConnectionAppById(ctx, user, request.GetId(), payload, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatConnectionAppResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatConnectionApp) DeleteChatConnectionAppById(ctx context.Context, request *pb.DeleteChatConnectionAppRequest) (*pb.DeleteChatConnectionAppResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err := service.ChatConnectionAppService.DeleteChatConnectionAppById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteChatConnectionAppResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
