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
		User        map[string]chan model.User
		mux         *sync.RWMutex
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
	chatApp, err := CheckConfigAppCache(ctx, data.AppId)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	var connectionCache model.ChatConnectionApp
	connectionCache, err = GetConfigConnectionAppCache(ctx, data.AppId, data.OaId, data.MessageType)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
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
		IsEcho:              data.IsEcho,
	}
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
			message.Content = val.Payload.Name
		}
	}

	var isNew bool
	var conversation model.ConversationView
	var user model.User
	userChan := make(chan []model.User)
	errChan := make(chan error)
	done := make(chan bool)
	// defer close(errChan)
	go func(userChan <-chan []model.User) {
		for _, item := range <-userChan {
			user = item
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
					errChan <- err
					return
				} else if len(*connectionQueueExists) < 1 {
					log.Errorf("connection queue " + connectionCache.Id + " not found")
					errChan <- errors.New("connection queue " + connectionCache.Id + " not found")
					return
				}

				filterChatManageQueueUser := model.ChatManageQueueUserFilter{
					QueueId: (*connectionQueueExists)[0].QueueId,
				}
				_, manageQueueUser, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, 1, 0)
				if err != nil {
					log.Error(err)
					errChan <- err
					return
				}
				if len(*manageQueueUser) > 0 {
					user.QueueId = (*manageQueueUser)[0].QueueId
					user.ConnectionId = (*manageQueueUser)[0].ConnectionId
					user.IsOk = true
				}
			}

			if user.IsOk {
				conversationTmp, isNewTmp, errConv := s.UpSertConversation(ctx, user.ConnectionId, data)
				if errConv != nil {
					log.Error(errConv)
					errChan <- errConv
					return
				}
				if err := util.ParseAnyToAny(conversationTmp, &conversation); err != nil {
					log.Error(err)
					errChan <- err
					return
				}
				isNew = isNewTmp

				if len(conversation.ConversationId) > 0 {
					message.ConversationId = conversation.ConversationId
					message.IsRead = "deactive"
					if message.IsEcho {
						message.Direction = variables.DIRECTION["send"]
						message.Avatar = conversation.OaAvatar
						message.SupporterId = ""
						message.SupporterName = "Admin OA"
					}
					if errMsg := InsertES(ctx, data.TenantId, ES_INDEX, conversation.AppId, docId, message); errMsg != nil {
						log.Error(errMsg)
						errChan <- errMsg
						return
					}
				} else {
					log.Error("conversation " + conversation.ConversationId + " not found")
					errChan <- errors.New("conversation " + conversation.ConversationId + " not found")
					return
				}
			} else if user.PreviousAssign != nil {
				// TODO: insert or update allocate user
				user.PreviousAssign.UserId = user.AuthUser.UserId
				user.PreviousAssign.MainAllocate = "active"
				user.PreviousAssign.UpdatedAt = time.Now()
				user.PreviousAssign.AllocatedTimestamp = time.Now().UnixMilli()
				if err := repository.UserAllocateRepo.Update(ctx, repository.DBConn, *user.PreviousAssign); err != nil {
					log.Error(err)
					errChan <- err
					return
				}
			} else if len(user.QueueId) < 1 && err != nil {
				log.Error(err)
				errChan <- err
				return
			}
			if len(conversation.ConversationId) < 1 {
				log.Error("conversation " + conversation.ConversationId + " not found")
				errChan <- errors.New("conversation " + conversation.ConversationId + " not found")
				return
			}

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
			if len(user.QueueId) > 0 {
				manageQueueUser, err := GetManageQueueUser(ctx, user.QueueId)
				if err != nil {
					log.Error(err)
					errChan <- err
					return
				} else if manageQueueUser == nil {
					log.Error("queue " + user.QueueId + " not found")
					errChan <- errors.New("queue " + user.QueueId + " not found")
					return
				}
				if len(manageQueueUser.ConnectionId) < 1 {
					manageQueueUser.ConnectionId = connectionCache.Id
				}

				// TODO: if user not found then assign conversation for manager
				filter := model.UserAllocateFilter{
					AppId:          conversation.AppId,
					ConversationId: conversation.ConversationId,
					MainAllocate:   "active",
				}
				_, userAllocates, err := repository.UserAllocateRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
				if err != nil {
					log.Error(err)
					errChan <- err
					return
				}
				if len(*userAllocates) < 1 {
					if len(conversation.ConversationId) > 0 {
						_, conversationDeactiveExist, err := repository.UserAllocateRepo.GetAllocateUsers(ctx, repository.DBConn, model.UserAllocateFilter{
							AppId:          conversation.AppId,
							ConversationId: conversation.ConversationId,
						}, -1, 0)
						if err != nil {
							log.Error(err)
							errChan <- err
							return
						}
						if len(*conversationDeactiveExist) > 0 {
							if err := repository.UserAllocateRepo.DeleteAllocateUsers(ctx, repository.DBConn, *conversationDeactiveExist); err != nil {
								log.Error(err)
								errChan <- err
								return
							}
						}
					}
					userAllocation := model.AllocateUser{
						Base:               model.InitBase(),
						TenantId:           conversation.TenantId,
						ConversationId:     conversation.ConversationId,
						AppId:              message.AppId,
						OaId:               message.OaId,
						UserId:             manageQueueUser.UserId,
						ConnectionQueueId:  connectionCache.ConnectionQueueId,
						QueueId:            manageQueueUser.QueueId,
						AllocatedTimestamp: time.Now().UnixMilli(),
						MainAllocate:       "active",
						ConnectionId:       manageQueueUser.ConnectionId,
					}
					if err := repository.UserAllocateRepo.Insert(ctx, repository.DBConn, userAllocation); err != nil {
						log.Error(err)
						errChan <- err
						return
					}
					if err := cache.RCache.Set(USER_ALLOCATE+"_"+conversation.ConversationId, userAllocation, USER_ALLOCATE_EXPIRE); err != nil {
						log.Error(err)
						errChan <- err
						return
					}
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
			}
			done <- true

			if ENABLE_CHAT_AUTO_SCRIPT_REPLY {
				if err = ExecutePlannedAutoScript(ctx, user, message, &conversation); err != nil {
					log.Error(err)
					errChan <- err
					return
				}
			}
		}
	}(userChan)

	// TODO: check queue setting
	go s.CheckChatSetting(ctx, message, *chatApp, userChan, errChan)

	select {
	case <-done:
		log.Debug("receive ott message done")
		return response.OKResponse()
	case err = <-errChan:
		return response.ServiceUnavailableMsg(err.Error())
	case <-ctx.Done():
		log.Debug("context timeout")
		return response.ServiceUnavailableMsg(errors.New("context timeout"))
	}
}

