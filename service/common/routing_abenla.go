package common

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/model"
)

func HandleDeliveryMessageAbenla(ctx context.Context, id string, routingConfig model.RoutingConfig, inboxMarketingRequest model.InboxMarketingRequest) (int, model.AbenlaSendMessageResponse, error) {
	result := model.AbenlaSendMessageResponse{}
	url := routingConfig.RoutingOption.Abenla.ApiUrl
	if len(url) < 1 {
		return 0, result, errors.New("api url is empty")
	}
	url += "/api/SendSms"
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
		return res.StatusCode(), result, err
	}

	if err := json.Unmarshal(res.Body(), &result); err != nil {
		return res.StatusCode(), result, err
	}

	return res.StatusCode(), result, nil
}
