package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IFpt interface {
		FptWebhook(ctx context.Context, routingConfigUuid string, data model.FptWebhook) (int, any)
	}
	FptWebhook struct {
		Index string
	}
)

var (
	HOOK_FPT     = "fpt"
	TTL_HOOK_FPT = 1 * time.Minute
)

func (s *Webhook) FptWebhook(ctx context.Context, routingConfigUuid string, data model.FptWebhook) (int, any) {
	routingConfig := model.RoutingConfig{}
	// Caching
	routingConfigCache := cache.NewMemCache().Get(INFO_ROUTING + "_" + routingConfigUuid)
	if routingConfigCache != nil {
		routing, ok := routingConfigCache.(*model.RoutingConfig)
		if !ok {
			log.Info("routing not found in system")
			return response.OK(map[string]any{
				"status": 0,
			})
		}
		routingConfig = *routing
	} else {
		routing, err := repository.RoutingConfigRepo.GetRoutingConfigById(ctx, routingConfigUuid)
		if err != nil {
			log.Error(err)
			return response.OK(map[string]any{
				"status": 0,
			})
		} else if routing == nil {
			log.Info("routing not found in system")
			return response.OK(map[string]any{
				"status": 0,
			})
		}
		cache.NewMemCache().Set(INFO_ROUTING+"_"+routingConfigUuid, routingConfig, EXPIRE_EXTERNAL_ROUTING)
		routingConfig = *routing
	}

	externalMessageId := strconv.Itoa(data.SmsId)
	logWebhookExist, err := repository.InboxMarketingESRepo.GetDocByRoutingExternalMessageId(ctx, s.Index, externalMessageId)
	if err != nil {
		log.Error(err)
		return response.OK(map[string]any{
			"status": 0,
		})
	} else if len(logWebhookExist.Id) < 1 {
		log.Error("message not found in system")
		return response.OK(map[string]any{
			"status": 0,
		})
	}

	// Map network
	telcoId := util.MapNetworkPlugin(data.Telco)
	telcoStr, _ := strconv.Atoi(telcoId)

	logWebhookExist.StatusHook = common.MapStatusFpt(data.Status)
	logWebhookExist.ChannelHook = "sms"
	logWebhookExist.ErrorCode = data.Error
	logWebhookExist.Quantity = data.MtCount
	logWebhookExist.TelcoId = telcoStr

	auditLogModel := model.LogInboxMarketing{
		Id:                logWebhookExist.Id,
		RoutingConfigUuid: logWebhookExist.RoutingConfigUuid,
		ExternalMessageId: externalMessageId,
		Plugin:            logWebhookExist.Plugin,
		Status:            common.MapStatusFpt(data.Status),
		ChannelHook:       logWebhookExist.ChannelHook,
		ErrorCodeHook:     logWebhookExist.ErrorCodeHook,
		TelcoId:           telcoStr,
		Quantity:          data.MtCount,
		Channel:           "sms",
		IsCheck:           logWebhookExist.IsCheck,
		Code:              0,
		CountAction:       logWebhookExist.CountAction + 1,
	}
	logAction, err := common.HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		log.Error(err)
		return response.OK(map[string]any{
			"status": 0,
		})
	}
	logWebhookExist.Log = append(logWebhookExist.Log, logAction)

	esDoc := map[string]any{}
	tmpByte, err := json.Marshal(logWebhookExist)
	if err != nil {
		log.Error(err)
		return response.OK(map[string]any{
			"status": 0,
		})
	}
	if err := json.Unmarshal(tmpByte, &esDoc); err != nil {
		log.Error(err)
		return response.OK(map[string]any{
			"status": 0,
		})
	}

	if err := repository.ESRepo.UpdateDocById(ctx, s.Index, logWebhookExist.Id, esDoc); err != nil {
		log.Error(err)
		return response.OK(map[string]any{
			"status": 0,
		})
	}

	// Set cache
	cache.NewMemCache().Set(HOOK_FPT+"_"+routingConfigUuid, routingConfig, TTL_HOOK_FPT)

	return response.OK(map[string]any{
		"status": 1,
	})
}
