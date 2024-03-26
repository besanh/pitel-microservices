package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

func CheckChatSetting(ctx context.Context, message model.Message) (model.User, error) {
	var user model.User
	var authInfo model.AuthUser
	var userLives []Subscriber
	var User model.UserAllocate
	var isOk bool

	newConversationId := GenerateConversationId(message.AppId, message.ExternalUserId)
	UserAllocationCache := cache.RCache.Get(USER_ALLOCATE + "_" + newConversationId)
	if UserAllocationCache != nil {
		UserTmp := model.UserAllocate{}
		if err := json.Unmarshal([]byte(UserAllocationCache.(string)), &UserTmp); err != nil {
			log.Error(err)
			user.AuthUser = &authInfo
			user.IsOk = isOk
			return user, err
		}
		User = UserTmp
		authInfo.TenantId = User.TenantId
		authInfo.UserId = User.UserId
		authInfo.Username = User.Username
		user.AuthUser = &authInfo
		user.IsOk = true
		user.ConnectionId = User.ConnectionId
		user.QueueId = User.QueueId

		return user, nil
	} else {
		filter := model.UserAllocateFilter{
			ConversationId: newConversationId,
			MainAllocate:   "active",
		}
		total, UserAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return user, err
		}
		if total > 0 {
			authInfo.TenantId = (*UserAllocations)[0].TenantId
			authInfo.UserId = (*UserAllocations)[0].UserId
			user.AuthUser = &authInfo
			user.IsOk = true
			user.ConnectionId = (*UserAllocations)[0].ConnectionId
			user.QueueId = (*UserAllocations)[0].QueueId

			if err := cache.RCache.Set(USER_ALLOCATE+"_"+newConversationId, (*UserAllocations)[0], USER_ALLOCATE_EXPIRE); err != nil {
				log.Error(err)
				return user, err
			}

			return user, nil
		} else {
			// Get connection
			connectionType := ""
			if message.MessageType == "zalo" {
				connectionType = "zalo"
			} else if message.MessageType == "face" {
				connectionType = "facebook"
			}
			connectionFilter := model.ChatConnectionAppFilter{
				ConnectionType: connectionType,
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
				totalConnectionQueue, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, 1, 0)
				if err != nil {
					log.Error(err)
					return user, err
				}
				if totalConnectionQueue > 0 {
					filterUserAllocation := model.UserAllocateFilter{
						ConversationId: newConversationId,
						QueueId:        (*connectionQueues)[0].QueueId,
						MainAllocate:   "active",
					}
					totalUserAllocation, UserAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, filterUserAllocation, -1, 0)
					if err != nil {
						log.Error(err)
						return user, err
					}
					if totalUserAllocation > 0 {
						authInfo.TenantId = (*UserAllocations)[0].TenantId
						authInfo.UserId = (*UserAllocations)[0].UserId

						for s := range WsSubscribers.Subscribers {
							if s.UserId == authInfo.UserId && (s.Level == "user" || s.Level == "User") {
								User.UserId = s.UserId
							}
						}
						if err := cache.RCache.Set(USER_ALLOCATE+"_"+newConversationId, User, USER_ALLOCATE_EXPIRE); err != nil {
							log.Error(err)
							return user, err
						}
						user.IsOk = true
						user.AuthUser = &authInfo
						user.ConnectionId = (*UserAllocations)[0].ConnectionId
						user.QueueId = (*UserAllocations)[0].QueueId

						return user, nil
					} else {
						// Connection prevent duplicate
						// Meaning: 1 connection with page A in 1 app => only recieve one queue
						queue, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, (*connectionQueues)[0].QueueId)
						if err != nil {
							log.Error(err)
							return user, err
						} else if queue == nil {
							log.Error("queue " + (*connectionQueues)[0].QueueId + " not found")
							return user, errors.New("queue " + (*connectionQueues)[0].QueueId + " not found")
						}

						chatRouting, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, queue.ChatRoutingId)
						if err != nil {
							log.Error(err)
							return user, err
						} else if chatRouting == nil {
							log.Error("chat routing " + queue.ChatRoutingId + " not found")
							return user, errors.New("chat routing " + queue.ChatRoutingId + " not found")
						}

						filterQueueUser := model.ChatQueueUserFilter{
							QueueId: []string{queue.Id},
						}
						totalQueueUsers, queueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, filterQueueUser, -1, 0)
						if err != nil {
							log.Error(err)
							return user, err
						}
						if totalQueueUsers > 0 {
							if strings.ToLower(chatRouting.RoutingAlias) == "random" {
								UserUuids := []string{}
								for _, item := range *queueUsers {
									UserUuids = append(UserUuids, item.UserId)
								}
								for s := range WsSubscribers.Subscribers {
									if slices.Contains[[]string](UserUuids, s.UserId) && s.Level == "user" || s.Level == "User" {
										userLives = append(userLives, *s)
									}
								}

								// Pick random
								if len(userLives) > 0 {
									rand.NewSource(time.Now().UnixNano())
									randomIndex := rand.Intn(len(userLives))
									tmp := userLives[randomIndex]
									User.TenantId = tmp.TenantId
									User.UserId = tmp.UserId
									User.Username = tmp.Username

									authInfo.TenantId = User.TenantId
									authInfo.UserId = User.UserId
								}
							} else if strings.ToLower(chatRouting.RoutingAlias) == "round_robin_online" {
								UserTmp, err := RoundRobinUserOnline(ctx, GenerateConversationId(message.AppId, message.ExternalUserId), queueUsers)
								if err != nil {
									log.Error(err)
									return user, err
								}
								userLives = append(userLives, *UserTmp)
								User.TenantId = UserTmp.TenantId
								User.UserId = UserTmp.UserId
								User.Username = UserTmp.Username

								authInfo.TenantId = User.TenantId
							}

							if len(userLives) > 0 {
								newConversationId = GenerateConversationId(message.AppId, message.ExternalUserId)
								UserAllocation := model.UserAllocate{
									Base:               model.InitBase(),
									TenantId:           (*connectionApps)[0].TenantId,
									ConversationId:     newConversationId,
									AppId:              message.AppId,
									UserId:             User.UserId,
									QueueId:            queue.Id,
									AllocatedTimestamp: time.Now().Unix(),
									MainAllocate:       "active",
									ConnectionId:       (*connectionQueues)[0].ConnectionId,
								}
								log.Infof("conversation %s allocated to User %s", newConversationId, User.Username)
								if err := repository.UserAllocateRepo.Insert(ctx, repository.DBConn, UserAllocation); err != nil {
									log.Error(err)
									return user, err
								}

								if err := cache.RCache.Set(USER_ALLOCATE+"_"+newConversationId, UserAllocation, USER_ALLOCATE_EXPIRE); err != nil {
									log.Error(err)
									return user, err
								}

								user.IsOk = true
								user.AuthUser = &authInfo
								user.ConnectionId = (*connectionQueues)[0].ConnectionId
								user.QueueId = (*connectionQueues)[0].QueueId

								return user, nil
							} else {
								log.Error("User not available")
								user.IsOk = true
								return user, errors.New("User not available")
							}
						} else {
							log.Error("queue User not found")
							return user, errors.New("queue User not found")
						}
					}
				} else {
					log.Error("queue not found")
					return user, errors.New("queue not found")
				}
			} else {
				log.Error("connect for conversation " + newConversationId + " not found")
				return user, errors.New("connect for conversation " + newConversationId + " not found")
			}
		}
	}
}

