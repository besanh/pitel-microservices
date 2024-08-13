package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IMessage interface {
		SendMessageToOTT(ctx context.Context, authUser *model.AuthUser, data model.MessageRequest, file *multipart.FileHeader) (result model.Message, err error)
		// SendMessageToOTTAsync(ctx context.Context, authUser *model.AuthUser, data model.MessageRequest, file *multipart.FileHeader) (int, any)
		GetMessages(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int) (total int, messages *[]model.Message, err error)
		GetMessagesWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit int, scrollId string) (total int, messages []*model.Message, respScrollId string, err error)
		MarkReadMessages(ctx context.Context, authUser *model.AuthUser, data model.MessageMarkRead) (result model.ReadMessageResponse, err error)
		ShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfo) (result model.OttCodeChallenge, err error)
		GetMessageMediasWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit int, scrollId string) (total int, messages []*model.MessageAttachmentsDetails, respScrollId string, err error)
		PostTicketToMessage(ctx context.Context, authUser *model.AuthUser, data model.MessagePostTicket) (err error)
	}
	Message struct{}
)

var MessageService IMessage

func NewMessage() IMessage {
	return &Message{}
}

func (s *Message) SendMessageToOTT(ctx context.Context, authUser *model.AuthUser, data model.MessageRequest, file *multipart.FileHeader) (message model.Message, err error) {
	conversation := model.Conversation{}
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.ConversationId)
	if conversationCache != nil {
		if err = json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return
		}
	} else {
		conversationTmp, errTmp := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, data.AppId, data.ConversationId)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if conversationTmp == nil {
			err = errors.New("conversation: " + data.ConversationId + " not found")
			log.Error(err)
			return
		}
		conversation = *conversationTmp
		if err = cache.RCache.Set(CONVERSATION+"_"+conversation.ConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}

	if ENABLE_CHAT_POLICY_SETTINGS {
		// block messages sent outside of chat window
		conversationTime := conversation.UpdatedAt
		if len(conversationTime) < 1 {
			conversationTime = conversation.CreatedAt
		}
		if err = CheckOutOfChatWindowTime(ctx, conversation.TenantId, conversation.ConversationType, conversationTime); err != nil {
			return
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
		var fileUrl string
		if file == nil {
			if len(data.Url) < 1 {
				err = errors.New("url or file is required")
				return
			}
			// TODO: validate url
			fileUrl = data.Url
		} else {
			fileUrlTmp, errTmp := UploadDoc(ctx, data.AppId, data.OaId, file)
			if errTmp != nil {
				err = errTmp
				log.Error(err)
				return
			}
			fileUrl = fileUrlTmp
		}
		if len(fileUrl) < 1 {
			err = errors.New("file url is required")
			return
		}

		att := model.OttAttachments{
			Payload: &model.OttPayloadMedia{
				Url: fileUrl,
			},
			AttType: data.EventName,
		}
		attachments = append(attachments, &att)
	}
	content := data.Content
	// if eventName != "text" {
	// 	if file != nil {
	// 		content = file.Filename
	// 	}
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
		Text:          content,
	}

	log.Info("message to ott: ", ottMessage)

	resOtt, err := sendMessageToOTT(ottMessage, attachments)
	if err != nil {
		log.Error(err)
		return
	}

	// Store ES
	message = model.Message{
		TenantId:               conversation.TenantId,
		ParentExternalMsgId:    "",
		MessageId:              docId,
		MessageType:            conversation.ConversationType,
		ConversationId:         conversation.ConversationId,
		ExternalMsgId:          resOtt.Data.MsgId,
		EventName:              eventName,
		Direction:              variables.DIRECTION["send"],
		AppId:                  conversation.AppId,
		OaId:                   conversation.OaId,
		Avatar:                 conversation.Avatar,
		SupporterId:            authUser.UserId,
		SupporterName:          authUser.Fullname,
		SendTime:               time.Now(),
		SendTimestamp:          timestampTmp,
		Content:                content,
		Attachments:            attachments,
		ExternalConversationId: conversation.ExternalConversationId,
	}
	log.Info("message to es: ", message)

	// Should to queue
	if err = InsertMessage(ctx, conversation.TenantId, ES_INDEX_MESSAGE, message.AppId, docId, message); err != nil {
		log.Error(err)
		return
	}

	// TODO: update conversation => after refresh page, this conversation will appearance on top the list
	conversation.UpdatedAt = time.Now().Format(time.RFC3339)
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		log.Error(err)
		return
	}
	esDoc := map[string]any{}
	if err = json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return
	}

	conversationQueue := model.ConversationQueue{
		DocId:        conversation.ConversationId,
		Conversation: conversation,
	}
	if err = PublishPutConversationToChatQueue(ctx, conversationQueue); err != nil {
		log.Error(err)
		return
	}

	// TODO: send message to admin/manager/user
	// Exclude: if user send message, remove user, if admin send message, remove admin, ...
	var queueId string
	filter := model.ChatConnectionAppFilter{
		AppId:          message.AppId,
		OaId:           message.OaId,
		ConnectionType: conversation.ConversationType,
	}
	_, connection, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*connection) < 1 {
		err = errors.New("connection " + (*connection)[0].Id + " not found")
		log.Errorf(err.Error())
		return
	} else {
		message.TenantId = (*connection)[0].TenantId
	}

	// TODO: find connection_queue
	connectionQueueExist, err := repository.ConnectionQueueRepo.GetById(ctx, repository.DBConn, (*connection)[0].ConnectionQueueId)
	if err != nil {
		log.Error(err)
		return
	} else if connectionQueueExist == nil {
		err = errors.New("connection queue not found")
		log.Error(err)
		return
	}

	filterChatManageQueueUser := model.ChatManageQueueUserFilter{
		QueueId: connectionQueueExist.QueueId,
	}
	_, manageQueueUser, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*manageQueueUser) > 0 {
		queueId = (*manageQueueUser)[0].QueueId
	}

	if authUser.Level != "manager" {
		if len(queueId) > 0 {
			if err = SendEventToManage(ctx, authUser, message, queueId); err != nil {
				log.Error(err)
				return
			}
		} else {
			log.Errorf("queue %s not found in send event to manage", queueId)
		}
	}

	return
}

