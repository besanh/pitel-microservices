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
	chatQueue := model.ChatQueue{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}
	routingExist, err := repository.ChatRoutingRepo.GetById(ctx, dbCon, data.ChatRoutingId)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	} else if routingExist == nil {
		err = errors.New(response.ERR_DATA_NOT_FOUND)
		return chatQueue.Base.GetId(), err
	}

	connectionAgents := []model.ConnectionQueue{}
	if len(data.ConnectionId) > 0 {
		for _, item := range data.ConnectionId {
			connectionAgent := model.ConnectionQueue{
				Base:         model.InitBase(),
				ConnectionId: item,
				QueueId:      chatQueue.Base.GetId(),
				Status:       data.Status,
			}
			connectionAgents = append(connectionAgents, connectionAgent)
		}
	}

	if len(connectionAgents) > 0 {
		if err = repository.ConnectionQueueRepo.BulkInsert(ctx, dbCon, connectionAgents); err != nil {
			log.Error(err)
			return chatQueue.Base.GetId(), err
		}
	}

	chatQueue.QueueName = data.QueueName
	chatQueue.Description = data.Description
	chatQueue.ChatRoutingId = data.ChatRoutingId
	chatQueue.Status = data.Status

	err = repository.ChatQueueRepo.Insert(ctx, dbCon, chatQueue)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}
	return chatQueue.Base.GetId(), nil
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

	filter := model.ConnectionQueueFilter{
		QueueId: queueExist.Id,
	}
	total, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}

	if total > 0 {
		for _, item := range *connectionQueues {
			if err := repository.ConnectionQueueRepo.Delete(ctx, dbCon, item.Id); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	connectionAgents := []model.ConnectionQueue{}
	for _, item := range data.ConnectionId {
		connectionAgent := model.ConnectionQueue{
			Base:         model.InitBase(),
			ConnectionId: item,
			QueueId:      queueExist.Id,
			Status:       data.Status,
		}
		connectionAgents = append(connectionAgents, connectionAgent)
	}
	if err = repository.ConnectionQueueRepo.BulkInsert(ctx, dbCon, connectionAgents); err != nil {
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
