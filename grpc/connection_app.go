package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/connection_app"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCConnectionApp struct {
	pb.UnimplementedConnectionAppServer
}

func NewGRPCConnectionApp() *GRPCConnectionApp {
	return &GRPCConnectionApp{}
}

func (s *GRPCConnectionApp) PostConnectionApp(ctx context.Context, req *pb.ConnectionAppBodyRequest) (result *pb.ConnectionAppResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.ConnectionAppRequest
	if err = util.ParseAnyToAny(req, &payload); err != nil {
		result = &pb.ConnectionAppResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	id, err := service.NewConnectionApp().InsertConnectionApp(ctx, authUser, payload)
	if err != nil {
		result = &pb.ConnectionAppResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	result = &pb.ConnectionAppResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return
}

func (s *GRPCConnectionApp) GetConnectionApp(ctx context.Context, req *pb.ConnectionAppRequest) (result *pb.ConnectionAppStructResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	status := sql.NullBool{}
	if len(req.Status) > 0 {
		status.Valid = true
		statusTmp, _ := strconv.ParseBool(req.Status)
		status.Bool = statusTmp
	}
	filter := model.ConnectionAppFilter{
		ConnectionName: req.ConnectionName,
		ConnectionType: req.ConnectionType,
		Status:         status,
	}
	limit := util.ParseLimit(req.GetLimit())
	offset := util.ParseOffset(req.GetOffset())

	total, apps, err := service.NewConnectionApp().GetConnectionApp(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.ConnectionAppStructResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	var data []*structpb.Struct
	if total > 0 {
		for _, item := range *apps {
			element := map[string]any{
				"connection_name": item.ConnectionName,
				"connection_type": item.ConnectionType,
				"id":              item.Id,
				"created_at":      item.CreatedAt,
				"updated_at":      item.UpdatedAt,
			}
			tmp, _ := util.ToStructPb(element)
			data = append(data, tmp)
		}
	}

	result = &pb.ConnectionAppStructResponse{
		Code:    "OK",
		Message: "ok",
		Total:   int32(total),
		Data:    data,
	}

	return
}
