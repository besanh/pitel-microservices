package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatUser interface {
		InsertChatUser(ctx context.Context, request *model.ChatUserRequest) (id string, err error)
		UpdateChatUserStatusById(ctx context.Context, authUser *model.AuthUser, status string) (err error)
		GetChatUserStatusById(ctx context.Context, authUser *model.AuthUser) (status string, err error)
	}
	ChatUser struct {
	}
)

var ChatUserService IChatUser

func NewChatUser() IChatUser {
	return &ChatUser{}
}

func (s *ChatUser) InsertChatUser(ctx context.Context, request *model.ChatUserRequest) (id string, err error) {
	chatUser := model.ChatUser{
		Base: model.InitBase(),
		Salt: uuid.NewString(),
	}

	_, chatUsers, err := repository.ChatUserRepo.GetChatUsers(ctx, repository.DBConn, model.ChatUserFilter{
		Username: request.Username,
	}, 1, 0)
	if err != nil {
		log.Error(err)
		return
	} else if len(*chatUsers) > 0 {
		err = fmt.Errorf("username already exists")
		log.Error(err)
		return
	}

	roleExist, err := repository.ChatRoleRepo.GetById(ctx, repository.DBConn, request.RoleId)
	if err != nil {
		log.Error(err)
		return
	} else if roleExist == nil {
		err = errors.New("role " + request.RoleId + " not found")
		log.Error(err)
		return
	}

	chatUser.Username = request.Username
	chatUser.Fullname = request.Fullname
	dbPassword := []byte(chatUser.Salt + request.Password)
	chatUser.Password = fmt.Sprintf("%x", md5.Sum(dbPassword))
	chatUser.Email = request.Email
	chatUser.Avatar = request.Avatar
	chatUser.Level = request.Level
	chatUser.Status = request.Status
	chatUser.RoleId = request.RoleId
	chatUser.CreatedAt = time.Now()
	if err = repository.ChatUserRepo.Insert(ctx, repository.DBConn, chatUser); err != nil {
		log.Error(err)
		return
	}

	id = chatUser.GetId()
	return
}

func (s *ChatUser) UpdateChatUserStatusById(ctx context.Context, authUser *model.AuthUser, status string) (err error) {
	subscriber, err := cache.RCache.HGet(BSS_SUBSCRIBERS, authUser.UserId)
	if err != nil && errors.Is(err, redis.Nil) {
		err = fmt.Errorf("subscriber %s not found", authUser.UserId)
		log.Error(err)
		return
	} else if err != nil {
		log.Error(err)
		return
	}

	// update in cache
	sub := Subscriber{}
	if err = json.Unmarshal([]byte(subscriber), &sub); err != nil {
		log.Error(err)
		return
	}
	sub.Status = status
	jsonByteSubscriber, err := json.Marshal(&sub)
	if err != nil {
		log.Error(err)
		return
	}
	if err = cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, authUser.UserId, string(jsonByteSubscriber)); err != nil {
		log.Error(err)
		return
	}

	// update in memory
	for wsSubscriber := range WsSubscribers.Subscribers {
		if wsSubscriber.Id == authUser.UserId {
			wsSubscriber.Status = status
		}
	}

	// send event to ws
	eventMessage := map[string]string{
		"status":  status,
		"user_id": authUser.UserId,
	}
	go PublishWsEventToOneUser(variables.EVENT_USER_STATUS["user_status_updated"], "data", authUser.UserId, []string{authUser.UserId}, true, eventMessage)
	return
}

func (s *ChatUser) GetChatUserStatusById(ctx context.Context, authUser *model.AuthUser) (status string, err error) {
	for wsSubscriber := range WsSubscribers.Subscribers {
		if wsSubscriber.Id == authUser.UserId {
			status = wsSubscriber.Status
			return
		}
	}

	subscriber, err := cache.RCache.HGet(BSS_SUBSCRIBERS, authUser.UserId)
	if err != nil && errors.Is(err, redis.Nil) {
		err = fmt.Errorf("subscriber %s not found", authUser.UserId)
		log.Error(err)
		return
	} else if err != nil {
		log.Error(err)
		return
	}

	sub := Subscriber{}
	if err = json.Unmarshal([]byte(subscriber), &sub); err != nil {
		log.Error(err)
		return
	}
	status = sub.Status
	return
}
