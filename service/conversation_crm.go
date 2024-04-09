package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

func (s *Conversation) GetConversationsByManage(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any) {
	filter.TenantId = authUser.TenantId
	if authUser.Source == "authen" {
		var queueUuids string
		if authUser.Level == "manager" {
			filterManageQueue := model.ChatManageQueueUserFilter{
				ManageId: authUser.UserId,
			}
			totalManageQueue, manageQueues, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterManageQueue, -1, 0)
			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
			if totalManageQueue > 0 {
				queueUuids = (*manageQueues)[0].QueueId
			}
		}
		total, conversations, err := getConversationByFilter(ctx, queueUuids, filter, limit, offset)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		return response.Pagination(conversations, total, limit, offset)
	} else {
		return response.Pagination(nil, 0, limit, offset)
	}
}

func getConversationByFilter(ctx context.Context, queueUuids string, filter model.ConversationFilter, limit, offset int) (total int, conversations *[]model.ConversationView, err error) {
	conversationIds := []string{}
	conversationFilter := model.UserAllocateFilter{
		TenantId: filter.TenantId,
		QueueId:  queueUuids,
	}
	if filter.IsDone.Bool {
		conversationFilter.MainAllocate = "deactive"
	} else {
		conversationFilter.MainAllocate = "active"
	}

	total, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, conversationFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return total, nil, err
	}
	if total > 0 {
		for _, item := range *userAllocations {
			conversationIds = append(conversationIds, item.ConversationId)
		}
	}
	if len(conversationIds) < 1 {
		log.Error("list conversation not found")
		return total, nil, err
	}
	filter.ConversationId = conversationIds
	total, conversations, err = repository.ConversationESRepo.GetConversations(ctx, "", ES_INDEX_CONVERSATION, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return total, nil, err
	}
	if total > 0 {
		for k, conv := range *conversations {
			filter := model.MessageFilter{
				TenantId:       conv.TenantId,
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

			(*conversations)[k] = conv
		}
	}
	return
}
