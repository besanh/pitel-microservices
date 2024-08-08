package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tel4vn/fins-microservices/common/response"
	authMdw "github.com/tel4vn/fins-microservices/middleware/auth"
	fpMdw "github.com/tel4vn/fins-microservices/middleware/fingerprint"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type (
	APIIBKAuth struct{}
)

func RegisterAPIIBKAuth(api hureg.APIGen) {
	handler := &APIIBKAuth{}

	group := api.AddBasePath("/inbox-marketing/v1")
	tags := []string{"IBK Auth"}
	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "login",
		Middlewares: huma.Middlewares{fpMdw.FingerprintMiddleware},
	}, handler.Login)

	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "logout",
		Method:      http.MethodPost,
		Path:        "logout",
		Middlewares: huma.Middlewares{fpMdw.FingerprintMiddleware},
	}, handler.Logout)

	hureg.Register(group, huma.Operation{
		Tags:        tags,
		OperationID: "refresh",
		Method:      http.MethodPost,
		Path:        "/auth/refresh",
		Security:    authMdw.DefaultAuthSecurity,
	}, handler.RefreshToken)
}

func (h *APIIBKAuth) Login(c context.Context, req *struct {
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

func (h *APIIBKAuth) Logout(c context.Context, req *struct{}) (res *response.GenericResponse[any], err error) {
	fp, _ := c.Value("fingerprint").(string)
	token, _ := c.Value("token").(string)
	if err = service.AuthService.Logout(c, token, fp); err != nil {
		err = response.HandleError(err)
		return
	}
	res = response.OK[any](nil)
	return
}

func (h *APIIBKAuth) RefreshToken(c context.Context, req *struct {
	Token string `header:"Authorization"`
}) (res *response.GenericResponse[any], err error) {
	authUser, ok := authMdw.GetUserFromContext(c)
	if !ok {
		err = response.ErrUnauthorized()
		return
	}
	req.Token = strings.Replace(req.Token, "Bearer ", "", -1)
	err = service.AuthService.RefreshAuthData(c, authUser, req.Token)
	if err != nil {
		err = response.HandleError(err)
		return
	}
	res = response.OK[any](nil)
	return
}
