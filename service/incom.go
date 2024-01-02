package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IIncom interface {
		IncomWebhook(ctx context.Context, data model.WebhookIncom) (int, any)
	}
	IncomWebhook struct {
		Index string
	}
)

const (
	HOOK_INCOM                     = "hook_incom"
	INFO_ROUTING                   = "info_routing"
	INFO_EXTERNAL_PLUGIN_CONNECT   = "INFO_EXTERNAL_PLUGIN_CONNECT"
	EXPIRE_EXTERNAL_PLUGIN_CONNECT = 1 * time.Minute

	// token
	ABENLA_TOKEN            = "abenla_token"
	INCOM_TOKEN             = "incom_token"
	FPT_TOKEN               = "fpt_token"
	EXPIRE_EXTERNAL_ROUTING = 1 * time.Minute
	TTL_HOOK_INCOM          = 1 * time.Minute
)

func (s *Webhook) IncomWebhook(ctx context.Context, data model.WebhookIncom) (int, any) {
	log.Info(ES_INDEX)
	logWebhookExist, err := repository.InboxMarketingESRepo.GetDocByRoutingExternalMessageId(ctx, ES_INDEX, data.IdOmniMess)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(logWebhookExist.Id) < 1 {
		log.Info(data.IdOmniMess + " is not exist")
		return response.ServiceUnavailableMsg(data.IdOmniMess + " is not exist")
	}

	// Map network
	telcoTmp := strconv.Itoa(data.TelcoId)
	telcoId := util.MapNetworkPlugin(telcoTmp)
	telcoStr, _ := strconv.Atoi(telcoId)

	logWebhookExist.StatusHook = strings.ToLower(data.Status)
	channel := strings.ToLower(data.Channel)
	logWebhookExist.ChannelHook = strings.ReplaceAll(channel, "brandnamesms", "sms")
	logWebhookExist.StatusHook = strings.ToLower(data.Status)
	logWebhookExist.ErrorCode = data.ErrorCode
	logWebhookExist.Quantity = data.Quantity
	logWebhookExist.TelcoId = telcoStr
	logWebhookExist.IsChargedZns = data.IsChargedZns

	auditLogModel := model.LogInboxMarketing{
		Id:                logWebhookExist.Id,
		RoutingConfigUuid: logWebhookExist.RoutingConfigUuid,
		ExternalMessageId: data.IdOmniMess,
		Plugin:            logWebhookExist.Plugin,
		Status:            data.Status,
		ChannelHook:       logWebhookExist.ChannelHook,
		ErrorCodeHook:     logWebhookExist.ErrorCodeHook,
		TelcoId:           telcoStr,
		IsChargedZns:      data.IsChargedZns,
		Quantity:          data.Quantity,
		Channel:           data.Channel,
		IsCheck:           logWebhookExist.IsCheck,
		Code:              0,
		CountAction:       logWebhookExist.CountAction + 1,
	}
	logAction, err := common.HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	logWebhookExist.Log = append(logWebhookExist.Log, logAction)

	esDoc := map[string]any{}
	tmpByte, err := json.Marshal(logWebhookExist)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if err := json.Unmarshal(tmpByte, &esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX, logWebhookExist.Id, esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]any{
		"message": "success",
	})
}
