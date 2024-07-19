package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

/**
* Chia patch cho phan loop
 */
func (s *OttMessage) CheckChatSetting(ctx context.Context, message model.Message, chatApp model.ChatApp, userChan chan<- []model.User, errChan chan<- error) {
	var authInfo model.AuthUser
	var user model.User
	var err error
	userAllocate := &model.AllocateUser{}
	tenants := []string{}
	chatIntegrateSystems := []model.ChatIntegrateSystem{}
	chatAppIntegrateSystems := []model.ChatAppIntegrateSystem{}
	conversationId := GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId)

	chatAppIntegrateSystemCache := cache.RCache.Get(CHAT_APP_INTEGRATE_SYSTEM + "_" + conversationId)
	if chatAppIntegrateSystemCache != nil {
		if err = json.Unmarshal([]byte(chatAppIntegrateSystemCache.(string)), &chatAppIntegrateSystems); err != nil {
			log.Error(err)
			errChan <- err
			return
		}
	} else {
		filterChatAppIntegrateSystem := model.ChatAppIntegrateSystemFilter{
			ChatAppId: chatApp.Id,
		}
		_, tmp, err := repository.ChatAppIntegrateSystemRepo.GetChatAppIntegrateSystems(ctx, repository.DBConn, filterChatAppIntegrateSystem, -1, 0)
		if err != nil {
			log.Error(err)
			errChan <- err
			return
		}
		chatAppIntegrateSystems = *tmp

		if err = cache.RCache.Set(CHAT_APP_INTEGRATE_SYSTEM+"_"+conversationId, chatAppIntegrateSystems, CHAT_APP_INTEGRATE_SYSTEM_EXPIRE); err != nil {
			log.Error(err)
			errChan <- err
			return
		}
	}

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
	if len(tenants) > 0 {
		for _, item := range tenants {
			userAllocatesCache, err := cache.RCache.HGetAll(USER_ALLOCATE + "_" + item + "_" + conversationId)
			if err != nil {
				log.Error(err)
				errChan <- err
				break
			} else if len(userAllocatesCache) > 0 {
				for _, userAllocateCache := range userAllocatesCache {
					if err = json.Unmarshal([]byte(userAllocateCache), userAllocate); err != nil {
						log.Error(err)
						errChan <- err
						break
					}
					// Exist user allocate
					if userAllocate.TenantId == item && userAllocate.ConversationId == conversationId {
						if userAllocate.MainAllocate == "active" {
							authInfo.TenantId = userAllocate.TenantId
							authInfo.UserId = userAllocate.UserId
							user.AuthUser = &authInfo
							user.IsOk = true
							user.ConnectionId = userAllocate.ConnectionId
							user.QueueId = userAllocate.QueueId
							user.ConnectionQueueId = userAllocate.ConnectionQueueId

							log.Infof("conversation %s allocated to username %s, id: %s, domain: %s, source: %s", conversationId, user.AuthUser.Fullname, user.AuthUser.UserId, user.AuthUser.TenantId, user.AuthUser.Source)
							userChan <- []model.User{user}
						} else {
							authInfo.TenantId = userAllocate.TenantId
							authInfo.UserId = userAllocate.UserId
							user.AuthUser = &authInfo
							user.ConnectionId = userAllocate.ConnectionId
							user.QueueId = userAllocate.QueueId
							user.ConnectionQueueId = userAllocate.ConnectionQueueId

							log.Info("user first: ", userAllocate.UserId)

							user, err := s.CheckAllSetting(ctx, GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), message, true, userAllocate, chatApp)
							if err != nil {
								log.Error(err)
								errChan <- err
								return
							}
							log.Info("user after: ", user.AuthUser.UserId, userAllocate.UserId)
							if user.AuthUser.UserId == userAllocate.UserId {
								user.IsReassignSame = true
							} else {
								user.IsReassignNew = true
								user.UserIdRemove = userAllocate.UserId
							}

							user.IsOk = true
							log.Infof("conversation %s allocated to username %s, id: %s, domain: %s, source: %s", conversationId, user.AuthUser.Fullname, user.AuthUser.UserId, user.AuthUser.TenantId, user.AuthUser.Source)
							userChan <- []model.User{user}
						}
					} else {
						user, err := s.CheckAllSetting(ctx, conversationId, message, false, nil, chatApp)
						if err != nil {
							log.Error(err)
							errChan <- err
							return
						}

						filter := model.UserAllocateFilter{
							ConversationId: conversationId,
							MainAllocate:   "deactive",
						}
						_, userAllocations, err := repository.UserAllocateRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
						if err != nil {
							log.Error(err)
							errChan <- err
							return
						}
						if len(*userAllocations) > 0 {
							if user.AuthUser.UserId == (*userAllocations)[0].UserId {
								user.IsReassignSame = true
							} else {
								user.IsReassignNew = true
								user.UserIdRemove = (*userAllocations)[0].UserId
							}
						}

						log.Infof("conversation %s allocated to username %s, id: %s", conversationId, user.AuthUser.Fullname, user.AuthUser.UserId)
						userChan <- []model.User{user}

						// Set to cache
						jsonByte, err := json.Marshal(&userAllocate)
						if err != nil {
							log.Error(err)
							errChan <- err
							return
						}
						if err = cache.RCache.HSetRaw(ctx, USER_ALLOCATE+"_"+item+"_"+conversationId, conversationId, string(jsonByte)); err != nil {
							log.Error(err)
							errChan <- err
							return
						}
					}
				}
			} else {
				user, err := s.CheckAllSetting(ctx, conversationId, message, false, nil, chatApp)
				if err != nil {
					log.Error(err)
					errChan <- err
					return
				}

				filter := model.UserAllocateFilter{
					ConversationId: conversationId,
					MainAllocate:   "deactive",
				}
				_, userAllocations, err := repository.UserAllocateRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
				if err != nil {
					log.Error(err)
					errChan <- err
					return
				}
				if len(*userAllocations) > 0 {
					if user.AuthUser.UserId == (*userAllocations)[0].UserId {
						user.IsReassignSame = true
					} else {
						user.IsReassignNew = true
						user.UserIdRemove = (*userAllocations)[0].UserId
					}
				}

				log.Infof("conversation %s allocated to username %s, id: %s", conversationId, user.AuthUser.Fullname, user.AuthUser.UserId)
				userChan <- []model.User{user}
				return
			}
		}
	}
}

