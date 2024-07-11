package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_integrate_system"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
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
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}
	filter := model.ChatIntegrateSystemFilter{
		SystemName: request.GetSystemName(),
		VendorName: request.GetVendorName(),
		Status:     status,
	}

	total, chatIntegrateSystems, err := service.ChatIntegrateSystemService.GetChatIntegrateSystems(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.GetChatIntegrateSystemResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}

	var data []*structpb.Struct
	if len(*chatIntegrateSystems) > 0 {
		for _, item := range *chatIntegrateSystems {
			var itm model.ChatIntegrateSystem
			if err := util.ParseAnyToAny(item, &itm); err != nil {
				result = &pb.GetChatIntegrateSystemResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
					Message: err.Error(),
				}
				return result, nil
			}

			tmp, _ := util.ToStructPb(itm)
			data = append(data, tmp)
		}
	}

	result = &pb.GetChatIntegrateSystemResponse{
		Code:    response.MAP_ERR_RESPONSE[response.SUCCESS].Code,
		Message: response.MAP_ERR_RESPONSE[response.SUCCESS].Message,
		Total:   int32(total),
		Data:    data,
	}

	return
}
