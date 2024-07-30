package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/internal/queue"
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
	OttMessage struct {
		NumConsumer int
		Users       map[string]chan []model.User
	}
)

var OttMessageService IOttMessage

func NewOttMessage() IOttMessage {
	s := &OttMessage{}
	s.InitQueueRequest()
	return s
}

/**
* PROBLEM: chat vao thi sao biet user thuoc db nao
* Khi chuyen qua fins thi lam sao biet setting nay cua db nao
* IMPROVE: add channel, mutex
 */
func (s *OttMessage) GetOttMessage(ctx context.Context, data model.OttMessage) (int, any) {
	// Check and cache the chat app configuration
	chatApp, err := CheckConfigAppCache(ctx, data.AppId)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}

	timestamp := time.Unix(0, data.Timestamp*int64(time.Millisecond))
	messageTmp := s.createMessage(data, timestamp)
	if slices.Contains[[]string](variables.EVENT_READ_MESSAGE, data.EventName) {
		messageTmp.ReadTime = timestamp
		messageTmp.ReadTimestamp = data.Timestamp
	}

	externalConversationId := GenerateConversationId(data.AppId, data.OaId, data.ExternalUserId)

	// Retrieve and cache chat app integrate systems
	chatAppIntegrateSystems, err := s.getChatAppIntegrateSystems(ctx, chatApp.Id, externalConversationId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// Process tenants
	tenants := s.proccessIntegrateSystems(chatAppIntegrateSystems)

	var isNew bool
	var conversation model.ConversationView
	var user model.User
	userChan := make(chan []model.User, len(tenants)) // Buffered channel to avoid blocking
	doneChan := make(chan bool, 1)
	errChan := make(chan error, 1)

	// Process users
	var wg sync.WaitGroup
	wg.Add(len(tenants))
	go func(timeStamp time.Time, userChan chan []model.User, doneChan chan<- bool) {
		defer close(userChan)
		for users := range userChan {
			if len(users) > 0 {
				for k, item := range users {
					user = item
					docId := uuid.NewString()
					connectionCache, err := GetConfigConnectionAppCache(ctx, tenants[k], data.AppId, data.OaId, data.MessageType)
					if err != nil {
						log.Error(err)
						return
					}
					message := s.createMessage(data, timestamp)
					if slices.Contains[[]string](variables.EVENT_READ_MESSAGE, data.EventName) {
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
									return
								}
								attachmentFile.Url = strings.ReplaceAll(attachmentFile.Url, "u0026", "&")
								attachmentDetail.Payload = &payload
							} else {
								if err := util.ParseAnyToAny(val.Payload, &payload); err != nil {
									log.Error(err)
									return
								}
								attachmentMedia.Url = strings.ReplaceAll(attachmentMedia.Url, "u0026", "&")
								attachmentDetail.Payload = &payload
							}
							message.Attachments = append(message.Attachments, &attachmentDetail)
							message.Content = val.Payload.Name
						}
					}

					if user.AuthUser != nil {
						data.TenantId = user.AuthUser.TenantId
						message.TenantId = user.AuthUser.TenantId
					} else {
						data.TenantId = connectionCache.TenantId
						message.TenantId = connectionCache.TenantId

						// TODO: find connection_queue
						connectionQueueFilter := model.ConnectionQueueFilter{
							TenantId:     connectionCache.TenantId,
							ConnectionId: connectionCache.Id,
						}
						_, connectionQueueExists, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, connectionQueueFilter, 1, 0)
						if err != nil {
							log.Error(err)
							return
						} else if len(*connectionQueueExists) < 1 {
							log.Errorf("connection queue " + connectionCache.Id + " not found")
							return
						}

						filterChatManageQueueUser := model.ChatManageQueueUserFilter{
							QueueId: (*connectionQueueExists)[0].QueueId,
						}
						_, manageQueueUser, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, 1, 0)
						if err != nil {
							log.Error(err)
							return
						}
						if len(*manageQueueUser) > 0 {
							user.QueueId = (*manageQueueUser)[0].QueueId
							user.ConnectionId = (*manageQueueUser)[0].ConnectionId
							user.IsOk = true
						}
					}

					if user.IsOk {
						conversationTmp, isNewTmp, errConv := s.UpSertConversation(ctx, user.ConnectionId, user.ConversationId, data)
						if errConv != nil {
							return
						}
						if err := util.ParseAnyToAny(conversationTmp, &conversation); err != nil {
							log.Error(err)
							return
						}
						isNew = isNewTmp

						if len(conversation.ConversationId) > 0 {
							message.IsRead = "deactive"
							if message.IsEcho {
								message.Direction = variables.DIRECTION["send"]
								message.Avatar = conversation.OaAvatar
								message.SupporterId = ""
								message.SupporterName = "Admin OA"
							}
						} else {
							log.Error("conversation " + conversation.ConversationId + " not found with app id " + data.AppId)
							return
						}
					} else if user.PreviousAssign != nil {
						// TODO: insert or update allocate user
						user.PreviousAssign.UserId = user.AuthUser.UserId
						user.PreviousAssign.MainAllocate = "active"
						user.PreviousAssign.UpdatedAt = time.Now()
						user.PreviousAssign.AllocatedTimestamp = time.Now().UnixMilli()
						if err := repository.AllocateUserRepo.Update(ctx, repository.DBConn, *user.PreviousAssign); err != nil {
							log.Error(err)
							return
						}
					}
					if len(conversation.ConversationId) < 1 {
						err = errors.New("conversation " + conversation.ConversationId + " not found")
						log.Error(err)
						return
					}

					// Insert message
					if len(conversation.ConversationId) > 0 {
						// Parsing conversation_id
						message.ConversationId = conversation.ConversationId
						message.MessageId = docId
						if err = InsertMessage(ctx, data.TenantId, ES_INDEX_MESSAGE, data.AppId, docId, message); err != nil {
							log.Error(err)
							return
						}
					}

					wg.Done()

					subscribers := []*Subscriber{}
					subscriberAdmins := []string{}
					subscriberManagers := []string{}
					for s := range WsSubscribers.Subscribers {
						if (user.AuthUser != nil && s.TenantId == user.AuthUser.TenantId) || (conversation.TenantId == s.TenantId) {
							subscribers = append(subscribers, s)
							if s.Level == "admin" {
								subscriberAdmins = append(subscriberAdmins, s.Id)
							}
							if s.Level == "manager" {
								subscriberManagers = append(subscriberManagers, s.Id)
							}
						}
					}

					if user.AuthUser != nil {
						go handlePublishEvent(false, &user, nil, subscriberAdmins, subscribers, conversation, message, isNew)
					}

					// TODO: publish message to manager
					manageQueueUser, err := GetManageQueueUser(ctx, user.QueueId)
					if err != nil {
						log.Error(err)
						return
					} else if manageQueueUser == nil {
						log.Error("queue " + user.QueueId + " not found")
						return
					}
					if len(manageQueueUser.ConnectionId) < 1 {
						manageQueueUser.ConnectionId = connectionCache.Id
					}

					// TODO: publish message to manager
					isExist := BinarySearchSlice(manageQueueUser.UserId, subscriberManagers)
					if isExist {
						if (user.AuthUser != nil && user.AuthUser.UserId != manageQueueUser.UserId) || (user.AuthUser == nil && len(manageQueueUser.UserId) > 0) {
							go handlePublishEvent(false, &user, manageQueueUser, subscriberAdmins, subscribers, conversation, message, isNew)
						}
					}

					// TODO: publish to admin
					if ENABLE_PUBLISH_ADMIN {
						go handlePublishEvent(true, &user, manageQueueUser, subscriberAdmins, subscribers, conversation, message, isNew)
					}

					if ENABLE_CHAT_AUTO_SCRIPT_REPLY {
						if err = ExecutePlannedAutoScript(ctx, user, message, &conversation); err != nil {
							log.Error(err)
							return
						}
					}
				}
			} else {
				wg.Done()
			}
		}
	}(timestamp, userChan, doneChan)

	// TODO: check queue setting
	go s.CheckChatSetting(ctx, externalConversationId, messageTmp, *chatApp, userChan, errChan, tenants)

	// Wait for all tenants to be processed
	go func() {
		wg.Wait() // Wait for all goroutines to finish processing
		doneChan <- true
	}()

	select {
	case <-doneChan:
		log.Debug("receive ott message done")
		return response.OKResponse()
	case err = <-errChan:
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	case <-ctx.Done():
		log.Debug("context timeout")
		return response.ServiceUnavailableMsg(errors.New("context timeout"))
	}
}

