package grpc

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_vendor"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatVendor struct {
	pb.UnimplementedChatVendorServiceServer
}

func NewGRPCChatVendor() *GRPCChatVendor {
	return &GRPCChatVendor{}
}

func (g *GRPCChatVendor) PostChatVendor(ctx context.Context, req *pb.PostChatVendorRequest) (*pb.PostChatVendorResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if (user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	payload := model.ChatVendorRequest{
		VendorName: req.GetVendorName(),
		VendorType: req.GetVendorType(),
		Status:     req.GetStatus(),
	}
	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatVendorService.InsertChatVendor(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.PostChatVendorResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}, nil
}

func (g *GRPCChatVendor) ListChatVendors(ctx context.Context, req *pb.ListChatVendorRequest) (*pb.ListChatVendorResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if (user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	limit, offset := req.GetLimit(), req.GetOffset()
	filter := model.ChatVendorFilter{
		VendorName: req.GetVendorName(),
		VendorType: req.GetVendorType(),
		Status: sql.NullBool{
			Bool:  req.GetStatus(),
			Valid: true,
		},
	}

	total, result, err := service.ChatVendorService.GetChatVendors(ctx, user, filter, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	data := make([]*pb.ChatVendorConfiguration, 0)
	if err = util.ParseAnyToAny(result, &data); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListChatVendorResponse{
		Code:    "OK",
		Message: "ok",
		Total:   int32(total),
		Data:    data,
	}, nil
}

func (g *GRPCChatVendor) GetChatVendor(ctx context.Context, req *pb.GetChatVendorRequest) (*pb.GetChatVendorResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if (user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	result, err := service.ChatVendorService.GetChatVendorById(ctx, user, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var data *pb.ChatVendorConfiguration
	if err = util.ParseAnyToAny(result, &data); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.GetChatVendorResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data,
	}, nil
}

func (g *GRPCChatVendor) UpdateChatVendor(ctx context.Context, req *pb.UpdateChatVendorRequest) (*pb.UpdateChatVendorResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if (user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	payload := model.ChatVendorRequest{
		VendorName: req.GetRequest().GetVendorName(),
		VendorType: req.GetRequest().GetVendorType(),
		Status:     req.GetRequest().GetStatus(),
	}
	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatVendorService.UpdateChatVendor(ctx, user, req.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.UpdateChatVendorResponse{
		Code:    "OK",
		Message: "ok",
	}, nil
}

func (g *GRPCChatVendor) DeleteChatVendor(ctx context.Context, req *pb.DeleteChatVendorRequest) (*pb.DeleteChatVendorResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if (user.GetLevel() != "superadmin") && len(user.SecretKey) < 1 {
		return nil, status.Errorf(codes.PermissionDenied, response.ERR_PERMISSION_DENIED)
	}

	err := service.ChatVendorService.DeleteChatVendor(ctx, user, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.DeleteChatVendorResponse{
		Code:    "OK",
		Message: "ok",
	}, nil
}
