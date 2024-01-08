package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatQueue interface {
		InsertChatQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequest) (string, error)
	}
	ChatQueue struct{}
)

func NewChatQueue() IChatQueue {
	return &ChatQueue{}
}

func (s *ChatQueue) InsertChatQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequest) (string, error) {
	chatQueue := model.ChatQueue{
		Base: model.InitBase(),
	}
	dbConn, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		return chatQueue.Base.GetId(), err
	}
	routingExist, err := repository.ChatRoutingRepo.GetById(ctx, dbConn, data.ChatRoutingId)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	} else if routingExist == nil {
		err = errors.New(response.ERR_DATA_NOT_FOUND)
		return chatQueue.Base.GetId(), err
	}
	chatQueue.QueueName = data.QueueName
	chatQueue.Description = data.Description
	chatQueue.ChatRoutingId = data.ChatRoutingId
	chatQueue.Status = data.Status
	err = repository.ChatQueueRepo.Insert(ctx, dbConn, chatQueue)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}
	return chatQueue.Base.GetId(), nil
}
