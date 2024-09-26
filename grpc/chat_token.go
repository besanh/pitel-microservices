package grpc

import (
	"context"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/util"
	pb "github.com/tel4vn/pitel-microservices/gen/proto/chat_auth"
	"github.com/tel4vn/pitel-microservices/internal/goauth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatToken struct {
	pb.UnimplementedChatTokenServiceServer
}

func NewGRPCChatToken() *GRPCChatToken {
	return &GRPCChatToken{}
}

func (s *GRPCChatToken) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (res *pb.VerifyTokenResponse, err error) {
	token := req.GetToken()
	if len(token) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, response.ERR_TOKEN_IS_EMPTY)
	}
	var user *goauth.AuthUser
	user, err = service.ChatAuthService.VerifyToken(ctx, token)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	} else if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	tokenData := &model.TokenData{}
	if err := util.ParseAnyToAny(user.Data, tokenData); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.VerifyTokenResponse{
		Code:    "OK",
		Message: "ok",
		Data: &pb.AuthUser{
			TenantId: tokenData.TenantId,
			UserId:   tokenData.UserId,
			Username: tokenData.Username,
			Level:    tokenData.Level,
			RoleId:   tokenData.RoleId,
		},
	}, nil
}

func (s *GRPCChatToken) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	body := model.RefreshTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	}
	result, err := service.ChatAuthService.RefreshToken(ctx, body)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	return &pb.RefreshTokenResponse{
		Code:    "OK",
		Message: "ok",
		Data: &pb.RefreshTokenResponseData{
			RefreshToken: result.RefreshToken,
			Token:        result.Token,
			ExpiredIn:    int32(result.ExpiredIn),
		},
	}, nil
}
