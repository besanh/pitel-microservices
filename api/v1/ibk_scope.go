package api

import (
	"context"
	"net/http"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/response"
	authMdw "github.com/tel4vn/fins-microservices/middleware/auth"
)

type APIIBKScope struct {
}

func RegisterAPIIBKScope(api hureg.APIGen) {
	handler := &APIIBKScope{}

	apiGroup := api.AddBasePath("/inbox-marketing/v1/scope")
	tags := []string{"IBK Scopes"}
	hureg.Register(apiGroup, huma.Operation{
		Tags:        tags,
		OperationID: "Get Scopes",
		Method:      http.MethodGet,
		Path:        "",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.GetScopes)

}

func (h *APIIBKScope) GetScopes(ctx context.Context, req *struct {
	Limit  int `query:"limit" default:"50" min:"1" max:"2999"`
	Offset int `query:"offset" default:"0" min:"0" max:"2999"`
}) (res *response.PaginationResponse[[]constants.API_SCOPE], err error) {
	res = response.Pagination(len(constants.API_SCOPES), constants.API_SCOPES)
	return
}
