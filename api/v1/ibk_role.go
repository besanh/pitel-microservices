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

type APIIBKRole struct{}

func RegisterAPIIBKRole(api hureg.APIGen) {
	handler := &APIIBKRole{}

	group := api.AddBasePath("/inbox-marketing/v1/role")
	tags := []string{"IBK Role"}

	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Get Roles",
		Method:      http.MethodGet,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.GetRoles)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Create Role",
		Method:      http.MethodPost,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.CreateRole)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Get Role By Id",
		Method:      http.MethodGet,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.GetRoleById)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Update Role",
		Method:      http.MethodPut,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.PutRole)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Delete Role",
		Method:      http.MethodDelete,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.DeleteRole)
}

func (api *APIIBKRole) GetRoles(c context.Context, req *struct {
	model.IBKRoleQueryParam
	Limit  int `query:"limit" default:"50" min:"1" max:"2999"`
	Offset int `query:"offset" default:"0" min:"0" max:"2999"`
}) (res *response.PaginationResponse[[]*model.IBKRoleInfo], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}

	total, roles, err := service.IBKRoleService.GetRoles(c, authUser, req.IBKRoleQueryParam, req.Limit, req.Offset)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.Pagination(total, roles)
	return
}

func (api *APIIBKRole) CreateRole(c context.Context, req *struct {
	model.IBKRoleBody
}) (res *response.IdResponse, err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}

	id, err := service.IBKRoleService.CreateRole(c, authUser, &req.IBKRoleBody)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = &response.IdResponse{
		Id: id,
	}
	return
}

func (api *APIIBKRole) GetRoleById(c context.Context, req *struct {
	Id string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[*model.IBKRoleInfo], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}

	role, err := service.IBKRoleService.GetRoleById(c, authUser, req.Id)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.OK(role)
	return
}

func (api *APIIBKRole) PutRole(c context.Context, req *struct {
	Body model.IBKRoleBody
	Id   string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[any], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	err = service.IBKRoleService.PutRole(c, authUser, req.Id, &req.Body)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	return
}

func (api *APIIBKRole) DeleteRole(c context.Context, req *struct {
	Id string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[any], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	err = service.IBKRoleService.DeleteRole(c, authUser, req.Id)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	return
}
