package service

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

func (s *Conversation) GetConversationsByManager(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any) {
	if authUser.Source == "authen" {
		if authUser.Level != "manager" && authUser.Level != "admin" {
			return response.Pagination(nil, 0, limit, offset)
		}
		url := API_CRM + "/v1/crm/user-crm"
		if authUser.Level == "manager" {
			url += "?level=user&unit_uuid=" + authUser.UnitUuid + "&limit=-1&offset=0"
		} else {
			url += "?limit=-1&offset=0"
		}
		client := resty.New()
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+authUser.Token).
			// SetHeader("Authorization", "Bearer "+token).
			Get(url)

		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		if res.StatusCode() == 200 {
			var responseData model.ResponseData
			err = json.Unmarshal(res.Body(), &responseData)
			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}

			userUuids := []string{
				authUser.UserId,
			}

			for _, item := range responseData.Data {
				userUuid, ok := item["user_uuid"].(string)
				if !ok {
					log.Error("user_uuid not found or not a string")
					continue
				}
				userUuids = append(userUuids, userUuid)
			}
			if len(userUuids) < 1 {
				log.Error("list user not found")
				return response.Pagination(nil, 0, limit, offset)
			}

			conversationIds := []string{}
			conversationFilter := model.AgentAllocateFilter{
				AgentId:      userUuids,
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

			return response.Pagination(conversations, total, limit, offset)
		} else {
			return response.ServiceUnavailableMsg("Can not get user crm")
		}
	} else {
		return response.Pagination(nil, 0, limit, offset)
	}
}
