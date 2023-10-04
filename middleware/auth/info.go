package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/model"
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
