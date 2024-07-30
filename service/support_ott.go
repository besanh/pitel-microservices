package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

/**
* TODO: improve workflow in new version
* We need to write new workflow that can support multi tenant with one or multi chat app
* We should use transaction that can rollback between each workflow, posgre and elasticsearch
* Always remember that if no user online, we will assign conversaton for manager of queue
* If customer chatting with us, we need to check conversation exist in database and elasticsearch to determine it is new conversation or not
* If it's new conversation, we should create new conversation in database and elasticsearch
* If it's not new conversation, we should update conversation in database and elasticsearch
* If user make done conversation then customer chat with us again, we should check rule allocate again to determine assign for old user or not
* If customer chatting with user A then admin assign for user B, we will transfer for that user
* Finally we should prepare data ready to push to wss, include conversation, message, allocated user
 */

/**
* Chia patch cho phan loop
 */
func (s *OttMessage) CheckChatSetting(ctx context.Context, externalConversationId string, message model.Message, chatApp model.ChatApp, userChan chan<- []model.User, errChan chan<- error, tenants []string) {
	for _, item := range tenants {
		if len(message.ExternalMsgId) > 0 {
			isExistMessage, err := s.checkMessageEcho(ctx, item, message)
			if err != nil || isExistMessage {
				userChan <- []model.User{}
				continue
			}
		}
		var authInfo model.AuthUser
		var user model.User
		allocateUser := &model.AllocateUser{}
		userUuidExcludes := []string{}
		conversationId := uuid.NewString()

		allocateUsersCache, err := cache.RCache.HGetAll(USER_ALLOCATE + "_" + item + "_" + externalConversationId)
		if err != nil {
			log.Error(err)
			userChan <- []model.User{}
			continue
		}
		if len(allocateUsersCache) > 0 {
			for _, AllocateUserCache := range allocateUsersCache {
				if err = json.Unmarshal([]byte(AllocateUserCache), allocateUser); err != nil {
					log.Error(err)
					userChan <- []model.User{}
					continue
				}

				// Exist user allocate
				if allocateUser.TenantId == item && allocateUser.ConversationId == externalConversationId {
					if allocateUser.MainAllocate == "active" {
						authInfo.TenantId = allocateUser.TenantId
						authInfo.UserId = allocateUser.UserId
						user.AuthUser = &authInfo
						user.IsOk = true
						user.ConnectionId = allocateUser.ConnectionId
						user.QueueId = allocateUser.QueueId
						user.ConnectionQueueId = allocateUser.ConnectionQueueId
						user.ConversationId = allocateUser.ConversationId

						log.Infof("conversation %s allocated to username %s, id: %s, domain: %s", externalConversationId, user.AuthUser.Fullname, user.AuthUser.UserId, user.AuthUser.TenantId)
						userChan <- []model.User{user}
					} else {
						authInfo.TenantId = allocateUser.TenantId
						authInfo.UserId = allocateUser.UserId
						user.AuthUser = &authInfo
						user.ConnectionId = allocateUser.ConnectionId
						user.QueueId = allocateUser.QueueId
						user.ConnectionQueueId = allocateUser.ConnectionQueueId

						userUuidExcludes = append(userUuidExcludes, allocateUser.UserId)
						user, err := s.CheckAllSetting(ctx, item, conversationId, GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), message, true, allocateUser, chatApp, userUuidExcludes)
						if err != nil {
							userChan <- []model.User{}
							continue
						}

						if user.AuthUser.UserId == allocateUser.UserId {
							user.IsReassignSame = true
						} else {
							user.IsReassignNew = true
							user.UserIdRemove = allocateUser.UserId
						}

						user.IsOk = true
						user.ConversationId = allocateUser.ConversationId
						log.Infof("conversation %s allocated to username %s, id: %s, domain: %s", externalConversationId, user.AuthUser.Fullname, user.AuthUser.UserId, user.AuthUser.TenantId)
						userChan <- []model.User{user}
					}
				} else {
					user, err := s.CheckAllSetting(ctx, item, conversationId, externalConversationId, message, false, nil, chatApp, userUuidExcludes)
					if err != nil {
						log.Error(err)
						userChan <- []model.User{}
						continue
					}

					filter := model.AllocateUserFilter{
						TenantId:               item,
						ExternalConversationId: externalConversationId,
						MainAllocate:           "deactive",
					}
					_, userAllocations, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
					if err != nil {
						log.Error(err)
						userChan <- []model.User{}
						continue
					}
					if len(*userAllocations) > 0 {
						if user.AuthUser.UserId == (*userAllocations)[0].UserId {
							user.IsReassignSame = true
						} else {
							user.IsReassignNew = true
							user.UserIdRemove = (*userAllocations)[0].UserId
						}
						user.ConversationId = (*userAllocations)[0].ConversationId
					} else {
						log.Error(errors.New("conversation " + externalConversationId + " not found"))
						continue
					}

					log.Infof("conversation %s allocated to username %s, id: %s, domain: %s", externalConversationId, user.AuthUser.Fullname, user.AuthUser.UserId, user.AuthUser.TenantId)
					userChan <- []model.User{user}

					// Set to cache
					jsonByte, err := json.Marshal(&allocateUser)
					if err != nil {
						log.Error(err)
						userChan <- []model.User{}
						continue
					}
					if err = cache.RCache.HSetRaw(ctx, USER_ALLOCATE+"_"+item+"_"+externalConversationId, externalConversationId, string(jsonByte)); err != nil {
						log.Error(err)
						userChan <- []model.User{}
						continue
					}
				}
			}
		} else {
			user, err := s.CheckAllSetting(ctx, item, conversationId, externalConversationId, message, false, nil, chatApp, userUuidExcludes)
			if err != nil {
				userChan <- []model.User{}
				continue
			}

			filter := model.AllocateUserFilter{
				TenantId:               user.AuthUser.TenantId,
				ExternalConversationId: externalConversationId,
				MainAllocate:           "deactive",
			}
			_, userAllocations, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
			if err != nil {
				log.Error(err)
				userChan <- []model.User{}
				continue
			}
			if len(*userAllocations) > 0 && user.AuthUser != nil {
				if user.AuthUser.UserId == (*userAllocations)[0].UserId {
					user.IsReassignSame = true
				} else {
					user.IsReassignNew = true
					user.UserIdRemove = (*userAllocations)[0].UserId
				}
				user.ConversationId = (*userAllocations)[0].ConversationId
			}

			log.Infof("conversation %s allocated to username %s, id: %s, domain: %s", externalConversationId, user.AuthUser.Fullname, user.AuthUser.UserId, user.AuthUser.TenantId)
			userChan <- []model.User{user}
		}
	}
}