func UpSertConversation(ctx context.Context, connectionId string, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
	newConversationId := GenerateConversationId(data.AppId, data.ExternalUserId)
	conversation = model.Conversation{
		TenantId:         data.TenantId,
		ConversationId:   newConversationId,
		AppId:            data.AppId,
		ConversationType: data.MessageType,
		Username:         data.Username,
		Avatar:           data.Avatar,
		OaId:             data.OaId,
		ShareInfo:        data.ShareInfo,
		ExternalUserId:   data.ExternalUserId,
		CreatedAt:        time.Now().Format(time.RFC3339),
	}
	shareInfo := data.ShareInfo

	isExisted := false
	// conversationCache := cache.RCache.Get(CONVERSATION + "_" + newConversationId)
	// if conversationCache != nil {
	// 	isExisted = true
	// 	if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
	// 		log.Error(err)
	// 		return conversation, isNew, err
	// 	}
	// 	if err := UpdateESAndCache(ctx, data.TenantId, data.AppId, data.ExternalUserId, connectionId, *conversation.ShareInfo); err != nil {
	// 		log.Error(err)
	// 		return conversation, isNew, err
	// 	}
	// 	return conversation, isNew, nil
	// } else {
	filter := model.ConversationFilter{
		ConversationId: []string{newConversationId},
		AppId:          []string{data.AppId},
	}
	total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, "", ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return conversation, isNew, err
	}
	if total > 0 {
		conversation.TenantId = (*conversations)[0].TenantId
		conversation.ConversationId = (*conversations)[0].ConversationId
		conversation.ConversationType = (*conversations)[0].ConversationType
		conversation.AppId = (*conversations)[0].AppId
		conversation.OaId = (*conversations)[0].OaId
		conversation.OaName = (*conversations)[0].OaName
		conversation.OaAvatar = (*conversations)[0].OaAvatar
		conversation.ExternalUserId = (*conversations)[0].ExternalUserId
		conversation.Username = (*conversations)[0].Username
		conversation.Username = (*conversations)[0].Username
		conversation.Avatar = (*conversations)[0].Avatar
		conversation.IsDone = (*conversations)[0].IsDone
		conversation.IsDoneBy = (*conversations)[0].IsDoneBy

		conversation.ShareInfo = shareInfo
		if len(connectionId) > 0 {
			conversation, err = CacheConnection(ctx, connectionId, conversation)
			if err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
		}

		tmpBytes, err := json.Marshal(conversation)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		esDoc := map[string]any{}
		if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		newConversationId := GenerateConversationId(conversation.AppId, conversation.ExternalUserId)
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversation.AppId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isExisted = true
		return conversation, isNew, nil
	}
	// }

	if !isExisted {
		id, err := InsertConversation(ctx, conversation, connectionId)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		conversation.ConversationId = id
		if len(connectionId) > 0 {
			conversation, err = CacheConnection(ctx, connectionId, conversation)
			if err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
		}
		if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isNew = true
	}

	return conversation, isNew, nil
}

