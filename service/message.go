package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
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
)

type (
	IMessage interface {
		SendMessageToOTT(ctx context.Context, authUser *model.AuthUser, data model.MessageRequest, file *multipart.FileHeader) (int, any)
		GetMessages(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int) (int, any)
		MarkReadMessages(ctx context.Context, authUser *model.AuthUser, data model.MessageMarkRead) (int, any)
		ShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfo) (int, any)
	}
	Message struct{}
)

func NewMessage() IMessage {
	return &Message{}
}

func (s *Message) SendMessageToOTT(ctx context.Context, authUser *model.AuthUser, data model.MessageRequest, file *multipart.FileHeader) (int, any) {
	conversation := model.Conversation{}
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.ConversationId)
	if conversationCache != nil {
		if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	} else {
		if len(data.ConversationId) < 1 {
			return response.BadRequestMsg("conversation " + data.ConversationId + " is required")
		}
		filter := model.ConversationFilter{
			ConversationId: []string{data.ConversationId},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if total > 0 {
			if err := util.ParseAnyToAny((*conversations)[0], &conversation); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
			if err := cache.RCache.Set(CONVERSATION+"_"+conversation.ConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		} else {
			return response.BadRequestMsg("conversation: " + data.ConversationId + " not found")
		}
	}

	timestampTmp := time.Now().UnixMilli()
	timestamp := fmt.Sprintf("%d", timestampTmp)
	eventName := "text"
	if len(data.Attachments) > 0 {
		for _, item := range data.Attachments {
			eventNameTmp, ok := variables.ATTACHMENT_TYPE_MAP[item.AttachmentType]
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

	attachments := []model.OttAttachments{}

	// Upload to Docs
	// fileUrl, err := s.UploadDoc(authUser, data, file)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// }
	// if len(data.Attachments) > 0 && (len(data.EventName) > 0 && data.EventName != "text") {
	// 	data.Attachments[0].AttachmentFile.Url = fileUrl
	// 	data.AttType = eventName
	// 	Payload: data.Attachments
	// }

	// Send to OTT
	ottMessage := model.SendMessageToOtt{
		Type:          conversation.ConversationType,
		EventName:     eventName,
		AppId:         conversation.AppId,
		OaId:          conversation.OaId,
		Uid:           conversation.ExternalUserId,
		SupporterId:   authUser.UserId,
		SupporterName: authUser.Username,
		Timestamp:     timestamp,
		Text:          data.Content,
	}

	log.Info("message to ott: ", ottMessage)

	resOtt, err := s.sendMessageToOTT(ottMessage, &attachments)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// Store ES
	message := model.Message{
		TenantId:            conversation.TenantId,
		ParentExternalMsgId: "",
		Id:                  docId,
		MessageType:         conversation.ConversationType,
		ConversationId:      conversation.ConversationId,
		ExternalMsgId:       resOtt.Data.MsgId,
		EventName:           eventName,
		Direction:           variables.DIRECTION["send"],
		AppId:               conversation.AppId,
		OaId:                conversation.OaId,
		Avatar:              conversation.Avatar,
		SupporterId:         authUser.UserId,
		SupporterName:       authUser.Username,
		SendTime:            time.Now(),
		SendTimestamp:       timestampTmp,
		Content:             data.Content,
		Attachments:         data.Attachments,
	}
	log.Info("message to es: ", message)

	// Should to queue
	if err := InsertES(ctx, conversation.TenantId, ES_INDEX, message.AppId, docId, message); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.Created(message)
}

func (s *Message) GetMessages(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int) (int, any) {
	total, messages, err := repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.Pagination(messages, total, limit, offset)
}

func (s *Message) MarkReadMessages(ctx context.Context, authUser *model.AuthUser, data model.MessageMarkRead) (int, any) {
	conversation := model.Conversation{}
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.ConversationId)
	if conversationCache != nil {
		if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
	} else {
		conversationFilter := model.ConversationFilter{
			ConversationId: []string{data.ConversationId},
			TenantId:       authUser.TenantId,
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, conversationFilter, 1, 0)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		if total > 0 {
			conversation := (*conversations)[0]
			if err := cache.RCache.Set(CONVERSATION+"_"+conversation.ConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
	}

	totalSuccess := 0
	totalFail := 0
	listMessageIdSuccess := map[string]string{}
	listMessageIdFail := map[string]string{}

	if data.ReadAll {
		filter := model.MessageFilter{
			IsRead:         "deactive",
			ConversationId: data.ConversationId,
		}
		total, messages, err := repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX, filter, -1, 0)
		if err != nil {
			log.Error(err)
			totalSuccess -= 1
			totalFail += 1
			return response.ServiceUnavailableMsg(err.Error())
		}
		totalSuccess = total
		if total > 0 {
			for _, item := range *messages {
				item.IsRead = "active"
				item.ReadBy = append(item.ReadBy, authUser.UserId)
				item.ReadTimestamp = time.Now().Unix()
				item.UpdatedAt = time.Now()
				item.ReadTime = time.Now()

				tmpBytes, err := json.Marshal(item)
				if err != nil {
					log.Error(err)
					totalSuccess -= 1
					totalFail += 1
					return response.ServiceUnavailableMsg(err.Error())
				}
				esDoc := map[string]any{}
				if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
					log.Error(err)
					totalSuccess -= 1
					totalFail += 1
					return response.ServiceUnavailableMsg(err.Error())
				}
				if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX, item.AppId, item.Id, esDoc); err != nil {
					log.Error(err)
					totalSuccess -= 1
					totalFail += 1
					return response.ServiceUnavailableMsg(err.Error())
				}
			}
		}

		return response.OK(map[string]any{
			"total_success": totalSuccess,
			"total_fail":    totalFail,
			"list_fail":     listMessageIdFail,
			"list_success":  listMessageIdSuccess,
		})
	} else {
		for _, item := range data.MessageIds {
			// Need tracking message stil not read ?
			message, err := repository.MessageESRepo.GetMessageById(ctx, "", ES_INDEX, item)
			if err != nil {
				log.Error(err)
				totalSuccess -= 1
				totalFail += 1
				listMessageIdFail[item] = err.Error()
				continue
			} else if len(message.Id) < 1 {
				totalSuccess -= 1
				totalFail += 1
				log.Errorf("message %s not found", item)
				listMessageIdFail[item] = "message " + item + " not found"
				continue
			}

			message.IsRead = "active"
			message.ReadBy = append(message.ReadBy, authUser.UserId)
			message.ReadTimestamp = time.Now().Unix()
			message.UpdatedAt = time.Now()
			message.ReadTime = time.Now()

			tmpBytes, err := json.Marshal(message)
			if err != nil {
				log.Error(err)
				totalSuccess -= 1
				totalFail += 1
				listMessageIdFail[item] = err.Error()
				continue
			}

			esDoc := map[string]any{}
			if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
				log.Error(err)
				totalSuccess -= 1
				totalFail += 1
				listMessageIdFail[item] = err.Error()
			}
			newConversationId := GenerateConversationId(data.AppId, item)
			if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX, data.AppId, newConversationId, esDoc); err != nil {
				log.Error(err)
				totalSuccess -= 1
				totalFail += 1
				listMessageIdFail[item] = err.Error()
			}
			listMessageIdSuccess[item] = "success"
		}

		return response.OK(map[string]any{
			"total_success": totalSuccess,
			"total_fail":    totalFail,
			"list_fail":     listMessageIdFail,
			"list_success":  listMessageIdSuccess,
		})
	}
}

func (s *Message) ShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfo) (int, any) {
	body := map[string]string{
		"name":     data.Fullname,
		"phone":    data.PhoneNumber,
		"address":  data.Address,
		"city":     data.City,
		"district": data.District,
	}

	url := OTT_URL
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(url)
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
