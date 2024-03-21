package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IConversation interface {
		InsertConversation(ctx context.Context, conversation model.Conversation) (id string, err error)
		GetConversations(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any)
		GetConversationsByManager(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any)
		UpdateConversationById(ctx context.Context, authUser *model.AuthUser, appId, id string, data model.ShareInfo) (int, any)
		UpdateStatusConversation(ctx context.Context, authUser *model.AuthUser, appId, id, updatedBy, status string) error
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
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX, conversation.TenantId); err != nil {
		log.Error(err)
		return docId, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX, conversation.TenantId); err != nil {
			log.Error(err)
			return docId, err
		}
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.TenantId, ES_INDEX, conversation.AppId, docId, esDoc); err != nil {
		log.Error(err)
		return docId, err
	}

	return docId, nil
}

func (s *Conversation) GetConversations(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any) {
	conversationIds := []string{}
	conversationFilter := model.AgentAllocationFilter{
		AgentId:      []string{authUser.UserId},
		MainAllocate: "active",
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
	filter.TenantId = authUser.TenantId
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
				EventNameExlucde: []string{
					"received",
					"seen",
				},
			}
			total, _, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX, filter, -1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			conv.TotalUnRead = int64(total)

			filterMessage := model.MessageFilter{
				ConversationId: conv.ConversationId,
			}
			totalTmp, message, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX, filterMessage, 1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			if totalTmp > 0 {
				if slices.Contains[[]string](variables.ATTACHMENT_TYPE, (*message)[0].EventName) {
					conv.LatestMessageContent = (*message)[0].EventName
				} else {
					conv.LatestMessageContent = (*message)[0].Content
				}
			}
			conv.LatestMessageDirection = (*message)[0].Direction

			(*conversations)[k] = conv
		}
	}
	return response.Pagination(conversations, total, limit, offset)
}

func (s *Conversation) UpdateConversationById(ctx context.Context, authUser *model.AuthUser, appId, id string, data model.ShareInfo) (int, any) {
	newConversationId := GenerateConversationId(appId, id)
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, appId, newConversationId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(conversationExist.ConversationId) < 1 {
		return response.NotFoundMsg("conversation " + newConversationId + " not found")
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
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OKResponse()
}

func (s *Conversation) UpdateStatusConversation(ctx context.Context, authUser *model.AuthUser, appId, conversationId, updatedBy, status string) error {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, appId, conversationId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found", conversationId)
		return errors.New("conversation " + conversationId + " not found")
	}

	if status != "reopen" {
		if conversationExist.IsDone {
			log.Errorf("conversation %s is done", conversationId)
			return errors.New("conversation " + conversationId + " is done")
		}
	}

	statusAllocate := "active"
	if status == "reopen" {
		statusAllocate = "deactive"
	}

	// Update agent allocate
	filter := model.AgentAllocationFilter{
		AppId:          appId,
		ConversationId: conversationId,
		MainAllocate:   statusAllocate,
	}
	total, agentAllocate, err := repository.AgentAllocationRepo.GetAgentAllocations(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total < 1 {
		log.Errorf("conversation %s not found with active agent", conversationId)
		return errors.New("conversation " + conversationId + " not found with active agent")
	}

	agentAllocateTmp := (*agentAllocate)[0]

	if status == "done" {
		agentAllocateTmp.MainAllocate = "deactive"
		agentAllocateTmp.AllocatedTimestamp = time.Now().Unix()
		agentAllocateTmp.UpdatedAt = time.Now()
		if err := repository.AgentAllocationRepo.Update(ctx, repository.DBConn, agentAllocateTmp); err != nil {
			log.Error(err)
			return err
		}
		conversationExist.IsDone = true
		conversationExist.IsDoneBy = updatedBy
		conversationExist.IsDoneAt = time.Now()
	} else if status == "reopen" {
		agentAllocateTmp.MainAllocate = "active"
		agentAllocateTmp.AllocatedTimestamp = time.Now().Unix()
		agentAllocateTmp.UpdatedAt = time.Now()
		if err := repository.AgentAllocationRepo.Update(ctx, repository.DBConn, agentAllocateTmp); err != nil {
			log.Error(err)
			return err
		}

		conversationExist.IsDone = false
		conversationExist.IsDoneBy = ""
		isDoneAt, err := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
		if err != nil {
			log.Error(err)
			return err
		}
		conversationExist.IsDoneAt = isDoneAt
	}

	if slices.Contains([]string{"done", "reopen"}, status) {
		tmpBytes, err := json.Marshal(conversationExist)
		if err != nil {
			log.Error(err)
			return err
		}
		esDoc := map[string]any{}
		if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
			log.Error(err)
			return err
		}
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, appId, conversationId, esDoc); err != nil {
			log.Error(err)
			if err := repository.AgentAllocationRepo.Update(ctx, repository.DBConn, (*agentAllocate)[0]); err != nil {
				log.Error(err)
			}
			return err
		}
	}

	// Event to manager
	manageQueueAgent, err := GetManageQueueAgent(ctx, agentAllocateTmp.QueueId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(manageQueueAgent.Id) < 1 {
		log.Error("queue " + agentAllocateTmp.QueueId + " not found")
		return errors.New("queue " + agentAllocateTmp.QueueId + " not found")
	}

	for s := range WsSubscribers.Subscribers {
		if s.Id == manageQueueAgent.AgentId {
			// TODO: publish message to manager
			event := map[string]any{
				"event_name": variables.EVENT_CHAT["conversation_done"],
				"event_data": map[string]any{
					"message": conversationExist,
				},
			}
			if err := PublishMessageToOne(manageQueueAgent.AgentId, event); err != nil {
				log.Error(err)
				return err
			}
			break
		}
	}

	return nil
}
