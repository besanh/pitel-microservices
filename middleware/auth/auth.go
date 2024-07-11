package auth

import (
	"context"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/token"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/goauth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

var cacheObj libcache.Cache
var strategy union.Union
var tokenStrategy auth.Strategy

var SECRET_TOKEN string

type IAuthMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
}

var AuthMdw IAuthMiddleware

func AuthMiddleware() gin.HandlerFunc {
	return AuthMdw.AuthMiddleware()
}

type LocalAuthMiddleware struct {
}

func NewLocalAuthMiddleware() IAuthMiddleware {
	return &LocalAuthMiddleware{}
}

func SetupGoGuardian() {
	cacheObj = libcache.FIFO.New(0)
	cacheObj.SetTTL(time.Minute * 10)
	tokenStrategy = token.New(validateTokenAuth, cacheObj)
	strategy = union.New(tokenStrategy)
}

func (auth *LocalAuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := strategy.AuthenticateRequest(c.Request)
		if err != nil {
			log.Error("invalid credentials")
			c.JSON(
				http.StatusUnauthorized,
				map[string]interface{}{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			c.Abort()
			return
		}
		c.Set("user", user)

	}
}

func validateToken(ctx context.Context, tokenString string) (authInfo auth.Info, t time.Time, err error) {
	var user *goauth.AuthUser
	user, err = service.ChatAuthService.VerifyToken(ctx, tokenString)
	if err != nil {
		return nil, time.Time{}, err
	}
	data := &model.TokenData{}
	if err := util.ParseAnyToAny(user.Data, data); err != nil {
		log.Error(err)
		return nil, time.Time{}, err
	}
	authInfo = NewGoAuthUser(data.UserId, data.Username, data.TenantId, data.RoleId, data.Level)
	return authInfo, time.Now(), nil
}

func validateTokenAuth(ctx context.Context, r *http.Request, tokenString string) (auth.Info, time.Time, error) {
	return validateToken(ctx, tokenString)
}

func GetAuthorizationHeader(c *gin.Context) string {
	authorizationHeader := c.GetHeader("Authorization")
	return authorizationHeader
}
