package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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
		_ = keysContent
		_ = keysTemplate
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
	} else if plugin == "fpt" {
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
		accessToken, err := common.GetAccessTokenFpt(ctx, dbCon)
		if err != nil {
			log.Error(err)
			return "access token error", err
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
		}
		if res.Code < 1 {
			return res.Status, nil
		} else {
			return res.Status, nil
		}
	}
	return err.Error(), err
}

// func (s *InboxMarketing) HandleUpdateStatusMessCurrentday() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
// 	defer cancel()

// 	index := "pitel_bss_inbox_marketing"

// 	t := time.Now()
// 	startCurrentDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Local().Format(time.RFC3339)
// 	endCurrentDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local).Local().Format(time.RFC3339)

// 	total, logData, err := repository.InboxMarketingESRepo.GetLogCurrentDay(ctx, index, "incom", startCurrentDay, endCurrentDay)
// 	if err != nil {
// 		log.Error(err)
// 	} else if total > 0 {
// 		for _, v := range logData {
// 			if len(v.ExternalMessageId) > 0 {
// 				esDocTmp := map[string]any{}
// 				bytesTmp, err := json.Marshal(v)
// 				if err != nil {
// 					log.Error(err)
// 					continue
// 				}
// 				if err := json.Unmarshal(bytesTmp, &esDocTmp); err != nil {
// 					log.Error(err)
// 					continue
// 				}

// 				// Get cache for config user
// 				pluginInfo := model.PluginInfo{}
// 				pluginInfoCache, err := cache.MCache.Get(PLUGIN_INCOM)
// 				if err != nil {
// 					log.Error(err)
// 					if err := repository.ESRepo.UpdateDocById(ctx, v.DomainUuid, index, v.InboxMarketingLogUuid, esDocTmp); err != nil {
// 						log.Error(err)
// 					}
// 					continue
// 				} else if pluginInfoCache == nil {
// 					pluginExist, err := repository.ExternalPluginRepo.GetExternalPluginByIdOrPlugin(ctx, v.DomainUuid, v.ExternalPluginUuid, "")
// 					if err != nil {
// 						log.Error(err)
// 						continue
// 					} else if len(pluginExist.ExternalPluginUuid) == 0 {
// 						log.Error("plugin " + pluginInfo.ExternalPluginUuid + " not exist")
// 						continue
// 					}
// 					pluginInfo.DomainUuid = pluginExist.DomainUuid
// 					pluginInfo.ExternalPluginUuid = pluginExist.ExternalPluginUuid
// 					pluginInfo.Username = pluginExist.Username
// 					pluginInfo.Password = pluginExist.Password
// 					pluginInfo.ApiUrl = pluginExist.ApiUrl
// 					pluginInfo.Status = pluginExist.Status
// 				} else {
// 					pluginInfo = pluginInfoCache.(model.PluginInfo)
// 				}
// 				bodyStatusMess := model.IncomBodyStatus{
// 					Username:   pluginInfo.Username,
// 					Password:   pluginInfo.Password,
// 					IdOmniMess: v.ExternalId,
// 				}
// 				logExist, err := common.HandleGetStatusMessage(ctx, pluginInfo, v.InboxMarketingLogUuid, bodyStatusMess)
// 				if err != nil {
// 					log.Error(err)
// 					if err := repository.ESRepo.UpdateDocById(ctx, v.DomainUuid, index, v.InboxMarketingLogUuid, esDocTmp); err != nil {
// 						log.Error(err)
// 					}
// 					continue
// 				}
// 				countAction := logExist.CountAction + 1
// 				auditLogModel := model.LogInboxMarketing{
// 					DomainUuid:         logExist.DomainUuid,
// 					Id:                 logExist.ExternalPluginUuid,
// 					Plugin:             v.Plugin,
// 					ExternalPluginUuid: logExist.ExternalPluginUuid,
// 					ExternalId:         logExist.ExternalId,
// 					Channel:            logExist.Channel,
// 					ChannelHook:        v.ChannelHook,
// 					ErrorCode:          logExist.ErrorCode,
// 					ErrorCodeHook:      v.ErrorCodeHook,
// 					Status:             logExist.Status,
// 					TelcoId:            logExist.TelcoId,
// 					IsChargedZns:       logExist.IsChargedZns,
// 					IsCheck:            logExist.IsCheck,
// 					Code:               logExist.Code,
// 					CountAction:        countAction,
// 				}
// 				auditLog, err := common.HandleAuditLogInboxMarketing(auditLogModel)
// 				if err != nil {
// 					log.Error(err)
// 					continue
// 				}
// 				logExist.Log = append(logExist.Log, auditLog)
// 				logExist.CountAction = countAction

// 				esDoc := map[string]any{}
// 				tmpBytes, err := json.Marshal(logExist)
// 				if err != nil {
// 					log.Error(err)
// 					if err := repository.ESRepo.UpdateDocById(ctx, v.DomainUuid, index, v.InboxMarketingLogUuid, esDocTmp); err != nil {
// 						log.Error(err)
// 					}
// 					continue
// 				}
// 				if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
// 					log.Error(err)
// 					if err := repository.ESRepo.UpdateDocById(ctx, v.DomainUuid, index, v.InboxMarketingLogUuid, esDocTmp); err != nil {
// 						log.Error(err)
// 					}
// 					continue
// 				}
// 				log.Infof("data logExist - UpdateStatusMessInDay: updatedAt: %s - id %s - quantity update: %d, isChargedZns: %t", logExist.UpdatedAt, v.ExternalPluginUuid, logExist.Quantity, logExist.IsChargedZns)

// 				if err := repository.ESRepo.UpdateDocById(ctx, pluginInfo.DomainUuid, index, v.InboxMarketingLogUuid, esDoc); err != nil {
// 					log.Error(err)
// 					if err := repository.ESRepo.UpdateDocById(ctx, v.DomainUuid, index, v.InboxMarketingLogUuid, esDoc); err != nil {
// 						log.Error(err)
// 					}
// 					continue
// 				}
// 				// Set cache
// 				if err := cache.MCache.SetTTL(PLUGIN_INCOM, pluginInfo, TTL_PLUGIN_INCOM); err != nil {
// 					log.Error(err)
// 					if err := repository.ESRepo.UpdateDocById(ctx, v.DomainUuid, index, v.InboxMarketingLogUuid, esDoc); err != nil {
// 						log.Error(err)
// 					}
// 					continue
// 				}
// 			}
// 		}
// 	}
// }

func (s *InboxMarketing) GetReportInboxMarketing(ctx context.Context, authUser *model.AuthUser, filter model.InboxMarketingFilter, limit, offset int) (total int, result []model.InboxMarketingLogReport, err error) {
	total, data, err := repository.InboxMarketingESRepo.GetReport(ctx, authUser.TenantId, authUser.DatabaseEsIndex, limit, offset, filter)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, data, nil
}
