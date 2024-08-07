package auth

import (
	"context"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/goauth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

var SECRET_TOKEN string

func validateToken(ctx context.Context, tokenString string) (authInfo auth.Info, err error) {
	var user *goauth.AuthUser
	if len(tokenString) > 0 && tokenString == SECRET_TOKEN {

	} else {
		user, err = service.AuthService.VerifyToken(ctx, tokenString)
		if err != nil {
			return
		}
	}
	data := &model.AuthUserData{}
	if err = util.ParseAnyToAny(user.Data, data); err != nil {
		log.Error(err)
		return
	}

	authInfo = NewGoAuthUser(data.UserId, data.Username, data.TenantId, data.BusinessUnitId, user.Level, user.RoleId, data.Scopes)
	return
}
