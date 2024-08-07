package v1

import (
	"context"
	"net/http"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tel4vn/fins-microservices/common/response"
	fpMdw "github.com/tel4vn/fins-microservices/middleware/fingerprint"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type (
	APIAuth struct{}
)

func RegisterAPIAuth(api hureg.APIGen) {
	handler := &APIAuth{}

	group := api.AddBasePath("bss-inbox-marketing/v1")
	hureg.Register(group, huma.Operation{
		Tags:        []string{"auth"},
		OperationID: "Auth",
		Method:      http.MethodPost,
		Path:        "auth",
		Middlewares: huma.Middlewares{fpMdw.FingerprintMiddleware},
	}, handler.Login)
}

func (h *APIAuth) Login(c context.Context, req *struct {
	Body model.LoginRequest
}) (res *response.GenericResponse[model.LoginResponse], err error) {
	fp, _ := c.Value("fingerprint").(string)
	userAgent, _ := c.Value("user_agent").(string)
	loginRequest := model.LoginRequest{
		Username:    req.Body.Username,
		Password:    req.Body.Password,
		UserAgent:   userAgent,
		Fingerprint: fp,
	}
	var result *model.LoginResponse
	result, err = service.AuthService.Login(c, loginRequest)
	if err != nil {
		err = response.HandleError(err)
		return
	}
	res = response.OK(*result)
	return
}
