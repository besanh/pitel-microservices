package auth

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
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

type IAuthMiddleware interface {
	GRPCAuthMiddleware(ctx context.Context) (context.Context, error)
}

var AuthMdw IAuthMiddleware

type GatewayAuthMiddleware struct {
	env string // dev, prod
}

const (
	DEV  = "dev"
	PROD = "prod"
)

func NewGatewayAuthMiddleware(env string) IAuthMiddleware {
	mdw := &GatewayAuthMiddleware{
		env: DEV,
	}
	if env == PROD {
		mdw.env = PROD
	}
	return mdw
}

// AuthFunc is a middleware (interceptor) that extracts token from header
func (mdw *GatewayAuthMiddleware) GRPCAuthMiddleware(ctx context.Context) (context.Context, error) {
	isAuthenticaed, ok := ctx.Value(AUTHENTICATED).(bool)
	if !ok || !isAuthenticaed {
		var user *model.AuthUser
		if mdw.env == DEV {
			user = ParseHeaderToUserDev(ctx)
		} else {
			user = ParseHeaderToUser(ctx)
		}
		if user == nil || len(user.UserId) < 1 {
			return nil, status.Errorf(codes.Unauthenticated, ERR_TOKEN_IS_INVALID)
		}
		ctx = context.WithValue(ctx, AUTHENTICATED, true)
		ctx = context.WithValue(ctx, USER, user)
	}
	return ctx, nil
}

func ParseHeaderToUser(ctx context.Context) *model.AuthUser {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Debugf("no metadata")
		return nil
	}
	return &model.AuthUser{
		TenantId:         GetFirstOfSlices(md["tenant_id"]),
		BusinessUnitId:   GetFirstOfSlices(md["business_unit_id"]),
		UserId:           GetFirstOfSlices(md["user_id"]),
		Username:         GetFirstOfSlices(md["username"]),
		Services:         md["services"],
		DatabaseName:     GetFirstOfSlices(md["database_name"]),
		DatabasePort:     util.ParseInt(GetFirstOfSlices(md["database_port"])),
		DatabaseHost:     GetFirstOfSlices(md["database_host"]),
		DatabaseUser:     GetFirstOfSlices(md["database_user"]),
		DatabasePassword: GetFirstOfSlices(md["database_password"]),
	}
}

func GetFirstOfSlices(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

// TENANT_ID = 8d264455-956c-4450-9338-673748fc07aa
// TENANT_NAME = dev
// BUSINESS_UNIT_ID = d7ee56b1-dc9d-4e23-9847-9c99c6361137
// BUSINESS_UNIT_NAME = dev
// ROLE_ID = 1368f569-c943-44ba-92e7-b9c9be851205
// USER_ID = 4755c226-f404-4df9-9b4c-27fddf7d1418
// USERNAME = fins_dev
// DATABASE_NAME = dev_fins_dev
// DATABASE_HOST = 42.96.44.195
// DATABASE_PORT = 9000
// DATABASE_USER = tel4vnDBAdmin
// DATABASE_PASSWORD = Tel4vn@PsWrd#202399
func ParseHeaderToUserDev(ctx context.Context) *model.AuthUser {
	return &model.AuthUser{
		TenantId:         "8d264455-956c-4450-9338-673748fc07aa",
		BusinessUnitId:   "d7ee56b1-dc9d-4e23-9847-9c99c6361137",
		UserId:           "4755c226-f404-4df9-9b4c-27fddf7d1418",
		Username:         "fins_dev",
		Services:         []string{},
		DatabaseName:     "dev_fins_collection",
		DatabasePort:     5432,
		DatabaseHost:     "103.56.162.66",
		DatabaseUser:     "fins_api",
		DatabasePassword: "FinS#1D1B1#!2023",
	}
}

// Authorization unary interceptor function to handle authorize per RPC call
func GRPCAuthInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	log.Infof("method name: %s", info.FullMethod)
	start := time.Now()
	if util.InArrayContains(info.FullMethod, []string{"proto.inbox_marketing_incom"}) {
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
