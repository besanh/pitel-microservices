package service

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	cacheUtil "github.com/tel4vn/fins-microservices/common/cache"
	httpUtil "github.com/tel4vn/fins-microservices/common/http"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"

	"github.com/go-resty/resty/v2"
)

const (
	INFO_ROUTING                   = "info_routing"
	INFO_EXTERNAL_PLUGIN_CONNECT   = "INFO_EXTERNAL_PLUGIN_CONNECT"
	EXPIRE_EXTERNAL_PLUGIN_CONNECT = 1 * time.Minute

	// token
	ABENLA_TOKEN            = "abenla_token"
	INCOM_TOKEN             = "incom_token"
	FPT_TOKEN               = "fpt_token"
	EXPIRE_EXTERNAL_ROUTING = 1 * time.Minute
)

func GetInfoPlugin(ctx context.Context, db sqlclient.ISqlClientConn, authUser *model.AuthUser, routingConfigUuid string) (*model.RoutingConfig, error) {
	routingConfigCache := cacheUtil.NewMemCache().Get(INFO_ROUTING + "_" + routingConfigUuid)
	if routingConfigCache != nil {
		routing := routingConfigCache.(*model.RoutingConfig)
		return routing, nil
	} else {
		routingConfig, err := repository.RoutingConfigRepo.GetById(ctx, db, routingConfigUuid)
		if err != nil {
			return nil, err
		} else if routingConfig == nil {
			return nil, errors.New("routing not found in system")
		}
		cacheUtil.NewMemCache().Set(INFO_ROUTING+"_"+routingConfigUuid, routingConfig, EXPIRE_EXTERNAL_ROUTING)
		return routingConfig, nil
	}
}

func CheckConnectionWithExternalPlugin(ctx context.Context, routingConfig model.RoutingConfig) error {
	params := map[string]string{}
	if routingConfig.RoutingOption.Abenla.Status {
		connectionAlive := cacheUtil.NewMemCache().Get(ABENLA_TOKEN)
		if connectionAlive != nil {
			return nil
		} else {
			if len(routingConfig.RoutingOption.Abenla.Username) > 0 {
				params["loginName"] = routingConfig.RoutingOption.Abenla.Username
			}
			if len(routingConfig.RoutingOption.Abenla.Password) > 0 {
				hasher := md5.New()
				hasher.Write([]byte(routingConfig.RoutingOption.Abenla.Password))
				params["sign"] = hex.EncodeToString(hasher.Sum(nil))
			}
			client := resty.New()
			client.SetTimeout(time.Second * 3)
			client.SetTLSClientConfig(&tls.Config{
				InsecureSkipVerify: true,
			})
			url := routingConfig.RoutingOption.Abenla.ApiAuthUrl
			res, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetQueryParams(params).
				Get(url)
			if err != nil {
				return err
			}
			result := model.AbenlaCheckConnectionResponse{}
			if err := json.Unmarshal(res.Body(), &result); err != nil {
				return err
			}
			if result.Code != 106 {
				return errors.New(result.Message)
			}
			cacheUtil.NewMemCache().Set(ABENLA_TOKEN, result, EXPIRE_EXTERNAL_ROUTING)

			return nil
		}
	} else if routingConfig.RoutingOption.Incom.Status {
		return nil
	} else if routingConfig.RoutingOption.Fpt.Status {
		return nil
	}

	return nil
}

func HandlePushRMQ(ctx context.Context, index, docId string, authUser *model.AuthUser, routingConfig model.RoutingConfig, tmpBytes []byte) error {
	esDoc := make(map[string]any)
	err := json.Unmarshal(tmpBytes, &esDoc)
	if err != nil {
		return err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX, authUser.TenantId); err != nil {
		return err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX, authUser.TenantId); err != nil {
			return err
		}
	}

	_, err = repository.ESRepo.CreateDocRabbitMQ(ctx, ES_INDEX, authUser.TenantId, authUser.TenantId, docId, esDoc)
	if err != nil {
		return err
	}

	return nil
}

func GetAccessTokenFpt(ctx context.Context, routingConfig model.RoutingConfig) (token string, err error) {
	tokenCache := cacheUtil.NewMemCache().Get(FPT_TOKEN + "_" + routingConfig.Id)
	if tokenCache != nil {
		return tokenCache.(string), nil
	}
	hasher := md5.New()
	hasher.Write([]byte(uuid.NewString()))

	body := map[string]any{
		"client_id":     routingConfig.RoutingOption.Fpt.ClientId,
		"client_secret": routingConfig.RoutingOption.Fpt.ClientSecret,
		"scope":         routingConfig.RoutingOption.Fpt.Scope,
		"session_id":    hex.EncodeToString(hasher.Sum(nil)),
		"grant_type":    routingConfig.RoutingOption.Fpt.GrantType,
	}
	url := fmt.Sprintf(routingConfig.RoutingOption.Fpt.ApiAuthUrl)
	res, err := httpUtil.Post(url, body)
	if err != nil {
		log.Error(err)
		return "", err
	} else if res.StatusCode() != http.StatusOK {
		loginResponse := model.FptResponseError{}
		err = json.Unmarshal([]byte(res.Body()), &loginResponse)
		if err != nil {
			log.Error(err)
			return "", err
		}
		return loginResponse.ErrorDescription, errors.New(loginResponse.ErrorDescription)
	}
	loginResponse := model.FptGetTokenResponseSuccess{}
	err = json.Unmarshal([]byte(res.Body()), &loginResponse)
	if err != nil {
		log.Error(err)
		return "", err
	}
	cacheUtil.NewMemCache().Set(FPT_TOKEN+"_"+routingConfig.Id, loginResponse.AccessToken, EXPIRE_EXTERNAL_ROUTING)
	return loginResponse.AccessToken, nil
}
