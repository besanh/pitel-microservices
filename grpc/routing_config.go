package grpc

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/routing_config"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCRoutingConfig struct {
	pb.UnimplementedRoutingConfigServer
}

func NewGRPCRoutingConfig() *GRPCRoutingConfig {
	return &GRPCRoutingConfig{}
}

func (g *GRPCRoutingConfig) PostRoutingConfig(ctx context.Context, request *pb.RoutingConfigBodyRequest) (result *pb.RoutingConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.RoutingConfig
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}
	// log.Info("payload -->", payload)

	if err := payload.Validate(); err != nil {
		log.Error(err)
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	id, err := service.NewRoutingConfig().InsertRoutingConfig(ctx, authUser, payload)
	if err != nil {
		log.Error(err)
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.RoutingConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Id:      id,
	}

	return result, nil
}

func (g *GRPCRoutingConfig) GetRoutingConfigs(ctx context.Context, request *pb.RoutingConfigRequest) (result *pb.RoutingConfigResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit := util.ParseLimit(request.GetLimit())
	offset := util.ParseOffset(request.GetOffset())
	status := sql.NullBool{}
	if len(request.Status) > 0 {
		status.Valid = true
		statusTmp, _ := strconv.ParseBool(request.Status)
		status.Bool = statusTmp
	}
	filter := model.RoutingConfigFilter{
		RoutingName: request.GetRoutingName(),
		RoutingType: request.GetRoutingType(),
		Brandname:   request.GetBrandName(),
		Status:      status,
	}

	total, routingConfigs, err := service.NewRoutingConfig().GetRoutingConfigs(ctx, authUser, filter, limit, offset)
	if err != nil {
		log.Error(err)
		result = &pb.RoutingConfigResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	var data []*structpb.Struct

	if total > 0 {
		for _, val := range *routingConfigs {
			var item model.RoutingConfig
			if err := util.ParseAnyToAny(val, &item); err != nil {
				log.Error(err)
				result = &pb.RoutingConfigResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
					Message: err.Error(),
				}
				return result, nil
			}

			tmp, _ := util.ToStructPb(item)
			data = append(data, tmp)
		}
	}

	result = &pb.RoutingConfigResponse{
		Code:    "OK",
		Message: "success",
		Data:    data,
		Total:   int32(total),
	}

	return result, nil
}

func (g *GRPCRoutingConfig) GetRoutingConfigById(ctx context.Context, request *pb.RoutingConfigByIdRequest) (result *pb.RoutingConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	routingConfig, err := service.NewRoutingConfig().GetRoutingConfigById(ctx, authUser, id)
	if err != nil {
		log.Error(err)
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	var res model.RoutingConfig

	if err := util.ParseAnyToAny(routingConfig, &res); err != nil {
		log.Error(err)
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	data, _ := util.ToStructPb(res)

	result = &pb.RoutingConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    data,
	}

	return result, nil
}

func (g *GRPCRoutingConfig) PutRoutingConfigById(ctx context.Context, request *pb.RoutingConfigPutRequest) (result *pb.RoutingConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	var routingConfig model.RoutingConfig
	if err := util.ParseAnyToAny(request, &routingConfig); err != nil {
		log.Error(err)
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	err = service.NewRoutingConfig().PutRoutingConfigById(ctx, authUser, id, routingConfig)
	if err != nil {
		log.Error(err)
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.RoutingConfigByIdResponse{
		Code:    "OK",
		Message: "success",
	}

	return result, nil
}

func (g *GRPCRoutingConfig) DeleteRoutingConfigById(ctx context.Context, request *pb.RoutingConfigByIdRequest) (result *pb.RoutingConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	err = service.NewRoutingConfig().DeleteRoutingConfigById(ctx, authUser, id)
	if err != nil {
		log.Error(err)
		result = &pb.RoutingConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.RoutingConfigByIdResponse{
		Code:    "OK",
		Message: "success",
	}

	return result, nil
}
