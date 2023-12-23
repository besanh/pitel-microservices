package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IIncom interface {
		WebhookIncom(ctx context.Context, routingConfigUuid string, authUser *model.AuthUser, data model.WebhookIncom) (err error)
	}
	IncomWebhook struct{}
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

func NewIncom() IIncom {
	return &IncomWebhook{}
}

func (s *IncomWebhook) WebhookIncom(ctx context.Context, routingConfigUuid string, authUser *model.AuthUser, data model.WebhookIncom) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return err
	}
	_ = dbCon
	routingConfig := model.RoutingConfig{}
	// Caching
	routingConfigCache := cache.NewMemCache().Get(INFO_ROUTING + "_" + routingConfigUuid)
	if routingConfigCache != nil {
		routing, ok := routingConfigCache.(*model.RoutingConfig)
		if !ok {
			log.Error(err)
			return err
		}
		routingConfig = *routing
	} else {
		// routing, err := repository.RoutingConfigRepo.GetById(ctx, dbCon, routingConfigUuid)
		// if err != nil {
		// 	log.Error(err)
		// 	return err
		// } else if routing == nil {
		// 	return errors.New("routing not found in system")
		// }
		// cache.NewMemCache().Set(INFO_ROUTING+"_"+routingConfigUuid, routingConfig, EXPIRE_EXTERNAL_ROUTING)
		// routingConfig = *routing
	}

	logWebhookExist, err := repository.InboxMarketingESRepo.GetDocByRoutingExternalMessageId(ctx, authUser.TenantId, authUser.DatabaseEsIndex, data.IdOmniMess)
	if err != nil {
		log.Error(err)
		return err
	} else if len(logWebhookExist.Id) < 1 {
		return fmt.Errorf("%s is not exist", data.IdOmniMess)
	}

	// Map network
	telcoTmp := strconv.Itoa(data.TelcoId)
	telcoId := util.MapNetworkPlugin(telcoTmp)
	telcoStr, _ := strconv.Atoi(telcoId)

	logWebhookExist.StatusHook = strings.ToLower(data.Status)
	channel := strings.ToLower(data.Channel)
	logWebhookExist.ChannelHook = strings.ReplaceAll(channel, "brandnamesms", "sms")
	logWebhookExist.ErrorCode = data.ErrorCode
	logWebhookExist.Quantity = data.Quantity
	logWebhookExist.TelcoId = telcoStr
	logWebhookExist.IsChargedZns = data.IsChargedZns

	auditLogModel := model.LogInboxMarketing{
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		Services:          authUser.Services,
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
		return err
	}
	logWebhookExist.Log = append(logWebhookExist.Log, logAction)

	esDoc := map[string]any{}
	tmpByte, err := json.Marshal(logWebhookExist)
	if err != nil {
		log.Error(err)
		return err
	}
	if err := json.Unmarshal(tmpByte, &esDoc); err != nil {
		log.Error(err)
		return err
	}

	if err := repository.ESRepo.UpdateDocById(ctx, authUser.DatabaseEsIndex, logWebhookExist.Id, esDoc); err != nil {
		log.Error(err)
		return err
	}

	// Set cache
	cache.NewMemCache().Set(HOOK_INCOM+"_"+routingConfigUuid, routingConfig, TTL_HOOK_INCOM)

	return err
}
