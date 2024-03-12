package api

import (
	"context"
	"encoding/json"
	"errors"
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

func AuthMiddleware(c *gin.Context) *model.AAAResponse {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth_url"),
		Source:  c.Query("source"),
	}

	if len(c.GetHeader("validator_header")) > 0 {
		bssAuthRequest = model.BssAuthRequest{
			Token:   c.GetHeader("token"),
			AuthUrl: c.GetHeader("auth_url"),
			Source:  c.GetHeader("source"),
		}
	}

	res := AAAMiddleware(c, bssAuthRequest)

	return res
}

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

func AAAMiddleware(ctx *gin.Context, bssAuthRequest model.BssAuthRequest) (result *model.AAAResponse) {
	if len(bssAuthRequest.Token) < 10 {
		return nil
	}
	if bssAuthRequest.Source == "authen" {
		result, err := RequestAuthen(ctx, bssAuthRequest)
		if err != nil {
			return nil
		}
		return result
	} else if bssAuthRequest.Source == "aaa" {
		result, err := RequestAAA(ctx, bssAuthRequest)
		if err != nil {
			return nil
		}
		return result
	}
	return
}

func RequestAAA(ctx *gin.Context, bssAuthRequest model.BssAuthRequest) (result *model.AAAResponse, err error) {
	body := map[string]string{
		"token": bssAuthRequest.Token,
	}
	// https://api.dev.fins.vn/aaa/v1/token/verify
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+bssAuthRequest.Token).
		SetBody(body).
		SetResult(result).
		Post(bssAuthRequest.AuthUrl)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, err
	}
	return result, nil
}

func RequestAuthen(ctx *gin.Context, bssAuthRequest model.BssAuthRequest) (result *model.AAAResponse, err error) {
	clientInfo := resty.New()
	urlInfo := bssAuthRequest.AuthUrl + "/v1/crm/auth/auth-info"
	res, err := clientInfo.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+bssAuthRequest.Token).
		Get(urlInfo)
	if err != nil {
		return nil, err
	}
	var resInfo map[string]any
	if err := json.Unmarshal(res.Body(), &resInfo); err != nil {
		return result, err
	}
	userUuid, _ := resInfo["user_uuid"].(string)
	if len(userUuid) < 1 {
		return nil, errors.New("invalid user uuid")
	}

	// Get Info agent
	agentInfo := model.AuthUserInfo{}
	agentInfoCache := cache.MCache.Get(AGENT_INFO + "_" + bssAuthRequest.Token)
	if agentInfoCache != nil {
		if err := util.ParseAnyToAny(agentInfoCache, &agentInfo); err != nil {
			log.Error(err)
			return nil, err
		}
	} else {
		url := bssAuthRequest.AuthUrl + "/v1/crm/user-crm/" + userUuid
		res, err := clientInfo.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+bssAuthRequest.Token).
			Get(url)
		if err != nil {
			return nil, err
		}
		var resp map[string]any
		if err := json.Unmarshal(res.Body(), &resp); err != nil {
			return result, err
		}

		agentInfo.UserUuid, _ = resp["user_uuid"].(string)
		agentInfo.DomainUuid, _ = resp["domain_uuid"].(string)
		agentInfo.Username, _ = resp["username"].(string)
		agentInfo.Password, _ = resp["password"].(string)
		agentInfo.ApiKey, _ = resp["api_key"].(string)
		agentInfo.UserEnabled, _ = resp["user_enabled"].(string)
		agentInfo.UserStatus, _ = resp["user_status"].(string)
		agentInfo.Level, _ = resp["level"].(string)
		agentInfo.LastName, _ = resp["last_name"].(string)
		agentInfo.MiddleName, _ = resp["middle_name"].(string)
		agentInfo.FirstName, _ = resp["first_name"].(string)
		agentInfo.UnitUuid, _ = resp["unit_uuid"].(string)
		agentInfo.UnitName, _ = resp["unit_name"].(string)
		agentInfo.RoleUuid, _ = resp["role_uuid"].(string)
		agentInfo.Extension, _ = resp["extension"].(string)
		agentInfo.ExtensionUuid, _ = resp["extension_uuid"].(string)

		cache.MCache.Set(AGENT_INFO+"_"+bssAuthRequest.Token, agentInfo, 1*time.Minute)
	}

	if len(agentInfo.UserUuid) > 1 {
		result = &model.AAAResponse{
			Data: &model.AuthUser{
				TenantId: agentInfo.DomainUuid,
				UserId:   agentInfo.UserUuid,
				Username: agentInfo.Username,
				Level:    agentInfo.Level,
				Source:   bssAuthRequest.Source,
				Token:    bssAuthRequest.Token,
				UnitUuid: agentInfo.UnitUuid,
			},
		}
	} else {
		cache.MCache.Del(AGENT_INFO + "_" + bssAuthRequest.Token)
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