// Check message echo to prevent duplicate message
func (s *OttMessage) checkMessageEcho(ctx context.Context, tenant string, message model.Message) (isExist bool, err error) {
	filter := model.MessageFilter{
		TenantId:          tenant,
		ExternalMessageId: message.ExternalMsgId,
	}
	_, messages, err := repository.MessageESRepo.GetMessages(ctx, tenant, ES_INDEX_MESSAGE, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*messages) > 0 {
		isExist = true
	}
	return
}

/**
* Check all setting to allocate conversation to user
 */
func (s *OttMessage) CheckAllSetting(ctx context.Context, tenantId, conversationId, externalConversationId string, message model.Message, isConversationExist bool, currentAllocateUser *model.AllocateUser, chatApp model.ChatApp, userUuidExlucdes []string) (user model.User, err error) {
	var authInfo model.AuthUser
	chatConnectionApp := model.ChatConnectionApp{}
	chatConnectionCache := cache.RCache.Get(CHAT_CONNECTION + "_" + tenantId + "_" + message.AppId + "_" + message.OaId)
	if chatConnectionCache != nil {
		if err = json.Unmarshal([]byte(chatConnectionCache.(string)), &chatConnectionApp); err != nil {
			log.Error(err)
			return
		}
	} else {
		connectionFilter := model.ChatConnectionAppFilter{
			TenantId:       tenantId,
			ConnectionType: message.MessageType,
			OaId:           message.OaId,
			AppId:          message.AppId,
			Status:         "active",
		}
		_, connectionApps, errTmp := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, connectionFilter, 1, 0)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		}
		if len(*connectionApps) > 0 {
			chatConnectionApp = (*connectionApps)[0]
			if err = cache.RCache.Set(CHAT_CONNECTION+"_"+tenantId+"_"+message.AppId+"_"+message.OaId, chatConnectionApp, CHAT_CONNECTION_EXPIRE); err != nil {
				log.Error(err)
				return
			}
		} else {
			err = errors.New("connect for conversation " + externalConversationId + " not found in tenant " + tenantId)
			log.Error(err)
			return
		}
	}

	connectionQueue, errTmp := repository.ConnectionQueueRepo.GetById(ctx, repository.DBConn, chatConnectionApp.ConnectionQueueId)
	if errTmp != nil {
		err = errTmp
		log.Error(err)
		return
	} else if connectionQueue == nil {
		err = errors.New("connection queue " + chatConnectionApp.ConnectionQueueId + " not found")
		log.Error(err)
		return
	}

	allocateUserFilter := model.AllocateUserFilter{
		TenantId:               tenantId,
		ConversationId:         conversationId,
		ExternalConversationId: externalConversationId,
		QueueId:                []string{connectionQueue.QueueId},
		MainAllocate:           "active",
	}
	_, allocateUsers, errTmp := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, allocateUserFilter, -1, 0)
	if errTmp != nil {
		err = errTmp
		log.Error(err)
		return
	}
	if len(*allocateUsers) > 0 {
		authInfo.TenantId = (*allocateUsers)[0].TenantId
		authInfo.UserId = (*allocateUsers)[0].UserId

		if err = cache.RCache.Set(USER_ALLOCATE+"_"+externalConversationId, (*allocateUsers)[0], USER_ALLOCATE_EXPIRE); err != nil {
			log.Error(err)
			return
		}
		user.IsOk = true
		user.AuthUser = &authInfo
		user.ConnectionId = (*allocateUsers)[0].ConnectionId
		user.QueueId = (*allocateUsers)[0].QueueId
		user.ConnectionQueueId = (*allocateUsers)[0].ConnectionQueueId
		user.ConversationId = (*allocateUsers)[0].ConversationId
	} else {
		// Connection prevent duplicate
		// Meaning: 1 connection with page A in 1 app => only recieve one queue
		var queue model.ChatQueue
		queueCache := cache.RCache.Get(CHAT_QUEUE + "_" + connectionQueue.QueueId)
		if queueCache != nil {
			if err = json.Unmarshal([]byte(queueCache.(string)), &queue); err != nil {
				log.Error(err)
				return
			}
		} else {
			queueTmp, errTmp := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, connectionQueue.QueueId)
			if errTmp != nil {
				err = errTmp
				log.Error(err)
				return
			} else if queueTmp == nil {
				err = errors.New("queue " + connectionQueue.QueueId + " not found")
				log.Error(err)
				return
			}
			queue = *queueTmp
		}

		// Get routing from cache
		chatRouting, errTmp := CacheChatRouting(ctx, queue.ChatRoutingId)
		if errTmp != nil {
			err = errTmp
			return
		}

		queueUserFilter := model.ChatQueueUserFilter{
			TenantId: tenantId,
			QueueId:  []string{queue.Id},
		}
		_, chatQueueUsers, errTmp := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, queueUserFilter, -1, 0)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		}
		chatManagerQueueFilter := model.ChatManageQueueUserFilter{
			TenantId:     tenantId,
			QueueId:      queue.Id,
			ConnectionId: connectionQueue.ConnectionId,
		}
		_, chatManangers, errTmp := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, chatManagerQueueFilter, -1, 0)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if len(*chatManangers) < 1 {
			err = errors.New("manage queue not found with queue id " + queue.Id + " and connection id " + connectionQueue.ConnectionId)
			log.Error(err)
			return
		}
		if len(*chatQueueUsers) > 0 {
			chatSetting := model.ChatSetting{
				ConnectionApp:       chatConnectionApp,
				ConnectionQueue:     *connectionQueue,
				QueueUser:           *chatQueueUsers,
				RoutingAlias:        chatRouting.RoutingAlias,
				Message:             message,
				ConnectionQueueUser: *chatQueueUsers,
				ManagerQueueUser:    (*chatManangers)[0],
			}

			userTmp, nextUserId, errTmp := s.GetAllocateUser(ctx, tenantId, conversationId, externalConversationId, chatSetting, isConversationExist, currentAllocateUser, userUuidExlucdes)
			if errTmp != nil {
				user.ConnectionId = connectionQueue.ConnectionId
				user.ConnectionQueueId = connectionQueue.Id
				err = errTmp
				return
			}
			if err = util.ParseAnyToAny(userTmp, &user); err != nil {
				log.Error(err)
				return
			}

			// Update main_allocate from deactive to active if customer chat again and assign to the same user deactive
			_, allocateUsers, errTmp := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, model.AllocateUserFilter{
				TenantId:               tenantId,
				ExternalConversationId: externalConversationId,
				MainAllocate:           "deactive",
			}, 1, 0)
			if errTmp != nil {
				err = errTmp
				log.Error(err)
				return
			}
			if len(*allocateUsers) > 0 {
				if (*allocateUsers)[0].UserId != nextUserId && len(nextUserId) > 0 {
					(*allocateUsers)[0].UserId = nextUserId
				}
				(*allocateUsers)[0].MainAllocate = "active"
				if err = repository.AllocateUserRepo.Update(ctx, repository.DBConn, (*allocateUsers)[0]); err != nil {
					log.Error(err)
					return
				}
			}

			user.QueueId = userTmp.QueueId
			user.ConnectionId = connectionQueue.ConnectionId
			user.ConnectionQueueId = connectionQueue.Id
			user.ConversationId = conversationId
		} else {
			err = errors.New("queue user not found")
			log.Error(err)
		}
	}
	return
}

