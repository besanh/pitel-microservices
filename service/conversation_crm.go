package service

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/variables"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
	"golang.org/x/exp/slices"
)

func (s *Conversation) GetConversationsByHighLevel(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (total int, conversations *[]model.ConversationView, err error) {
	filter.TenantId = authUser.TenantId
	var queueId []string
	if authUser.Level == "manager" {
		filterManageQueue := model.ChatManageQueueUserFilter{
			TenantId: authUser.TenantId,
			UserId:   authUser.UserId,
		}
		totalManageQueue, manageQueues, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterManageQueue, -1, 0)
		if err != nil {
			log.Error(err)
			return 0, nil, err
		}
		if totalManageQueue > 0 {
			for _, item := range *manageQueues {
				queueId = append(queueId, item.QueueId)
			}
		}
	}
	total, conversations, err = s.getConversationByFilter(ctx, queueId, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	if len(*conversations) > 0 {
		for k, item := range *conversations {
			if item.Labels != nil && !reflect.DeepEqual(item.Labels, "") {
				var labels []map[string]string
				if err = json.Unmarshal([]byte(item.Labels), &labels); err != nil {
					log.Error(err)
					return
				}
				chatLabelIds := []string{}
				if len(labels) > 0 {
					for _, item := range labels {
						chatLabelIds = append(chatLabelIds, item["label_id"])
					}
					if len(chatLabelIds) > 0 {
						_, chatLabelExist, err := repository.ChatLabelRepo.GetChatLabels(ctx, repository.DBConn, model.ChatLabelFilter{
							TenantId: authUser.TenantId,
							LabelIds: chatLabelIds,
						}, -1, 0)
						if err != nil {
							log.Error(err)
							return 0, nil, err
						}
						if len(*chatLabelExist) > 0 {
							tmp, err := json.Marshal(*chatLabelExist)
							if err != nil {
								log.Error(err)
								return 0, nil, err
							}
							(*conversations)[k].Labels = tmp
						} else {
							(*conversations)[k].Labels = []byte("[]")
						}
					} else {
						(*conversations)[k].Labels = []byte("[]")
					}
				} else {
					(*conversations)[k].Labels = []byte("[]")
				}
			} else {
				(*conversations)[k].Labels = []byte("[]")
			}
		}
	}

	return
}

func (s *Conversation) getConversationByFilter(ctx context.Context, queueUuids []string, filter model.ConversationFilter, limit, offset int) (total int, conversations *[]model.ConversationView, err error) {
	conversationIds := []string{}
	conversationFilter := model.AllocateUserFilter{
		TenantId: filter.TenantId,
	}
	if len(queueUuids) > 0 {
		conversationFilter.QueueId = queueUuids
	}
	if filter.IsDone.Bool {
		conversationFilter.MainAllocate = "deactive"
	} else {
		conversationFilter.MainAllocate = "active"
	}

	total, userAllocations, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, conversationFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return total, nil, err
	}
	if total > 0 {
		for _, item := range *userAllocations {
			conversationIds = append(conversationIds, item.ConversationId)
		}
	}
	filter.ConversationId = conversationIds
	total, conversations, err = repository.ConversationESRepo.GetConversations(ctx, filter.TenantId, ES_INDEX_CONVERSATION, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
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
			total, _, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX_MESSAGE, filter, -1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			conv.TotalUnRead = int64(total)

			filterMessage := model.MessageFilter{
				ConversationId: conv.ConversationId,
			}
			_, message, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX_MESSAGE, filterMessage, 1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			if len(*message) > 0 {
				if slices.Contains[[]string](variables.ATTACHMENT_TYPE, (*message)[0].EventName) {
					conv.LatestMessageContent = (*message)[0].EventName
				} else {
					conv.LatestMessageContent = (*message)[0].Content
				}
				conv.LatestMessageDirection = (*message)[0].Direction
			}

			(*conversations)[k] = conv
		}
	}
	return
}

