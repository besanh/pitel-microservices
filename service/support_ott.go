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

	userAllocationCache := cache.RCache.Get(USER_ALLOCATE + "_" + GenerateConversationId(message.AppId, message.ExternalUserId))
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
		filter := model.UserAllocateFilter{
			ConversationId: GenerateConversationId(message.AppId, message.ExternalUserId),
		}
		_, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return user, err
		}
		if len((*userAllocations)) > 0 {
			authInfo.TenantId = (*userAllocations)[0].TenantId
			authInfo.UserId = (*userAllocations)[0].UserId
			user.AuthUser = &authInfo
			user.IsOk = true
			user.ConnectionId = (*userAllocations)[0].ConnectionId
			user.QueueId = (*userAllocations)[0].QueueId

			user, err := CheckAllSetting(ctx, GenerateConversationId(message.AppId, message.ExternalUserId), message, (*userAllocations)[0])
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

			return user, nil
		} else {
			user, err := CheckAllSetting(ctx, GenerateConversationId(message.AppId, message.ExternalUserId), message, userAllocate)
			if err != nil {
				log.Error(err)
				return user, err
			}

			return user, nil
		}
	}
}

/**
* Check all setting to allocate conversation to user
 */
func CheckAllSetting(ctx context.Context, newConversationId string, message model.Message, userAllocate model.UserAllocate) (user model.User, err error) {
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

				for s := range WsSubscribers.Subscribers {
					if s.UserId == authInfo.UserId && (s.Level == "user" || s.Level == "agent") {
						userAllocate.UserId = s.UserId
					}
				}
				if err := cache.RCache.Set(USER_ALLOCATE+"_"+newConversationId, userAllocate, USER_ALLOCATE_EXPIRE); err != nil {
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
						ConnectionApp:   &(*connectionApps)[0],
						ConnectionQueue: &(*connectionQueues)[0],
						QueueUser:       chatQueueUsers,
						RoutingAlias:    chatRouting.RoutingAlias,
						Message:         &message,
						UserAllocate:    &userAllocate,
					}

					userTmp, err := GetAllocateUser(ctx, chatSetting)
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
			log.Error("queue not found")
			return user, fmt.Errorf("queue not found")
		}
	} else {
		log.Error("connect for conversation " + newConversationId + " not found")
		return user, fmt.Errorf("connect for conversation " + newConversationId + " not found")
	}
}

/**
* Pick user
 */
func GetAllocateUser(ctx context.Context, chatSetting model.ChatSetting) (user model.User, err error) {
	var userAllocate model.UserAllocate
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
		userTmp, err := RoundRobinUserOnline(ctx, GenerateConversationId(chatSetting.Message.AppId, chatSetting.Message.ExternalUserId), chatSetting.QueueUser)
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

	if len(userLives) > 0 {
		conversationId := GenerateConversationId(chatSetting.Message.AppId, chatSetting.Message.ExternalUserId)
		if len(chatSetting.UserAllocate.UserId) > 0 {
			chatSetting.UserAllocate.MainAllocate = "active"
			chatSetting.UserAllocate.AllocatedTimestamp = time.Now().Unix()
			chatSetting.UserAllocate.UpdatedAt = time.Now()
			if err = UpdateConversationById(ctx, chatSetting.UserAllocate.TenantId, *chatSetting.UserAllocate, *chatSetting.Message); err != nil {
				log.Error(err)
				return
			}
			userAllocateTmp := chatSetting.UserAllocate
			userAllocateTmp.Username = userAllocate.Username
			userAllocate = *userAllocateTmp
		} else {
			userAllocate = model.UserAllocate{
				Base:               model.InitBase(),
				TenantId:           chatSetting.ConnectionApp.TenantId,
				ConversationId:     conversationId,
				AppId:              chatSetting.Message.AppId,
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
		}

		log.Infof("conversation %s allocated to user %s, id: %s", conversationId, userAllocate.Username, userAllocate.Id)

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

func UpdateConversationById(ctx context.Context, tenantId string, userAllocate model.UserAllocate, message model.Message) error {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, message.AppId, GenerateConversationId(message.AppId, message.ExternalUserId))
	if err != nil {
		log.Error(err)
		return err
	}

	conversationExist.IsDone = false
	conversationExist.IsDoneBy = ""
	isDoneAt, err := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
	if err != nil {
		log.Error(err)
		return err
	}
	conversationExist.IsDoneAt = isDoneAt
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		log.Error(err)
		return err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return err
	}
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, message.AppId, GenerateConversationId(message.AppId, message.ExternalUserId), esDoc); err != nil {
		log.Error(err)
		return err
	}
	if err := repository.UserAllocateRepo.Update(ctx, repository.DBConn, userAllocate); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
