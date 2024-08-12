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
		OperationID: "Get User",
		Method:      http.MethodGet,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.GetUser)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Create User",
		Method:      http.MethodPost,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.CreateUser)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Get User By Id",
		Method:      http.MethodGet,
		Path:        "{id}",
	}, handler.GetUserById)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Update User",
		Method:      http.MethodPut,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.PutUser)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Delete User",
		Method:      http.MethodDelete,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.DeleteUser)
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

func (h *APIIBKUser) GetUserById(c context.Context, req *struct {
	Id string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[*model.IBKUserInfo], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	user, err := service.IBKUserService.GetUserById(c, authUser, req.Id)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.OK(user)
	return
}

func (h *APIIBKUser) CreateUser(c context.Context, req *struct {
	Body model.IBKUserBody
}) (res *response.IdResponse, err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	id, err := service.IBKUserService.CreateUser(c, authUser, &req.Body)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}

	res = &response.IdResponse{
		Id: id,
	}

	return
}

func (h *APIIBKUser) PutUser(c context.Context, req *struct {
	Body model.IBKUserBody
	Id   string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[any], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	err = service.IBKUserService.PutUser(c, authUser, req.Id, &req.Body)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.OK[any](nil)
	return
}

func (h *APIIBKUser) DeleteUser(c context.Context, req *struct {
	Id string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[any], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	err = service.IBKUserService.DeleteUser(c, authUser, req.Id)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.OK[any](nil)
	return
}
