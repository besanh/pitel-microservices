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
		UpdateChatQueueAgentById(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueAgentRequest) (*model.ChatQueueAgentUpdateResponse, error)
	}
	ChatQueueAgent struct{}
)

func NewChatQueueAgent() IChatQueueAgent {
	return &ChatQueueAgent{}
}

func (s *ChatQueueAgent) InsertChatQueueAgent(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueAgentRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		return err
	}
	_, err = repository.ChatQueueAgentRepo.GetById(ctx, dbCon, data.QueueId)
	if err != nil {
		log.Error(err)
		return err
	}

	filter := model.ChatQueueAgentFilter{
		QueueId: []string{data.QueueId},
	}
	total, chatQueueAgents, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}

	for _, item := range data.AgentId {
		chatQueueAgent := model.ChatQueueAgent{
			Base:     model.InitBase(),
			TenantId: authUser.TenantId,
			QueueId:  data.QueueId,
			AgentId:  item,
			Source:   authUser.Source,
		}
		err = repository.ChatQueueAgentRepo.Insert(ctx, dbCon, chatQueueAgent)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if total > 0 {
		for _, item := range *chatQueueAgents {
			if err = repository.ChatQueueAgentRepo.Delete(ctx, dbCon, item.Id); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}

func (s *ChatQueueAgent) UpdateChatQueueAgentById(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueAgentRequest) (*model.ChatQueueAgentUpdateResponse, error) {
	result := model.ChatQueueAgentUpdateResponse{}
	totalSuccess := len(data.AgentId)
	totalFail := 0
	dbCon, err := HandleGetDBConSource(authUser)
	if err == ERR_EMPTY_CONN {
		err = errors.New(response.ERR_EMPTY_CONN)
		totalSuccess -= 1
		totalFail += 1
		return nil, err
	}

	filter := model.ChatQueueAgentFilter{
		QueueId: []string{data.QueueId},
	}
	total, chatQueueAgents, err := repository.ChatQueueAgentRepo.GetChatQueueAgents(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, item := range data.AgentId {
		_, err := repository.ChatQueueAgentRepo.GetById(ctx, dbCon, item)
		if err != nil {
			log.Error(err)
			totalSuccess -= 1
			totalFail += 1
			result.ListFail = append(result.ListFail, item)
			return nil, err
		}

		chatQueueAgent := model.ChatQueueAgent{
			Base:     model.InitBase(),
			TenantId: authUser.TenantId,
			QueueId:  data.QueueId,
			AgentId:  item,
			Source:   authUser.Source,
		}
		if err = repository.ChatQueueAgentRepo.Insert(ctx, dbCon, chatQueueAgent); err != nil {
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
		for _, item := range *chatQueueAgents {
			if err = repository.ChatQueueAgentRepo.Delete(ctx, dbCon, item.Id); err != nil {
				log.Error(err)
				return nil, err
			}
		}
	}
	result.TotalSuccess = totalSuccess
	result.TotalFail = totalFail

	return &result, nil
}
