package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IAssignConversation interface {
		GetUserAssigned(ctx context.Context, authUser *model.AuthUser, conversationId string, status string) (*model.AllocateUser, error)
		GetUserInQueue(ctx context.Context, authUser *model.AuthUser, data model.UserInQueueFilter) (int, []model.ChatQueueUserView, error)
		AllocateConversation(ctx context.Context, authUser *model.AuthUser, data *model.AssignConversation) error
	}
	AssignConversation struct{}
)

var AssignConversationService IAssignConversation

func NewAssignConversation() IAssignConversation {
	return &AssignConversation{}
}

func (s *AssignConversation) GetUserInQueue(ctx context.Context, authUser *model.AuthUser, data model.UserInQueueFilter) (total int, result []model.ChatQueueUserView, err error) {
	filter := model.ChatConnectionAppFilter{
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
			result = append(result, model.ChatQueueUserView{
				TenantId: item.TenantId,
				QueueId:  item.QueueId,
				UserId:   item.UserId,
			})
		}
	}
	if authUser.Source == "authen" {
		if authUser.Level == "manager" || authUser.Level == "admin" {
			chatManageQueueUserFiler := model.ChatManageQueueUserFilter{
				TenantId: authUser.TenantId,
				QueueId:  connectionQueueExist.QueueId,
			}
			_, chatManageQueueUsers, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, chatManageQueueUserFiler, 1, 0)
			if err != nil {
				log.Error(err)
				return 0, nil, err
			}
			if len(*chatManageQueueUsers) > 0 {
				result = append(result, model.ChatQueueUserView{
					TenantId: (*chatManageQueueUsers)[0].TenantId,
					QueueId:  (*chatManageQueueUsers)[0].QueueId,
					UserId:   (*chatManageQueueUsers)[0].UserId,
				})
			}
		}
	}

	return len(result), result, nil
}

func (s *AssignConversation) GetUserAssigned(ctx context.Context, authUser *model.AuthUser, conversationId string, status string) (result *model.AllocateUser, err error) {
	filter := model.ConversationFilter{
		ConversationId: []string{conversationId},
		TenantId:       authUser.TenantId,
	}
	_, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}

	if len(*conversations) < 1 {
		log.Errorf("conversation %s not found", conversationId)
		err = fmt.Errorf("conversation %s not found", conversationId)
		return
	}

	conversationFilter := model.AllocateUserFilter{
		ConversationId: (*conversations)[0].ConversationId,
		MainAllocate:   status,
	}
	_, userAllocates, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, conversationFilter, -1, 0)

	if err != nil {
		log.Error(err)
		return
	}

	if len(*userAllocates) < 1 {
		return nil, nil
	}
	return &(*userAllocates)[0], nil
}

func (s *AssignConversation) AllocateConversation(ctx context.Context, authUser *model.AuthUser, data *model.AssignConversation) (err error) {
	filter := model.ConversationFilter{
		ConversationId: []string{data.ConversationId},
		TenantId:       authUser.TenantId,
	}
	_, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*conversations) < 1 {
		log.Errorf("conversation not found")
		err = errors.New("conversation not found")
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
		}
		conv.LatestMessageDirection = (*message)[0].Direction

		(*conversations)[k] = conv
	}

	allocateFilter := model.AllocateUserFilter{
		ConversationId: (*conversations)[0].ConversationId,
		MainAllocate:   data.Status,
	}
	_, userAllocates, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, allocateFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}

	if len(*userAllocates) < 1 {
		userAllocate := model.AllocateUser{
			Base:               model.InitBase(),
			TenantId:           (*conversations)[0].TenantId,
			AppId:              (*conversations)[0].AppId,
			OaId:               (*conversations)[0].OaId,
			UserId:             data.UserId,
			QueueId:            data.QueueId,
			AllocatedTimestamp: time.Now().UnixMilli(),
			MainAllocate:       "active",
			ConnectionId:       (*conversations)[0].ConversationId,
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
		var conversationEvent model.ConversationView
		if err = util.ParseAnyToAny((*conversations)[0], &conversationEvent); err != nil {
			log.Error(err)
			return
		}
		PublishConversationToOneUser(variables.EVENT_CHAT["conversation_unassigned"], userIdAssigned, subscribers, true, &conversationEvent)

		// TODO: publish message
		// filterMessage := model.MessageFilter{
		// 	TenantId:       (*conversations)[0].TenantId,
		// 	ConversationId: (*conversations)[0].ConversationId,
		// }
		// _, messages, err := repository.MessageESRepo.GetMessages(ctx, (*conversations)[0].TenantId, ES_INDEX, filterMessage, 1, 0)
		// if err != nil {
		// 	log.Error(err)
		// 	return response.ServiceUnavailableMsg(err.Error())
		// }
		// if len(*messages) > 0 {
		// 	PublishMessageToOneUser(variables.EVENT_CHAT["message_created"], userIdAssigned, subscribers, &(*messages)[0])
		// }
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
		conversationEvent := model.ConversationView{}
		if err = util.ParseAnyToAny((*conversations)[0], &conversationEvent); err != nil {
			log.Error(err)
			return
		}

		filterMessage := model.MessageFilter{
			TenantId:       conversationEvent.TenantId,
			ConversationId: conversationEvent.ConversationId,
		}
		_, message, err := repository.MessageESRepo.GetMessages(ctx, conversationEvent.TenantId, ES_INDEX, filterMessage, 1, 0)
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

		PublishConversationToManyUser(variables.EVENT_CHAT["conversation_assigned"], userUuids, true, &conversationEvent)

		// TODO: publish message
		// filterMessage := model.MessageFilter{
		// 	TenantId:       (*conversations)[0].TenantId,
		// 	ConversationId: (*conversations)[0].ConversationId,
		// }
		// _, messages, err := repository.MessageESRepo.GetMessages(ctx, (*conversations)[0].TenantId, ES_INDEX, filterMessage, 1, 0)
		// if err != nil {
		// 	log.Error(err)
		// 	return response.ServiceUnavailableMsg(err.Error())
		// }
		// if len(*messages) > 0 {
		// 	PublishMessageToManyUser(variables.EVENT_CHAT["message_created"], userUuids, &(*messages)[0])
		// }
	}

	return
}
