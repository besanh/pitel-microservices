package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/tel4vn/fins-microservices/model"
	"google.golang.org/grpc/metadata"
)

func GetUser(c *gin.Context) (*model.AuthUser, bool) {
	tmp, isExist := c.Get("user")
	if isExist {
		user, ok := tmp.(*model.AuthUser)
		return user, ok
	} else {
		return nil, false
	}
}

func GetUserFromContext(ctx context.Context) (*model.AuthUser, bool) {
	tmp := ctx.Value("user")
	if tmp != nil {
		user, ok := tmp.(*model.AuthUser)
		return user, ok
	} else {
		return nil, false
	}
}

func HandleMetadata(ctx context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)
	md["user_id"] = r.Header.Get("user-id")
	md["tenant_id"] = r.Header.Get("tenant-id")
	md["username"] = r.Header.Get("username")
	md["level"] = r.Header.Get("level")
	md["token"] = parseTokenFromAuthorization(r.Header.Get("token"))
	md["role_id"] = r.Header.Get("role-id")
	md["secret_key"] = r.Header.Get("secret-key")
	md["system_id"] = r.Header.Get("system-key")

	return metadata.New(md)
}

type GoAuthInfo interface {
	auth.Info
}

func NewGoAuthUser(userId, username, tenantId, roleId string, level string) GoAuthInfo {
	user := &model.AuthUser{
		TenantId: tenantId,
		UserId:   userId,
		Username: username,
		RoleId:   roleId,
		Level:    level,
	}
	return user
}
