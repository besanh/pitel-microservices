package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IOttMessage interface {
		GetOttMessage(ctx context.Context, data model.OttMessage) (int, any)
		GetCodeChallenge(ctx context.Context, authUser *model.AuthUser, appId string) (int, any)
		PostShareInfoEvent(ctx context.Context, authUser *model.AuthUser, data model.ShareInfo) (int, any)
	}
	OttMessage struct{}
)

func NewOttMessage() IOttMessage {
	return &OttMessage{}
}

/**
* PROBLEM: chat vao thi sao biet user thuoc db nao
* Khi chuyen qua fins thi lam sao biet setting nay cua db nao
 */
func (s *OttMessage) GetOttMessage(ctx context.Context, data model.OttMessage) (int, any) {
	var connectionCache model.ChatConnectionApp
	connectionAppCache := cache.RCache.Get(CHAT_CONNECTION + "_" + data.AppId + "_" + data.OaId)
	if connectionAppCache != nil {
		if err := json.Unmarshal([]byte(connectionAppCache.(string)), &connectionCache); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	} else {
		filter := model.ChatConnectionAppFilter{
			AppId:          data.AppId,
			OaId:           data.OaId,
			ConnectionType: data.MessageType,
		}
		total, connections, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if total < 1 {
			log.Error("connect for app_id: " + data.AppId + ", oa_id: " + data.OaId + " not found")
			return response.ServiceUnavailableMsg("connect for app_id: " + data.AppId + ", oa_id: " + data.OaId + " not found")
		}
		connectionCache = (*connections)[0]
		if err := cache.RCache.Set(CHAT_CONNECTION+"_"+data.AppId+"_"+data.OaId, (*connections)[0], CHAT_CONNECTION_EXPIRE); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	}

	docId := uuid.NewString()
	timestamp := time.Unix(0, data.Timestamp*int64(time.Millisecond))
	message := model.Message{
		Id:                  docId,
		ParentExternalMsgId: "",
		ExternalMsgId:       data.MsgId,
		MessageType:         data.MessageType,
		EventName:           data.EventName,
		Direction:           variables.DIRECTION["receive"],
		AppId:               data.AppId,
		OaId:                data.OaId,
		UserIdByApp:         data.UserIdByApp,
		ExternalUserId:      data.ExternalUserId,
		Avatar:              data.Avatar,
		SendTime:            timestamp,
		SendTimestamp:       data.Timestamp,
		Content:             data.Content,
		UserAppname:         data.Username,
		CreatedAt:           time.Now(),
		ShareInfo:           data.ShareInfo,
	}
	if slices.Contains[[]string](variables.EVENT_READ_MESSAGE, data.EventName) {
		timestamp := time.Unix(0, data.Timestamp*int64(time.Millisecond))
		message.ReadTime = timestamp
		message.ReadTimestamp = data.Timestamp
	}
	if data.Attachments != nil {
		for _, val := range *data.Attachments {
			var attachmentFile model.OttPayloadFile
			var attachmentMedia model.OttPayloadMedia
			var attachmentDetail model.OttAttachments
			var payload model.OttPayloadMedia
			attachmentDetail.AttType = val.AttType
			if val.AttType == "file" {
				if err := util.ParseAnyToAny(val.Payload, &payload); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				attachmentFile.Url = strings.ReplaceAll(attachmentFile.Url, "u0026", "&")
				attachmentDetail.Payload = &payload
			} else {
				if err := util.ParseAnyToAny(val.Payload, &payload); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				attachmentMedia.Url = strings.ReplaceAll(attachmentMedia.Url, "u0026", "&")
				attachmentDetail.Payload = &payload
			}
			message.Attachments = append(message.Attachments, &attachmentDetail)
		}
	}

	var isNew bool
	var conversation model.Conversation

	// TODO: check queue setting
	user, err := CheckChatSetting(ctx, message)
	if user.AuthUser != nil {
		data.TenantId = user.AuthUser.TenantId
		message.TenantId = user.AuthUser.TenantId
	} else {
		data.TenantId = connectionCache.TenantId
		message.TenantId = connectionCache.TenantId
		filterChatManageQueueUser := model.ChatManageQueueUserFilter{
			QueueId: connectionCache.QueueId,
		}
		totalManageQueueUser, manageQueueUser, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if totalManageQueueUser > 0 {
			user.QueueId = (*manageQueueUser)[0].QueueId
			user.ConnectionId = (*manageQueueUser)[0].ConnectionId
			user.IsOk = true
		}
	}
	if user.IsOk {
		conversationTmp, isNewTmp, errConv := UpSertConversation(ctx, user.ConnectionId, data)
		if errConv != nil {
			log.Error(errConv)
			return response.ServiceUnavailableMsg(errConv.Error())
		}
		conversation = conversationTmp
		isNew = isNewTmp

		// TODO: add rabbitmq message
		if len(conversation.ConversationId) > 0 {
			message.ConversationId = conversation.ConversationId
			message.IsRead = "deactive"
			if errMsg := InsertES(ctx, data.TenantId, ES_INDEX, conversation.AppId, docId, message); errMsg != nil {
				log.Error(errMsg)
				return response.ServiceUnavailableMsg(errMsg.Error())
			}
		}
	} else if len(user.QueueId) < 1 && err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if user.AuthUser != nil {
		// TODO: publish to rmq
		// if err := HandlePushRMQ(ctx, ES_INDEX, docId, message, tmpBytes); err != nil {
		// 	log.Error(err)
		// 	return response.ServiceUnavailableMsg(err.Error())
		// }

		if isNew {
			event := model.Event{
				EventName: variables.EVENT_CHAT[4],
				EventData: &model.EventData{
					Conversation: conversation,
				},
			}
			for s := range WsSubscribers.Subscribers {
				if s.Id == user.AuthUser.UserId && (s.Level == "agent" || s.Level == "user") {
					if err := PublishMessageToOne(user.AuthUser.UserId, event); err != nil {
						log.Error(err)
						return response.ServiceUnavailableMsg(err.Error())
					}
					break
				}
			}
		}
		event := model.Event{
			EventName: variables.EVENT_CHAT[3],
			EventData: &model.EventData{
				Message: message,
			},
		}

		for s := range WsSubscribers.Subscribers {
			if s.Id == user.AuthUser.UserId && (s.Level == "agent" || s.Level == "user") {
				if err := PublishMessageToOne(user.AuthUser.UserId, event); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				break
			}
		}
	}

	if len(conversation.ConversationId) < 1 {
		log.Error("conversation not found")
		return response.NotFoundMsg("conversation not found")
	}
	if len(conversation.TenantId) < 1 {
		conversation.TenantId = connectionCache.TenantId
	}

	// TODO: publish message to manager
	if len(user.QueueId) > 0 {
		manageQueueUser, err := GetManageQueueUser(ctx, user.QueueId)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if len(manageQueueUser.Id) < 1 {
			log.Error("queue " + user.QueueId + " not found")
			return response.NotFoundMsg("queue " + user.QueueId + " not found")
		}
		if len(manageQueueUser.ConnectionId) < 1 {
			manageQueueUser.ConnectionId = connectionCache.Id
		}

		filter := model.UserAllocateFilter{
			AppId:          conversation.AppId,
			ConversationId: conversation.ConversationId,
			MainAllocate:   "active",
		}
		totalUserAllocate, _, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if totalUserAllocate < 1 {
			userAllocation := model.UserAllocate{
				Base:               model.InitBase(),
				TenantId:           conversation.TenantId,
				ConversationId:     conversation.ConversationId,
				AppId:              message.AppId,
				UserId:             manageQueueUser.ManageId,
				QueueId:            manageQueueUser.QueueId,
				AllocatedTimestamp: time.Now().Unix(),
				MainAllocate:       "active",
				ConnectionId:       manageQueueUser.ConnectionId,
			}
			if err := repository.UserAllocateRepo.Insert(ctx, repository.DBConn, userAllocation); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}

			if err := cache.RCache.Set(USER_ALLOCATE+"_"+conversation.ConversationId, userAllocation, USER_ALLOCATE_EXPIRE); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}

		// TODO: publish message to manager
		if isNew {
			event := model.Event{
				EventName: variables.EVENT_CHAT[0],
				EventData: &model.EventData{
					Conversation: conversation,
				},
			}
			if err := PublishMessageToOne(manageQueueUser.ManageId, event); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
		event := model.Event{
			EventName: variables.EVENT_CHAT[3],
			EventData: &model.EventData{
				Message: message,
			},
		}
		for s := range WsSubscribers.Subscribers {
			if s.Id == manageQueueUser.ManageId {
				if err := PublishMessageToOne(manageQueueUser.ManageId, event); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
			}
		}

		// TODO: publish to admin
		if ENABLE_PUBLISH_ADMIN {
			userUuids := []string{}
			for s := range WsSubscribers.Subscribers {
				if s.TenantId == conversation.TenantId && s.Level == "admin" {
					userUuids = append(userUuids, s.Id)
				}
			}

			if len(userUuids) > 0 {
				if isNew {
					eventConversation := model.Event{
						EventName: variables.EVENT_CHAT[4],
						EventData: &model.EventData{
							Conversation: conversation,
						},
					}
					if err := PublishMessageToMany(userUuids, eventConversation); err != nil {
						log.Error(err)
						return response.ServiceUnavailableMsg(err.Error())
					}
				}

				if err := PublishMessageToMany(userUuids, event); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
			}
		}
	}

	return response.OKResponse()
}

func (s *OttMessage) GetCodeChallenge(ctx context.Context, authUser *model.AuthUser, appId string) (int, any) {
	url := OTT_URL + "/ott/v1/zalo/code-challenge/" + appId
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	if resp.StatusCode() == 200 {
		var result model.OttCodeChallenge
		if err := json.Unmarshal([]byte(resp.Body()), &result); err != nil {
			return response.ServiceUnavailableMsg(err.Error())
		}
		return response.OK(result)
	} else {
		return response.ServiceUnavailableMsg(resp.String())
	}
}

func (s *OttMessage) PostShareInfoEvent(ctx context.Context, authUser *model.AuthUser, data model.ShareInfo) (int, any) {
	// Because submit info is used sometimes, not to use struct
	event := map[string]any{
		"event_name": "share_info",
		"event_data": map[string]any{
			"share_info": data,
		},
	}
	if err := PublishMessageToOne(authUser.UserId, event); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OKResponse()
}
