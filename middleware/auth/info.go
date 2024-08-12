package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
)

type GoAuthInfo interface {
	auth.Info
}

var ENV string

func GetUser(c *gin.Context) (user *model.AuthUser, ok bool) {
	tmp, isExist := c.Get("user")
	if isExist {
		user, ok = tmp.(*model.AuthUser)
		return
	} else {
		return
	}
}

func GetUserFromContext(ctx context.Context) (user *model.AuthUser, ok bool) {
	if ENV == "local" {
		return &model.AuthUser{
			AuthUserData: &model.AuthUserData{
				Level:          env.GetStringENV("USER_LEVEL", constants.SUPERADMIN),
				TenantId:       env.GetStringENV("USER_TENANT_ID", "6b84fd80-0dad-4755-9a43-14bf4fb70011"),
				BusinessUnitId: env.GetStringENV("USER_BUSINESS_UNIT_ID", "deca75f7-686a-4694-9b7c-e33518280ad5"),
				UserId:         env.GetStringENV("USER_ID", "8a4e8ee4-3cb1-4ff6-a355-1520d58408fe"),
				Username:       "usertest",
			},
		}, true
	}
	tmp := ctx.Value("user")
	if tmp != nil {
		user, ok = tmp.(*model.AuthUser)
		return
	} else {
		return
	}
}

func NewGoAuthUser(userId, username, tenantId, businessUnitId, level, roleId string, scopes []string) (user GoAuthInfo) {
	return &model.AuthUser{
		AuthUserData: &model.AuthUserData{
			Level:          level,
			TenantId:       tenantId,
			BusinessUnitId: businessUnitId,
			UserId:         userId,
			Username:       username,
			Scopes:         scopes,
			RoleId:         roleId,
		},
	}
}

func ValidateTokenQueryGin(c *gin.Context) {
	token := c.Query("token")
	if len(token) < 1 {
		c.JSON(http.StatusUnauthorized, response.ErrUnauthorized())
		c.Abort()
		return
	}
	// validate token user
	authInfo, err := validateToken(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrUnauthorized())
		c.Abort()
		return
	}
	c.Set("user", authInfo)
	c.Next()
}

func ValidateTokenGin(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) < 1 {
		c.JSON(http.StatusUnauthorized, response.ErrUnauthorized())
		c.Abort()
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	// validate token user
	authInfo, err := validateToken(c, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrUnauthorized())
		c.Abort()
		return
	}
	c.Set("user", authInfo)
	c.Next()
}
