package service

import (
	"context"
	"encoding/json"
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

func CheckChatSetting(ctx context.Context, message model.Message) (model.User, error) {
	var user model.User
	var authInfo model.AuthUser
	var userAllocate model.UserAllocate
	var isOk bool

	userAllocationCache := cache.RCache.Get(USER_ALLOCATE + "_" + GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId))
	if userAllocationCache != nil {
		if err := json.Unmarshal([]byte(userAllocationCache.(string)), &userAllocate); err != nil {
			log.Error(err)
			user.AuthUser = &authInfo
			user.IsOk = isOk
			return user, err
		}
		authInfo.TenantId = userAllocate.TenantId
		authInfo.UserId = userAllocate.UserId
		authInfo.Username = userAllocate.Username
		user.AuthUser = &authInfo
		user.IsOk = true
		user.ConnectionId = userAllocate.ConnectionId
		user.QueueId = userAllocate.QueueId

		return user, nil
	} else {
		/**
		* TODO: check conversation is exist or not
		* if not exist, create
		* if exist, then check status is active or not
		* if status is active, then get user
		* if status is not active, then check setting to reassign
		 */
		filter := model.UserAllocateFilter{
			ConversationId: GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId),
		}
		_, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return user, err
		}
		if len((*userAllocations)) > 0 {
			if (*userAllocations)[0].MainAllocate == "deactive" {
				authInfo.TenantId = (*userAllocations)[0].TenantId
				authInfo.UserId = (*userAllocations)[0].UserId
				user.AuthUser = &authInfo
				user.IsOk = true
				user.ConnectionId = (*userAllocations)[0].ConnectionId
				user.QueueId = (*userAllocations)[0].QueueId

				user, err := CheckAllSetting(ctx, GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), message, true, &(*userAllocations)[0])
				if err != nil {
					log.Error(err)
					return user, err
				}
				if user.AuthUser.UserId == (*userAllocations)[0].UserId {
					user.IsReassignSame = true
				} else {
					user.IsReassignNew = true
					user.UserIdRemove = (*userAllocations)[0].UserId
				}

				log.Infof("conversation %s allocated to username %s, id: %s", GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), user.AuthUser.Fullname, user.AuthUser.UserId)
				return user, nil
			} else {
				authInfo.TenantId = (*userAllocations)[0].TenantId
				authInfo.UserId = (*userAllocations)[0].UserId
				user.AuthUser = &authInfo
				user.IsOk = true
				user.ConnectionId = (*userAllocations)[0].ConnectionId
				user.QueueId = (*userAllocations)[0].QueueId

				log.Infof("conversation %s allocated to username %s, id: %s", GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), user.AuthUser.Fullname, user.AuthUser.UserId)
				return user, nil
			}
		} else {
			user, err := CheckAllSetting(ctx, GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), message, false, nil)
			if err != nil {
				log.Error(err)
				return user, err
			}

			filter := model.UserAllocateFilter{
				ConversationId: GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId),
				MainAllocate:   "deactive",
			}
			_, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, filter, 1, 0)
			if err != nil {
				log.Error(err)
				return user, err
			}
			if len((*userAllocations)) > 0 {
				if user.AuthUser.UserId == (*userAllocations)[0].UserId {
					user.IsReassignSame = true
				} else {
					user.IsReassignNew = true
					user.UserIdRemove = (*userAllocations)[0].UserId
				}
			}

			log.Infof("conversation %s allocated to username %s, id: %s", GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), user.AuthUser.Fullname, user.AuthUser.UserId)
			return user, nil
		}
	}
}

/**
* Check all setting to allocate conversation to user
 */
