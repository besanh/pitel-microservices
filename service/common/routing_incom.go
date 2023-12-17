package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func HandleDeliveryMessageIncom(ctx context.Context, id string, routingConfig model.RoutingConfig, templateCode string, inboxMarketing model.InboxMarketingLogInfo, inboxMarketingRequest model.InboxMarketingRequest) (int, model.ResponseInboxMarketing, error) {
	inboxMarketing.RouteRule = removeDuplicateValues(inboxMarketing.RouteRule)
	resultInternal := model.ResponseInboxMarketing{}
	result := model.IncomSendMessageResponse{}

	url := routingConfig.RoutingOption.Incom.ApiUrl
	if len(url) < 1 {
		return 0, resultInternal, errors.New("api url is empty")
	}
	// url += "/api/OmniMessage/SendMessage"

	body := map[string]any{
		"username":     routingConfig.RoutingOption.Incom.Username,
		"password":     routingConfig.RoutingOption.Incom.Password,
		"phonenumber":  inboxMarketingRequest.PhoneNumber,
		"templatecode": templateCode,
		"routerule":    inboxMarketing.RouteRule,
		"list_param":   inboxMarketing.ListParam,
	}
	client := resty.New()
	client.SetTimeout(time.Second * 5)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		resultInternal.Status = strings.ToLower(result.Code)
		var status int
		var err error
		if len(result.Status) > 0 {
			status, err = strconv.Atoi(result.Status)
		}
		resultInternal.Code = status

		return res.StatusCode(), resultInternal, err
	}
	var r interface{}
	if err := json.Unmarshal(res.Body(), &r); err != nil {
		resultInternal.Status = strings.ToLower(result.Code)
		var status int
		var err error
		if len(result.Status) > 0 {
			status, err = strconv.Atoi(result.Status)
		}
		resultInternal.Code = status

		return res.StatusCode(), resultInternal, err
	}
	if err := util.ParseAnyToAny(r, &result); err != nil {
		return res.StatusCode(), resultInternal, err
	}

	resultInternal.Id = result.IdOmniMess
	resultInternal.Status = strings.ToLower(result.Code)
	status, _ := strconv.Atoi(result.Status)
	resultInternal.Code = status

	return res.StatusCode(), resultInternal, nil
}

func HandleGetStatusMessage(ctx context.Context, dbCon sqlclient.ISqlClientConn, authUser *model.AuthUser, routingConfig model.RoutingConfig, docId string, body model.IncomBodyStatus) (model.InboxMarketingLogInfo, error) {
	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	url := fmt.Sprintf("%s/api/OmniReport/GetStatusOmni", routingConfig.RoutingOption.Incom.ApiUrl)
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&model.IncomStatusMess{}).
		Post(url)
	if err != nil {
		return model.InboxMarketingLogInfo{}, err
	}
	result := &model.IncomStatusMess{}
	if err := json.Unmarshal(res.Body(), result); err != nil {
		return model.InboxMarketingLogInfo{}, err
	}
	log.Info("HandleGetStatusMessage", *result)
	logExist, err := repository.InboxMarketingESRepo.GetDocById(ctx, authUser.TenantId, authUser.DatabaseEsIndex, docId)
	if err != nil {
		return model.InboxMarketingLogInfo{}, err
	} else if len(logExist.Id) < 1 {
		return model.InboxMarketingLogInfo{}, errors.New("log not found")
	}
	logExist.PhoneNumber = result.PhoneNumber
	logExist.ListParam = result.ListParam
	logExist.SendTime = result.CreateDatetime
	logExist.TemplateCode = strings.ToLower(result.TemplateCode)
	logExist.Status = strings.ToLower(result.Status)
	channel := strings.ToLower(result.Channel)
	logExist.Channel = strings.ReplaceAll(channel, "brandnamesms", "sms")
	logExist.ErrorCode = strings.ToLower(result.ErrorCode)
	logExist.IsCheck = true
	logExist.UpdatedAt = time.Now()
	logExist.Quantity, _ = strconv.Atoi(result.MtCount)

	// Map network
	// telcoTmp := util.MapNetworkPlugin(result.TelcoId)
	// telco, _ := strconv.Atoi(telcoTmp)
	// logExist.TelcoId = telco

	logExist.IsChargedZns = false
	if result.Ischarged == "True" || result.Ischarged == "true" {
		logExist.IsChargedZns = true
	}

	return *logExist, nil
}