func (s *OttMessage) createMessage(data model.OttMessage, timestamp time.Time) model.Message {
	return model.Message{
		MessageId:           uuid.NewString(),
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
		IsEcho:              data.IsEcho,
	}
}

func (s *OttMessage) getChatAppIntegrateSystems(ctx context.Context, chatAppId, externalConversationId string) (result []model.ChatAppIntegrateSystem, err error) {
	chatAppIntegrateSystemCache := cache.RCache.Get(CHAT_APP_INTEGRATE_SYSTEM + "_" + externalConversationId)
	if chatAppIntegrateSystemCache != nil {
		if err = json.Unmarshal([]byte(chatAppIntegrateSystemCache.(string)), &result); err != nil {
			log.Error(err)
			return
		}
	} else {
		filterChatAppIntegrateSystem := model.ChatAppIntegrateSystemFilter{
			ChatAppId: chatAppId,
		}
		_, tmp, errTmp := repository.ChatAppIntegrateSystemRepo.GetChatAppIntegrateSystems(ctx, repository.DBConn, filterChatAppIntegrateSystem, -1, 0)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		}
		result = *tmp

		if err = cache.RCache.Set(CHAT_APP_INTEGRATE_SYSTEM+"_"+externalConversationId, result, CHAT_APP_INTEGRATE_SYSTEM_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}
	return
}

