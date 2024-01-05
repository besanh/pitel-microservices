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
	IChatQueueAgent interface {
		InsertChatQueueAgent(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueAgentRequest) error
	}
	ChatQueueAgent struct{}
)

func NewChatQueueAgent() IChatQueueAgent {
	return &ChatQueueAgent{}
}

func (s *ChatQueueAgent) InsertChatQueueAgent(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueAgentRequest) error {
	dbConn, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		return err
	}
	_, err = repository.ChatQueueAgentRepo.GetById(ctx, dbConn, data.QueueId)
	if err != nil {
		log.Error(err)
		return err
	}
	chatQueueAgent := model.ChatQueueAgent{
		Base:    model.InitBase(),
		QueueId: data.QueueId,
		AgentId: data.AgentId,
		Source:  data.Source,
	}
	err = repository.ChatQueueAgentRepo.Insert(ctx, dbConn, chatQueueAgent)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
