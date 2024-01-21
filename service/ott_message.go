package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"golang.org/x/exp/slices"
)

type (
	IOttMessage interface {
		GetOttMessage(ctx context.Context, data model.OttMessage) (int, any)
	}
	OttMessage struct{}
)

func NewOttMessage() IOttMessage {
	return &OttMessage{}
}

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
			var attachmentDetail model.AttachmentsDetails
			attachmentDetail.AttachmentType = val.AttType
			if val.AttType == variables.ATTACHMENT_TYPE["file"] {
				if err := util.ParseAnyToAny(val.Payload, &attachmentFile); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				attachmentFile.Url = strings.ReplaceAll(attachmentFile.Url, "u0026", "&")
				attachmentDetail.AttachmentFile = &attachmentFile
			} else {
				if err := util.ParseAnyToAny(val.Payload, &attachmentMedia); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				attachmentMedia.Url = strings.ReplaceAll(attachmentMedia.Url, "u0026", "&")
				attachmentDetail.AttachmentMedia = &attachmentMedia
			}
			message.Attachments = append(message.Attachments, &attachmentDetail)
		}
	}

	// TODO: check conversation and add message
	conversation, isNew, err := UpSertConversation(ctx, data)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	message.ConversationId = conversation.ConversationId

	//  TODO: add rabbitmq message
	if err := InsertES(ctx, data.AppId, ES_INDEX, docId, message); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// if err := HandlePushRMQ(ctx, ES_INDEX, docId, message, tmpBytes); err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// }

	var agentId string

	// TODO: check queue setting
	filter := model.QueueFilter{
		AppId: data.AppId,
	}
	agentId, err = CheckChatQueueSetting(ctx, filter, data.ExternalUserId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(agentId) > 0 {
		if isNew {
			event := map[string]any{
				"event_name": "conversation_created",
				"event_data": map[string]any{
					"conversation": conversation,
				},
			}
			PublishMessageToOne(agentId, event)
		}
		event := map[string]any{
			"event_name": "message_created",
			"event_data": map[string]any{
				"message": message,
			},
		}
		PublishMessageToOne(agentId, event)
	}

	return response.OKResponse()
}