/**
* Check all setting to allocate conversation to user
 */
func (s *OttMessage) CheckAllSetting(ctx context.Context, newConversationId string, message model.Message, isConversationExist bool, currentUserAllocate *model.AllocateUser, chatApp model.ChatApp) (user model.User, err error) {
	var authInfo model.AuthUser
	connectionFilter := model.ChatConnectionAppFilter{
		ConnectionType: message.MessageType,
		OaId:           message.OaId,
		AppId:          message.AppId,
		Status:         "active",
	}
	_, connectionApps, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, connectionFilter, 1, 0)
	if err != nil {
		log.Error(err)
		return user, err
	}
	if len(*connectionApps) > 0 {
		connectionQueue, err := repository.ConnectionQueueRepo.GetById(ctx, repository.DBConn, (*connectionApps)[0].ConnectionQueueId)
		if err != nil {
			log.Error(err)
			return user, err
		} else if connectionQueue == nil {
			log.Error("connection queue " + (*connectionApps)[0].ConnectionQueueId + " not found")
			return user, errors.New("connection queue " + (*connectionApps)[0].ConnectionQueueId + " not found")
		}

		filterUserAllocation := model.UserAllocateFilter{
			ConversationId: newConversationId,
			QueueId:        []string{connectionQueue.QueueId},
			MainAllocate:   "active",
		}
		_, userAllocations, err := repository.UserAllocateRepo.GetAllocateUsers(ctx, repository.DBConn, filterUserAllocation, -1, 0)
		if err != nil {
			log.Error(err)
			return user, err
		}
		if len(*userAllocations) > 0 {
			authInfo.TenantId = (*userAllocations)[0].TenantId
			authInfo.UserId = (*userAllocations)[0].UserId

			if err := cache.RCache.Set(USER_ALLOCATE+"_"+newConversationId, (*userAllocations)[0], USER_ALLOCATE_EXPIRE); err != nil {
				log.Error(err)
				return user, err
			}
			user.IsOk = true
			user.AuthUser = &authInfo
			user.ConnectionId = (*userAllocations)[0].ConnectionId
			user.QueueId = (*userAllocations)[0].QueueId
			user.ConnectionQueueId = (*userAllocations)[0].ConnectionQueueId

			return user, nil
		} else {
			// Connection prevent duplicate
			// Meaning: 1 connection with page A in 1 app => only recieve one queue
			var queue model.ChatQueue
			queueCache := cache.RCache.Get(CHAT_QUEUE + "_" + connectionQueue.QueueId)
			if queueCache != nil {
				var tmp model.ChatQueue
				if err := json.Unmarshal([]byte(queueCache.(string)), &tmp); err != nil {
					log.Error(err)
					return user, err
				}
				queue = tmp
			} else {
				queueTmp, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, connectionQueue.QueueId)
				if err != nil {
					log.Error(err)
					return user, err
				} else if queueTmp == nil {
					log.Error("queue " + connectionQueue.QueueId + " not found")
					return user, fmt.Errorf("queue " + connectionQueue.QueueId + " not found")
				}
				queue = *queueTmp
			}

			chatRouting, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, queue.ChatRoutingId)
			if err != nil {
				log.Error(err)
				return user, err
			} else if chatRouting == nil {
				log.Error("chat routing " + queue.ChatRoutingId + " not found")
				return user, fmt.Errorf("chat routing " + queue.ChatRoutingId + " not found")
			}

			filterQueueUser := model.ChatQueueUserFilter{
				QueueId: []string{queue.Id},
			}
			_, chatQueueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, filterQueueUser, -1, 0)
			if err != nil {
				log.Error(err)
				return user, err
			}
			if len(*chatQueueUsers) > 0 {
				chatSetting := model.ChatSetting{
					ConnectionApp:       (*connectionApps)[0],
					ConnectionQueue:     *connectionQueue,
					QueueUser:           *chatQueueUsers,
					RoutingAlias:        chatRouting.RoutingAlias,
					Message:             message,
					ConnectionQueueUser: *chatQueueUsers,
					ManagerQueueUser:    (*chatQueueUsers)[0],
				}

				userTmp, err := s.GetAllocateUser(ctx, chatSetting, isConversationExist, currentUserAllocate)
				if err != nil {
					user.ConnectionId = connectionQueue.ConnectionId
					user.ConnectionQueueId = connectionQueue.Id
					return user, err
				}
				if err := util.ParseAnyToAny(userTmp, &user); err != nil {
					log.Error(err)
					return user, err
				}
				user.QueueId = userTmp.QueueId
				user.ConnectionId = connectionQueue.ConnectionId
				user.ConnectionQueueId = connectionQueue.Id

				return user, nil
			} else {
				log.Error("queue user not found")
				return user, fmt.Errorf("queue user not found")
			}
		}
	} else {
		log.Error("connect for conversation " + newConversationId + " not found")
		return user, fmt.Errorf("connect for conversation " + newConversationId + " not found")
	}
}