func InsertConversation(ctx context.Context, conversation model.Conversation, connectionId string) (id string, err error) {
	id = GenerateConversationId(conversation.AppId, conversation.ExternalUserId)
	if len(connectionId) > 0 {
		conversation, err = CacheConnection(ctx, connectionId, conversation)
		if err != nil {
			log.Error(err)
			return id, err
		}
	}
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		log.Error(err)
		return id, err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return id, err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX_CONVERSATION, conversation.TenantId); err != nil {
		log.Error(err)
		return id, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX_CONVERSATION, conversation.TenantId); err != nil {
			log.Error(err)
			return id, err
		}
	}
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, conversation.AppId, id)
	if err != nil {
		log.Error(err)
		return id, err
	} else if len(conversationExist.ExternalUserId) > 0 {
		log.Errorf("conversation %s not found", id)
		return id, errors.New("conversation " + id + " not found")
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.TenantId, ES_INDEX_CONVERSATION, conversation.AppId, id, esDoc); err != nil {
		log.Error(err)
		return id, err
	}

	return id, nil
}

func CheckConversationInUser(userId string, allocationUser []*model.UserAllocate) bool {
	for _, item := range allocationUser {
		if item.ConversationId == userId {
			return true
		}
	}
	return false
}

/**
* Update ES and Cache
* API get conversation can get from redis, and here can caching to descrese the number of api calls to ES
 */
func UpdateESAndCache(ctx context.Context, tenantId, appId, conversationId, connectionId string, shareInfo model.ShareInfo) error {
	var isUpdate bool
	newConversationId := GenerateConversationId(appId, conversationId)
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, appId, newConversationId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(conversationExist.ExternalUserId) < 1 {
		isUpdate = true
		// Use when routing is pitel_bss_conversation_
		conversationExistSecond, err := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, appId, newConversationId)
		if err != nil {
			log.Error(err)
			return err
		} else if len(conversationExistSecond.ExternalUserId) < 1 {
			log.Errorf("conversation %s not found", newConversationId)
			return errors.New("conversation " + newConversationId + " not found")
		}
		if err := util.ParseAnyToAny(conversationExistSecond, &conversationExist); err != nil {
			log.Error(err)
			return err
		}
		conversationExist = conversationExistSecond
	}

	conversationExist.ShareInfo = &shareInfo
	conversationExist.UpdatedAt = time.Now().Format(time.RFC3339)
	conversationExist.TenantId = tenantId
	if len(connectionId) > 0 {
		*conversationExist, err = CacheConnection(ctx, connectionId, *conversationExist)
		if err != nil {
			log.Error(err)
			return err
		}
	}
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

	// PROBLEM: Use for conv have not authUser and then having authUser
	// TODO: insert new conv
	if isUpdate {
		if err := repository.ESRepo.DeleteById(ctx, ES_INDEX_CONVERSATION, newConversationId); err != nil {
			log.Error(err)
			return err
		}
		if err := repository.ESRepo.InsertLog(ctx, tenantId, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return err
		}
	} else {
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return err
		}
	}

	if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversationExist, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
