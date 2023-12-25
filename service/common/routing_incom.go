package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func HandleDeliveryMessageIncom(ctx context.Context, id string, routingConfig model.RoutingConfig, templateCode string, inboxMarketing model.InboxMarketingLogInfo, inboxMarketingRequest model.InboxMarketingRequest) (int, model.ResponseInboxMarketing, error) {
	resultStandard := model.ResponseInboxMarketing{}
	result := model.IncomSendMessageResponse{}

	url := routingConfig.RoutingOption.Incom.ApiSendMessageUrl
	if len(url) < 1 {
		resultStandard.Code = "3" // fail
		resultStandard.Message = constants.MESSAGE_TEL4VN["fail"]
		return 0, resultStandard, errors.New("api url is empty")
	}

	inboxMarketing.RouteRule = removeDuplicateValues(inboxMarketing.RouteRule)
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
		resultStandard := HandleMapResponsePlugin("incom", id, 0, result)
		return res.StatusCode(), resultStandard, err
	}
	var r any
	if err := json.Unmarshal(res.Body(), &r); err != nil {
		resultStandard := HandleMapResponsePlugin("incom", id, 0, result)
		return res.StatusCode(), resultStandard, err
	}
	if err := util.ParseAnyToAny(r, &result); err != nil {
		return res.StatusCode(), resultStandard, err
	}
	resultStandard = HandleMapResponsePlugin("incom", id, 0, result)
	statusCode := 0
	if resultStandard.Status == "success" {
		statusCode = 200
	} else if resultStandard.Status == "fail" {
		statusCode = 400
	}
	return statusCode, resultStandard, err
}

func HandleGetStatusMessage(ctx context.Context, dbCon sqlclient.ISqlClientConn, authUser *model.AuthUser, routingConfig model.RoutingConfig, docId string, body model.IncomBodyStatus) (model.InboxMarketingLogInfo, error) {
	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	url := routingConfig.RoutingOption.Incom.WebhookUrl
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

	logExist.IsChargedZns = false
	if result.Ischarged == "True" || result.Ischarged == "true" {
		logExist.IsChargedZns = true
	}

	return *logExist, nil
}
