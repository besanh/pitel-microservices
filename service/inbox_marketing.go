package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IInboxMarketing interface {
		SendInboxMarketing(ctx context.Context, authUser *model.AuthUser, data model.InboxMarketingRequest) (response model.ResponseInboxMarketing, err error)
		// HandleUpdateStatusMessCurrentday()
		GetReportInboxMarketing(ctx context.Context, authUser *model.AuthUser, filter model.InboxMarketingFilter, limit, offset int) (total int, result []model.InboxMarketingLogReport, err error)
		PostExportReportInboxMarketing(ctx context.Context, authUser *model.AuthUser, fileType string, filter model.InboxMarketingFilter) (int, any)
	}
	InboxMarketing struct{}
	ColumnInfo     struct {
		Index int
		Name  string
	}
)

var (
	PLUGIN_INCOM      = "plugin_incom"
	RECIPIENT_NETWORK = "recipient_network"
	TTL_PLUGIN_INCOM  = 1 * time.Minute
)

func NewInboxMarketing() IInboxMarketing {
	return &InboxMarketing{}
}

func (s *InboxMarketing) SendInboxMarketing(ctx context.Context, authUser *model.AuthUser, data model.InboxMarketingRequest) (res model.ResponseInboxMarketing, err error) {
	res = model.ResponseInboxMarketing{
		Status: "failed",
		Code:   constants.STANDARD_CODE["3"],
	}
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		res.Message = err.Error()
		return res, errors.New(response.ERR_EMPTY_CONN)
	}

	routingConfig, err := common.GetInfoPlugin(ctx, dbCon, authUser, data.RoutingConfig)
	if err != nil {
		log.Error(err)
		res.Message = "get rounting info error"
		return res, err
	} else if len(routingConfig.RoutingName) < 1 {
		res.Message = "routing not found in system"
		return res, errors.New("routing not found in system")
	}

	network := util.HandleNetwork(data.PhoneNumber)
	ok := handleCheckNetworkWithRecipient(ctx, dbCon, network, routingConfig.RoutingFlow.FlowUuid)
	if !ok {
		res.Message = fmt.Sprintf("network %s not match with routing %s", network, routingConfig.RoutingName)
		return res, errors.New("network " + network + " not match with routing " + routingConfig.RoutingName)
	}

	// Handle content match template
	template, keysContent, keysTemplate, err := common.HandleCheckContentMatchTemplate(ctx, dbCon, authUser, data.Template, data.Content)
	if err != nil {
		log.Error(err)
		res.Message = "check content match template error"
		return res, err
	} else if len(keysContent) != len(keysTemplate) {
		res.Message = "content not match template"
		return res, errors.New("content not match template")
	}

	// Check connection to plugin
	err = common.CheckConnectionWithExternalPlugin(ctx, *routingConfig)
	if err != nil {
		log.Error(err)
		res.Message = "check connection error"
		return res, err
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
	content := strings.ReplaceAll(data.Content, "{{", "")
	content = strings.ReplaceAll(content, "}}", "")
	inboxMarketing := model.InboxMarketingLogInfo{
		Id:             docId,
		TenantId:       authUser.TenantId,
		BusinessUnitId: authUser.BusinessUnitId,
		UserId:         authUser.UserId,
		Username:       authUser.Username,
		// Services:          authUser.Services,
		FlowType:          routingConfig.RoutingFlow.FlowType,
		FlowUuid:          routingConfig.RoutingFlow.FlowUuid,
		RoutingConfigUuid: routingConfig.Id,
		Plugin:            plugin,
		Channel:           routingConfig.RoutingType,
		TemplateUuid:      data.Template,
		TemplateCode:      template.TemplateCode,
		Message:           content,
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
			res.Message = err.Error()
			return res, err
		}
		listParam = append(listParam, tmpByte...)
		inboxMarketing.Message = data.Content
		inboxMarketing.ListParam = listParam

		// Get rout rule from flow_uuid
		routeRules := []string{}
		if routingConfig.RoutingFlow.FlowType == "recipient" {
			routeRule, err := common.HandleGetRouteRule(ctx, dbCon, routingConfig.RoutingFlow.FlowUuid)
			if err != nil {
				log.Error(err)
				res.Message = err.Error()
				return res, err
			}
			routeRules = append(routeRules, routeRule)
		}
		if len(routeRules) < 1 {
			res.Message = "route rule not found"
			return res, errors.New("route rule not found")
		}
		inboxMarketing.RouteRule = routeRules
	} else if plugin == "fpt" {
		inboxMarketing.ExternalMessageId = docId
		_ = keysContent
		_ = keysTemplate
	}

	auditLogModel := model.LogInboxMarketing{
		Id:                docId,
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		FlowType:          routingConfig.RoutingFlow.FlowType,
		FlowUuid:          routingConfig.RoutingFlow.FlowUuid,
		RoutingConfigUuid: routingConfig.Id,
		Status:            "",
		Plugin:            plugin,
		Quantity:          0,
		TelcoId:           0,
		IsChargedZns:      false,
		IsCheck:           false,
		Code:              0,
		CountAction:       1,
		UpdatedBy:         authUser.UserId,
	}
	if plugin == "abenla" || plugin == "fpt" {
		auditLogModel.ExternalMessageId = smsGuid
	}

	auditLog, err := common.HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		log.Error(err)
		res.Message = err.Error()
		return res, err
	}
	inboxMarketing.Log = append(inboxMarketing.Log, auditLog)
	tmpBytes, err := json.Marshal(inboxMarketing)
	if err != nil {
		log.Error(err)
		res.Message = err.Error()
		return res, err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		res.Message = err.Error()
		return res, err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, authUser.DatabaseEsIndex, authUser.TenantId); err != nil {
		log.Error(err)
		res.Message = err.Error()
		return res, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, authUser.DatabaseEsIndex, authUser.TenantId); err != nil {
			log.Error(err)
			res.Message = err.Error()
			return res, err
		}
	}
	if err := repository.ESRepo.InsertLog(ctx, authUser.TenantId, authUser.DatabaseEsIndex, docId, esDoc); err != nil {
		log.Error(err)
		res.Message = err.Error()
		return res, err
	}

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
	// Check pool to delivery for channel or delivery to plugin immediately
	if plugin == "abenla" {
		inboxMarketingBasic.ExternalMessageId = smsGuid
		_, err := HandleMainInboxMarketingAbenla(ctx, authUser, inboxMarketingBasic, *routingConfig, inboxMarketing, data)
		if err != nil {
			log.Error(err)
			res.Message = err.Error()
			return res, err
		}
		res.Message = err.Error()
		return res, err
	} else if plugin == "incom" {
		res, err := HandleMainInboxMarketingIncom(ctx, dbCon, inboxMarketingBasic, *routingConfig, authUser, inboxMarketing, data)
		if err != nil {
			log.Error(err)
			return res, err
		}
		return res, nil
	} else if plugin == "fpt" {
		accessToken, err := common.GetAccessTokenFpt(ctx, *routingConfig)
		if err != nil {
			log.Error(err)
			res.Message = err.Error()
			return res, err
		} else if accessToken == "" {
			res.Message = err.Error()
			return res, errors.New("access token not found")
		}
		hasher := md5.New()
		hasher.Write([]byte(routingConfig.RoutingOption.Fpt.ClientSecret))
		sessionId := hex.EncodeToString(hasher.Sum(nil))
		fpt := model.FptRequireRequest{
			AccessToken: accessToken,
			SessionId:   sessionId,
		}
		res, err := HandleMainInboxMarketingFpt(ctx, authUser, inboxMarketingBasic, *routingConfig, inboxMarketing, data, fpt)
		if err != nil {
			log.Error(err)
			res.Message = err.Error()
			return res, err
		}
		return res, nil
	}
	res.Message = err.Error()
	return res, err
}

