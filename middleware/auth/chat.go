package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/pitel-microservices/common/cache"
	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
	"github.com/tel4vn/pitel-microservices/service"
)

const (
	AUTHEN_TOKEN = "authen_token"
	USER_INFO    = "user_info"
)

func ChatMiddleware(ctx context.Context, token, systemId string) (result *model.ChatResponse) {
	var integrateType string
	chatIntegrateSystem := &model.ChatIntegrateSystem{}
	chatIScache := cache.RCache.Get(service.CHAT_INTEGRATE_SYSTEM + "_" + systemId)
	if chatIScache != nil {
		if err := json.Unmarshal([]byte(chatIScache.(string)), chatIntegrateSystem); err != nil {
			log.Error(err)
			return
		}
	} else {
		_, chatIntegrateSystems, errTmp := repository.ChatIntegrateSystemRepo.GetIntegrateSystems(ctx, repository.DBConn, model.ChatIntegrateSystemFilter{
			SystemId: systemId}, 1, 0)
		if errTmp != nil {
			log.Error(errTmp)
			return
		} else if len(*chatIntegrateSystems) < 1 {
			log.Error("invalid system id " + systemId)
			return
		}

		chatIntegrateSystem = &(*chatIntegrateSystems)[0]

		if err := cache.RCache.Set(service.CHAT_INTEGRATE_SYSTEM+"_"+systemId, chatIntegrateSystem, service.CHAT_INTEGRATE_SYSTEM_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}
	integrateType = chatIntegrateSystem.InfoSystem.AuthType

	bssAuthRequest := model.BssAuthRequest{
		ApiUrl:        chatIntegrateSystem.InfoSystem.ApiUrl,
		AuthUrl:       chatIntegrateSystem.InfoSystem.ApiAuthUrl,
		Token:         token,
		UserDetailUrl: chatIntegrateSystem.InfoSystem.ApiGetUserDetailUrl,
		ServerName:    chatIntegrateSystem.InfoSystem.ServerName,
	}

	switch integrateType {
	case "pitel_crm":
		result, err := CrmMiddleware(ctx, token, systemId, bssAuthRequest)
		if err != nil {
			return nil
		}
		return result
	default:
		return
	}
}

func CrmMiddleware(ctx context.Context, token, systemId string, bssAuthRequest model.BssAuthRequest) (result *model.ChatResponse, err error) {
	// Get Info user
	clientInfo := resty.New()
	urlInfo := bssAuthRequest.AuthUrl
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
		url := bssAuthRequest.UserDetailUrl + "/" + userUuid
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
				ApiUrl:   bssAuthRequest.ApiUrl,
				SystemId: systemId,
				Server:   bssAuthRequest.ServerName,
			},
		}
	} else {
		cache.RCache.Del([]string{USER_INFO + "_" + bssAuthRequest.Token})
		return nil, errors.New("failed to get user info")
	}

	return result, nil
}