func (s *OttMessage) RegisterUserChannel(key string, resultChan chan model.User) (err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.User[key]; !ok {
		return errors.New("channel " + key + " not found")
	}

	s.User[key] = resultChan
	return
}

func (s *OttMessage) UnregisterUserChannel(key string) {
	s.mux.Lock()
	defer s.mux.Unlock()

	close(s.User[key])
	delete(s.User, key)
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
	isExist := queue.RMQ.Server.IsHasQueue(BSS_CHAT_QUEUE_NAME)
	if isExist {
		queue.RMQ.Server.RemoveQueue(BSS_CHAT_QUEUE_NAME)
	}

	if s.NumConsumer < 1 {
		s.NumConsumer = 1
	}

	if err := queue.RMQ.Server.AddQueue(BSS_CHAT_QUEUE_NAME, s.handleEsConversationQueue, 1); err != nil {
		log.Error(err)
	}
}

func (s *OttMessage) handleEsConversationQueue(d rmq.Delivery) {
	payload := model.Conversation{}
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
		if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
			log.Error(err)
			return
		}
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, payload.AppId, payload.ConversationId, esDoc); err != nil {
			log.Error(err)
			return
		}
		if err := cache.RCache.Del([]string{CONVERSATION + "_" + payload.ConversationId}); err != nil {
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
