package api

import (
	"context"
	"net/http"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/internal/goauth"
	authMdw "github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type (
	APIIBKToken struct {
	}
)

func RegisterAPIIBKToken(api hureg.APIGen) {
	handler := &APIIBKToken{}

	apiGroup := api.AddBasePath("/inbox-marketing/v1/token")

	hureg.Register(apiGroup, huma.Operation{
		Tags:        []string{"IBK Token"},
		OperationID: "Verify Token",
		Method:      http.MethodPost,
		Path:        "verify",
	}, handler.VerifyToken)

	hureg.Register(apiGroup, huma.Operation{
		Tags:        []string{"IBK Token"},
		OperationID: "Refresh Token",
		Method:      http.MethodPost,
		Path:        "refresh",
	}, handler.RefreshToken)
}

func (h *APIIBKToken) VerifyToken(ctx context.Context, req *struct {
	XTenantId string `header:"X-Tenant-Id"`
	Body      struct {
		Token string `json:"token" required:"true" minLength:"5"`
	}
}) (res *response.GenericResponse[goauth.AuthUser], err error) {

	token := req.Body.Token
	var user *goauth.AuthUser
	if token == authMdw.SECRET_TOKEN {
		user = &goauth.AuthUser{
			UserId:       "22d4859b-77f8-436a-a8d6-7fa61ba3dede",
			Token:        authMdw.SECRET_TOKEN,
			RefreshToken: "",
			Data: model.AuthUserData{
				TenantId:       "",
				BusinessUnitId: "",
				UserId:         "a08707a5-459a-466d-8e5c-9fcc676c867a",
			},
		}
		xTenantId := req.XTenantId
		if len(xTenantId) > 0 {
			var tenant *model.IBKTenant
			tenant, err = service.IBKTenantService.GetById(ctx, xTenantId)
			if err != nil {
				err = response.HandleError(err)
				return
			} else if tenant != nil {
				user.Data = model.AuthUserData{
					TenantId:       xTenantId,
					BusinessUnitId: "",
					UserId:         "a08707a5-459a-466d-8e5c-9fcc676c867a",
				}
			}
		}
	} else {
		user, err = service.AuthService.VerifyToken(ctx, token)
		if err != nil {
			err = response.HandleError(err)
			return
		}
	}
	res = response.OK(*user)
	return
}

func (h *APIIBKToken) RefreshToken(ctx context.Context, req *struct {
	Body model.RefreshTokenRequest
}) (res *response.GenericResponse[model.RefreshTokenResponse], err error) {
	var result *model.RefreshTokenResponse
	result, err = service.AuthService.RefreshToken(ctx, req.Body)
	if err != nil {
		err = response.HandleError(err)
		return
	}
	res = response.OK(*result)
	return
}
