package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatQueueUser interface {
		InsertChatQueueUser(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueUserRequest) error
		UpdateChatQueueUserById(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueUserRequest) (*model.ChatQueueUserUpdateResponse, error)
	}
	ChatQueueUser struct{}
)

func NewChatQueueUser() IChatQueueUser {
	return &ChatQueueUser{}
}

func (s *ChatQueueUser) InsertChatQueueUser(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueUserRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = repository.ChatQueueUserRepo.GetById(ctx, dbCon, data.QueueId)
	if err != nil {
		log.Error(err)
		return err
	}

	filter := model.ChatQueueUserFilter{
		QueueId: []string{data.QueueId},
	}
	total, chatQueueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}

	for _, item := range data.UserId {
		chatQueueUser := model.ChatQueueUser{
			Base:     model.InitBase(),
			TenantId: authUser.TenantId,
			QueueId:  data.QueueId,
			UserId:   item,
			Source:   authUser.Source,
		}
		err = repository.ChatQueueUserRepo.Insert(ctx, dbCon, chatQueueUser)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if total > 0 {
		for _, item := range *chatQueueUsers {
			if err = repository.ChatQueueUserRepo.Delete(ctx, dbCon, item.Id); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}

func (s *ChatQueueUser) UpdateChatQueueUserById(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueUserRequest) (*model.ChatQueueUserUpdateResponse, error) {
	result := model.ChatQueueUserUpdateResponse{}
	totalSuccess := len(data.UserId)
	totalFail := 0
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		totalSuccess -= 1
		totalFail += 1
		return nil, err
	}

	filter := model.ChatQueueUserFilter{
		QueueId: []string{data.QueueId},
	}
	total, chatQueueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, item := range data.UserId {
		_, err := repository.ChatQueueUserRepo.GetById(ctx, dbCon, item)
		if err != nil {
			log.Error(err)
			totalSuccess -= 1
			totalFail += 1
			result.ListFail = append(result.ListFail, item)
			return nil, err
		}

		chatQueueUser := model.ChatQueueUser{
			Base:     model.InitBase(),
			TenantId: authUser.TenantId,
			QueueId:  data.QueueId,
			UserId:   item,
			Source:   authUser.Source,
		}
		if err = repository.ChatQueueUserRepo.Insert(ctx, dbCon, chatQueueUser); err != nil {
			log.Error(err)
			totalSuccess -= 1
			totalFail += 1
			result.ListFail = append(result.ListFail, item)
			return nil, err
		}

		result.ListSuccess = append(result.ListSuccess, item)
	}

	// After insert success, remove old item
	if total > 0 {
		for _, item := range *chatQueueUsers {
			// TODO: clear cache
			chatQueueCache := cache.RCache.Get(CHAT_QUEUE + "_" + item.Id)
			if chatQueueCache != nil {
				if err = cache.RCache.Del([]string{CHAT_QUEUE + "_" + item.Id}); err != nil {
					log.Error(err)
					continue
				}
			}
			if err = repository.ChatQueueUserRepo.Delete(ctx, dbCon, item.Id); err != nil {
				log.Error(err)
				return nil, err
			}
		}
	}
	result.TotalSuccess = totalSuccess
	result.TotalFail = totalFail

	return &result, nil
}
