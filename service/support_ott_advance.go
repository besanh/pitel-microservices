package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func GetManageQueueUser(ctx context.Context, queueId string) (manageQueueUser *model.ChatManageQueueUser, err error) {
	manageQueueUserCache := cache.RCache.Get(MANAGE_QUEUE_USER + "_" + queueId)
	if manageQueueUserCache != nil {
		if err = json.Unmarshal([]byte(manageQueueUserCache.(string)), &manageQueueUser); err != nil {
			log.Error(err)
		}
	}
	filter := model.ChatManageQueueUserFilter{
		QueueId: queueId,
	}
	total, manageQueueUsers, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		manageQueueUser = &(*manageQueueUsers)[0]
		if err = cache.RCache.Set(MANAGE_QUEUE_USER+"_"+queueId, manageQueueUser, MANAGE_QUEUE_USER_EXPIRE); err != nil {
			log.Error(err)
			return
		}
	}
	return manageQueueUser, nil
}

func GenerateConversationId(appId, conversationId string) (newConversationId string) {
	newConversationId = appId + "_" + conversationId
	return
}

func RoundRobinUserOnline(ctx context.Context, conversationId string, queueUsers *[]model.ChatQueueUser) (*Subscriber, error) {
	userLive := Subscriber{}
	userLives := []Subscriber{}
	subscribers, err := cache.RCache.HGetAll(BSS_SUBSCRIBERS)
	if err != nil {
		log.Error(err)
		return &userLive, err
	}
	for _, item := range subscribers {
		s := Subscriber{}
		if err := json.Unmarshal([]byte(item), &s); err != nil {
			log.Error(err)
			return &userLive, err
		}
		if (s.Level == "user" || s.Level == "User") && CheckInLive(*queueUsers, s.Id) {
			userLives = append(userLives, s)
		}
	}
	if len(userLives) > 0 {
		index, userAllocate := GetUserIsRoundRobin(userLives)
		userLive = *userAllocate
		userLive.IsAssignRoundRobin = true
		userPrevious := Subscriber{}
		if index < len(userLives) {
			userPrevious = userLives[(index+1)%len(userLives)]
		} else {
			userPrevious = userLives[0]
		}
		userPrevious.IsAssignRoundRobin = false

		// Update current
		jsonByteUserLive, err := json.Marshal(&userLive)
		if err != nil {
			log.Error(err)
			return &userLive, err
		}
		if err := cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, userLive.Id, string(jsonByteUserLive)); err != nil {
			log.Error(err)
			return &userLive, err
		}

		// Update previous
		if userPrevious.Id != userLive.Id {
			jsonByteUserLivePrevious, err := json.Marshal(&userPrevious)
			if err != nil {
				log.Error(err)
				return &userLive, err
			}
			if err := cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, userPrevious.Id, string(jsonByteUserLivePrevious)); err != nil {
				log.Error(err)
				return &userLive, err
			}
		}
		return &userLive, nil
	} else {
		return &userLive, errors.New("no user online")
	}
}

func GetUserIsRoundRobin(userLives []Subscriber) (int, *Subscriber) {
	isOk := false
	index := 0
	userLive := Subscriber{}
	for i, item := range userLives {
		if item.IsAssignRoundRobin {
			if (i+1)%len(userLives) <= len(userLives) {
				userLive = userLives[(i+1)%len(userLives)]
				isOk = true
				index = (i + 1) % len(userLives)
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
	userLive = userLives[0]
	return index, &userLive
}

func CheckInLive(queueUsers []model.ChatQueueUser, id string) bool {
	for _, item := range queueUsers {
		if item.UserId == id {
			return true
		}
	}
	return false
}

// TODO: caching
func CacheConnection(ctx context.Context, connectionId string, conversation model.Conversation) (model.Conversation, error) {
	connectionExist, err := repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, connectionId)
	if err != nil {
		log.Error(err)
		return conversation, err
	}
	if connectionExist != nil {
		if connectionExist.ConnectionType == "zalo" {
			conversation.OaName = connectionExist.OaInfo.Zalo[0].OaName
			conversation.OaAvatar = connectionExist.OaInfo.Zalo[0].Avatar
		} else if connectionExist.ConnectionType == "facebook" {
			conversation.OaName = connectionExist.OaInfo.Facebook[0].OaName
			conversation.OaAvatar = connectionExist.OaInfo.Facebook[0].Avatar
		}
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

func CheckConfigApp(ctx context.Context, appId string) (isExist bool, err error) {
	chatAppCache := cache.RCache.Get(CHAT_APP + "_" + appId)
	if chatAppCache != nil {
		isExist = true
	} else {
		filter := model.AppFilter{
			AppId: appId,
		}
		total, chatApp, err := repository.ChatAppRepo.GetChatApp(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			return isExist, err
		} else if total > 0 {
			isExist = true
			if err = cache.RCache.Set(CHAT_APP+"_"+appId, chatApp, CHAT_APP_EXPIRE); err != nil {
				log.Error(err)
				return isExist, err
			}
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
		total, connections, errConnection := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
		if err != nil {
			log.Error(err)
			err = errConnection
			return
		}
		if total < 1 {
			log.Error("connect for app_id: " + appId + ", oa_id: " + oaId + " not found")
			err = errors.New("connect for app_id: " + appId + ", oa_id: " + oaId + " not found")
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