func (s *OttMessage) proccessIntegrateSystems(chatAppIntegrateSystems []model.ChatAppIntegrateSystem) (tenants []string) {
	chatIntegrateSystems := []model.ChatIntegrateSystem{}
	if len(chatAppIntegrateSystems) > 0 {
		for _, integrateSystem := range chatAppIntegrateSystems {
			if integrateSystem.ChatIntegrateSystem[0].Status {
				chatIntegrateSystems = append(chatIntegrateSystems, *integrateSystem.ChatIntegrateSystem[0])
			}
		}
	}

	if len(chatIntegrateSystems) > 0 {
		for _, integrateSystem := range chatIntegrateSystems {
			if len(integrateSystem.TenantDefaultId) > 0 {
				tenants = append(tenants, integrateSystem.TenantDefaultId)
			}
		}
	}
	return
}

func handlePublishEvent(isPublishToAdmin bool, user *model.User, manageQueueUser *model.ChatManageQueueUser, subscriberAdmins []string, subscribers []*Subscriber, conversation model.ConversationView, message model.Message, isNew bool) {
	var userId string
	var userIdRemove string
	if manageQueueUser != nil {
		userId = manageQueueUser.UserId
		userIdRemove = manageQueueUser.UserId
	} else {
		userId = user.AuthUser.UserId
		userIdRemove = user.UserIdRemove
	}

	if isPublishToAdmin {
		if user.IsReassignSame {
			PublishConversationToManyUser(variables.EVENT_CHAT["conversation_reopen"], subscriberAdmins, true, &conversation)
			PublishMessageToManyUser(variables.EVENT_CHAT["message_created"], subscriberAdmins, &message)
		} else if user.IsReassignNew {
			PublishConversationToManyUser(variables.EVENT_CHAT["conversation_removed"], subscriberAdmins, true, &conversation)
			PublishConversationToManyUser(variables.EVENT_CHAT["conversation_created"], subscriberAdmins, isNew, &conversation)
			PublishMessageToManyUser(variables.EVENT_CHAT["message_created"], subscriberAdmins, &message)
		} else {
			PublishConversationToManyUser(variables.EVENT_CHAT["conversation_created"], subscriberAdmins, isNew, &conversation)
			PublishMessageToManyUser(variables.EVENT_CHAT["message_created"], subscriberAdmins, &message)
		}
	} else {
		if user.IsReassignSame {
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_reopen"], userId, subscribers, true, &conversation)
			PublishMessageToOneUser(variables.EVENT_CHAT["message_created"], userId, subscribers, &message)
		} else if user.IsReassignNew {
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_removed"], userIdRemove, subscribers, true, &conversation)
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_created"], userId, subscribers, isNew, &conversation)
			PublishMessageToOneUser(variables.EVENT_CHAT["message_created"], userId, subscribers, &message)
		} else {
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_created"], userId, subscribers, isNew, &conversation)
			PublishMessageToOneUser(variables.EVENT_CHAT["message_created"], userId, subscribers, &message)
		}
	}
}

