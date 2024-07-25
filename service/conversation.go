package service

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IConversation interface {
		InsertConversation(ctx context.Context, conversation model.Conversation) (id string, err error)
		GetConversations(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (total int, result []model.ConversationCustomView, err error)
		GetConversationsWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit int, scrollId string) (int, any)
		GetConversationsByHighLevel(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (int, any)
		GetConversationsByHighLevelWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit int, scrollId string) (int, any)
		UpdateConversationById(ctx context.Context, authUser *model.AuthUser, appId, oaId, id string, data model.ShareInfo) (int, any)
		UpdateStatusConversation(ctx context.Context, authUser *model.AuthUser, appId, id, updatedBy, status string) error
		GetConversationById(ctx context.Context, authUser *model.AuthUser, appId, conversationId string) (int, any)
		UpdateUserPreferenceConversation(ctx context.Context, authUser *model.AuthUser, preferRequest model.ConversationPreferenceRequest) error
	}
	Conversation struct{}
)

var ConversationService IConversation

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

func (s *Conversation) GetConversations(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit, offset int) (total int, result []model.ConversationCustomView, err error) {
	conversationIds := []string{}
	conversationFilter := model.AllocateUserFilter{
		TenantId: authUser.TenantId,
		UserId:   []string{authUser.UserId},
	}
	if filter.IsDone.Valid && !filter.IsDone.Bool {
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
	if len(conversationIds) < 1 {
		log.Error("list conversation not found")
		return
	}
	filter.ConversationId = conversationIds
	filter.TenantId = authUser.TenantId
	_, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	var conversationCustomViews []model.ConversationCustomView
	if len(*conversations) > 0 {
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
			_, messages, errTmp := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX, filter, -1, 0)
			if errTmp != nil {
				log.Error(errTmp)
				break
			}
			conv.TotalUnRead = int64(len(*messages))

			filterMessage := model.MessageFilter{
				TenantId:       conv.TenantId,
				ConversationId: conv.ConversationId,
			}
			_, message, errTmp := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX, filterMessage, 1, 0)
			if errTmp != nil {
				log.Error(errTmp)
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

			// TODO: parse label
			var conversationCustomView model.ConversationCustomView
			if err = util.ParseAnyToAny((*conversations)[k], &conversationCustomView); err != nil {
				log.Error(err)
				return
			}
			labels := []any{}
			if err = json.Unmarshal([]byte((*conversations)[k].Label), &labels); err != nil {
				log.Error(err)
				return
			}
			if len(labels) > 0 {
				chatLabelIds := []string{}
				for _, item := range labels {
					var tmp map[string]string
					if err := util.ParseAnyToAny(item, &tmp); err != nil {
						log.Error(err)
						continue
					}
					chatLabelIds = append(chatLabelIds, tmp["label_id"])
				}
				if len(chatLabelIds) > 0 {
					_, chatLabelExist, errTmp := repository.ChatLabelRepo.GetChatLabels(ctx, repository.DBConn, model.ChatLabelFilter{
						LabelIds: chatLabelIds,
					}, -1, 0)
					if errTmp != nil {
						log.Error(errTmp)
						return
					}
					if len(*chatLabelExist) > 0 {
						conversationCustomView.Label = chatLabelExist
					}
				}
			}
			conversationCustomViews = append(conversationCustomViews, conversationCustomView)
		}
	}
	result = conversationCustomViews
	return
}

