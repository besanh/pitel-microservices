package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
)

type GoAuthInfo interface {
	auth.Info
}

var ENV string

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

func NewGoAuthUser(userId, username, tenantId, businessUnitId, level, roleId string, scopes []string) GoAuthInfo {
	user := &model.AuthUser{
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
	return user
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
