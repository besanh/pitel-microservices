package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/queue"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func GenerateConversationId(appId, oaId, conversationId string) (newConversationId string) {
	if len(appId) == 0 {
		log.Error("app id cannot be empty")
		return
	} else if len(oaId) == 0 {
		log.Error("oa id cannot be empty")
		return
	} else if len(conversationId) == 0 {
		log.Error("conversation id cannot be empty")
		return
	}
	newConversationId = appId + "_" + oaId + "_" + conversationId
	return
}

func GetManageQueueUser(ctx context.Context, queueId string) (manageQueueUser *model.ChatManageQueueUser, err error) {
	manageQueueUserCache := cache.RCache.Get(MANAGE_QUEUE_USER + "_" + queueId)
	if manageQueueUserCache != nil {
		if err = json.Unmarshal([]byte(manageQueueUserCache.(string)), &manageQueueUser); err != nil {
			log.Error(err)
			return
		}
	}
	filter := model.ChatManageQueueUserFilter{
		QueueId: queueId,
	}
	_, manageQueueUsers, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*manageQueueUsers) > 0 {
		manageQueueUser = &(*manageQueueUsers)[0]
		if err = cache.RCache.Set(MANAGE_QUEUE_USER+"_"+queueId, manageQueueUser, MANAGE_QUEUE_USER_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}
	return manageQueueUser, nil
}

func RoundRobinUserOnline(ctx context.Context, tenantId, conversationId string, queueUsers *[]model.ChatQueueUser) (userLive *Subscriber, err error) {
	userLives := []Subscriber{}
	subscribers, err := cache.RCache.HGetAll(BSS_SUBSCRIBERS)
	if err != nil {
		log.Error(err)
		return
	}
	for _, item := range subscribers {
		s := Subscriber{}
		if err = json.Unmarshal([]byte(item), &s); err != nil {
			log.Error(err)
			return
		}
		if (s.Level == "user" || s.Level == "agent") && CheckInLive(*queueUsers, s.Id) {
			userLives = append(userLives, s)
		}
	}
	if len(userLives) > 0 {
		index, userAllocate := GetUserIsRoundRobin(tenantId, userLives)
		userLive = userAllocate
		userLive.IsAssignRoundRobin = true
		// Check len > 1 for assign online 1 user
		if len(userLives) > 1 {
			userPrevious := Subscriber{}
			if index == 0 {
				userPrevious = userLives[len(userLives)-1]
			} else {
				userPrevious = userLives[index-1]
			}
			userPrevious.IsAssignRoundRobin = false

			// Update current
			jsonByteUserLive, errTmp := json.Marshal(&userLive)
			if errTmp != nil {
				log.Error(errTmp)
				err = errTmp
				return
			}
			if err = cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, userLive.Id, string(jsonByteUserLive)); err != nil {
				log.Error(err)
				return
			}

			// Update previous
			if userPrevious.Id != userLive.Id {
				jsonByteUserLivePrevious, errTmp := json.Marshal(&userPrevious)
				if errTmp != nil {
					log.Error(errTmp)
					err = errTmp
					return
				}
				if err = cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, userPrevious.Id, string(jsonByteUserLivePrevious)); err != nil {
					log.Error(err)
					return
				}
			}
		}
		return
	}

	// Because if no user online, conversation will assign to manager
	return
}

func GetUserIsRoundRobin(tenantId string, userLives []Subscriber) (int, *Subscriber) {
	isOk := false
	index := 0
	userLive := Subscriber{}
	for i, item := range userLives {
		if item.IsAssignRoundRobin && item.TenantId == tenantId {
			if (i + 1) < len(userLives) {
				userLive = userLives[(i + 1)]
				isOk = true
				index = (i + 1)
				break
			} else {
				isOk = true
				userLive = userLives[0]
				break
			}
		}
	}
	if isOk {
		return index, &userLive
	}
	if len(userLives) > 0 {
		userLive = userLives[0]
	}
	return index, &userLive
}

