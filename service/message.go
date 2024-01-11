package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IMessage interface {
		SendMessageToOTT(ctx context.Context, authUser *model.AuthUser, data model.MessageRequest) (int, any)
		GetMessages(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int) (int, any)
	}
	Message struct {
		OttSendMessageUrl string
	}
)

func NewMessage(OttSendMessageUrl string) IMessage {
	return &Message{
		OttSendMessageUrl: OttSendMessageUrl,
	}
}

func (s *Message) SendMessageToOTT(ctx context.Context, authUser *model.AuthUser, data model.MessageRequest) (int, any) {
	conversation := model.Conversation{}
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.UserIdByApp)
	if conversationCache != nil {
		if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	} else {
		filter := model.ConversationFilter{}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, data.AppId, ES_INDEX_CONVERSATION, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if total > 0 {
			conversation = (*conversations)[0]
			if err := cache.RCache.Set(CONVERSATION+"_"+data.UserIdByApp, conversation, CONVERSATION_EXPIRE); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
	}

	timestampTmp := time.Now().UnixMilli()
	timestamp := fmt.Sprintf("%d", timestampTmp)
	eventName := "text"
	if len(data.Attachments) > 0 {
		for _, item := range data.Attachments {
			eventNameTmp, ok := variables.ATTACHMENT_TYPE[item.AttachmentType]
			if !ok {
				break
			}
			eventName = eventNameTmp
		}
	}
	if len(eventName) < 1 {
		log.Errorf("event name %s not found", eventName)
		return response.BadRequestMsg("event name " + eventName + " not found")
	}

	docId := uuid.NewString()

	// Store ES
	message := model.Message{
		ParentExternalMsgId: "",
		Id:                  docId,
		MessageType:         conversation.ConversationType,
		EventName:           eventName,
		Direction:           variables.DIRECTION["send"],
		AppId:               conversation.AppId,
		OaId:                conversation.OaId,
		UserIdByApp:         conversation.UserIdByApp,
		Avatar:              conversation.Avatar,
		SupporterId:         authUser.UserId,
		SupporterName:       authUser.Username,
		SendTime:            time.Now(),
		SendTimestamp:       timestampTmp,
		Content:             data.Content,
		Attachments:         data.Attachments,
	}
	log.Info(message)
	if err := InsertES(ctx, conversation.AppId, ES_INDEX, docId, message); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// Send to OTT
	ottMessage := model.SendMessageToOtt{
		Type:          conversation.ConversationType,
		EventName:     eventName,
		AppId:         conversation.AppId,
		UserIdByApp:   conversation.UserIdByApp,
		Uid:           conversation.Uid,
		SupporterId:   authUser.UserId,
		SupporterName: authUser.Username,
		Timestamp:     timestamp,
		Text:          data.Content,
	}

	resOtt, err := s.sendMessageToOTT(ctx, ottMessage)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	message.ExternalMsgId = resOtt.MsgId

	// Update msgId to ES
	tmpBytes, err := json.Marshal(message)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX, docId, esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OKResponse()
}

func (s *Message) GetMessages(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int) (int, any) {
	total, messages, err := repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(messages, total, limit, offset)
}