func CheckAllSetting(ctx context.Context, newConversationId string, message model.Message, isConversationExist bool, currentUserAllocate *model.UserAllocate) (user model.User, err error) {
	var authInfo model.AuthUser
	connectionFilter := model.ChatConnectionAppFilter{
		ConnectionType: message.MessageType,
		OaId:           message.OaId,
		AppId:          message.AppId,
		Status:         "active",
	}
	totalConnection, connectionApps, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, connectionFilter, 1, 0)
	if err != nil {
		log.Error(err)
		return user, err
	}
	if totalConnection > 0 {
		filter := model.ConnectionQueueFilter{
			ConnectionId: (*connectionApps)[0].Id,
		}
		_, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return user, err
		}
		if len(*connectionQueues) > 0 {
			filteruserAllocation := model.UserAllocateFilter{
				ConversationId: newConversationId,
				QueueId:        (*connectionQueues)[0].QueueId,
				MainAllocate:   "active",
			}
			_, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, filteruserAllocation, -1, 0)
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

				return user, nil
			} else {
				// Connection prevent duplicate
				// Meaning: 1 connection with page A in 1 app => only recieve one queue
				var queue model.ChatQueue
				queueCache := cache.RCache.Get(CHAT_QUEUE + "_" + (*connectionQueues)[0].QueueId)
				if queueCache != nil {
					var tmp model.ChatQueue
					if err := json.Unmarshal([]byte(queueCache.(string)), &tmp); err != nil {
						log.Error(err)
						return user, err
					}
					queue = tmp
				} else {
					queueTmp, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, (*connectionQueues)[0].QueueId)
					if err != nil {
						log.Error(err)
						return user, err
					} else if queueTmp == nil {
						log.Error("queue " + (*connectionQueues)[0].QueueId + " not found")
						return user, fmt.Errorf("queue " + (*connectionQueues)[0].QueueId + " not found")
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
				_, chatQueueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, filterQueueUser, 1, 0)
				if err != nil {
					log.Error(err)
					return user, err
				}
				if len(*chatQueueUsers) > 0 {
					chatSetting := model.ChatSetting{
						ConnectionApp:   (*connectionApps)[0],
						ConnectionQueue: (*connectionQueues)[0],
						QueueUser:       *chatQueueUsers,
						RoutingAlias:    chatRouting.RoutingAlias,
						Message:         message,
					}

					userTmp, err := GetAllocateUser(ctx, chatSetting, isConversationExist, currentUserAllocate)
					if err != nil {
						user.QueueId = (*chatQueueUsers)[0].QueueId
						user.ConnectionId = (*connectionQueues)[0].ConnectionId
						return user, err
					}
					if err := util.ParseAnyToAny(userTmp, &user); err != nil {
						log.Error(err)
						return user, err
					}
					user.QueueId = (*chatQueueUsers)[0].QueueId
					user.ConnectionId = (*connectionQueues)[0].ConnectionId

					return user, nil
				} else {
					log.Error("queue user not found")
					return user, fmt.Errorf("queue user not found")
				}
			}
		} else {
			log.Errorf("queue not found with connection %s", (*connectionApps)[0].Id)
			return user, fmt.Errorf("queue not found with connection %s", (*connectionApps)[0].Id)
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
func GetAllocateUser(ctx context.Context, chatSetting model.ChatSetting, isConversationExist bool, currentUserAllocate *model.UserAllocate) (user model.User, err error) {
	userAllocate := model.UserAllocate{}
	var authInfo model.AuthUser
	var userLives []Subscriber

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
			userAllocate.TenantId = tmp.TenantId
			userAllocate.UserId = tmp.UserId
			userAllocate.Username = tmp.Username

			authInfo.TenantId = userAllocate.TenantId
			authInfo.UserId = userAllocate.UserId
		}
	} else if strings.ToLower(chatSetting.RoutingAlias) == "round_robin_online" {
		userTmp, err := RoundRobinUserOnline(ctx, GenerateConversationId(chatSetting.Message.AppId, chatSetting.Message.OaId, chatSetting.Message.ExternalUserId), &chatSetting.QueueUser)
		if err != nil {
			log.Error(err)
			return user, err
		}
		userLives = append(userLives, *userTmp)
		userAllocate.TenantId = userTmp.TenantId
		userAllocate.UserId = userTmp.UserId
		userAllocate.Username = userTmp.Username

		authInfo.TenantId = userTmp.TenantId
		authInfo.UserId = userTmp.UserId
		authInfo.Username = userTmp.Username
		authInfo.Level = userTmp.Level
	}

	conversationId := GenerateConversationId(chatSetting.Message.AppId, chatSetting.Message.OaId, chatSetting.Message.ExternalUserId)
	if isConversationExist {
		currentUserAllocate.UserId = userAllocate.UserId
		currentUserAllocate.MainAllocate = "active"
		currentUserAllocate.AllocatedTimestamp = time.Now().Unix()
		currentUserAllocate.UpdatedAt = time.Now()
		tenantId := userAllocate.TenantId
		if len(tenantId) < 1 {
			tenantId = (*currentUserAllocate).TenantId
		}
		if err = UpdateConversationById(ctx, tenantId, *currentUserAllocate, chatSetting.Message); err != nil {
			log.Error(err)
			return user, err
		}

		authInfo.TenantId = userAllocate.TenantId
		authInfo.UserId = userAllocate.UserId
		authInfo.Username = userAllocate.Username

		user.AuthUser = &authInfo
		user.ConnectionId = chatSetting.ConnectionQueue.ConnectionId
		user.QueueId = chatSetting.ConnectionQueue.QueueId
		if len(userAllocate.UserId) < 1 {
			user.IsOk = false
		}

		return user, nil
	} else {
		if len(userLives) > 0 {
			total, conversationDeactiveExist, errConv := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, model.UserAllocateFilter{
				AppId:          chatSetting.Message.AppId,
				ConversationId: chatSetting.Message.ConversationId,
			}, -1, 0)
			if errConv != nil {
				log.Error(errConv)
				return user, errConv
			}
			if total > 0 {
				if err := repository.UserAllocateRepo.DeleteUserAllocates(ctx, repository.DBConn, *conversationDeactiveExist); err != nil {
					log.Error(err)
					return user, err
				}
			}
			userAllocate = model.UserAllocate{
				Base:               model.InitBase(),
				TenantId:           chatSetting.ConnectionApp.TenantId,
				ConversationId:     conversationId,
				AppId:              chatSetting.Message.AppId,
				OaId:               chatSetting.Message.OaId,
				UserId:             userAllocate.UserId,
				QueueId:            chatSetting.ConnectionQueue.QueueId,
				AllocatedTimestamp: time.Now().Unix(),
				MainAllocate:       "active",
				ConnectionId:       chatSetting.ConnectionQueue.ConnectionId,
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

			return user, nil
		} else {
			log.Error("user not available")
			return user, fmt.Errorf("user not available")
		}
	}
}

func UpdateConversationById(ctx context.Context, tenantId string, userAllocate model.UserAllocate, message model.Message) error {
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
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, message.AppId, GenerateConversationId(message.AppId, message.OaId, message.ExternalUserId), esDoc); err != nil {
		return err
	}

	if err := repository.UserAllocateRepo.Update(ctx, repository.DBConn, userAllocate); err != nil {
		return err
	}
	return nil
}