func (s *Conversation) GetConversationsWithScrollAPI(ctx context.Context, authUser *model.AuthUser, filter model.ConversationFilter, limit int, scrollId string) (int, any) {
	conversationIds := []string{}
	conversationFilter := model.AllocateUserFilter{
		TenantId: authUser.TenantId,
		UserId:   []string{authUser.UserId},
	}
	if filter.IsDone.Valid {
		conversationFilter.MainAllocate = "deactive"
	} else {
		conversationFilter.MainAllocate = "active"
	}
	total, userAllocations, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, conversationFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if total > 0 {
		for _, item := range *userAllocations {
			conversationIds = append(conversationIds, item.ConversationId)
		}
	}
	if len(conversationIds) < 1 {
		log.Error("list conversation not found")
		result := map[string]any{
			"conversations": nil,
			"scroll_id":     "",
		}
		return response.Pagination(result, 0, limit, 0)
	}
	filter.ExternalConversationId = conversationIds
	filter.TenantId = authUser.TenantId
	_, conversations, respScrollId, err := repository.ConversationESRepo.SearchWithScroll(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, limit, scrollId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	conversationCustomViews := make([]model.ConversationCustomView, 0)
	if len(conversations) > 0 {
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
			_, messages, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX, filter, -1, 0)
			if err != nil {
				log.Error(err)
				break
			}
			conv.TotalUnRead = int64(len(*messages))

			filterMessage := model.MessageFilter{
				TenantId:       conv.TenantId,
				ConversationId: conv.ConversationId,
			}
			_, message, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX, filterMessage, 1, 0)
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

			// TODO: parse label
			var conversationCustomView model.ConversationCustomView
			if err := util.ParseAnyToAny(conversations[k], &conversationCustomView); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}

			if !reflect.DeepEqual(conversations[k].Label, "") {
				var labels []map[string]string
				if err = json.Unmarshal([]byte(conversations[k].Label), &labels); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				chatLabelIds := []string{}
				if len(labels) > 0 {
					for _, item := range labels {
						chatLabelIds = append(chatLabelIds, item["label_id"])
					}
					if len(chatLabelIds) > 0 {
						_, chatLabelExist, err := repository.ChatLabelRepo.GetChatLabels(ctx, repository.DBConn, model.ChatLabelFilter{
							LabelIds: chatLabelIds,
						}, -1, 0)
						if err != nil {
							log.Error(err)
							return response.ServiceUnavailableMsg(err.Error())
						}
						if len(*chatLabelExist) > 0 {
							conversationCustomView.Label = chatLabelExist
						}
					}
				}
			}
			conversationCustomViews = append(conversationCustomViews, conversationCustomView)
		}
	}
	result := map[string]any{
		"conversations": conversationCustomViews,
		"scroll_id":     respScrollId,
	}
	return response.Pagination(result, len(conversationCustomViews), limit, 0)
}

