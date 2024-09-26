package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/tel4vn/pitel-microservices/model"
)

type GoAuthInfo interface {
	auth.Info
}

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

func NewGoAuthUser(userId, username, tenantId, roleId, level, systemId, secretKey string) GoAuthInfo {
	user := &model.AuthUser{
		TenantId:  tenantId,
		UserId:    userId,
		Username:  username,
		RoleId:    roleId,
		Level:     level,
		SystemId:  systemId,
		SecretKey: secretKey,
	}
	return user
}
