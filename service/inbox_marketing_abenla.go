package service

import (
	"context"
	"encoding/json"
	"errors"
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
	IInboxMarketingLogAbenla interface {
		WebhookReceiveSmsStatus(ctx context.Context, authUser *model.AuthUser, routingConfig model.RoutingConfig, data model.WebhookReceiveSmsStatus) (int, any)
	}
	InboxMarketingAbenla struct{}
)

const (
	SIGNATURE = "signature"
	TTL       = 1 * time.Minute
)

func NewInboxMarketingAbenla() IInboxMarketingLogAbenla {
	return &InboxMarketingAbenla{}
}

func (s *InboxMarketingAbenla) WebhookReceiveSmsStatus(ctx context.Context, authUser *model.AuthUser, routingConfig model.RoutingConfig, data model.WebhookReceiveSmsStatus) (int, any) {
	// Caching
	_, err := cache.MCache.Get(SIGNATURE + "_" + routingConfig.Id)
	if err != nil {
		log.Error(err)
		return response.ResponseXml("Status", "1")
	}

	// Check signature
	// pluginInfo, err := repository.PluginConfigRepo.GetExternalPluginConfigByUsernameOrSignature(ctx, dbCon, "", data.SercretSign)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ResponseXml("Status", "1")
	// }

	// Update ES
	logExist, err := repository.InboxMarketingESRepo.GetDocById(ctx, authUser.TenantId, authUser.DatabaseEsIndex, data.SmsGuid)
	if err != nil {
		log.Error(err)
		return response.ResponseXml("Status", "1")
	} else if len(logExist.Id) < 1 {
		return response.ResponseXml("Status", "1")
	}
	auditLogModel := model.LogInboxMarketing{
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		Services:          authUser.Services,
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
		return response.ResponseXml("Status", "1")
	}
	logExist.Log = append(logExist.Log, logAction)

	status := strconv.Itoa(data.Status)
	// Send to hook
	if len(routingConfig.RoutingOption.Abenla.WebhookUrl) > 0 {
		dataHook := model.WebhookSendData{
			Id:     data.SmsGuid,
			Status: status,
		}
		errArr := common.HandleWebhook(ctx, routingConfig, dataHook)
		if len(errArr) > 0 {
			log.Error(err)
			return response.ResponseXml("Status", "1")
		}
	}
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
		return response.ResponseXml("Status", "1")
	}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return response.ResponseXml("Status", "1")
	}
	if err := repository.ESRepo.UpdateDocById(ctx, authUser.DatabaseEsIndex, logExist.Id, esDoc); err != nil {
		log.Error(err)
		return response.ResponseXml("Status", "1")
	}

	// Set cache
	if err := cache.MCache.SetTTL(SIGNATURE+"_"+routingConfig.Id, routingConfig.RoutingOption.Abenla.Signature, TTL); err != nil {
		log.Error(err)
		return response.ResponseXml("Status", "1")
	}

	return response.ResponseXml("Status", "2")
}

func HandleMainInboxMarketingAbenla(ctx context.Context, authUser *model.AuthUser, inboxMarketingBasic model.InboxMarketingBasic, routingConfig model.RoutingConfig, inboxMarketing model.InboxMarketingLogInfo, inboxMarketingRequest model.InboxMarketingRequest) (model.ResponseInboxMarketing, error) {
	res := model.ResponseInboxMarketing{
		Id: inboxMarketingBasic.DocId,
	}
	dataUpdate := map[string]any{}

	_, result, err := common.HandleDeliveryMessageAbenla(ctx, inboxMarketingBasic.DocId, routingConfig, inboxMarketingRequest)
	if err != nil {
		return res, err
	}

	// Send to hook
	// if len(pluginInfo.WebhookUrl) > 0 {
	// 	dataHook := model.WebhookSendData{
	// 		Id:       docId,
	// 		Quantity: result.SmsPerMessage,
	// 		Status:   result.Message,
	// 	}
	// 	errArr := common.HandleWebhook(ctx, pluginInfo, dataHook)
	// 	if len(errArr) > 0 {
	// 		return response.ServiceUnavailableMsg(errArr)
	// 	}
	// }

	// Find in ES to avoid 404 not found
	dataExist, err := repository.InboxMarketingESRepo.GetDocById(ctx, inboxMarketingBasic.TenantId, authUser.DatabaseEsIndex, inboxMarketingBasic.Id)
	if err != nil {
		return res, err
	} else if len(dataExist.Id) < 1 {
		return res, errors.New("log is not exist")
	}

	inboxMarketing.Quantity = result.Quantity
	code, _ := strconv.Atoi(result.Code)
	// log
	auditLogModel := model.LogInboxMarketing{
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		Services:          authUser.Services,
		Id:                inboxMarketing.Id,
		RoutingConfigUuid: routingConfig.Id,
		ExternalMessageId: inboxMarketingBasic.ExternalMessageId,
		Status:            "",
		Quantity:          0,
		TelcoId:           0,
		IsChargedZns:      false,
		IsCheck:           false,
		Code:              code,
		CountAction:       2,
		UpdatedBy:         inboxMarketingBasic.UpdatedBy,
	}
	auditLog, err := common.HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		return res, err
	}
	inboxMarketing.Log = append(inboxMarketing.Log, auditLog)
	inboxMarketing.UpdatedAt = time.Now()

	tmpBytesUpdate, err := json.Marshal(inboxMarketing)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(tmpBytesUpdate, &dataUpdate); err != nil {
		return res, err
	}

	if err := repository.ESRepo.UpdateDocById(ctx, inboxMarketingBasic.Index, inboxMarketingBasic.DocId, dataUpdate); err != nil {
		return res, err
	}

	res.Status = result.Message

	return res, nil
}
