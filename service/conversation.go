package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IConversation interface {
		InsertConversation(ctx context.Context, conversation model.Conversation) (id string, err error)
		GetConversations(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any)
		UpdateConversationById(ctx context.Context, authUser *model.AuthUser, id string, data model.ShareInfo) (int, any)
	}
	Conversation struct {
	}
)

func NewConversation() IConversation {
	return &Conversation{}
}

func (s *Conversation) InsertConversation(ctx context.Context, conversation model.Conversation) (id string, err error) {
	docId := uuid.NewString()
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		log.Error(err)
		return docId, err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return docId, err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX, conversation.AppId); err != nil {
		log.Error(err)
		return docId, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX, conversation.AppId); err != nil {
			log.Error(err)
			return docId, err
		}
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.AppId, ES_INDEX, docId, esDoc); err != nil {
		log.Error(err)
		return docId, err
	}

	return docId, nil
}

func (s *Conversation) GetConversations(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any) {
	conversationIds := []string{}
	conversationFilter := model.AgentAllocationFilter{
		AgentId: []string{authUser.UserId},
	}
	total, agentAllocations, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, conversationFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if total > 0 {
		for _, item := range *agentAllocations {
			conversationIds = append(conversationIds, item.ConversationId)
		}
	}
	if len(conversationIds) < 1 {
		log.Error("list conversation not found")
		return response.Pagination(nil, 0, limit, offset)
	}
	filter.ConversationId = conversationIds
	total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, "", ES_INDEX_CONVERSATION, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if total > 0 {
		for k, conv := range *conversations {
			filter := model.MessageFilter{
				ConversationId: conv.ConversationId,
				IsRead:         "deactive",
			}
			total, _, err := repository.MessageESRepo.GetMessages(ctx, conv.AppId, ES_INDEX, filter, -1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			conv.TotalUnRead = int64(total)

			filterMessage := model.MessageFilter{
				ConversationId: conv.ConversationId,
			}
			total, message, err := repository.MessageESRepo.GetMessages(ctx, conv.AppId, ES_INDEX, filterMessage, 1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			if total > 0 {
				conv.LatestMessageContent = (*message)[0].Content
			}

			(*conversations)[k] = conv
		}
	}
	return response.Pagination(conversations, total, limit, offset)
}

func (s *Conversation) UpdateConversationById(ctx context.Context, authUser *model.AuthUser, id string, data model.ShareInfo) (int, any) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, id)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if conversationExist == nil {
		return response.NotFoundMsg("conversation not found")
	}
	conversationExist.ShareInfo = &data
	conversationExist.UpdatedAt = time.Now().Format(time.RFC3339)
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, id, esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OKResponse()
}
