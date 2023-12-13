package grpc

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/balance_config"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCBalanceConfig struct {
	pb.UnimplementedBalanceConfigServer
}

func NewGRPCBalanceConfig() *GRPCBalanceConfig {
	return &GRPCBalanceConfig{}
}

func (g *GRPCBalanceConfig) PostBalanceConfig(ctx context.Context, request *pb.BalanceConfigBodyRequest) (result *pb.BalanceConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.BalanceConfigBodyRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}
	log.Println("payload -->", payload)

	if err := payload.Validate(); err != nil {
		log.Error(err)
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	if err := service.NewBalanceConfig().InsertBalanceConfig(ctx, authUser, payload); err != nil {
		log.Error(err)
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.BalanceConfigByIdResponse{
		Code:    "OK",
		Message: "success",
	}

	return result, nil
}

func (s *GRPCBalanceConfig) GetBalanceConfigs(ctx context.Context, request *pb.BalanceConfigRequest) (result *pb.BalanceConfigResponse, err error) {
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
	filter := model.BalanceConfigFilter{
		Provider: request.GetProvider(),
		Priority: request.GetPriority(),
		Status:   status,
	}

	total, Balances, err := service.NewBalanceConfig().GetBalanceConfigs(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.BalanceConfigResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}

	var data []*structpb.Struct
	if total > 0 {
		for _, val := range *Balances {
			var item model.BalanceConfigView
			if err = util.ParseAnyToAny(val, &item); err != nil {
				result = &pb.BalanceConfigResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
					Message: err.Error(),
				}
				return result, err
			}
			tmp, _ := util.ToStructPb(item)
			data = append(data, tmp)
		}
	}

	result = &pb.BalanceConfigResponse{
		Code:    "OK",
		Message: "success",
		Total:   int32(total),
		Data:    data,
	}

	return result, nil
}

func (s *GRPCBalanceConfig) GetBalanceConfigById(ctx context.Context, request *pb.BalanceConfigByIdRequest) (result *pb.BalanceConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	BalanceConfig, err := service.NewBalanceConfig().GetBalanceConfigById(ctx, authUser, id)
	if err != nil {
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return
	}
	var res *model.BalanceConfig
	if err = util.ParseAnyToAny(BalanceConfig, &res); err != nil {
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	data, _ := util.ToStructPb(res)

	result = &pb.BalanceConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    data,
	}

	return result, nil
}

func (s *GRPCBalanceConfig) PutBalanceConfigById(ctx context.Context, request *pb.BalanceConfigPutRequest) (result *pb.BalanceConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	var payload model.BalanceConfigPutBodyRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	err = service.NewBalanceConfig().PutBalanceConfigById(ctx, authUser, id, payload)
	if err != nil {
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_PUT_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}

	result = &pb.BalanceConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    nil,
	}
	return result, nil
}

func (s *GRPCBalanceConfig) DeleteBalanceConfigById(ctx context.Context, request *pb.BalanceConfigByIdRequest) (result *pb.BalanceConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	err = service.NewBalanceConfig().DeleteBalanceConfigById(ctx, authUser, id)
	if err != nil {
		result = &pb.BalanceConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DELETE_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}

	result = &pb.BalanceConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    nil,
	}
	return result, nil
}
