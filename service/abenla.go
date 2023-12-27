package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IAbenla interface {
		AbenlaWebhook(ctx context.Context, routingConfigUuid string, data model.WebhookReceiveSmsStatus) (int, any)
	}
	Abenla struct {
		Index string
	}
)

const (
	HOOK_ABENLA = "hook_abenal"

	TTL_HOOK_ABENLA = 1 * time.Minute
)

var AbenlaService IAbenla

func (s *Webhook) AbenlaWebhook(ctx context.Context, routingConfigUuid string, data model.WebhookReceiveSmsStatus) (int, any) {
	statusCode := "1"
	routingConfig := model.RoutingConfig{}
	// Caching
	routingConfigCache := cache.NewMemCache().Get(HOOK_ABENLA + "_" + data.SercretSign)
	if routingConfigCache != nil {
		routing, ok := routingConfigCache.(*model.RoutingConfig)
		if !ok {
			return response.ResponseXml("Status", statusCode)
		}
		routingConfig = *routing
	} else {
		routing, err := repository.RoutingConfigRepo.GetRoutingConfigById(ctx, data.SercretSign)
		if err != nil {
			log.Error(err)
			return response.ResponseXml("Status", statusCode)
		} else if routing == nil {
			return response.ResponseXml("Status", statusCode)
		}
		cache.NewMemCache().Set(INFO_ROUTING+"_"+data.SercretSign, routingConfig, EXPIRE_EXTERNAL_ROUTING)
		routingConfig = *routing
	}

	logExist, err := repository.InboxMarketingESRepo.GetDocByRoutingExternalMessageId(ctx, s.Index, data.SmsGuid)
	if err != nil {
		log.Error(err)
		return response.ResponseXml("Status", statusCode)
	} else if len(logExist.Id) < 1 {
		return response.ResponseXml("Status", statusCode)
	}
	auditLogModel := model.LogInboxMarketing{
		Id:                logExist.Id,
		RoutingConfigUuid: logExist.RoutingConfigUuid,
		ExternalMessageId: data.SmsGuid,
		Plugin:            logExist.Plugin,
		Status:            strconv.Itoa(data.Status),
		Quantity:          1,
		IsCheck:           logExist.IsCheck,
		Code:              0,
		CountAction:       logExist.CountAction + 1,
	}
	logAction, err := common.HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		log.Error(err)
		return response.ResponseXml("Status", statusCode)
	}
	logExist.Log = append(logExist.Log, logAction)

	status := strconv.Itoa(data.Status)

	statusStr := ""
	if status == "3" {
		statusStr = "success"
	} else if status == "4" {
		statusStr = "sent_fail"
	} else if status == "5" {
		statusStr = "wrong_phone_number"
	} else if status == "6" {
		statusStr = "account_expired"
	} else if status == "7" {
		statusStr = "amount_zero"
	} else if status == "8" {
		statusStr = "not_price"
	} else if status == "9" {
		statusStr = "can_not_sent"
	} else if status == "10" {
		statusStr = "deny_phone_number"
	} else if status == "13" {
		statusStr = "wrong_sender_name"
	}
	logExist.Status = statusStr
	esDoc := map[string]any{}
	tmpBytes, err := json.Marshal(logExist)
	if err != nil {
		log.Error(err)
		return response.ResponseXml("Status", statusCode)
	}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return response.ResponseXml("Status", statusCode)
	}
	if err := repository.ESRepo.UpdateDocById(ctx, "", logExist.Id, esDoc); err != nil {
		log.Error(err)
		return response.ResponseXml("Status", statusCode)
	}

	// Set cache
	cache.NewMemCache().Set(HOOK_ABENLA+"_"+data.SmsGuid, routingConfig, TTL_HOOK_ABENLA)

	statusCode = "2"
	return response.ResponseXml("Status", statusCode)
}
