package auth

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	AUTHENTICATED        = "authenticated"
	USER                 = "user"
	ERR_TOKEN_IS_EMPTY   = "token is empty"
	ERR_TOKEN_IS_INVALID = "token is invalid"
	ERR_TOKEN_IS_EXPIRED = "token is expired"
)

// AuthFunc is a middleware (interceptor) that extracts token from header
func GRPCAuthMiddleware(ctx context.Context) (context.Context, error) {
	isAuthenticaed, ok := ctx.Value(AUTHENTICATED).(bool)
	if !ok || !isAuthenticaed {
		var user *model.AuthUser = ParseHeaderToUser(ctx)
		if len(user.SecretKey) > 0 {
			encryptPassword := []byte(service.SECRET_KEY_SUPERADMIN)
			password := fmt.Sprintf("%x", md5.Sum(encryptPassword))
			if user.SecretKey != password {
				log.Error("invalid secret key")
				return nil, status.Errorf(codes.Unauthenticated, "invalid secret key")
			}
			ctx = context.WithValue(ctx, AUTHENTICATED, true)
			ctx = context.WithValue(ctx, USER, user)
		} else {
			// if
			token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
			if err != nil {
				return nil, err
			}
			authUser, _, err := validateToken(ctx, token)
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, err.Error())
			}
			ctx = context.WithValue(ctx, AUTHENTICATED, true)
			ctx = context.WithValue(ctx, USER, authUser)
		}
	}
	return ctx, nil
}

// Authorization unary interceptor function to handle authorize per RPC call
func GRPCAuthInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	log.Infof("method name: %s", info.FullMethod)
	start := time.Now()
	if util.InArrayContains(info.FullMethod, []string{"proto.chatAuth"}) {
		ctx = context.WithValue(ctx, AUTHENTICATED, true)
	}
	h, err := handler(ctx, req)
	logrus.WithFields(logrus.Fields{
		"method":  info.FullMethod,
		"latency": time.Since(start).Microseconds(),
		"error":   err,
	}).Info("gRPC Request")
	return h, err
}

func ParseHeaderToUser(ctx context.Context) *model.AuthUser {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Debugf("no metadata")
		return nil
	}
	return &model.AuthUser{
		TenantId:  GetFirstOfSlices(md["tenant_id"]),
		UserId:    GetFirstOfSlices(md["user_id"]),
		Username:  GetFirstOfSlices(md["username"]),
		Token:     parseTokenFromAuthorization(GetFirstOfSlices(md["authorization"])),
		SecretKey: GetFirstOfSlices(md["secret_key"]),
		Level:     GetFirstOfSlices(md["level"]),
		RoleId:    GetFirstOfSlices(md["role_id"]),
		SystemId:  GetFirstOfSlices(md["system_id"]),
	}
}

func GetFirstOfSlices(s []string) (str string) {
	if len(s) > 0 {
		str = s[0]
	}
	return
}

func StringToSlice(str string) []string {
	return []string{str}
}

func parseTokenFromAuthorization(authorizationHeader string) string {
	return strings.Replace(authorizationHeader, "Bearer ", "", 1)
}

func ParseHeaderToUserDev(ctx context.Context) *model.AuthUser {
	return &model.AuthUser{}
}

func HandleAuthWithChat(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	if GetFirstOfSlices(md["is_http_request"]) == "true" {
		return ctx, nil
	}
	if len(md["secret_key"]) > 0 {
		return ctx, nil
	}
	authorizationHeader := parseTokenFromAuthorization(GetFirstOfSlices(md["authorization"]))
	if len(authorizationHeader) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	ctx = metadata.NewIncomingContext(ctx, md)
	return ctx, nil
}