func (s *OttMessage) GetCodeChallenge(ctx context.Context, authUser *model.AuthUser, appId string) (int, any) {
	url := OTT_URL + "/ott/" + OTT_VERSION + "/zalo/code-challenge/" + appId
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
	event := model.Event{
		EventName: "share_info",
		EventData: &model.EventData{
			ShareInfo: &data,
		},
	}
	if err := PublishMessageToOne(authUser.UserId, event); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OKResponse()
}

func (s *OttMessage) InitQueueRequest() {
	isExistConversationQueue := queue.RMQ.Server.IsHasQueue(BSS_CHAT_CONVERSATION_QUEUE_NAME)
	if isExistConversationQueue {
		queue.RMQ.Server.RemoveQueue(BSS_CHAT_CONVERSATION_QUEUE_NAME)
	}

	isExistMessageQueue := queue.RMQ.Server.IsHasQueue(BSS_CHAT_MESSAGE_QUEUE_NAME)
	if isExistMessageQueue {
		queue.RMQ.Server.RemoveQueue(BSS_CHAT_MESSAGE_QUEUE_NAME)
	}

	if s.NumConsumer < 2 {
		s.NumConsumer = 2
	}

	if err := queue.RMQ.Server.AddQueue(BSS_CHAT_CONVERSATION_QUEUE_NAME, s.handleEsConversationQueue, 1); err != nil {
		log.Error(err)
	}

	if err := queue.RMQ.Server.AddQueue(BSS_CHAT_MESSAGE_QUEUE_NAME, s.handleEsMessageQueue, 1); err != nil {
		log.Error(err)
	}
}

func (s *OttMessage) handleEsConversationQueue(d rmq.Delivery) {
	payload := model.ConversationQueue{}
	if err := util.ParseStringToAny(d.Payload(), &payload); err != nil {
		log.Error(err)
		d.Reject()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	done := make(chan struct{})

	go func() {
		defer func() {
			close(done)
		}()
		tmpBytes, err := json.Marshal(payload.Conversation)
		if err != nil {
			log.Error(err)
			return
		}
		esDoc := map[string]any{}
		if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
			log.Error(err)
			return
		}
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, payload.Conversation.AppId, payload.DocId, esDoc); err != nil {
			log.Error(err)
			return
		}
		if err := cache.RCache.Del([]string{CONVERSATION + "_" + payload.Conversation.ConversationId}); err != nil {
			log.Error(err)
			return
		}
	}()

	select {
	case <-done:
		// job have been success
		d.Ack()
	case <-ctx.Done():
		// exceeded timeout
		log.Errorf("handleEsConversationQueue exceeded timeout, msg=%s", d.Payload())
		d.Reject()
	}
}

func (s *OttMessage) handleEsMessageQueue(d rmq.Delivery) {
	payload := model.Message{}
	if err := util.ParseStringToAny(d.Payload(), &payload); err != nil {
		log.Error(err)
		d.Reject()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	done := make(chan struct{})

	go func() {
		defer func() {
			close(done)
		}()
		tmpBytes, err := json.Marshal(payload)
		if err != nil {
			log.Error(err)
			return
		}
		esDoc := map[string]any{}
		if err = json.Unmarshal(tmpBytes, &esDoc); err != nil {
			log.Error(err)
			return
		}
		if err = repository.ESRepo.UpdateDocById(ctx, ES_INDEX_MESSAGE, payload.AppId, payload.MessageId, esDoc); err != nil {
			log.Error(err)
			return
		}
	}()

	select {
	case <-done:
		// job have been success
		d.Ack()
	case <-ctx.Done():
		// exceeded timeout
		log.Errorf("handleEsMessageQueue exceeded timeout, msg=%s", d.Payload())
		d.Reject()
	}
}
