package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/util"
	pb "github.com/tel4vn/pitel-microservices/gen/proto/chat_routing"
	"github.com/tel4vn/pitel-microservices/middleware/auth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatRouting struct{}

func NewGRPCChatRouting() pb.ChatRoutingServiceServer {
	return &GRPCChatRouting{}
}

func (g *GRPCChatRouting) InsertChatRouting(ctx context.Context, request *pb.PostChatRoutingRequest) (result *pb.PostChatRoutingResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	statusTmp, _ := strconv.ParseBool(request.GetStatus())
	payload := model.ChatRoutingRequest{
		RoutingName:  request.GetRoutingName(),
		RoutingAlias: request.GetRoutingAlias(),
		Status:       statusTmp,
	}

	if err = payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatRoutingService.InsertChatRouting(ctx, user, &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PostChatRoutingResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return
}

func (g *GRPCChatRouting) GetChatRoutings(ctx context.Context, request *pb.GetChatRoutingsRequest) (result *pb.GetChatRoutingsResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var statusTmp sql.NullBool
	if len(request.GetStatus()) > 0 {
		tmp, _ := strconv.ParseBool(request.GetStatus())
		statusTmp.Valid = true
		statusTmp.Bool = tmp
	}
	filter := model.ChatRoutingFilter{
		RoutingName:  request.GetRoutingName(),
		RoutingAlias: request.GetRoutingAlias(),
		Status:       statusTmp,
	}
	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.ChatRoutingService.GetChatRoutings(ctx, user, filter, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatRouting, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			var tmp pb.ChatRouting
			tmp.CreatedAt = timestamppb.New(item.CreatedAt)
			tmp.UpdatedAt = timestamppb.New(item.UpdatedAt)

			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				result = &pb.GetChatRoutingsResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, nil
			}
			resultData = append(resultData, &tmp)
		}
	}

	result = &pb.GetChatRoutingsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}

func (g *GRPCChatRouting) GetChatRoutingById(ctx context.Context, request *pb.GetChatRoutingByIdRequest) (result *pb.GetChatRoutingByIdResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatRoutingService.GetChatRoutingById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatRouting{}
	tmp.CreatedAt = timestamppb.New(data.CreatedAt)
	tmp.UpdatedAt = timestamppb.New(data.UpdatedAt)
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.GetChatRoutingByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return
}

func (g *GRPCChatRouting) UpdateChatRoutingById(ctx context.Context, request *pb.PutChatRoutingRequest) (result *pb.PutChatRoutingResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	statusTmp, _ := strconv.ParseBool(request.GetStatus())
	payload := model.ChatRoutingRequest{
		RoutingName:  request.GetRoutingName(),
		RoutingAlias: request.GetRoutingAlias(),
		Status:       statusTmp,
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ChatRoutingService.UpdateChatRoutingById(ctx, user, request.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutChatRoutingResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCChatRouting) DeleteChatRoutingById(ctx context.Context, request *pb.DeleteChatRoutingRequest) (result *pb.DeleteChatRoutingResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if len(request.GetId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	err = service.ChatRoutingService.DeleteChatRoutingById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.DeleteChatRoutingResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
