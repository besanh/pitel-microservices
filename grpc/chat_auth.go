package grpc

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/hash"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCChatAuth struct {
	pb.UnimplementedChatAuthServiceServer
}

func NewGRPCChatAuth() *GRPCChatAuth {
	return &GRPCChatAuth{}
}

func (s *GRPCChatAuth) Login(ctx context.Context, req *pb.LoginRequest) (res *pb.LoginResponse, err error) {
	fingerprint := GetFingerprintFromMD(ctx)
	userAgent := GetUserAgentFromMD(ctx)
	body := model.LoginRequest{
		Username:    req.GetUsername(),
		Password:    req.GetPassword(),
		Fingerprint: hash.HashMD5(fingerprint),
		UserAgent:   userAgent,
	}
	var result *model.LoginResponse
	result, err = service.ChatAuthService.Login(ctx, body)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res = &pb.LoginResponse{
		Code:    "OK",
		Message: "ok",
		Data: &pb.LoginResponseData{
			UserId:       result.UserId,
			Username:     result.Username,
			Token:        result.Token,
			RefreshToken: result.RefreshToken,
			TenantId:     result.TenantId,
			ExpiredIn:    int32(result.ExpiredIn),
			Fullname:     result.Fullname,
			TenantName:   result.TenantName,
			TenantLogo:   result.TenantLogo,
			Level:        result.Level,
			RoleId:       result.RoleId,
		},
	}
	return
}

func GetUserAgentFromMD(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	userAgent, ok := md["user_agent"]
	if !ok {
		return ""
	}
	return userAgent[0]
}

func GetFingerprintFromMD(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	fingerprint, ok := md["fingerprint"]
	if !ok {
		return ""
	}
	return fingerprint[0]
}

func GetTokenFromMD(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	token, ok := md["token"]
	if !ok {
		return ""
	}
	return token[0]
}

// func (s *GRPCChatAuth) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (res *pb.ForgotPasswordResponse, err error) {
// 	validator, err := validator.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = validator.Validate(req)
// 	if err != nil {
// 		res = &pb.ForgotPasswordResponse{
// 			Code:    response.MAP_ERR_RESPONSE[response.ERR_REQUEST_INVALID].Code,
// 			Message: err.Error(),
// 		}
// 		if isStr, r := response.HandleValidatorPBError(err); isStr {
// 			res.ErrorCode = &pb.ForgotPasswordResponse_Error{Error: r.(string)}
// 		} else {
// 			var tmp *structpb.Struct
// 			tmp, err = util.ToStructPb(r)
// 			if err != nil {
// 				return nil, status.Error(codes.Unavailable, err.Error())
// 			}
// 			res.ErrorCode = &pb.ForgotPasswordResponse_Errors{Errors: tmp}
// 		}
// 		return res, nil
// 	}
// 	userAgent := GetUserAgentFromMD(ctx)

// 	body := model.ForgotPassword{
// 		State:     req.GetState(),
// 		Email:     req.GetEmail(),
// 		Code:      req.GetCode(),
// 		Password:  req.GetPassword(),
// 		UserAgent: userAgent,
// 		Token:     req.GetToken(),
// 	}
// 	code := "OK"
// 	message := "ok"
// 	result, err := service.ChatAuthService.ForgotPassword(ctx, body)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unavailable, err.Error())
// 	} else if result != nil && !result.IsValid {
// 		code = result.Code
// 		message = response.MAP_ERR_RESPONSE[result.Code].Message
// 	}
// 	log.Debugf("result: %v", result)

// 	data, err := util.ToStructPb(result)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unavailable, err.Error())
// 	}

// 	res = &pb.ForgotPasswordResponse{
// 		Code:    code,
// 		Message: message,
// 		Data:    data,
// 	}
// 	return
// }

func (s *GRPCChatAuth) Logout(ctx context.Context, req *pb.LogoutRequest) (res *pb.LogoutResponse, err error) {
	fingerprint := GetFingerprintFromMD(ctx)
	token := GetTokenFromMD(ctx)
	if len(token) < 1 {
		err = errors.New(response.ERR_REQUEST_INVALID)
		return
	}
	err = service.ChatAuthService.Logout(ctx, token, fingerprint)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res = &pb.LogoutResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
