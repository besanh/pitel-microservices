package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/share_info"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCShareInfo struct {
	pb.UnimplementedShareInfoServiceServer
}

func NewGRPCShareInfo() *GRPCShareInfo {
	return &GRPCShareInfo{}
}

// Send to ott share info
func (s *GRPCShareInfo) PostShareInfo(ctx context.Context, req *pb.PostShareInfoRequest) (result *pb.PostShareInfoResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		err = status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
		return
	}

	shareInfo := model.ShareInfoFormSubmitRequest{
		ShareType:      req.GetShareType(),
		EventName:      req.EventName,
		AppId:          req.GetAppId(),
		OaId:           req.GetOaId(),
		ExternalUserId: req.GetExternalUserId(),
	}

	if err = shareInfo.Validate(); err != nil {
		log.Error(err)
		return
	}

	err = service.ShareInfoService.PostRequestShareInfo(ctx, authUser, shareInfo)
	if err != nil {
		err = status.Errorf(codes.Internal, err.Error())
		return
	}
	result = &pb.PostShareInfoResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (s *GRPCShareInfo) GetShareInfos(ctx context.Context, req *pb.GetShareInfoRequest) (result *pb.GetShareInfoResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		err = status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
		return
	}

	limit, offset := util.ParseLimit(req.Limit), util.ParseOffset(req.Offset)
	filter := model.ShareInfoFormFilter{
		OaId:      req.GetOaId(),
		ShareType: req.GetShareType(),
		AppId:     req.GetAppId(),
	}
	total, data, err := service.ShareInfoService.GetShareInfos(ctx, authUser, filter, limit, offset)
	if err != nil {
		err = status.Errorf(codes.Internal, err.Error())
		return
	}
	tmp := make([]*pb.ShareInfo, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			itm := &pb.ShareInfo{}
			if err = util.ParseAnyToAny(item, itm); err != nil {
				log.Error(err)
				return
			}
			tmp = append(tmp, itm)
		}
	}

	result = &pb.GetShareInfoResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}

func (s *GRPCShareInfo) GetShareInfoById(ctx context.Context, req *pb.GetShareInfoByIdRequest) (result *pb.GetShareInfoByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		err = status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
		return
	}

	shareInfo, err := service.NewShareInfo().GetShareInfoById(ctx, authUser, req.GetId())
	if err != nil {
		err = status.Errorf(codes.Internal, err.Error())
		return
	}
	tmp := &pb.ShareInfo{}
	if err = util.ParseAnyToAny(shareInfo, tmp); err != nil {
		log.Error(err)
		return
	}
	result = &pb.GetShareInfoByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return
}

func (s *GRPCShareInfo) DeleteShareInfoById(ctx context.Context, req *pb.DeleteShareInfoRequest) (result *pb.DeleteShareInfoResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		err = status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
		return
	}

	err = service.ShareInfoService.DeleteShareInfoById(ctx, authUser, req.GetId())
	if err != nil {
		err = status.Errorf(codes.Internal, err.Error())
		return
	}
	result = &pb.DeleteShareInfoResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