func (s *InboxMarketing) GetReportInboxMarketing(ctx context.Context, authUser *model.AuthUser, filter model.InboxMarketingFilter, limit, offset int) (total int, result []model.InboxMarketingLogReport, err error) {
	total, data, err := repository.InboxMarketingESRepo.GetReport(ctx, authUser.TenantId, authUser.DatabaseEsIndex, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, data, nil
}

func handleCheckNetworkWithRecipient(ctx context.Context, dbCon sqlclient.ISqlClientConn, network string, recipientUuid string) bool {
	var ok bool
	recipientConfigCache, err := cache.MCache.Get(RECIPIENT_NETWORK + "_" + recipientUuid)
	if err != nil {
		return false
	} else if recipientConfigCache != nil {
		recipientConfig := recipientConfigCache.(*model.RecipientConfig)
		if recipientConfig.Recipient == network {
			ok = true
		}
		return ok
	} else {
		recipientConfig, err := repository.RecipientConfigRepo.GetById(ctx, dbCon, recipientUuid)
		if err != nil {
			log.Error(err)
			return false
		} else if recipientConfig != nil {
			if recipientConfig.Recipient == network {
				ok = true
			}
			if err := cache.MCache.SetTTL(RECIPIENT_NETWORK+"_"+recipientUuid, recipientConfig, common.EXPIRE_RECIPIENT); err != nil {
				log.Error(err)
				return false
			}
			return ok
		}
		return !ok
	}
}