/**
* Get user
* if isConversationExist = true,  it means conversation is exist, and we can get user from user_allocate
* if isConversationExist = false, it means conversation is not exist, we need to get user from chat_setting
 */
func (s *OttMessage) GetAllocateUser(ctx context.Context, tenantId, conversationId, externalConversationId string, chatSetting model.ChatSetting, isConversationExist bool, currentAllocateUser *model.AllocateUser, userUuidExclude []string) (user model.User, nextUserId string, err error) {
	allocateUser := model.AllocateUser{}
	var authInfo model.AuthUser
	var userLives []Subscriber
	var isUserAccept bool

	if strings.ToLower(chatSetting.RoutingAlias) == "random" {
		for s := range WsSubscribers.Subscribers {
			if (s.Level == "user" || s.Level == "agent") && s.TenantId == tenantId {
				userLives = append(userLives, *s)
			}
		}

		// Pick random
		if len(userLives) > 0 {
			rand.NewSource(time.Now().UnixNano())
			randomIndex := rand.Intn(len(userLives))
			tmp := userLives[randomIndex]

			// TODO: check user exist in queue
			if len(chatSetting.ConnectionQueueUser) > 0 {
				for _, item := range chatSetting.ConnectionQueueUser {
					if item.UserId == tmp.UserId && item.TenantId == tenantId {
						allocateUser.TenantId = tmp.TenantId
						allocateUser.UserId = tmp.UserId
						allocateUser.Username = tmp.Username

						authInfo.TenantId = allocateUser.TenantId
						authInfo.UserId = allocateUser.UserId
						isUserAccept = true
						break
					}
				}
			}
		}
	} else if strings.ToLower(chatSetting.RoutingAlias) == "round_robin_online" {
		userTmp, errTmp := RoundRobinUserOnline(ctx, tenantId, GenerateConversationId(chatSetting.Message.AppId, chatSetting.Message.OaId, chatSetting.Message.ExternalUserId), &chatSetting.QueueUser)
		if errTmp != nil {
			err = errTmp
			return
		} else if userTmp != nil {
			// TODO: check user exist in queue
			if len(chatSetting.ConnectionQueueUser) > 0 {
				for _, item := range chatSetting.ConnectionQueueUser {
					if item.UserId == userTmp.UserId && item.TenantId == tenantId {
						allocateUser.TenantId = userTmp.TenantId
						allocateUser.UserId = userTmp.UserId
						allocateUser.Username = userTmp.Username

						authInfo.TenantId = userTmp.TenantId
						authInfo.UserId = userTmp.UserId
						authInfo.Username = userTmp.Username
						authInfo.Level = userTmp.Level
						isUserAccept = true
						break
					}
				}
			}
		}
	}

	// If no user online, assign for manager
	if !isUserAccept {
		allocateUser.UserId = chatSetting.ManagerQueueUser.UserId
		authInfo.UserId = chatSetting.ManagerQueueUser.UserId
		authInfo.TenantId = chatSetting.ManagerQueueUser.TenantId
	}

	if isConversationExist {
		currentAllocateUser.ConversationId = conversationId
		currentAllocateUser.UserId = allocateUser.UserId
		currentAllocateUser.MainAllocate = "active"
		currentAllocateUser.AllocatedTimestamp = time.Now().UnixMilli()
		currentAllocateUser.UpdatedAt = time.Now()
		if err = s.UpdateConversationById(ctx, tenantId, *currentAllocateUser, chatSetting.Message); err != nil {
			log.Error(err)
			return
		}

		authInfo.TenantId = tenantId
		authInfo.UserId = allocateUser.UserId
		authInfo.Username = allocateUser.Username
	} else {
		if len(conversationId) > 0 {
			total, conversationDeactiveExist, errConv := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, model.AllocateUserFilter{
				AppId:          chatSetting.Message.AppId,
				ConversationId: conversationId,
			}, -1, 0)
			if errConv != nil {
				err = errConv
				log.Error(err)
				return
			}

			if total > 0 {
				if err = repository.AllocateUserRepo.DeleteAllocateUsers(ctx, repository.DBConn, *conversationDeactiveExist); err != nil {
					log.Error(err)
					return
				}
			}
		}
	}

	_, allocateUsers, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, model.AllocateUserFilter{
		TenantId:               tenantId,
		ExternalConversationId: externalConversationId,
	}, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*allocateUsers) < 1 {
		allocateUser = model.AllocateUser{
			Base:                   model.InitBase(),
			TenantId:               chatSetting.ConnectionApp.TenantId,
			ConversationId:         conversationId,
			ExternalConversationId: externalConversationId,
			AppId:                  chatSetting.Message.AppId,
			OaId:                   chatSetting.Message.OaId,
			UserId:                 authInfo.UserId,
			QueueId:                chatSetting.ConnectionQueue.QueueId,
			AllocatedTimestamp:     time.Now().UnixMilli(),
			MainAllocate:           "active",
			ConnectionId:           chatSetting.ConnectionQueue.ConnectionId,
			ConnectionQueueId:      chatSetting.ConnectionApp.ConnectionQueueId,
		}

		if err = repository.AllocateUserRepo.Insert(ctx, repository.DBConn, allocateUser); err != nil {
			log.Error(err)
			return
		}

		if err = cache.RCache.Set(USER_ALLOCATE+"_"+externalConversationId, allocateUser, USER_ALLOCATE_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}

	user.IsOk = true
	user.AuthUser = &authInfo
	user.ConnectionId = chatSetting.ConnectionQueue.ConnectionId
	user.QueueId = chatSetting.ConnectionQueue.QueueId
	user.ConnectionQueueId = chatSetting.ConnectionApp.ConnectionQueueId
	user.ConversationId = allocateUser.ConversationId
	nextUserId = allocateUser.UserId

	return
}

func (s *OttMessage) UpdateConversationById(ctx context.Context, tenantId string, allocateUser model.AllocateUser, message model.Message) (err error) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, message.AppId, allocateUser.ConversationId)
	if err != nil {
		return
	} else if conversationExist == nil {
		err = errors.New("conversation " + allocateUser.ConversationId + " not found")
		return
	}

	conversationExist.IsDone = false
	conversationExist.IsDoneBy = ""
	isDoneAt, err := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	if err != nil {
		return
	}
	conversationExist.IsDoneAt = isDoneAt
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		return
	}

	esDoc := map[string]any{}
	if err = json.Unmarshal(tmpBytes, &esDoc); err != nil {
		return
	}
	conversationQueue := model.ConversationQueue{
		DocId:        conversationExist.ConversationId,
		Conversation: *conversationExist,
	}
	if err = PublishPutConversationToChatQueue(ctx, conversationQueue); err != nil {
		log.Error(err)
		return
	}

	if err = repository.AllocateUserRepo.Update(ctx, repository.DBConn, allocateUser); err != nil {
		return
	}
	return
}
