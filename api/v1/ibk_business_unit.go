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

type (
	APIIBKBusinessUnit struct {
	}
)

func RegisterAPIIBKBusinessUnit(api hureg.APIGen) {
	handler := &APIIBKBusinessUnit{}

	group := api.AddBasePath("/inbox-marketing/v1/business-unit")
	tags := []string{"IBK Business Unit"}

	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Get Business Units",
		Method:      http.MethodGet,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.GetBusinessUnits)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Create Business Unit",
		Method:      http.MethodPost,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.CreateBusinessUnit)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Get Business Unit By Id",
		Method:      http.MethodGet,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.GetBusinessUnitById)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Update Business Unit",
		Method:      http.MethodPut,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.PutBusinessUnit)
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "Delete Business Unit",
		Method:      http.MethodDelete,
		Path:        "{id}",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.DeleteBusinessUnit)
}

func (h *APIIBKBusinessUnit) GetBusinessUnits(c context.Context, req *struct {
	model.IBKBusinessUnitQueryParam
	Limit  int `query:"limit" default:"50" min:"1" max:"2999"`
	Offset int `query:"offset" default:"0" min:"0" max:"2999"`
}) (res *response.PaginationResponse[[]*model.IBKBusinessUnit], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}

	total, businessUnits, err := service.IBKBusinessUnitService.GetBusinessUnits(c, authUser, req.IBKBusinessUnitQueryParam, req.Limit, req.Offset)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.Pagination(total, businessUnits)
	return
}

func (h *APIIBKBusinessUnit) CreateBusinessUnit(c context.Context, req *struct {
	Body model.IBKBusinessUnitBody
}) (res *response.IdResponse, err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	id, err := service.IBKBusinessUnitService.CreateBusinessUnit(c, authUser, &req.Body)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = &response.IdResponse{
		Id: id,
	}
	return
}

func (h *APIIBKBusinessUnit) PutBusinessUnit(c context.Context, req *struct {
	Body model.IBKBusinessUnitBody
	Id   string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[any], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}

	err = service.IBKBusinessUnitService.PutBusinessUnit(c, authUser, req.Id, &req.Body)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.OK[any](nil)
	return
}

func (h *APIIBKBusinessUnit) GetBusinessUnitById(c context.Context, req *struct {
	Id string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[*model.IBKBusinessUnit], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	businessUnit, err := service.IBKBusinessUnitService.GetBusinessUnitById(c, authUser, req.Id)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.OK(businessUnit)
	return
}
func (h *APIIBKBusinessUnit) DeleteBusinessUnit(c context.Context, req *struct {
	Id string `path:"id" required:"true" format:"uuid"`
}) (res *response.GenericResponse[any], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	err = service.IBKBusinessUnitService.DeleteBusinessUnit(c, authUser, req.Id)
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	res = response.OK[any](nil)
	return
}
