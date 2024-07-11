package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/response"
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
