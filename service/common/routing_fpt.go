package common

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	httpUtil "github.com/tel4vn/fins-microservices/common/http"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
)

// const (
// 	INFO_EXTERNAL_PLUGIN_CONNECT   = "INFO_EXTERNAL_PLUGIN_CONNECT"
// 	EXPIRE_EXTERNAL_PLUGIN_CONNECT = 1 * time.Minute
// )

// func GetExternalPluginConnectFromCache(ctx context.Context, dbCon sqlclient.ISqlClientConn, externalPluginConnectType string) (*model.ExternalPluginConnect, error) {
// 	externalPluginConnectCache, err := cacheUtil.MCache.Get(INFO_EXTERNAL_PLUGIN_CONNECT + "_" + externalPluginConnectType)
// 	if err != nil {
// 		return nil, err
// 	} else if externalPluginConnectCache != nil {
// 		externalPluginConnect := externalPluginConnectCache.(*model.ExternalPluginConnect)
// 		return externalPluginConnect, nil
// 	} else {
// 		externalPluginConnect, err := repository.ExternalPluginConnectRepo.GetExternalPluginByType(ctx, dbCon, externalPluginConnectType)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err := cacheUtil.MCache.SetTTL(INFO_EXTERNAL_PLUGIN_CONNECT+"_"+externalPluginConnectType, externalPluginConnect, EXPIRE_EXTERNAL_PLUGIN_CONNECT); err != nil {
// 			return nil, err
// 		}
// 		return externalPluginConnect, nil
// 	}
// }

func HandleDeliveryMessageFpt(ctx context.Context, id string, routingConfig model.RoutingConfig, inboxMarketingRequest model.InboxMarketingRequest, fpt model.FptRequireRequest) (int, *model.ResponseInboxMarketing, *model.FptSendMessageResponse, error) {
	resultStandard := model.ResponseInboxMarketing{}
	body := map[string]any{
		"access_token": fpt.AccessToken,
		"session_id":   fpt.SessionId,
		"BrandName":    routingConfig.RoutingOption.Fpt.BrandName,
		"Phone":        inboxMarketingRequest.PhoneNumber,
		"Message":      base64.StdEncoding.EncodeToString([]byte(inboxMarketingRequest.Content)),
		"RequestId":    id,
	}
	url := fmt.Sprintf(routingConfig.RoutingOption.Fpt.ApiSendMessageUrl)
	res, err := httpUtil.Post(url, body)
	if err != nil {
		log.Error(err)
		resultStandard = HandleMapResponsePlugin("fpt", id, 0, res)
		return 0, &resultStandard, nil, err
	} else if res.StatusCode() != http.StatusOK {
		resErr := model.FptResponseError{}
		err = json.Unmarshal([]byte(res.Body()), &resErr)
		if err != nil {
			log.Error(err)
			resultStandard = HandleMapResponsePlugin("fpt", id, 0, resErr)
			return 0, &resultStandard, nil, err
		}
		resultStandard = HandleMapResponsePlugin("fpt", id, 0, resErr)
		return 0, &resultStandard, nil, errors.New(resErr.ErrorDescription)
	}
	resSuccess := model.FptSendMessageResponse{}
	err = json.Unmarshal([]byte(res.Body()), &resSuccess)
	if err != nil {
		log.Error(err)
		resultStandard = HandleMapResponsePlugin("fpt", id, 0, resSuccess)
		return 0, &resultStandard, nil, err
	}
	var tmp any
	if err := json.Unmarshal([]byte(res.Body()), &tmp); err != nil {
		log.Error(err)
		resultStandard = HandleMapResponsePlugin("fpt", id, 0, resSuccess)
		return 0, &resultStandard, nil, err
	}
	if len(resSuccess.MessageId) < 1 {
		resultStandard = HandleMapResponsePlugin("fpt", id, 0, resSuccess)
		return 0, &resultStandard, nil, err
	}
	resultStandard = HandleMapResponsePlugin("fpt", id, 200, resSuccess)
	statusCode := 0
	if resultStandard.Status == "Success" {
		statusCode = 200
	} else if resultStandard.Status == "Fail" {
		statusCode = 400
	}
	return statusCode, &resultStandard, &resSuccess, nil
}
