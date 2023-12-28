package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IInboxMarketingIncom interface {
		WebhookReceiveStatus(ctx context.Context, routingConfig *model.RoutingConfig, authUser *model.AuthUser, pluginId string, data model.WebhookIncom) error
	}
	InboxMarketingIncom struct{}
)

const (
	HOOK_INCOM     = "hook_incom"
	TTL_HOOK_INCOM = 1 * time.Minute
)

func NewInboxMarketingIncom() IInboxMarketingIncom {
	return &InboxMarketingIncom{}
}

func (s *InboxMarketingIncom) WebhookReceiveStatus(ctx context.Context, routingConfig *model.RoutingConfig, authUser *model.AuthUser, pluginId string, data model.WebhookIncom) error {
	// Check signature
	// pluginInfo, err := repository.ExternalPluginRepo.GetExternalPluginByIdOrPlugin(ctx, "", pluginId, "")
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// }
	cache.NewMemCache().Set(HOOK_INCOM+"_"+pluginId, pluginId, TTL_HOOK_INCOM)
	logWebhookExist, err := repository.InboxMarketingESRepo.GetDocByRoutingExternalMessageId(ctx, authUser.TenantId, ES_INDEX, data.IdOmniMess)
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
	logAction, err := HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		log.Error(err)
		return err
	}
	logWebhookExist.Log = append(logWebhookExist.Log, logAction)
	// Send to hook
	if len(routingConfig.RoutingOption.Incom.WebhookUrl) > 0 {
		dataHook := model.WebhookSendData{
			Id:           logWebhookExist.Id,
			Status:       data.Status,
			Channel:      data.Channel,
			ErrorCode:    data.ErrorCode,
			Quantity:     data.Quantity,
			TelcoId:      telcoStr,
			IsChargedZns: data.IsChargedZns,
		}
		errArr := HandleWebhook(ctx, *routingConfig, dataHook)
		if len(errArr) > 0 {
			return errors.New("send hook error")
		}
	}

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

	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX, logWebhookExist.Id, esDoc); err != nil {
		log.Error(err)
		return err
	}

	// Set cache
	cache.NewMemCache().Set(HOOK_INCOM+"_"+pluginId, pluginId, TTL_HOOK_INCOM)

	return nil
}

func HandleMainInboxMarketingIncom(ctx context.Context, dbCon sqlclient.ISqlClientConn, inboxMarketingBasic model.InboxMarketingBasic, routingConfig model.RoutingConfig, authUser *model.AuthUser, inboxMarketing model.InboxMarketingLogInfo, inboxMarketingRequest model.InboxMarketingRequest) (model.ResponseInboxMarketing, error) {
	res := model.ResponseInboxMarketing{
		Id: inboxMarketingBasic.DocId,
	}
	dataUpdate := map[string]any{}
	time.Sleep(3 * time.Second)
	// Find in ES to avoid 404 not found
	dataExist, err := repository.InboxMarketingESRepo.GetDocById(ctx, authUser.TenantId, ES_INDEX, inboxMarketingBasic.DocId)
	if err != nil {
		return res, err
	} else if len(dataExist.Id) < 1 {
		return res, errors.New("log is not exist")
	}
	template, err := GetTemplate(ctx, dbCon, inboxMarketingRequest.Template)
	if err != nil {
		return res, err
	}
	statusCode, result, err := HandleDeliveryMessageIncom(ctx, inboxMarketingBasic.DocId, routingConfig, template.TemplateCode, inboxMarketing, inboxMarketingRequest)
	if err != nil {
		return res, err
	} else if statusCode != 200 {
		return result, errors.New(result.Message)
	}

	// Update id to ES
	inboxMarketing.ExternalMessageId = result.Id
	// inboxMarketing.Status = result.Status
	if len(constants.ROUTERULE[inboxMarketing.RouteRule[0]]) > 0 {
		inboxMarketing.ChannelHook = constants.ROUTERULE[inboxMarketing.RouteRule[0]]
	}
	countAction := inboxMarketing.CountAction + 1
	// log
	code, _ := strconv.Atoi(result.Code)
	auditLogModel := model.LogInboxMarketing{
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		Services:          authUser.Services,
		Id:                inboxMarketingBasic.DocId,
		RoutingConfigUuid: routingConfig.Id,
		ExternalMessageId: result.Id,
		Plugin:            routingConfig.RoutingType,
		Status:            result.Status,
		Quantity:          0,
		TelcoId:           0,
		IsChargedZns:      false,
		IsCheck:           false,
		Code:              code,
		CountAction:       countAction,
		UpdatedBy:         inboxMarketingBasic.UpdatedBy,
	}
	auditLog, err := HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		return res, err
	}
	inboxMarketing.Log = append(inboxMarketing.Log, auditLog)
	inboxMarketing.CountAction = countAction
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

	return result, err
}

func GetTemplate(ctx context.Context, dbCon sqlclient.ISqlClientConn, id string) (*model.TemplateBss, error) {
	templateCache := cache.NewMemCache().Get(INFO_TEMPLATE + "_" + id)
	if templateCache != nil {
		template := templateCache.(*model.TemplateBss)
		return template, nil
	} else {
		template, err := repository.TemplateBssRepo.GetById(ctx, dbCon, id)
		if err != nil {
			return nil, err
		}
		cache.NewMemCache().Set(INFO_TEMPLATE+"_"+id, template, EXPIRE_TEMPLATE)
		return template, nil
	}
}
