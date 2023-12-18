package common

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
	INFO_ROUTING = "info_routing"

	// token
	ABENLA_TOKEN            = "abenla_token"
	INCOM_TOKEN             = "incom_token"
	FPT_TOKEN               = "fpt_token"
	EXPIRE_EXTERNAL_ROUTING = 1 * time.Minute
)

func GetInfoPlugin(ctx context.Context, db sqlclient.ISqlClientConn, authUser *model.AuthUser, routingConfigUuid string) (*model.RoutingConfig, error) {
	routingConfigCache, err := cacheUtil.MCache.Get(INFO_ROUTING + "_" + routingConfigUuid)
	if err != nil {
		return nil, err
	} else if routingConfigCache != nil {
		routing := routingConfigCache.(*model.RoutingConfig)
		return routing, nil
	} else {
		routingConfig, err := repository.RoutingConfigRepo.GetById(ctx, db, routingConfigUuid)
		if err != nil {
			return nil, err
		} else if routingConfig == nil {
			return nil, errors.New("routing not found in system")
		}
		if err := cacheUtil.MCache.SetTTL(INFO_ROUTING+"_"+routingConfigUuid, routingConfig, EXPIRE_EXTERNAL_ROUTING); err != nil {
			return nil, err
		}
		return routingConfig, nil
	}
}

func CheckConnectionWithExternalPlugin(ctx context.Context, routingConfig model.RoutingConfig) error {
	params := map[string]string{}
	if routingConfig.RoutingOption.Abenla.Status {
		connectionAlive, err := cacheUtil.MCache.Get(ABENLA_TOKEN)
		if err != nil {
			return err
		} else if connectionAlive != nil {
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
			url := routingConfig.RoutingOption.Abenla.ApiUrl + "/api/CheckConnection"
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
			if err := cacheUtil.MCache.SetTTL(ABENLA_TOKEN, result, EXPIRE_EXTERNAL_ROUTING); err != nil {
				return err
			}

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
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, authUser.DatabaseEsIndex, authUser.TenantId); err != nil {
		return err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, authUser.DatabaseEsIndex, authUser.TenantId); err != nil {
			return err
		}
	}

	_, err = repository.ESRepo.CreateDocRabbitMQ(ctx, authUser.DatabaseEsIndex, authUser.TenantId, authUser.TenantId, docId, esDoc)
	if err != nil {
		return err
	}

	return nil
}

func GetAccessTokenFpt(ctx context.Context, dbCon *sqlclient.SqlClientConn) (token string, err error) {
	plugin, err := GetExternalPluginConnectFromCache(ctx, dbCon, "fpt")
	if err != nil {
		return "", err
	}
	hasher := md5.New()
	hasher.Write([]byte(uuid.NewString()))

	body := map[string]any{
		"client_id":     plugin.Config.FptConfig.ClientId,
		"client_sercet": plugin.Config.FptConfig.ClientSercet,
		"scope":         plugin.Config.FptConfig.Scope,
		"session_id":    hex.EncodeToString(hasher.Sum(nil)),
		"grant_type":    "client_credentials",
	}
	url := fmt.Sprintf(plugin.Config.FptConfig.Api)
	res, err := httpUtil.Post(url, body)
	if err != nil {
		log.Error(err)
		return "", err
	} else if res.StatusCode() != http.StatusOK {
		loginResponse := model.FptGetTokenResponseError{}
		err = json.Unmarshal([]byte(res.Body()), &loginResponse)
		if err != nil {
			log.Error(err)
			return "", err
		}
	}
	loginResponse := model.FptGetTokenResponseSuccess{}
	err = json.Unmarshal([]byte(res.Body()), &loginResponse)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if err := cacheUtil.MCache.SetTTL(FPT_TOKEN, loginResponse.AccessToken, 5*time.Minute); err != nil {
		log.Error(err)
	}
	return loginResponse.AccessToken, nil
}
