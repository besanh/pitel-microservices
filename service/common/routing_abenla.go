package common

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/model"
)

func HandleDeliveryMessageAbenla(ctx context.Context, id string, routingConfig model.RoutingConfig, inboxMarketingRequest model.InboxMarketingRequest) (int, model.ResponseInboxMarketing, error) {
	resultStandard := model.ResponseInboxMarketing{}
	result := model.AbenlaSendMessageResponse{}
	url := routingConfig.RoutingOption.Abenla.ApiSendMessageUrl
	hasher := md5.New()
	hasher.Write([]byte(routingConfig.RoutingOption.Abenla.Password))
	var serviceTypeId int
	serviceTypeId, _ = strconv.Atoi(routingConfig.RoutingOption.Abenla.ServiceTypeId)
	params := map[string]string{
		"loginName":     routingConfig.RoutingOption.Abenla.Username,
		"sign":          hex.EncodeToString(hasher.Sum(nil)),
		"serviceTypeId": strconv.Itoa(serviceTypeId),
		"phoneNumber":   inboxMarketingRequest.PhoneNumber,
		"message":       inboxMarketingRequest.Content,
		"brandName":     routingConfig.RoutingOption.Abenla.Brandname,
		"smsGuid":       id,
	}
	// if len(routingConfig.RoutingOption.Abenla.WebhookUrl) > 0 {
	// temporary get first hook
	// isCallbackTmp := false
	// if len(config.WebhookUrl) > 0 {
	// 	isCallbackTmp = true
	// }
	// isCallback := strconv.FormatBool(isCallbackTmp)
	// params["callback"] = isCallback
	// }
	params["callback"] = strconv.FormatBool(true)

	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		Get(url)
	if err != nil {
		resultStandard := HandleMapResponsePlugin("abenla", id, 0, result)
		return res.StatusCode(), resultStandard, err
	}

	if err := json.Unmarshal(res.Body(), &result); err != nil {
		return res.StatusCode(), resultStandard, err
	}
	resultStandard = HandleMapResponsePlugin("abenla", id, 0, result)
	return res.StatusCode(), resultStandard, nil
}
