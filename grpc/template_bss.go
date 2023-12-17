package grpc

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/template_bss"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCTemplateBss struct {
	pb.UnimplementedTemplateBssServer
}

func NewGRPCTemplateBss() *GRPCTemplateBss {
	return &GRPCTemplateBss{}
}

func (g *GRPCTemplateBss) PostTemplateBss(ctx context.Context, request *pb.TemplateBssBodyRequest) (result *pb.TemplateBssByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.TemplateBssBodyRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.TemplateBssByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	log.Println("payload -->", payload)
	if err = payload.Validate(); err != nil {
		log.Error(err)
		result = &pb.TemplateBssByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	err = service.NewTemplateBss().InsertTemplateBss(ctx, authUser, payload)
	if err != nil {
		log.Error(err)
		result = &pb.TemplateBssByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.TemplateBssByIdResponse{
		Code:    "OK",
		Message: "success",
	}

	return result, nil
}

func (g *GRPCTemplateBss) GetTemplateBsses(ctx context.Context, request *pb.TemplateBssRequest) (result *pb.TemplateBssResponse, err error) {
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
	filter := model.TemplateBssFilter{
		TemplateName: request.GetTemplateName(),
		TemplateCode: request.GetTemplateCode(),
		TemplateType: request.GetTemplateType(),
		Content:      request.GetContent(),
		Status:       status,
	}

	total, templateBsses, err := service.NewTemplateBss().GetTemplateBsses(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.TemplateBssResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	var data []*structpb.Struct
	if total > 0 {
		for _, val := range *templateBsses {
			item := model.TemplateBss{}
			if err := util.ParseAnyToAny(val, &item); err != nil {
				log.Error(err)
				result = &pb.TemplateBssResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
					Message: err.Error(),
				}
				return result, err
			}
			tmp, _ := util.ToStructPb(item)
			data = append(data, tmp)
		}
	}

	result = &pb.TemplateBssResponse{
		Code:    "OK",
		Message: "success",
		Data:    data,
		Total:   int32(total),
	}

	return result, nil
}

func (g *GRPCTemplateBss) GetTemplateBssById(ctx context.Context, request *pb.TemplateBssByIdRequest) (result *pb.TemplateBssByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	templateBss, err := service.NewTemplateBss().GetTemplateBssById(ctx, authUser, id)
	if err != nil {
		result = &pb.TemplateBssByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	data, _ := util.ToStructPb(*templateBss)
	result = &pb.TemplateBssByIdResponse{
		Code:    "OK",
		Message: "success",
		Data:    data,
	}
	return result, nil
}

func (g *GRPCTemplateBss) PutTemplateBssById(ctx context.Context, request *pb.TemplateBssPutRequest) (result *pb.TemplateBssByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	var payload model.TemplateBssBodyRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.TemplateBssByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	err = service.NewTemplateBss().PutTemplateBssById(ctx, authUser, id, payload)
	if err != nil {
		result = &pb.TemplateBssByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.TemplateBssByIdResponse{
		Code:    "OK",
		Message: "success",
	}

	return result, nil
}

func (g *GRPCTemplateBss) DeleteTemplateBssById(ctx context.Context, request *pb.TemplateBssByIdRequest) (result *pb.TemplateBssByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	id := request.GetId()
	if len(id) < 1 {
		return nil, errors.New("id is missing")
	}

	err = service.NewTemplateBss().DeleteTemplateBssById(ctx, authUser, id)
	if err != nil {
		result = &pb.TemplateBssByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	result = &pb.TemplateBssByIdResponse{
		Code:    "OK",
		Message: "success",
	}

	return result, nil
}
