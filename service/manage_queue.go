package service

import (
	"context"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IManageQueue interface {
		PostManageQueue(ctx context.Context, authUser *model.AuthUser, data model.ManageQueueAgentRequest) (id string, err error)
		UpdateManageQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ManageQueueAgentRequest) error
		DeleteManageQueueById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ManageQueue struct{}
)

func NewManageQueue() IManageQueue {
	return &ManageQueue{}
}

func (s *ManageQueue) PostManageQueue(ctx context.Context, authUser *model.AuthUser, data model.ManageQueueAgentRequest) (id string, err error) {
	manageQueue := model.ChatManageQueueAgent{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return manageQueue.GetId(), err
	}

	manageQueue.TenantId = authUser.TenantId
	manageQueue.QueueId = data.QueueId
	manageQueue.AgentId = data.AgentId

	if err := repository.ManageQueueRepo.Insert(ctx, dbCon, manageQueue); err != nil {
		log.Error(err)
		return manageQueue.GetId(), err
	}
	return manageQueue.GetId(), nil
}

func (s *ManageQueue) UpdateManageQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ManageQueueAgentRequest) (err error) {
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

	manageQueueExist.QueueId = data.QueueId
	manageQueueExist.UpdatedAt = time.Now()
	if err = repository.ManageQueueRepo.Update(ctx, dbCon, *manageQueueExist); err != nil {
		log.Error(err)
		return err
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
