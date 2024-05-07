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
	if len(data.EventName) > 0 {
		eventName = data.EventName
	}

	docId := uuid.NewString()

	attachments := []*model.OttAttachments{}

	// Upload to Docs
	if len(data.EventName) > 0 && data.EventName != "text" {
		fileUrl, err := s.UploadDoc(ctx, file)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		// var payload model.
		att := model.OttAttachments{
			Payload: &model.OttPayloadMedia{
				Url: fileUrl,
			},
			AttType: data.EventName,
		}
		attachments = append(attachments, &att)
	}
	content := data.Content
	if eventName != "text" {
		if file != nil {
			content = file.Filename
		}
	}

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
		Text:          content,
	}

	log.Info("message to ott: ", ottMessage)

	resOtt, err := s.sendMessageToOTT(ottMessage, attachments)
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
		SupporterName:       authUser.Fullname,
		SendTime:            time.Now(),
		SendTimestamp:       timestampTmp,
		Content:             data.Content,
		Attachments:         attachments,
	}
	log.Info("message to es: ", message)

	// Should to queue
	if err := InsertES(ctx, conversation.TenantId, ES_INDEX, message.AppId, docId, message); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// TODO: update conversation => after refresh page, this conversation will appearance on top the list
	conversation.UpdatedAt = time.Now().Format(time.RFC3339)
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversation.AppId, conversation.ConversationId, esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// TODO: send message to admin/manager/user
	// Exclude: if user send message, remove user, if admin send message, remove admin, ...
	var queueId string
	filter := model.ChatConnectionAppFilter{
		AppId:          message.AppId,
		OaId:           message.OaId,
		ConnectionType: conversation.ConversationType,
	}
	total, connection, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if total < 1 {
		log.Errorf("connection %s not found", (*connection)[0].Id)
		return response.ServiceUnavailableMsg("connection " + (*connection)[0].Id + " not found")
	} else {
		message.TenantId = (*connection)[0].TenantId
	}

	filterChatManageQueueUser := model.ChatManageQueueUserFilter{
		QueueId: (*connection)[0].QueueId,
	}
	totalManageQueueUser, manageQueueUser, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if totalManageQueueUser > 0 {
		queueId = (*manageQueueUser)[0].QueueId
	}

	if authUser.Level != "manager" {
		if len(queueId) > 0 {
			if err := SendEventToManage(ctx, authUser, message, queueId); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		} else {
			log.Errorf("queue %s not found in send event to manage", queueId)
		}
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

	var result model.ReadMessageResponse
	if data.ReadAll {
		filter := model.MessageFilter{
			IsRead:         "deactive",
			ConversationId: data.ConversationId,
		}
		total, messages, err := repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX, filter, -1, 0)
		if err != nil {
			log.Error(err)
			result.TotalSuccess -= 1
			result.TotalFail += 1
			return response.ServiceUnavailableMsg(err.Error())
		}
		result.TotalSuccess = total
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
					result.TotalSuccess -= 1
					result.TotalFail += 1
					return response.ServiceUnavailableMsg(err.Error())
				}
				esDoc := map[string]any{}
				if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
					log.Error(err)
					result.TotalSuccess -= 1
					result.TotalFail += 1
					return response.ServiceUnavailableMsg(err.Error())
				}
				if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX, item.AppId, item.Id, esDoc); err != nil {
					log.Error(err)
					result.TotalSuccess -= 1
					result.TotalFail += 1
					return response.ServiceUnavailableMsg(err.Error())
				}
			}
		}
	} else {
		for _, item := range data.MessageIds {
			// Need tracking message stil not read ?
			message, err := repository.MessageESRepo.GetMessageById(ctx, "", ES_INDEX, item)
			if err != nil {
				log.Error(err)
				result.TotalSuccess -= 1
				result.TotalFail += 1
				result.ListFail[item] = err.Error()
				continue
			} else if len(message.Id) < 1 {
				result.TotalSuccess -= 1
				result.TotalFail += 1
				log.Errorf("message %s not found", item)
				result.ListFail[item] = "message " + item + " not found"
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
				result.TotalSuccess -= 1
				result.TotalFail += 1
				result.ListFail[item] = err.Error()
				continue
			}

			esDoc := map[string]any{}
			if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
				log.Error(err)
				result.TotalSuccess -= 1
				result.TotalFail += 1
				result.ListFail[item] = err.Error()
			}
			newConversationId := GenerateConversationId(data.AppId, data.OaId, item)
			if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX, data.AppId, newConversationId, esDoc); err != nil {
				log.Error(err)
				result.TotalSuccess -= 1
				result.TotalFail += 1
				result.ListFail[item] = err.Error()
			}
			result.ListSuccess[item] = "success"
		}
	}
	return response.OK(result)
}

func (s *Message) ShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfo) (int, any) {
	body := model.ShareInfoSendToOtt{
		Name:     data.Fullname,
		Phone:    data.PhoneNumber,
		Address:  data.Address,
		City:     data.City,
		District: data.District,
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(OTT_URL)
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