func (s *Conversation) GetConversationsByHighLevelWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit int, scrollId string) (total int, conversations []*model.ConversationView, respScrollId string, err error) {
	filter.TenantId = authUser.TenantId
	var queueId []string
	if authUser.Level == "manager" {
		filterManageQueue := model.ChatManageQueueUserFilter{
			TenantId: authUser.TenantId,
			UserId:   authUser.UserId,
		}
		totalManageQueue, manageQueues, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterManageQueue, -1, 0)
		if err != nil {
			log.Error(err)
			return 0, nil, "", err
		}
		if totalManageQueue > 0 {
			for _, item := range *manageQueues {
				queueId = append(queueId, item.QueueId)
			}
		}
	}
	total, conversations, respScrollId, err = s.getConversationByFilterWithScrollAPI(ctx, queueId, filter, limit, scrollId)
	if err != nil {
		log.Error(err)
		return
	}

	for k, item := range conversations {
		if conversations[k] == nil {
			continue
		}
		if item.Labels != nil && !reflect.DeepEqual(item.Labels, "") {
			var labels []map[string]string
			if err = json.Unmarshal([]byte(item.Labels), &labels); err != nil {
				log.Error(err)
				return
			}
			chatLabelIds := []string{}
			if len(labels) > 0 {
				for _, item := range labels {
					chatLabelIds = append(chatLabelIds, item["label_id"])
				}
				if len(chatLabelIds) > 0 {
					_, chatLabelExist, err := repository.ChatLabelRepo.GetChatLabels(ctx, repository.DBConn, model.ChatLabelFilter{
						TenantId: authUser.TenantId,
						LabelIds: chatLabelIds,
					}, -1, 0)
					if err != nil {
						log.Error(err)
						return 0, nil, "", err
					}
					if len(*chatLabelExist) > 0 {
						tmp, err := json.Marshal(*chatLabelExist)
						if err != nil {
							log.Error(err)
							return 0, nil, "", err
						}
						conversations[k].Labels = tmp
					} else {
						conversations[k].Labels = []byte("[]")
					}
				} else {
					conversations[k].Labels = []byte("[]")
				}
			} else {
				conversations[k].Labels = []byte("[]")
			}
		} else {
			conversations[k].Labels = []byte("[]")
		}
	}

	return
}

func (s *Conversation) getConversationByFilterWithScrollAPI(ctx context.Context, queueUuids []string, filter model.ConversationFilter, limit int, scrollId string) (total int, conversations []*model.ConversationView, respScrollId string, err error) {
	conversationIds := []string{}
	conversationFilter := model.AllocateUserFilter{
		TenantId: filter.TenantId,
	}
	if len(queueUuids) > 0 {
		conversationFilter.QueueId = queueUuids
	}
	if filter.IsDone.Bool {
		conversationFilter.MainAllocate = "deactive"
	} else {
		conversationFilter.MainAllocate = "active"
	}

	total, userAllocations, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, conversationFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		for _, item := range *userAllocations {
			conversationIds = append(conversationIds, item.ConversationId)
		}
	}
	filter.ConversationId = conversationIds
	total, conversations, respScrollId, err = repository.ConversationESRepo.SearchWithScroll(ctx, filter.TenantId, ES_INDEX_CONVERSATION, filter, limit, scrollId)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		for k, conv := range conversations {
			filter := model.MessageFilter{
				TenantId:       conv.TenantId,
				ConversationId: conv.ConversationId,
				IsRead:         "deactive",
				EventNameExlucde: []string{
					"received",
					"seen",
				},
			}
			total, _, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX_MESSAGE, filter, -1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			conv.TotalUnRead = int64(total)

			filterMessage := model.MessageFilter{
				ConversationId: conv.ConversationId,
			}
			_, message, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX_MESSAGE, filterMessage, 1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			if len(*message) > 0 {
				if slices.Contains[[]string](variables.ATTACHMENT_TYPE, (*message)[0].EventName) {
					conv.LatestMessageContent = (*message)[0].EventName
				} else {
					conv.LatestMessageContent = (*message)[0].Content
				}
				conv.LatestMessageDirection = (*message)[0].Direction
			}

			conversations[k] = conv
		}
	}
	return
}
