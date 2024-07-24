package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatConnectionQueue interface {
		GetChatConnectionQueueById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ConnectionQueue, error)
	}
	ChatConnectionQueue struct{}
)

var ChatConnectionQueueService IChatConnectionQueue

func NewChatConnectionQueue() *ChatConnectionQueue {
	return &ChatConnectionQueue{}
}

func (s *ChatConnectionQueue) GetChatConnectionQueueById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ConnectionQueue, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}
	result, err = repository.ConnectionQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}