func CheckInLive(queueUsers []model.ChatQueueUser, id string) (isExist bool) {
	low := 0
	high := len(queueUsers) - 1
	mid := -1
	for low <= high {
		mid = (low + high) / 2
		if queueUsers[mid].Id == id {
			return true
		} else if queueUsers[mid].Id < id {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if mid != -1 {
		isExist = true
	}
	return
}

// TODO: caching
func CacheConnection(ctx context.Context, connectionId string, conversation model.Conversation) (model.Conversation, error) {
	connectionExist, err := repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, connectionId)
	if err != nil {
		log.Error(err)
		return conversation, err
	}
	if connectionExist != nil {
		if connectionExist.ConnectionType == "zalo" && conversation.ConversationType == "zalo" {
			conversation.OaName = connectionExist.OaInfo.Zalo[0].OaName
			conversation.OaAvatar = connectionExist.OaInfo.Zalo[0].Avatar
		} else if connectionExist.ConnectionType == "facebook" && conversation.ConversationType == "facebook" {
			conversation.OaName = connectionExist.OaInfo.Facebook[0].OaName
			conversation.OaAvatar = connectionExist.OaInfo.Facebook[0].Avatar
		}
	}

	if err := cache.RCache.Set(CONVERSATION+"_"+conversation.ConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return conversation, err
	}
	return conversation, nil
}

func GetProfile(ctx context.Context, appId, oaId, userId string) (result *model.ProfileResponse, err error) {
	params := map[string]string{
		"app_id": appId,
		"oa_id":  oaId,
		"uid":    userId,
	}
	url := OTT_URL + "/ott/v1/zalo/profile"
	client := resty.New()
	var res *resty.Response
	res, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		Get(url)
	if err != nil {
		log.Error(err)
		return
	}

	var resp model.ProfileResponse
	if err = json.Unmarshal([]byte(res.Body()), &resp); err != nil {
		log.Error(err)
		return
	}

	result = &resp

	return
}

func CheckConfigAppCache(ctx context.Context, appId string) (chatApp *model.ChatApp, err error) {
	chatAppCache := cache.RCache.Get(CHAT_APP + "_" + appId)
	if chatAppCache != nil {
		if err := json.Unmarshal([]byte(chatAppCache.(string)), &chatApp); err != nil {
			log.Error(err)
			return chatApp, err
		}
	} else {
		filter := model.ChatAppFilter{
			AppId:  appId,
			Status: "active",
		}
		total, chatApps, err := repository.ChatAppRepo.GetChatApp(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			return chatApp, err
		} else if total > 0 {
			chatApp = &(*chatApps)[0]
			if err = cache.RCache.Set(CHAT_APP+"_"+appId, chatApp, CHAT_APP_EXPIRE); err != nil {
				return chatApp, err
			}
		} else {
			return chatApp, fmt.Errorf("app %s not found", appId)
		}
	}
	return
}

func GetConfigConnectionAppCache(ctx context.Context, appId, oaId, connectionType string) (connectionApp model.ChatConnectionApp, err error) {
	connectionAppCache := cache.RCache.Get(CHAT_CONNECTION + "_" + appId + "_" + oaId)
	if connectionAppCache != nil {
		var tmp model.ChatConnectionApp
		if err = json.Unmarshal([]byte(connectionAppCache.(string)), &tmp); err != nil {
			log.Error(err)
			return
		}
		connectionApp = tmp
	} else {
		filter := model.ChatConnectionAppFilter{
			AppId:          appId,
			OaId:           oaId,
			ConnectionType: connectionType,
		}
		_, connections, errConnection := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
		if errConnection != nil {
			log.Error(err)
			err = errConnection
			return
		}
		if len(*connections) < 1 {
			err = fmt.Errorf("connect for app_id: %s, oa_id: %s, connection_type: %s not found", appId, oaId, connectionType)
			log.Error(err)
			return
		}

		if err = cache.RCache.Set(CHAT_CONNECTION+"_"+appId+"_"+oaId, (*connections)[0], CHAT_CONNECTION_EXPIRE); err != nil {
			log.Error(err)
			return
		}

		connectionApp = (*connections)[0]
	}
	return
}

func PublishPutConversationToChatQueue(ctx context.Context, conversation model.ConversationQueue) (err error) {
	var b []byte
	if b, err = json.Marshal(conversation); err != nil {
		log.Error(err)
		return
	}
	if err = queue.RMQ.Client.PublishBytes(BSS_CHAT_QUEUE_NAME, b); err != nil {
		log.Error(err)
		return
	}
	return
}

func CacheChatRouting(ctx context.Context, chatRoutingId string) (chatRouting *model.ChatRouting, err error) {
	chatRoutingCache := cache.RCache.Get(CHAT_ROUTING + "_" + chatRoutingId)
	if chatRoutingCache != nil {
		if err = json.Unmarshal([]byte(chatRoutingCache.(string)), &chatRouting); err != nil {
			log.Error(err)
			return
		}
	} else {
		chatRouting, err = repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, chatRoutingId)
		if err != nil {
			log.Error(err)
			return
		} else if chatRouting == nil {
			err = errors.New("chat routing " + chatRoutingId + " not found")
			log.Error(err)
			return
		}

		if err = cache.RCache.Set(CHAT_ROUTING+"_"+chatRoutingId, chatRouting, CHAT_ROUTING_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}

	return
}
