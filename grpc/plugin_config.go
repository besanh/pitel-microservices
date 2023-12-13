package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/plugin_config"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCPluginConfig struct {
	pb.UnimplementedPluginConfigServer
}

func NewGRPCPluginConfig() *GRPCPluginConfig {
	return &GRPCPluginConfig{}
}

func (g *GRPCPluginConfig) PostPluginConfig(ctx context.Context, request *pb.PluginConfigBodyRequest) (result *pb.PluginConfigResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	var payload model.PluginConfigRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.PluginConfigResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}
	log.Println("payload -->", payload)
	if err := payload.Validate(); err != nil {
		result = &pb.PluginConfigResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	if err := service.NewPluginConfig().PostPluginConfig(ctx, authUser, payload); err != nil {
		log.Error(err)
		result = &pb.PluginConfigResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.PluginConfigResponse{
		Code:    "OK",
		Message: "success",
	}
	return result, nil
}

func (g *GRPCPluginConfig) GetPluginConfigs(ctx context.Context, request *pb.PluginConfigRequest) (result *pb.PluginConfigResponse, err error) {
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
	filter := model.PluginConfigFilter{
		PluginName: request.GetPluginName(),
		PluginType: request.PluginType,
		Status:     status,
	}

	total, plugins, err := service.NewPluginConfig().GetPluginConfigs(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.PluginConfigResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}

	var data []*structpb.Struct
	if total > 0 {
		for _, val := range *plugins {
			var item model.PluginConfigView
			if err = util.ParseAnyToAny(val, &item); err != nil {
				result = &pb.PluginConfigResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
					Message: err.Error(),
				}
				return result, err
			}
			tmp, _ := util.ToStructPb(item)
			data = append(data, tmp)
		}
	}

	result = &pb.PluginConfigResponse{
		Code:    "OK",
		Message: "success",
		Data:    data,
		Total:   int32(total),
	}

	return result, nil
}
