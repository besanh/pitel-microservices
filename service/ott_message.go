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
			if val.AttType == variables.ATTACHMENT_TYPE_MAP["file"] {
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
		filter := model.ChatConnectionAppFilter{
			AppId: data.AppId,
			OaId:  data.OaId,
		}
		total, connection, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if total > 0 {
			data.TenantId = (*connection)[0].TenantId
			message.TenantId = (*connection)[0].TenantId
		}
		filterChatManageQueueAgent := model.ChatManageQueueAgentFilter{}
		totalManageQueueAgent, manageQueueAgent, err := repository.ManageQueueRepo.GetManageQueue(ctx, repository.DBConn, filterChatManageQueueAgent, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if totalManageQueueAgent > 0 {
			user.QueueId = (*manageQueueAgent)[0].QueueId
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
	} else if err != nil {
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
			event := map[string]any{
				"event_name": "conversation_created",
				"event_data": map[string]any{
					"conversation": conversation,
				},
			}
			if err := PublishMessageToOne(user.AuthUser.UserId, event); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
		event := map[string]any{
			"event_name": variables.EVENT_CHAT["message_created"],
			"event_data": map[string]any{
				"message": message,
			},
		}

		for s := range WsSubscribers.Subscribers {
			if s.Id == user.AuthUser.UserId {
				if err := PublishMessageToOne(user.AuthUser.UserId, event); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				break
			}
		}
	}

	// TODO: publish message to manager
	if len(user.QueueId) > 0 {
		manageQueueAgent, err := GetManageQueueAgent(ctx, user.QueueId)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		} else if len(manageQueueAgent.Id) < 1 {
			log.Error("queue " + user.QueueId + " not found")
			return response.NotFoundMsg("queue " + user.QueueId + " not found")
		}

		filter := model.AgentAllocateFilter{
			AppId:          conversation.AppId,
			ConversationId: conversation.ConversationId,
			MainAllocate:   "active",
		}
		totalAgentAllocate, _, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if totalAgentAllocate < 1 {
			agentAllocation := model.AgentAllocate{
				Base:               model.InitBase(),
				TenantId:           conversation.TenantId,
				ConversationId:     conversation.ConversationId,
				AppId:              message.AppId,
				AgentId:            manageQueueAgent.AgentId,
				QueueId:            manageQueueAgent.QueueId,
				AllocatedTimestamp: time.Now().Unix(),
				MainAllocate:       "active",
				ConnectionId:       manageQueueAgent.ConnectionId,
			}
			if err := repository.AgentAllocationRepo.Insert(ctx, repository.DBConn, agentAllocation); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}

			if err := cache.RCache.Set(AGENT_ALLOCATION+"_"+conversation.ConversationId, agentAllocation, AGENT_ALLOCATION_EXPIRE); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}

		// TODO: publish message to manager
		if isNew {
			event := map[string]any{
				"event_name": variables.EVENT_CHAT["conversation_created"],
				"event_data": map[string]any{
					"conversation": conversation,
				},
			}
			if err := PublishMessageToOne(manageQueueAgent.AgentId, event); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
		event := map[string]any{
			"event_name": variables.EVENT_CHAT["message_created"],
			"event_data": map[string]any{
				"message": message,
			},
		}
		if err := PublishMessageToOne(manageQueueAgent.AgentId, event); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		// TODO: publish to admin
		// go
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
