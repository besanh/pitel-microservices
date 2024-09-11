package service

import (
	"context"
	"errors"
	"time"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
)

type (
	IManageQueue interface {
		PostManageQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatManageQueueUserRequest) (id string, err error)
		UpdateManageQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatManageQueueUserRequest) error
		DeleteManageQueueById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ManageQueue struct{}
)

var ManageQueueService IManageQueue

func NewManageQueue() IManageQueue {
	return &ManageQueue{}
}

func (s *ManageQueue) PostManageQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatManageQueueUserRequest) (id string, err error) {
	manageQueue := model.ChatManageQueueUser{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return manageQueue.GetId(), err
	}

	manageQueue.TenantId = authUser.TenantId
	manageQueue.ConnectionId = data.ConnectionId
	manageQueue.QueueId = data.QueueId
	manageQueue.UserId = data.UserId

	queueExist, err := repository.ChatQueueRepo.GetById(ctx, dbCon, data.QueueId)
	if err != nil {
		log.Error(err)
		return manageQueue.GetId(), err
	} else if queueExist == nil {
		log.Error("queue not found")
		return manageQueue.GetId(), errors.New("queue " + data.QueueId + " not found")
	}
	queueExist.ManageQueueId = manageQueue.GetId()

	if err := repository.ManageQueueRepo.Insert(ctx, dbCon, manageQueue); err != nil {
		log.Error(err)
		return manageQueue.GetId(), err
	}
	if err := repository.ChatQueueRepo.Update(ctx, dbCon, *queueExist); err != nil {
		log.Error(err)
		return manageQueue.GetId(), err
	}

	return manageQueue.GetId(), nil
}

func (s *ManageQueue) UpdateManageQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatManageQueueUserRequest) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}
	manageQueueExist, err := repository.ManageQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if manageQueueExist == nil {
		log.Error("manage queue not found")
		return errors.New("manage queue " + id + " not found")
	}

	queueExist, err := repository.ChatQueueRepo.GetById(ctx, dbCon, data.QueueId)
	if err != nil {
		log.Error(err)
		return err
	} else if queueExist == nil {
		log.Error("queue not found")
		return errors.New("queue " + data.QueueId + " not found")
	}

	if data.IsNew {
		// TODO: check exist data
		// TODO: move to transaction
		// Remove old connection queue
		filter := model.ConnectionQueueFilter{
			ConnectionId: data.ConnectionId,
			QueueId:      queueExist.Id,
		}
		_, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, dbCon, filter, -1, 0)
		if err != nil {
			log.Error(err)
			return err
		}

		if len(*connectionQueues) > 0 {
			for _, item := range *connectionQueues {
				if err := repository.ConnectionQueueRepo.Delete(ctx, dbCon, item.Id); err != nil {
					log.Error(err)
					return err
				}
			}
		}

		manageQueueNew := model.ChatManageQueueUser{
			Base:         model.InitBase(),
			TenantId:     authUser.TenantId,
			ConnectionId: data.ConnectionId,
			QueueId:      data.QueueId,
			UserId:       data.UserId,
		}
		if err = repository.ManageQueueRepo.Insert(ctx, dbCon, manageQueueNew); err != nil {
			log.Error(err)
			return err
		}
	} else {
		// TODO: move to transaction
		manageQueueExist.QueueId = data.QueueId
		manageQueueExist.UserId = data.UserId
		manageQueueExist.UpdatedAt = time.Now()
		if err = repository.ManageQueueRepo.Update(ctx, dbCon, *manageQueueExist); err != nil {
			log.Error(err)
			return err
		}

		queueExist.ManageQueueId = manageQueueExist.GetId()
		if err := repository.ChatQueueRepo.Update(ctx, dbCon, *queueExist); err != nil {
			log.Error(err)
			return err
		}
	}

	return
}

func (s *ManageQueue) DeleteManageQueueById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = repository.ManageQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ManageQueueRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
