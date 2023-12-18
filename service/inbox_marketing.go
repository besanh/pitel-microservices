package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IInboxMarketing interface {
		SendInboxMarketing(ctx context.Context, authUser *model.AuthUser, data model.InboxMarketingRequest) (message string, err error)
	}
	InboxMarketing struct {
		StorageDialect string
	}
	ColumnInfo struct {
		Index int
		Name  string
	}
)

var (
	PLUGIN_INCOM     = "plugin_incom"
	TTL_PLUGIN_INCOM = 1 * time.Minute
)

func NewInboxMarketing() IInboxMarketing {
	return &InboxMarketing{}
}

func (s *InboxMarketing) SendInboxMarketing(ctx context.Context, authUser *model.AuthUser, data model.InboxMarketingRequest) (message string, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return err.Error(), errors.New(response.ERR_EMPTY_CONN)
	}

	routingConfig, err := common.GetInfoPlugin(ctx, dbCon, authUser, data.RoutingConfig)
	if err != nil {
		log.Error(err)
		return "get rounting info error", err
	} else if len(routingConfig.RoutingName) < 1 {
		return "routing not found in system", errors.New("routing not found in system")
	}

	// Force foreign phone numbers for abenla
	if util.HandleNetwork(data.PhoneNumber) == "dnc" && routingConfig.RoutingOption.Abenla.Status {
		routingConfig.RoutingName = "abenla"
		routingConfig.RoutingType = "abenla"
	}

	// Handle content match template
	keysContent, keysTemplate, err := common.HandleCheckContentMatchTemplate(ctx, dbCon, authUser, data.Template, data.Content)
	if err != nil {
		log.Error(err)
		return "check content match template error", err
	} else if len(keysContent) != len(keysTemplate) {
		return err.Error(), errors.New("content not match template")
	}

	// Check connection to plugin
	err = common.CheckConnectionWithExternalPlugin(ctx, *routingConfig)
	if err != nil {
		log.Error(err)
		return "check connection error", err
	}

	plugin := ""
	if routingConfig.RoutingOption.Incom.Status {
		plugin = "incom"
	}
	if routingConfig.RoutingOption.Abenla.Status {
		plugin = "abenla"
	}
	if routingConfig.RoutingOption.Fpt.Status {
		plugin = "fpt"
	}

	docId := uuid.NewString()
	smsGuid := docId
	// index := authUser.DatabaseEsIndex
	inboxMarketing := model.InboxMarketingLogInfo{
		Id:                docId,
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		Services:          authUser.Services,
		RoutingConfigUuid: routingConfig.Id,
		Plugin:            plugin,
		PhoneNumber:       data.PhoneNumber,
		CountAction:       1,
		CreatedAt:         time.Now(),
		CreatedBy:         authUser.UserId,
	}

	// Add param
	if plugin == "abenla" {
		inboxMarketing.Message = data.Content
		inboxMarketing.ExternalMessageId = smsGuid
		serviceTypeId, _ := strconv.Atoi(routingConfig.RoutingOption.Abenla.ServiceTypeId)
		inboxMarketing.ServiceTypeId = serviceTypeId
		_ = keysContent
		_ = keysTemplate
	} else if plugin == "incom" {
		// Handle map list param array with key in template
		listParam := json.RawMessage{}
		tmp := make(map[string]any)
		for i, item := range keysTemplate {
			tmp[item] = keysContent[i]
		}
		tmpByte, err := json.Marshal(tmp)
		if err != nil {
			log.Error(err)
			return err.Error(), err
		}
		listParam = append(listParam, tmpByte...)
		inboxMarketing.Message = data.Content
		inboxMarketing.ListParam = listParam

		// Get rout rule from flow_uuid
		routeRules := []string{}
		if routingConfig.RoutingFlow.FlowType == "recipient" {
			for _, val := range routingConfig.RoutingFlow.FlowUuid {
				routeRule, err := common.HandleGetRouteRule(ctx, dbCon, val)
				if err != nil {
					log.Error(err)
					return err.Error(), err
				}
				routeRules = append(routeRules, routeRule)
			}
		}
		if len(routeRules) < 1 {
			return err.Error(), errors.New("route rule not found")
		}
		inboxMarketing.RouteRule = routeRules
		// inboxMarketing.RouteRule = routingConfig.RoutingOption.Incom.RouteRule
	} else if plugin == "fpt" {

	}

	tmpBytes := []byte{}
	if plugin == "abenla" {
		// log
		auditLogModel := model.LogInboxMarketing{
			TenantId:          authUser.TenantId,
			BusinessUnitId:    authUser.BusinessUnitId,
			UserId:            authUser.UserId,
			Username:          authUser.Username,
			Services:          authUser.Services,
			Id:                docId,
			RoutingConfigUuid: routingConfig.Id,
			ExternalMessageId: smsGuid,
			Channel:           "sms",
			Status:            "",
			Quantity:          0,
			TelcoId:           0,
			IsChargedZns:      false,
			IsCheck:           false,
			Code:              0,
			CountAction:       1,
			UpdatedBy:         authUser.UserId,
		}
		auditLog, err := common.HandleAuditLogInboxMarketing(auditLogModel)
		if err != nil {
			log.Error(err)
			return err.Error(), err
		}
		inboxMarketing.Log = append(inboxMarketing.Log, auditLog)
		tmpBytes, err = json.Marshal(inboxMarketing)
		if err != nil {
			log.Error(err)
			return err.Error(), err
		}
	} else if plugin == "incom" {
		// log
		auditLogModel := model.LogInboxMarketing{
			TenantId:          authUser.TenantId,
			BusinessUnitId:    authUser.BusinessUnitId,
			UserId:            authUser.UserId,
			Username:          authUser.Username,
			Services:          authUser.Services,
			Id:                docId,
			RoutingConfigUuid: routingConfig.Id,
			ExternalMessageId: "",
			Status:            "",
			Quantity:          0,
			TelcoId:           0,
			IsChargedZns:      false,
			IsCheck:           false,
			Code:              0,
			CountAction:       1,
			UpdatedBy:         authUser.UserId,
		}
		auditLog, err := common.HandleAuditLogInboxMarketing(auditLogModel)
		if err != nil {
			log.Error(err)
			return err.Error(), err
		}
		inboxMarketing.Log = append(inboxMarketing.Log, auditLog)
		tmpBytes, err = json.Marshal(inboxMarketing)
		if err != nil {
			log.Error(err)
			return err.Error(), err
		}
	}

	// Push rabbitmq
	// err = common.HandlePushRMQ(ctx, index, docId, pluginInfo, tmpBytes)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err)
	// }
	// Sleep to rmq write to es
	// time.Sleep(2 * time.Second)
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return err.Error(), err
	}
	if err := repository.ESRepo.InsertLog(ctx, authUser.TenantId, authUser.DatabaseEsIndex, docId, esDoc); err != nil {
		log.Error(err)
		return err.Error(), err
	}

	// Check pool to delivery for channel or delivery to plugin immediately
	if plugin == "abenla" {
		inboxMarketingBasic := model.InboxMarketingBasic{
			TenantId:          authUser.TenantId,
			BusinessUnitId:    authUser.BusinessUnitId,
			UserId:            authUser.UserId,
			Username:          authUser.Username,
			Services:          authUser.Services,
			ExternalMessageId: smsGuid,
			DocId:             docId,
			Index:             authUser.DatabaseEsIndex,
			UpdatedBy:         authUser.UserId,
		}
		_, err := HandleMainInboxMarketingAbenla(ctx, authUser, inboxMarketingBasic, *routingConfig, inboxMarketing, data)
		if err != nil {
			log.Error(err)
			return err.Error(), err
		}

		return err.Error(), err
	} else if plugin == "incom" {
		inboxMarketingBasic := model.InboxMarketingBasic{
			TenantId:       authUser.TenantId,
			BusinessUnitId: authUser.BusinessUnitId,
			UserId:         authUser.UserId,
			Username:       authUser.Username,
			Services:       authUser.Services,
			DocId:          docId,
			Index:          authUser.DatabaseEsIndex,
			UpdatedBy:      authUser.UserId,
		}
		res, err := HandleMainInboxMarketingIncom(ctx, dbCon, inboxMarketingBasic, *routingConfig, authUser, inboxMarketing, data)
		if err != nil {
			log.Error(err)
			return err.Error(), err
		}

		if res.Code < 1 {
			return res.Status, nil
		} else {
			return res.Status, nil
		}
	}
	return err.Error(), err
}
