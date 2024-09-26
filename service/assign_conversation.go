package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/tel4vn/pitel-microservices/common/cache"
	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/util"
	"github.com/tel4vn/pitel-microservices/common/variables"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IAssignConversation interface {
		GetUserAssigned(ctx context.Context, authUser *model.AuthUser, conversationId string, status string) (*model.AllocateUser, error)
		GetUserInQueue(ctx context.Context, authUser *model.AuthUser, data model.UserInQueueFilter) (result []model.ChatQueueUserView, err error)
		AllocateConversation(ctx context.Context, authUser *model.AuthUser, data *model.AssignConversation) error
	}
	AssignConversation struct{}
)

var AssignConversationService IAssignConversation

func NewAssignConversation() IAssignConversation {
	return &AssignConversation{}
}

func (s *AssignConversation) GetUserInQueue(ctx context.Context, authUser *model.AuthUser, data model.UserInQueueFilter) (result []model.ChatQueueUserView, err error) {
	filter := model.ChatConnectionAppFilter{
		TenantId:       authUser.TenantId,
		AppId:          data.AppId,
		OaId:           data.OaId,
		ConnectionType: data.ConversationType,
	}
	_, connections, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*connections) < 1 {
		log.Errorf("connection not found")
		err = errors.New("connection not found")
		return
	}

	// TODO: find connection_queue
	connectionQueueExist, err := repository.ConnectionQueueRepo.GetById(ctx, repository.DBConn, (*connections)[0].ConnectionQueueId)
	if err != nil {
		log.Error(err)
		return
	} else if connectionQueueExist == nil {
		log.Errorf("connection queue not found")
		err = errors.New("connection queue not found")
		return
	}

	filterChatManageQueueUser := model.ChatManageQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  connectionQueueExist.QueueId,
	}
	_, manageQueueUsers, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}

	var queueIds []string
	if len(*manageQueueUsers) > 0 {
		for _, item := range *manageQueueUsers {
			queueIds = append(queueIds, item.QueueId)
		}
	}

	filterUserInQueue := model.ChatQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  queueIds,
	}

	_, userInQueues, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, filterUserInQueue, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}

	result = []model.ChatQueueUserView{}
	if len(*userInQueues) > 0 {
		for _, item := range *userInQueues {
			tmp := model.ChatQueueUserView{}
			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				return
			}
			tmp.Id = item.Id
			result = append(result, tmp)
		}
	}
	if authUser.Level == "manager" || authUser.Level == "admin" {
		chatManageQueueUserFiler := model.ChatManageQueueUserFilter{
			TenantId: authUser.TenantId,
			QueueId:  connectionQueueExist.QueueId,
		}
		_, chatManageQueueUsers, errTmp := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, chatManageQueueUserFiler, 1, 0)
		if errTmp != nil {
			log.Error(errTmp)
			return
		}
		if len(*chatManageQueueUsers) > 0 {
			tmp := model.ChatQueueUserView{}
			if err = util.ParseAnyToAny((*chatManageQueueUsers)[0], &tmp); err != nil {
				log.Error(err)
				return
			}
			tmp.Id = (*chatManageQueueUsers)[0].Id
			result = append(result, tmp)
		}
	}

	return
}

func (s *AssignConversation) GetUserAssigned(ctx context.Context, authUser *model.AuthUser, conversationId string, status string) (result *model.AllocateUser, err error) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, "", conversationId)
	if err != nil {
		log.Error(err)
		return
	} else if conversationExist == nil {
		err = errors.New("conversation " + conversationId + " not found")
		log.Error(err)
		return
	}

	conversationFilter := model.AllocateUserFilter{
		TenantId:       authUser.TenantId,
		ConversationId: conversationId,
		MainAllocate:   status,
	}
	_, userAllocates, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, conversationFilter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*userAllocates) > 0 {
		result = &(*userAllocates)[0]
	}
	return
}

