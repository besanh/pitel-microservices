package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_app"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCApp struct {
	pb.UnimplementedAppServer
}

func NewGRPCApp() *GRPCApp {
	return &GRPCApp{}
}

func (g *GRPCApp) InsertApp(ctx context.Context, req *pb.ChatAppBodyRequest) (result *pb.AppResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.ChatAppRequest
	if err = util.ParseAnyToAny(req, &payload); err != nil {
		result = &pb.AppResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	if err = payload.Validate(); err != nil {
		result = &pb.AppResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_VALIDATION_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	id, err := service.NewChatApp().InsertChatApp(ctx, authUser, payload)
	if err != nil {
		result = &pb.AppResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_INSERT_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	result = &pb.AppResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (s *GRPCApp) GetApp(ctx context.Context, req *pb.AppRequest) (result *pb.AppGetResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	filter := model.AppFilter{
		AppName: req.AppName,
		Status:  req.Status,
	}
	limit := util.ParseLimit(req.GetLimit())
	offset := util.ParseOffset(req.GetOffset())

	total, apps, err := service.NewChatApp().GetChatApp(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.AppGetResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return
	}

	var data []*structpb.Struct
	if total > 0 {
		for _, item := range *apps {
			element := map[string]any{
				"app_name":   item.AppName,
				"status":     item.Status,
				"id":         item.Id,
				"created_at": item.CreatedAt,
				"updated_at": item.UpdatedAt,
			}
			tmp, _ := util.ToStructPb(element)
			data = append(data, tmp)
		}
	}

	result = &pb.AppGetResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data,
		Total:   int32(total),
	}

	return
}
