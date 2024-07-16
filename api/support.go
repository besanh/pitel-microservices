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
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"nhooyr.io/websocket"
)

const (
	AUTHEN_TOKEN = "authen_token"
	USER_INFO    = "user_info"
)

func AuthMiddleware(c *gin.Context) *model.ChatResponse {
	bssAuthRequest := model.BssAuthRequest{
		Token:   c.Query("token"),
		AuthUrl: c.Query("auth-url"),
		Source:  c.Query("source"),
	}

	if len(c.GetHeader("validator-header")) > 0 {
		bssAuthRequest = model.BssAuthRequest{
			Token:   c.GetHeader("token"),
			AuthUrl: c.GetHeader("auth-url"),
			Source:  c.GetHeader("source"),
		}
	}

	res := AAAMiddleware(c, bssAuthRequest)
	if res == nil {
		c.JSON(http.StatusUnauthorized, map[string]any{
			"error": http.StatusText(http.StatusUnauthorized),
		})
		c.Abort()
		return nil
	}

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

func AAAMiddleware(ctx *gin.Context, bssAuthRequest model.BssAuthRequest) (result *model.ChatResponse) {
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

func RequestAAA(ctx *gin.Context, bssAuthRequest model.BssAuthRequest) (result *model.ChatResponse, err error) {
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
		Post(service.OTT_URL + "/aaa")
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, err
	}
	return result, nil
}

func RequestAuthen(ctx *gin.Context, bssAuthRequest model.BssAuthRequest) (result *model.ChatResponse, err error) {
	clientInfo := resty.New()
	urlInfo := bssAuthRequest.AuthUrl + "/v1/crm/auth/auth-info"
	res, err := clientInfo.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+bssAuthRequest.Token).
		Get(urlInfo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var resInfo map[string]any
	if err := json.Unmarshal(res.Body(), &resInfo); err != nil {
		log.Error(err)
		return result, err
	}

	userUuid, _ := resInfo["user_uuid"].(string)
	if len(userUuid) < 1 {
		log.Errorf("user_uuid %s is invalid", userUuid)
		return nil, errors.New("invalid user uuid")
	}

	// Get Info user
	userInfo := model.AuthUserInfo{}
	userInfoCache := cache.RCache.Get(USER_INFO + "_" + bssAuthRequest.Token)
	if userInfoCache != nil {
		if err := json.Unmarshal([]byte(userInfoCache.(string)), &userInfo); err != nil {
			log.Error(err)
			return result, err
		}
	} else {
		url := bssAuthRequest.AuthUrl + "/v1/crm/user-crm/" + userUuid
		res, err := clientInfo.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+bssAuthRequest.Token).
			Get(url)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		var resp map[string]any
		if err := json.Unmarshal(res.Body(), &resp); err != nil {
			log.Error(err)
			return result, err
		}

		userInfo.UserUuid, _ = resp["user_uuid"].(string)
		userInfo.DomainUuid, _ = resp["domain_uuid"].(string)
		userInfo.Username, _ = resp["username"].(string)
		userInfo.Password, _ = resp["password"].(string)
		userInfo.ApiKey, _ = resp["api_key"].(string)
		userInfo.UserEnabled, _ = resp["user_enabled"].(string)
		userInfo.UserStatus, _ = resp["user_status"].(string)
		userInfo.Level, _ = resp["level"].(string)
		userInfo.LastName, _ = resp["last_name"].(string)
		userInfo.MiddleName, _ = resp["middle_name"].(string)
		userInfo.FirstName, _ = resp["first_name"].(string)
		userInfo.UnitUuid, _ = resp["unit_uuid"].(string)
		userInfo.UnitName, _ = resp["unit_name"].(string)
		userInfo.RoleUuid, _ = resp["role_uuid"].(string)

		cache.RCache.Set(USER_INFO+"_"+bssAuthRequest.Token, userInfo, 3*time.Minute)
	}

	if len(userInfo.UserUuid) > 1 {
		result = &model.ChatResponse{
			Data: &model.AuthUser{
				TenantId: userInfo.DomainUuid,
				UserId:   userInfo.UserUuid,
				Username: userInfo.Username,
				Level:    userInfo.Level,
				Source:   bssAuthRequest.Source,
				Token:    bssAuthRequest.Token,
				Fullname: userInfo.FirstName + " " + userInfo.MiddleName + " " + userInfo.LastName,
			},
		}
	} else {
		cache.RCache.Del([]string{USER_INFO + "_" + bssAuthRequest.Token})
		return nil, errors.New("failed to get user info")
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

func AuthMiddlewareNew(c *gin.Context) (result *model.ChatResponse) {
	if len(c.Query("token")) < 1 || len(c.Query("system-key")) < 1 {
		c.JSON(http.StatusUnauthorized, map[string]any{
			"error": http.StatusText(http.StatusUnauthorized),
		})
		c.Abort()
		return
	}
	result = auth.ChatMiddleware(c, c.Query("token"), c.Query("system-key"))
	if result == nil {
		c.JSON(http.StatusUnauthorized, map[string]any{
			"error": http.StatusText(http.StatusUnauthorized),
		})
		c.Abort()
		return
	}
	return
}