func (s *AssignConversation) AllocateConversation(ctx context.Context, authUser *model.AuthUser, data *model.AssignConversation) (err error) {
	filter := model.ConversationFilter{
		TenantId:       authUser.TenantId,
		ConversationId: []string{data.ConversationId},
	}
	_, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*conversations) < 1 {
		err = fmt.Errorf("conversation %s not found", data.ConversationId)
		log.Error(err)
		return
	}

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
		_, messages, err := repository.MessageESRepo.GetMessages(ctx, conv.TenantId, ES_INDEX_MESSAGE, filter, -1, 0)
		if err != nil {
			log.Error(err)
			break
		}
		conv.TotalUnRead = int64(len(*messages))

		filterMessage := model.MessageFilter{
			TenantId:       conv.TenantId,
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
		}
		conv.LatestMessageDirection = (*message)[0].Direction

		(*conversations)[k] = conv
	}

	conversationEvent := &model.ConversationView{}
	if err = util.ParseAnyToAny((*conversations)[0], conversationEvent); err != nil {
		log.Error(err)
		return
	}
	if conversationEvent.Labels != nil && !reflect.DeepEqual(conversationEvent.Labels, "") {
		var labels []map[string]string
		if err = json.Unmarshal([]byte(conversationEvent.Labels), &labels); err != nil {
			log.Error(err)
			return
		}
		chatLabelIds := []string{}
		if len(labels) > 0 {
			for _, item := range labels {
				chatLabelIds = append(chatLabelIds, item["label_id"])
			}
			if len(chatLabelIds) > 0 {
				_, chatLabelExist, errTmp := repository.ChatLabelRepo.GetChatLabels(ctx, repository.DBConn, model.ChatLabelFilter{
					LabelIds: chatLabelIds,
				}, -1, 0)
				if errTmp != nil {
					err = errTmp
					log.Error(err)
					return
				}
				if len(*chatLabelExist) > 0 {
					tmp, errTmp := json.Marshal((*chatLabelExist))
					if errTmp != nil {
						err = errTmp
						log.Error(err)
						return
					}
					conversationEvent.Labels = tmp
				}
			}
		}
	}

	allocateFilter := model.AllocateUserFilter{
		TenantId:               authUser.TenantId,
		ExternalConversationId: (*conversations)[0].ExternalConversationId,
		MainAllocate:           data.Status,
	}
	_, userAllocates, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, allocateFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}

	if len(*userAllocates) < 1 {
		userAllocate := model.AllocateUser{
			Base:                   model.InitBase(),
			TenantId:               (*conversations)[0].TenantId,
			AppId:                  (*conversations)[0].AppId,
			OaId:                   (*conversations)[0].OaId,
			ConversationId:         (*conversations)[0].ConversationId,
			ExternalConversationId: (*conversations)[0].ExternalConversationId,
			UserId:                 data.UserId,
			QueueId:                data.QueueId,
			AllocatedTimestamp:     time.Now().UnixMilli(),
			MainAllocate:           "active",
			ConnectionId:           (*conversations)[0].ConversationId,
		}
		log.Infof("conversation %s allocated to user %s", (*conversations)[0].ConversationId, data.UserId)
		if err = repository.AllocateUserRepo.Insert(ctx, repository.DBConn, userAllocate); err != nil {
			log.Error(err)
			return
		}
	}

	userIdAssigned := (*userAllocates)[0].UserId
	(*userAllocates)[0].UserId = data.UserId
	if err = repository.AllocateUserRepo.Update(ctx, repository.DBConn, (*userAllocates)[0]); err != nil {
		log.Error(err)
		return
	}

	// TODO: clear cache
	userAllocateCache := cache.RCache.Get(USER_ALLOCATE + "_" + GenerateConversationId((*userAllocates)[0].AppId, (*userAllocates)[0].OaId, (*userAllocates)[0].ConversationId))
	if userAllocateCache != nil {
		if err = cache.RCache.Del([]string{USER_ALLOCATE + "_" + GenerateConversationId((*userAllocates)[0].AppId, (*userAllocates)[0].OaId, (*userAllocates)[0].ConversationId)}); err != nil {
			log.Error(err)
			return
		}
	}

	if authUser.UserId != userIdAssigned {
		var subscribers []*Subscriber
		for s := range WsSubscribers.Subscribers {
			if s.TenantId == authUser.TenantId && s.Id == userIdAssigned {
				subscribers = append(subscribers, s)
				break
			}
		}

		PublishConversationToOneUser(variables.EVENT_CHAT["conversation_unassigned"], userIdAssigned, subscribers, true, conversationEvent)
	}

	// Event user_assigned
	userUuids := []string{}
	manageQueueUser, err := GetManageQueueUser(ctx, (*userAllocates)[0].QueueId)
	if err != nil {
		log.Error(err)
		return
	} else if len(manageQueueUser.Id) < 1 {
		log.Error("queue " + (*userAllocates)[0].QueueId + " not found")
		err = errors.New("queue " + (*userAllocates)[0].QueueId + " not found")
		return
	}
	for s := range WsSubscribers.Subscribers {
		if s.TenantId == manageQueueUser.TenantId && s.Level == "admin" {
			userUuids = append(userUuids, s.Id)
		}
		if s.TenantId == manageQueueUser.TenantId && manageQueueUser.UserId == s.Id && s.Level == "manager" && s.Id != authUser.UserId {
			userUuids = append(userUuids, s.Id)
		}
		if s.TenantId == manageQueueUser.TenantId && s.Id == data.UserId {
			userUuids = append(userUuids, s.Id)
		}
	}

	if len(userUuids) > 0 {
		filterMessage := model.MessageFilter{
			TenantId:       conversationEvent.TenantId,
			ConversationId: conversationEvent.ConversationId,
		}
		_, message, err := repository.MessageESRepo.GetMessages(ctx, conversationEvent.TenantId, ES_INDEX_MESSAGE, filterMessage, 1, 0)
		if err != nil {
			log.Error(err)
		}
		if len(*message) > 0 {
			if slices.Contains[[]string](variables.ATTACHMENT_TYPE, (*message)[0].EventName) {
				conversationEvent.LatestMessageContent = (*message)[0].EventName
			} else {
				conversationEvent.LatestMessageContent = (*message)[0].Content
			}
			conversationEvent.LatestMessageDirection = (*message)[0].Direction
		}

		PublishConversationToManyUser(variables.EVENT_CHAT["conversation_assigned"], userUuids, true, conversationEvent)
	}

	return
}
