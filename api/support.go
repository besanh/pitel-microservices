package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/model"
	"nhooyr.io/websocket"
)

var (
	AAA_URL = "https://api.dev.fins.vn/aaa"
)

func MoveTokenToHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if len(token) < 10 {
			ctx.JSON(http.StatusUnauthorized, map[string]any{
				"error": http.StatusText(http.StatusUnauthorized),
			})
			ctx.Abort()
			return
		}
		ctx.Request.Header.Set("Authorization", "Bearer "+token)
		ctx.Next()
	}
}

func AAAMiddleware(ctx *gin.Context) (result *model.AAAResponse) {
	token := ctx.Query("token")
	if len(token) < 10 {
		return nil
	}
	if ctx.Query("source") == "authen" {
		result, err := RequestAuthen(ctx, token)
		if err != nil {
			return nil
		}
		return result
	} else if ctx.Query("source") == "aaa" {
		result, err := RequestAAA(ctx, token)
		if err != nil {
			return nil
		}
		return result
	}
	return
}

func RequestAAA(ctx *gin.Context, token string) (result *model.AAAResponse, err error) {
	body := map[string]string{
		"token": token,
	}
	// https://api.dev.fins.vn/aaa/v1/token/verify
	url := ctx.Query("auth_url")
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		SetResult(result).
		Post(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, err
	}
	return result, nil
}

func RequestAuthen(ctx *gin.Context, token string) (result *model.AAAResponse, err error) {
	body := map[string]string{
		"token": token,
	}
	resp := model.Authen{}
	// https://api-loadtest.tel4vn.com/v1/crm/auth
	url := ctx.Query("auth_url")
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		SetResult(resp).
		Post(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, err
	}
	result = &model.AAAResponse{
		Data: &model.AuthUser{
			TenantId: resp.DomainUuid,
			UserId:   resp.UserUuid,
			Username: resp.Username,
			Level:    resp.Level,
		},
	}
	return result, nil
}

func WriteTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := c.Write(ctx, websocket.MessageText, msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}