func (s *Message) GetMessages(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit, offset int) (total int, messages *[]model.Message, err error) {
	total, messages, err = repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX_MESSAGE, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Message) GetMessagesWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit int, scrollId string) (total int, messages []*model.Message, respScrollId string, err error) {
	total, messages, respScrollId, err = repository.MessageESRepo.SearchWithScroll(ctx, authUser.TenantId, ES_INDEX_MESSAGE, filter, limit, scrollId)
	if err != nil {
		log.Error(err)
		return
	}
	if messages == nil {
		messages = make([]*model.Message, 0)
	}
	return
}

func (s *Message) MarkReadMessages(ctx context.Context, authUser *model.AuthUser, data model.MessageMarkRead) (result model.ReadMessageResponse, err error) {
	conversation := model.Conversation{}
	conversationCache := cache.RCache.Get(CONVERSATION + "_" + data.ConversationId)
	if conversationCache != nil {
		if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
			log.Error(err)
			return result, err
		}
	} else {
		conversationFilter := model.ConversationFilter{
			TenantId:               authUser.TenantId,
			ExternalConversationId: []string{data.ConversationId},
		}
		total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, conversationFilter, 1, 0)
		if err != nil {
			log.Error(err)
			return result, err
		}
		if total > 0 {
			conversation := (*conversations)[0]
			if err := cache.RCache.Set(CONVERSATION+"_"+conversation.ConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
				log.Error(err)
				return result, err
			}
		}
	}

	if data.ReadAll {
		filter := model.MessageFilter{
			IsRead:         "deactive",
			ConversationId: data.ConversationId,
		}
		total, messages, err := repository.MessageESRepo.GetMessages(ctx, authUser.TenantId, ES_INDEX_MESSAGE, filter, -1, 0)
		if err != nil {
			log.Error(err)
			result.TotalSuccess -= 1
			result.TotalFail += 1
			return result, err
		}
		result.TotalSuccess = total
		if total > 0 {
			for _, item := range *messages {
				item.IsRead = "active"
				item.ReadBy = append(item.ReadBy, authUser.UserId)
				item.ReadTimestamp = time.Now().UnixMilli()
				item.UpdatedAt = time.Now()
				item.ReadTime = time.Now()

				// TODO: add queue to update
				if err := PublishPutMessageToChatQueue(ctx, item); err != nil {
					log.Error(err)
					result.TotalSuccess -= 1
					result.TotalFail += 1
					return result, err
				}
			}
		}
	} else {
		for _, item := range data.MessageIds {
			// Need tracking message still not read ?
			message, err := repository.MessageESRepo.GetMessageById(ctx, "", ES_INDEX_MESSAGE, item)
			if err != nil {
				log.Error(err)
				result.TotalSuccess -= 1
				result.TotalFail += 1
				result.ListFail[item] = err.Error()
				continue
			} else if len(message.MessageId) < 1 {
				result.TotalSuccess -= 1
				result.TotalFail += 1
				log.Errorf("message %s not found", item)
				result.ListFail[item] = "message " + item + " not found"
				continue
			}

			message.IsRead = "active"
			message.ReadBy = append(message.ReadBy, authUser.UserId)
			message.ReadTimestamp = time.Now().UnixMilli()
			message.UpdatedAt = time.Now()
			message.ReadTime = time.Now()

			if err := PublishPutMessageToChatQueue(ctx, *message); err != nil {
				log.Error(err)
				result.TotalSuccess -= 1
				result.TotalFail += 1
				result.ListFail[item] = err.Error()
			}
			result.ListSuccess[item] = "success"
		}
	}
	return
}

func (s *Message) ShareInfo(ctx context.Context, authUser *model.AuthUser, data model.ShareInfo) (result model.OttCodeChallenge, err error) {
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
		return
	}
	if resp.StatusCode() == 200 {
		if err := json.Unmarshal([]byte(resp.Body()), &result); err != nil {
			return result, err
		}
		return
	} else {
		err = errors.New(resp.String())
		return
	}
}

func (s *Message) GetMessageMediasWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.MessageFilter, limit int, scrollId string) (total int, messages []*model.MessageAttachmentsDetails, respScrollId string, err error) {
	total, messages, respScrollId, err = repository.MessageESRepo.GetMessageMediasWithScroll(ctx, authUser.TenantId, ES_INDEX_MESSAGE, filter, limit, scrollId)
	if err != nil {
		log.Error(err)
		return
	}
	if messages == nil {
		messages = make([]*model.MessageAttachmentsDetails, 0)
	}
	return
}

func (s *Message) PostTicketToMessage(ctx context.Context, authUser *model.AuthUser, data model.MessagePostTicket) (err error) {
	message, err := repository.MessageESRepo.GetMessageById(ctx, "", ES_INDEX_MESSAGE, data.MessageId)
	if err != nil {
		log.Error(err)
		return
	} else if len(message.MessageId) < 1 {
		log.Errorf("message %s not found", data.MessageId)
		return errors.New("message " + data.MessageId + " not found")
	}

	message.TicketId = data.TicketId

	if err = PublishPutMessageToChatQueue(ctx, *message); err != nil {
		return
	}
	return
}
