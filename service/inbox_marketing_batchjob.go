package service

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
