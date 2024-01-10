package service

import (
	"context"
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
		Uid:                 data.UserId,
		Avatar:              data.Avatar,
		SendTime:            timestamp,
		SendTimestamp:       data.Timestamp,
		Content:             data.Content,
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
			if val.AttType == variables.ATTACHMENT_TYPE["file"] {
				if err := util.ParseAnyToAny(val.Payload, &attachmentFile); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
			} else {
				if err := util.ParseAnyToAny(val.Payload, &attachmentMedia); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
			}
			message.Attachments = append(message.Attachments, &model.Attachments{
				Id:             uuid.NewString(),
				MsgId:          docId,
				AttachmentType: val.AttType,
				AttachmentsDetail: &model.AttachmentsDetail{
					AttachmentFile:  &attachmentFile,
					AttachmentMedia: &attachmentMedia,
				},
				SendTime:      timestamp,
				SendTimestamp: data.Timestamp,
			})
		}
	}

	// TODO: check conversation and add message
	conversationId := ""
	conversation, isExisted, err := GetConversationExist(ctx, data)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if !isExisted {
		id, err := InsertConversation(ctx, conversation)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		conversationId = id
	}

	message.ConversationId = conversationId

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
	agentId, err = CheckChatQueueSetting(ctx, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(agentId) > 0 {
		// PublishMessageToOne(agentId, message.Content)
		NewSubscriberService().PublishMessageToSubscriber(ctx, agentId, message)
	}

	// TODO: add to queue

	return response.OKResponse()
}
