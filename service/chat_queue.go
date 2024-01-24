package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatQueue interface {
		InsertChatQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequest) (string, error)
		GetChatQueues(ctx context.Context, authUser *model.AuthUser, filter model.QueueFilter, limit, offset int) (int, *[]model.ChatQueue, error)
		GetChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatQueue, error)
		UpdateChatQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatQueueRequest) error
		DeleteChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) error
	}
	ChatQueue struct{}
)

func NewChatQueue() IChatQueue {
	return &ChatQueue{}
}

func (s *ChatQueue) InsertChatQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequest) (string, error) {
	queueId := uuid.NewString()
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return queueId, err
	}
	routingExist, err := repository.ChatRoutingRepo.GetById(ctx, dbCon, data.ChatRoutingId)
	if err != nil {
		log.Error(err)
		return queueId, err
	} else if routingExist == nil {
		err = errors.New(response.ERR_DATA_NOT_FOUND)
		return queueId, err
	}

	chatQueues := []model.ChatQueue{}

	for _, item := range data.ConnectionId {
		connection, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, item)
		if err != nil {
			log.Error(err)
			return queueId, err
		} else if connection == nil {
			log.Error("connection " + item + " not found")
		}
		chatQueue := model.ChatQueue{
			Base:          model.InitBase(),
			QueueId:       queueId,
			QueueName:     data.QueueName,
			ConnectionId:  item,
			ChatRoutingId: data.ChatRoutingId,
			Description:   data.Description,
			ChatRouting:   routingExist,
			Status:        data.Status,
		}
		chatQueues = append(chatQueues, chatQueue)
	}
	err = repository.ChatQueueRepo.BulkInsert(ctx, dbCon, chatQueues)
	if err != nil {
		log.Error(err)
		return queueId, err
	}
	return queueId, nil
}

func (s *ChatQueue) GetChatQueues(ctx context.Context, authUser *model.AuthUser, filter model.QueueFilter, limit, offset int) (int, *[]model.ChatQueue, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	total, queues, err := repository.ChatQueueRepo.GetQueues(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, queues, nil
}

func (s *ChatQueue) GetChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatQueue, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err == ERR_EMPTY_CONN {
		log.Error(err)
		return nil, err
	}
	queue, err := repository.ChatQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return queue, nil
}

func (s *ChatQueue) UpdateChatQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatQueueRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	queueExist, err := repository.ChatQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	queueExist.QueueName = data.QueueName
	queueExist.Description = data.Description
	queueExist.ChatRoutingId = data.ChatRoutingId
	queueExist.Status = data.Status
	err = repository.ChatQueueRepo.Update(ctx, dbCon, *queueExist)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ChatQueue) DeleteChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = repository.ChatQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ChatQueueRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// Delete queue agent
	if err := repository.ChatQueueAgentRepo.DeleteChatQueueAgents(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