func (s *Conversation) UpdateConversationById(ctx context.Context, authUser *model.AuthUser, appId, oaId, conversationId string, data model.ShareInfo) (int, any) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, appId, conversationId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found with app_id %s", conversationId, appId)
		return response.NotFoundMsg("conversation " + conversationId + " not found")
	}
	conversationExist.Username = data.Fullname
	conversationExist.ShareInfo = &data
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	esDoc := map[string]any{}
	if err = json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	// IMPROVE: should we update direction es or use queue ?
	if err = repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, appId, conversationId, esDoc); err != nil {
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

	// Update User allocate
	filter := model.AllocateUserFilter{
		AppId:                  appId,
		ExternalConversationId: conversationId,
		MainAllocate:           statusAllocate,
	}
	_, allocateUser, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(*allocateUser) < 1 {
		log.Errorf("conversation %s not found with active user", conversationId)
		return errors.New("conversation " + conversationId + " not found with active user")
	}

	allocateUserTmp := (*allocateUser)[0]

	if status == "done" {
		allocateUserTmp.MainAllocate = "deactive"
		allocateUserTmp.AllocatedTimestamp = time.Now().UnixMilli()
		allocateUserTmp.UpdatedAt = time.Now()
		if err := repository.AllocateUserRepo.Update(ctx, repository.DBConn, allocateUserTmp); err != nil {
			log.Error(err)
			return err
		}
		conversationExist.IsDone = true
		conversationExist.IsDoneBy = updatedBy
		conversationExist.IsDoneAt = time.Now()
	} else if status == "reopen" {
		allocateUserTmp.MainAllocate = "active"
		allocateUserTmp.AllocatedTimestamp = time.Now().UnixMilli()
		allocateUserTmp.UpdatedAt = time.Now()
		if err := repository.AllocateUserRepo.Update(ctx, repository.DBConn, allocateUserTmp); err != nil {
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

		conversationQueue := model.ConversationQueue{
			DocId:        conversationExist.ConversationId,
			Conversation: *conversationExist,
		}
		if err = PublishPutConversationToChatQueue(ctx, conversationQueue); err != nil {
			log.Error(err)
			if err := repository.AllocateUserRepo.Update(ctx, repository.DBConn, (*allocateUser)[0]); err != nil {
				log.Error(err)
			}
			return err
		}
	}

	// TODO: clear cache
	allocateUserCache := cache.RCache.Get(USER_ALLOCATE + "_" + GenerateConversationId(conversationExist.AppId, conversationExist.OaId, conversationExist.ExternalUserId))
	if allocateUserCache != nil {
		if err = cache.RCache.Del([]string{USER_ALLOCATE + "_" + GenerateConversationId(conversationExist.AppId, conversationExist.OaId, conversationExist.ExternalUserId)}); err != nil {
			log.Error(err)
			return err
		}
	}

	conversationConverted := &model.ConversationView{}
	if err := util.ParseAnyToAny(conversationExist, conversationConverted); err != nil {
		log.Error(err)
		return err
	}

	// TODO: get message to display, otherwise use api get conversation to get latest message
	filterMessage := model.MessageFilter{
		TenantId:       conversationExist.TenantId,
		ConversationId: conversationExist.ConversationId,
		// IsRead:         "deactive",
		// EventNameExlucde: []string{
		// 	"received",
		// 	"seen",
		// },
	}
	_, messages, err := repository.MessageESRepo.GetMessages(ctx, conversationExist.TenantId, ES_INDEX, filterMessage, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(*messages) > 0 {
		conversationConverted.LatestMessageContent = (*messages)[0].Content
		for _, item := range *messages {
			if item.IsRead == "deactive" {
				conversationConverted.TotalUnRead += 1
			}
		}
	}

	manageQueueUser, err := GetManageQueueUser(ctx, allocateUserTmp.QueueId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(manageQueueUser.Id) < 1 {
		log.Error("queue " + allocateUserTmp.QueueId + " not found")
		return errors.New("queue " + allocateUserTmp.QueueId + " not found")
	}

	var subscribers []*Subscriber
	var subscriberAdmins []string
	var subscriberManagers []string
	for s := range WsSubscribers.Subscribers {
		if s.TenantId == authUser.TenantId {
			subscribers = append(subscribers, s)
			if s.Level == "admin" {
				subscriberAdmins = append(subscriberAdmins, s.Id)
			}
			if s.Level == "manager" {
				subscriberManagers = append(subscriberManagers, s.Id)
			}
		}
	}

	if slices.Contains(subscriberManagers, authUser.UserId) {
		if status == "done" {
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_done"], authUser.UserId, subscribers, true, conversationConverted)
		} else if status == "reopen" {
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_reopen"], authUser.UserId, subscribers, true, conversationConverted)
		}
	}

	// Event to manager
	isExist := BinarySearchSlice(manageQueueUser.UserId, subscriberManagers)
	if isExist && (manageQueueUser.UserId != conversationExist.IsDoneBy) {
		if status == "done" {
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_done"], manageQueueUser.UserId, subscribers, true, conversationConverted)
		} else if status == "reopen" {
			PublishConversationToOneUser(variables.EVENT_CHAT["conversation_reopen"], manageQueueUser.UserId, subscribers, true, conversationConverted)
		}

		// PublishMessageToOneUser(variables.EVENT_CHAT["message_created"], manageQueueUser.ManageId, subscribers, &(*messages)[0])
	}

	// Event to admin
	if ENABLE_PUBLISH_ADMIN && len(subscriberAdmins) > 0 {
		if status == "done" {
			PublishConversationToManyUser(variables.EVENT_CHAT["conversation_done"], subscriberAdmins, true, conversationConverted)
		} else if status == "reopen" {
			PublishConversationToManyUser(variables.EVENT_CHAT["conversation_reopen"], subscriberAdmins, true, conversationConverted)
		}

		// PublishMessageToManyUser(variables.EVENT_CHAT["message_created"], subscriberAdmins, &(*messages)[0])
	}

	return nil
}

func (s *Conversation) GetConversationById(ctx context.Context, authUser *model.AuthUser, appId, conversationId string) (int, any) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, appId, conversationId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(conversationExist.ConversationId) < 1 {
		err = errors.New("conversation " + conversationId + " with app_id " + appId + " not found")
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if !reflect.DeepEqual(conversationExist.Labels, "") {
		var labels []map[string]string
		if err = json.Unmarshal([]byte(conversationExist.Labels), &labels); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		chatLabelIds := []string{}
		if len(labels) > 0 {
			for _, item := range labels {
				chatLabelIds = append(chatLabelIds, item["label_id"])
			}
			if len(chatLabelIds) > 0 {
				_, chatLabelExist, err := repository.ChatLabelRepo.GetChatLabels(ctx, repository.DBConn, model.ChatLabelFilter{
					LabelIds: chatLabelIds,
				}, -1, 0)
				if err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				if len(*chatLabelExist) > 0 {
					tmp, err := json.Marshal(*chatLabelExist)
					if err != nil {
						log.Error(err)
						return response.ServiceUnavailableMsg(err.Error())
					}
					conversationExist.Labels = tmp
				} else {
					conversationExist.Labels = []byte("[]")
				}
			} else {
				conversationExist.Labels = []byte("[]")
			}
		} else {
			conversationExist.Labels = []byte("[]")
		}
	} else {
		conversationExist.Labels = []byte("[]")
	}

	return response.OK(conversationExist)
}

func (s *Conversation) UpdateUserPreferenceConversation(ctx context.Context, authUser *model.AuthUser, preferRequest model.ConversationPreferenceRequest) error {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, preferRequest.AppId, preferRequest.ConversationId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found", preferRequest.ConversationId)
		return errors.New("conversation " + preferRequest.ConversationId + " not found")
	}

	filter := model.AllocateUserFilter{
		AppId:                  preferRequest.AppId,
		ExternalConversationId: preferRequest.ConversationId,
		MainAllocate:           "active",
	}

	tmp, _ := strconv.ParseBool(preferRequest.PreferenceValue)
	switch preferRequest.PreferenceType {
	case "major":
		conversationExist.Major = tmp
	case "following":
		conversationExist.Following = tmp
	default:
	}

	conversationQueue := model.ConversationQueue{
		DocId:        conversationExist.ConversationId,
		Conversation: *conversationExist,
	}
	if err = PublishPutConversationToChatQueue(ctx, conversationQueue); err != nil {
		log.Error(err)
		return err
	}

	conversationConverted := &model.ConversationView{}
	if err = util.ParseAnyToAny(conversationExist, conversationConverted); err != nil {
		log.Error(err)
		return err
	}

	_, allocateUser, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(*allocateUser) > 0 {
		allocateUserTmp := (*allocateUser)[0]
		// Event to manager
		manageQueueUser, err := GetManageQueueUser(ctx, allocateUserTmp.QueueId)
		if err != nil {
			log.Error(err)
			return err
		} else if len(manageQueueUser.Id) < 1 {
			log.Error("queue " + allocateUserTmp.QueueId + " not found")
			return errors.New("queue " + allocateUserTmp.QueueId + " not found")
		}
		s.publishConversationEventToManagerAndAdmin(authUser, manageQueueUser, variables.PREFERENCE_EVENT[preferRequest.PreferenceType], conversationConverted)
	} else {
		s.publishConversationEventToManagerAndAdmin(authUser, nil, variables.PREFERENCE_EVENT[preferRequest.PreferenceType], conversationConverted)
	}

	return nil
}
