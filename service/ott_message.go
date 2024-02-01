package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
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
		GetCodeChallenge(ctx context.Context, authUser *model.AuthUser, appId string) (int, any)
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

	var agentId string

	// TODO: check queue setting
	agentId, err := CheckChatSetting(ctx, message)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(agentId) > 0 {
		// TODO: check conversation and add message
		conversation, isNew, err := UpSertConversation(ctx, data)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		//  TODO: add rabbitmq message
		if len(conversation.ConversationId) > 0 {
			message.ConversationId = conversation.ConversationId
			message.IsRead = "deactive"
			if err := InsertES(ctx, data.AppId, ES_INDEX, docId, message); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}

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
			PublishMessageToOne(agentId, event)
		}
		event := map[string]any{
			"event_name": "message_created",
			"event_data": map[string]any{
				"message": message,
			},
		}
		PublishMessageToOne(agentId, event)
	} else {
		// TODO: check conversation and add message
		conversation, _, err := UpSertConversation(ctx, data)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		// TODO: add rabbitmq message
		if len(conversation.ConversationId) > 0 {
			message.ConversationId = conversation.ConversationId
			message.IsRead = "deactive"
			if err := InsertES(ctx, data.AppId, ES_INDEX, docId, message); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}

		// if err := HandlePushRMQ(ctx, ES_INDEX, docId, message, tmpBytes); err != nil {
		// 	log.Error(err)
		// 	return response.ServiceUnavailableMsg(err.Error())
		// }
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
