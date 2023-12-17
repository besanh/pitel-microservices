package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
)

var (
	MethodActionHookAccept = []string{"post", "get"}
)

func HandleWebhook(ctx context.Context, routingConfig model.RoutingConfig, data model.WebhookSendData) []error {
	errArr := []error{}
	// dataHook := model.WebhookSendData{
	// 	Id:           data.Id,
	// 	Status:       data.Status,
	// 	Channel:      data.Channel,
	// 	ErrorCode:    data.ErrorCode,
	// 	Quantity:     data.Quantity,
	// 	TelcoId:      data.TelcoId,
	// 	IsChargedZns: data.IsChargedZns,
	// }

	// for _, val := range routingConfig.WebhookUrl {
	// 	if len(val.MethodAction) > 0 {
	// 		if slices.Contains[[]string](MethodActionHookAccept, strings.ToLower(val.MethodAction)) {
	// 			if strings.ToLower(val.MethodAction) == "get" {
	// 				_, err := HandleWebHookMethodGet(ctx, val, pluginInfo, dataHook)
	// 				if err != nil {
	// 					errArr = append(errArr, err)
	// 					continue
	// 				}
	// 			} else if strings.ToLower(val.MethodAction) == "post" {
	// 				_, err := HandleWebhookMethodPost(ctx, val, pluginInfo, dataHook)
	// 				if err != nil {
	// 					errArr = append(errArr, err)
	// 					continue
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return errArr
}

func HandleWebHookMethodGet(ctx context.Context, webhookInfo model.WebhookPlugin, routingConfig model.RoutingConfig, dataHook model.WebhookSendData) (int, error) {
	queryParams := map[string]string{}
	queryParams["id"] = dataHook.Id
	queryParams["quantity"] = strconv.Itoa(dataHook.Quantity)
	queryParams["status"] = dataHook.Status
	queryParams["telco_id"] = strconv.Itoa(dataHook.TelcoId)
	queryParams["is_charged_zns"] = strconv.FormatBool(dataHook.IsChargedZns)

	if len(dataHook.Channel) > 0 {
		queryParams["channel"] = dataHook.Channel
	}
	if len(dataHook.ErrorCode) > 0 {
		queryParams["error_code"] = dataHook.ErrorCode
	}
	if len(webhookInfo.Signature) > 0 {
		queryParams["signature"] = webhookInfo.Signature
	}

	url := webhookInfo.Url

	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	exe := client.R().
		SetHeader("Content-Type", "application/json")
	if len(webhookInfo.Token) > 0 {
		exe.SetHeader("Authorization", fmt.Sprintf("Bearer %s", webhookInfo.Token))
	} else if len(webhookInfo.Username) > 0 && len(webhookInfo.Password) > 0 {
		exe.SetBasicAuth(webhookInfo.Username, webhookInfo.Password)
	}
	exe.SetQueryParams(queryParams)

	res, err := exe.Get(url)
	if err != nil {
		log.Error(err)
		return res.StatusCode(), err
	}
	result := map[string]any{}
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		return http.StatusBadRequest, err
	}
	return res.StatusCode(), nil
}

func HandleWebhookMethodPost(ctx context.Context, webhookInfo model.WebhookPlugin, pluginInfo model.RoutingConfig, dataHook model.WebhookSendData) (int, error) {
	client := resty.New()
	client.SetTimeout(time.Second * 20)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	body := map[string]any{
		"id":             dataHook.Id,
		"quantity":       dataHook.Quantity,
		"status":         dataHook.Status,
		"telco_id":       dataHook.TelcoId,
		"is_charged_zns": dataHook.IsChargedZns,
	}
	if len(dataHook.Channel) > 0 {
		body["channel"] = dataHook.Channel
	}
	if len(dataHook.ErrorCode) > 0 {
		body["error_code"] = dataHook.ErrorCode
	}

	exe := client.R().
		SetHeader("Content-Type", "application/json")
	if len(webhookInfo.Token) > 0 {
		exe.SetHeader("Authorization", fmt.Sprintf("Bearer %s", webhookInfo.Token))
	}
	if len(webhookInfo.Username) > 0 && len(webhookInfo.Password) > 0 {
		exe.SetBasicAuth(webhookInfo.Username, webhookInfo.Password)
	}
	exe.SetBody(body)
	res, err := exe.Post(webhookInfo.Url)
	if err != nil {
		log.Error(err)
		return res.StatusCode(), err
	}
	result := map[string]any{}
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		return http.StatusBadRequest, err
	}
	return res.StatusCode(), nil
}
