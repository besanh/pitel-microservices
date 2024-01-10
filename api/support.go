package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"nhooyr.io/websocket"
)

const (
	AUTHEN_TOKEN = "authen_token"
	AGENT_INFO   = "agent_info"
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

func RequestAuthen(ctx *gin.Context, apiKey string) (result *model.AAAResponse, err error) {
	body := map[string]string{
		"api_key": apiKey,
	}
	var token string
	resp := model.Authen{}
	tokenCache := cache.MCache.Get(AUTHEN_TOKEN + "_" + apiKey)
	if tokenCache != nil {
		if err := util.ParseAnyToAny(tokenCache, &resp); err != nil {
			log.Error(err)
			return nil, err
		}
	} else {
		// https://api-loadtest.tel4vn.com/v3/auth/token
		url := ctx.Query("auth_url")
		client := resty.New()
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			SetResult(resp).
			Post(url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode() != 200 {
			return nil, err
		}
		if err := json.Unmarshal(res.Body(), &resp); err != nil {
			return result, err
		}
		cache.MCache.Set(AUTHEN_TOKEN+"_"+apiKey, resp, 1*time.Minute)
	}

	// Get Info agent
	agentInfo := model.AuthUserInfo{}
	agentInfoCache := cache.MCache.Get(AGENT_INFO + "_" + token)
	if agentInfoCache != nil {
		if err := util.ParseAnyToAny(agentInfoCache, &agentInfo); err != nil {
			log.Error(err)
			return nil, err
		}
	} else {
		// https://api-loadtest.tel4vn.com/crm/user-crm
		urlInfo := "https://api-loadtest.tel4vn.com/v1/crm/user-crm" + "/" + resp.UserId
		clientInfo := resty.New()
		res, err := clientInfo.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+resp.Token).
			Get(urlInfo)
		if err != nil {
			return nil, err
		}
		if res.StatusCode() != 200 {
			return nil, err
		}
		if err := json.Unmarshal(res.Body(), &agentInfo); err != nil {
			return result, err
		}
		cache.MCache.Set(AGENT_INFO+"_"+token, agentInfo, 1*time.Minute)
	}

	if len(agentInfo.UserUuid) > 1 {
		result = &model.AAAResponse{
			Data: &model.AuthUser{
				TenantId: agentInfo.DomainUuid,
				UserId:   agentInfo.UserUuid,
				Username: agentInfo.Username,
				Level:    agentInfo.Level,
			},
		}
	} else {
		return nil, fmt.Errorf("failed to get user info")
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
