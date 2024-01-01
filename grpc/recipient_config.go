package grpc

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/recipient_config"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCRecipientConfig struct {
	pb.UnimplementedRecipientConfigServer
}

func NewGRPCRecipientConfig() *GRPCRecipientConfig {
	return &GRPCRecipientConfig{}
}

func (s *GRPCRecipientConfig) PostRecipientConfig(ctx context.Context, req *pb.RecipientConfigBodyRequest) (result *pb.RecipientConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.RecipientConfigRequest
	if err := util.ParseAnyToAny(req, &payload); err != nil {
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}
	log.Println("payload -->", payload)

	if err := payload.Validate(); err != nil {
		log.Error(err)
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	id, err := service.NewRecipientConfig().InsertRecipientConfig(ctx, authUser, payload)
	if err != nil {
		log.Error(err)
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.RecipientConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Id:      id,
	}
	return result, nil
}

func (s *GRPCRecipientConfig) GetRecipientConfigs(ctx context.Context, request *pb.RecipientConfigRequest) (result *pb.RecipientConfigResponse, err error) {
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
	filter := model.RecipientConfigFilter{
		Provider: request.GetProvider(),
		Priority: request.GetPriority(),
		Status:   status,
	}

	total, recipients, err := service.NewRecipientConfig().GetRecipientConfigs(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.RecipientConfigResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}

	var data []*structpb.Struct
	if total > 0 {
		for _, val := range *recipients {
			var item model.RecipientConfigView
			if err = util.ParseAnyToAny(val, &item); err != nil {
				result = &pb.RecipientConfigResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
					Message: err.Error(),
				}
				return result, err
			}
			tmp, _ := util.ToStructPb(item)
			data = append(data, tmp)
		}
	}

	result = &pb.RecipientConfigResponse{
		Code:    "OK",
		Message: "success",
		Total:   int32(total),
		Data:    data,
	}

	return result, nil
}

func (s *GRPCRecipientConfig) GetRecipientConfigById(ctx context.Context, request *pb.RecipientConfigByIdRequest) (result *pb.RecipientConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	recipientConfig, err := service.NewRecipientConfig().GetRecipientConfigById(ctx, authUser, id)
	if err != nil {
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return
	}
	var res *model.RecipientConfig
	if err = util.ParseAnyToAny(recipientConfig, &res); err != nil {
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	data, _ := util.ToStructPb(res)

	result = &pb.RecipientConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    data,
	}

	return result, nil
}

func (s *GRPCRecipientConfig) PutRecipientConfigById(ctx context.Context, request *pb.RecipientConfigPutRequest) (result *pb.RecipientConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	var payload model.RecipientConfigPutRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	err = service.NewRecipientConfig().PutRecipientConfigById(ctx, authUser, id, payload)
	if err != nil {
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_PUT_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}

	result = &pb.RecipientConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    nil,
	}
	return result, nil
}

func (s *GRPCRecipientConfig) DeleteRecipientConfigById(ctx context.Context, request *pb.RecipientConfigByIdRequest) (result *pb.RecipientConfigByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	err = service.NewRecipientConfig().DeleteRecipientConfigById(ctx, authUser, id)
	if err != nil {
		result = &pb.RecipientConfigByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DELETE_FAILED].Code,
			Message: err.Error(),
		}
		return result, err
	}

	result = &pb.RecipientConfigByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    nil,
	}
	return result, nil
}