/**
* Get user
* if isConversationExist = true,  it means conversation is exist, and we can get user from user_allocate
* if isConversationExist = false, it means conversation is not exist, we need to get user from chat_setting
 */
func (s *OttMessage) GetAllocateUser(ctx context.Context, chatSetting model.ChatSetting, isConversationExist bool, currentUserAllocate *model.AllocateUser) (user model.User, err error) {
	userAllocate := model.AllocateUser{}
	var authInfo model.AuthUser
	var userLives []Subscriber
	var isUserAccept bool

	if strings.ToLower(chatSetting.RoutingAlias) == "random" {
		for s := range WsSubscribers.Subscribers {
			if s.Level == "user" || s.Level == "agent" {
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
					if item.UserId == tmp.UserId {
						userAllocate.TenantId = tmp.TenantId
						userAllocate.UserId = tmp.UserId
						userAllocate.Username = tmp.Username

						authInfo.TenantId = userAllocate.TenantId
						authInfo.UserId = userAllocate.UserId
						isUserAccept = true
						break
					}
				}
			}
		}
	} else if strings.ToLower(chatSetting.RoutingAlias) == "round_robin_online" {
		userTmp, isOk, err := RoundRobinUserOnline(ctx, chatSetting.ConnectionApp.TenantId, GenerateConversationId(chatSetting.Message.AppId, chatSetting.Message.OaId, chatSetting.Message.ExternalUserId), &chatSetting.QueueUser)
		if !isOk {
			log.Error(err)
			return user, err
		} else if userTmp != nil {
			// TODO: check user exist in queue
			if len(chatSetting.ConnectionQueueUser) > 0 {
				for _, item := range chatSetting.ConnectionQueueUser {
					if item.UserId == userTmp.UserId {
						userAllocate.TenantId = userTmp.TenantId
						userAllocate.UserId = userTmp.UserId
						userAllocate.Username = userTmp.Username

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

	if !isUserAccept {
		userAllocate.UserId = chatSetting.ManagerQueueUser.UserId
		authInfo.UserId = chatSetting.ManagerQueueUser.UserId
		authInfo.TenantId = chatSetting.ManagerQueueUser.TenantId
	}

	conversationId := GenerateConversationId(chatSetting.Message.AppId, chatSetting.Message.OaId, chatSetting.Message.ExternalUserId)
	if isConversationExist {
		currentUserAllocate.UserId = userAllocate.UserId
		currentUserAllocate.MainAllocate = "active"
		currentUserAllocate.AllocatedTimestamp = time.Now().UnixMilli()
		currentUserAllocate.UpdatedAt = time.Now()
		tenantId := userAllocate.TenantId
		if len(tenantId) < 1 {
			tenantId = (*currentUserAllocate).TenantId
		}
		if err = s.UpdateConversationById(ctx, tenantId, *currentUserAllocate, chatSetting.Message); err != nil {
			log.Error(err)
			return user, err
		}

		authInfo.TenantId = userAllocate.TenantId
		authInfo.UserId = userAllocate.UserId
		authInfo.Username = userAllocate.Username

		user.AuthUser = &authInfo
		user.ConnectionId = chatSetting.ConnectionQueue.ConnectionId
		user.QueueId = chatSetting.ConnectionQueue.QueueId
		user.ConnectionQueueId = chatSetting.ConnectionApp.ConnectionQueueId
		if len(userAllocate.UserId) < 1 {
			user.IsOk = false
		}

		return user, nil
	} else {
		if len(chatSetting.Message.ConversationId) > 0 {
			total, conversationDeactiveExist, errConv := repository.UserAllocateRepo.GetAllocateUsers(ctx, repository.DBConn, model.UserAllocateFilter{
				AppId:          chatSetting.Message.AppId,
				ConversationId: chatSetting.Message.ConversationId,
			}, -1, 0)
			if errConv != nil {
				log.Error(errConv)
				return user, errConv
			}
			if total > 0 {
				if err := repository.UserAllocateRepo.DeleteAllocateUsers(ctx, repository.DBConn, *conversationDeactiveExist); err != nil {
					log.Error(err)
					return user, err
				}
			}
		}

		userAllocate = model.AllocateUser{
			Base:               model.InitBase(),
			TenantId:           chatSetting.ConnectionApp.TenantId,
			ConversationId:     conversationId,
			AppId:              chatSetting.Message.AppId,
			OaId:               chatSetting.Message.OaId,
			UserId:             userAllocate.UserId,
			QueueId:            chatSetting.ConnectionQueue.QueueId,
			AllocatedTimestamp: time.Now().UnixMilli(),
			MainAllocate:       "active",
			ConnectionId:       chatSetting.ConnectionQueue.ConnectionId,
			ConnectionQueueId:  chatSetting.ConnectionApp.ConnectionQueueId,
			Username:           userAllocate.Username,
		}

		if err = repository.UserAllocateRepo.Insert(ctx, repository.DBConn, userAllocate); err != nil {
			log.Error(err)
			return
		}

		if err = cache.RCache.Set(USER_ALLOCATE+"_"+conversationId, userAllocate, USER_ALLOCATE_EXPIRE); err != nil {
			log.Error(err)
			return
		}

		user.IsOk = true
		user.AuthUser = &authInfo
		user.ConnectionId = chatSetting.ConnectionQueue.ConnectionId
		user.QueueId = chatSetting.ConnectionQueue.QueueId
		user.ConnectionQueueId = chatSetting.ConnectionApp.ConnectionQueueId

		return user, nil
	}
}

func (s *OttMessage) UpdateConversationById(ctx context.Context, tenantId string, userAllocate model.AllocateUser, message model.Message) error {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, message.AppId, GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId))
	if err != nil {
		return err
	}

	conversationExist.IsDone = false
	conversationExist.IsDoneBy = ""
	isDoneAt, err := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	if err != nil {
		return err
	}
	conversationExist.IsDoneAt = isDoneAt
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		return err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		return err
	}
	if err = PublishPutConversationToChatQueue(ctx, *conversationExist); err != nil {
		log.Error(err)
		return err
	}

	if err := repository.UserAllocateRepo.Update(ctx, repository.DBConn, userAllocate); err != nil {
		return err
	}
	return nil
}
