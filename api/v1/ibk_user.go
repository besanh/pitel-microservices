package api

import (
	"context"
	"net/http"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	authMdw "github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type APIIBKUser struct{}

func RegisterAPIIBKUser(api hureg.APIGen) {
	handler := &APIIBKUser{}

	group := api.AddBasePath("/inbox-marketing/v1/user")
	tags := []string{"IBK User"}
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.GetUser)
}

func (h *APIIBKUser) GetUser(c context.Context, req *struct {
	model.IBKUserQueryParam
	Limit  int `query:"limit" default:"50" min:"1" max:"2999"`
	Offset int `query:"offset" default:"0" min:"0" max:"2999"`
}) (res *response.PaginationResponse[[]*model.IBKUserInfo], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}

	total, users, err := service.IBKUserService.GetUsers(c, authUser, req.IBKUserQueryParam, req.Limit, req.Offset)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.Pagination(total, users)
	return
}
