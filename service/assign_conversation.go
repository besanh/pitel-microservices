package service

import (
	"context"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IAssignConversation interface {
		GetUserAssigned(ctx context.Context, authUser *model.AuthUser, conversationId string, status string) (int, any)
		GetUserInQueue(ctx context.Context, authUser *model.AuthUser, data model.UserInQueueFilter) (int, any)
		AllocateConversation(ctx context.Context, authUser *model.AuthUser, data *model.AssignConversation) (int, any)
	}
	AssignConversation struct{}
)

func NewAssignConversation() IAssignConversation {
	return &AssignConversation{}
}

func (s *AssignConversation) GetUserInQueue(ctx context.Context, authUser *model.AuthUser, data model.UserInQueueFilter) (int, any) {
	filter := model.ChatConnectionAppFilter{
		AppId:          data.AppId,
		OaId:           data.OaId,
		ConnectionType: data.ConversationType,
	}
	_, connections, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(*connections) < 1 {
		log.Errorf("connection not found")
		return response.ServiceUnavailableMsg("connection not found")
	}

	filterChatManageQueueUser := model.ChatManageQueueUserFilter{
		QueueId: (*connections)[0].QueueId,
	}
	_, manageQueueUsers, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	var queueId string
	if len(*manageQueueUsers) > 0 {
		queueId = (*manageQueueUsers)[0].QueueId
	}

	filterUserInQueue := model.ChatQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  []string{queueId},
	}

	if len(queueId) > 0 {
		filterChatManageQueueUser.QueueId = queueId
	}

	_, userInQueues, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, filterUserInQueue, -1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	result := []model.ChatQueueUserView{}
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
			conversationFilter := model.UserAllocateFilter{
				ConversationId: data.ConversationId,
				MainAllocate:   data.Status,
			}
			_, userAllocates, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, conversationFilter, -1, 0)

			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}

			if len(*userAllocates) < 1 {
				log.Errorf("conversation not found")
				return response.ServiceUnavailableMsg("conversation not found")
			}

			found := false
			for _, existing := range result {
				if (*userAllocates)[0].UserId == existing.UserId {
					found = true
					break
				}
			}
			if !found {
				result = append(result, model.ChatQueueUserView{
					TenantId: (*userAllocates)[0].TenantId,
					QueueId:  (*userAllocates)[0].QueueId,
					UserId:   (*userAllocates)[0].UserId,
				})
			}
		}
	}

	return response.OK(result)
}

func (s *AssignConversation) GetUserAssigned(ctx context.Context, authUser *model.AuthUser, conversationId string, status string) (int, any) {
	filter := model.ConversationFilter{
		ConversationId: []string{conversationId},
		TenantId:       authUser.TenantId,
	}
	_, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if len(*conversations) < 1 {
		log.Errorf("conversation %s not found", (*conversations)[0].ConversationId)
		return response.ServiceUnavailableMsg("conversation " + (*conversations)[0].ConversationId + " not found")
	}

	conversationFilter := model.UserAllocateFilter{
		ConversationId: (*conversations)[0].ConversationId,
		MainAllocate:   status,
	}
	_, userAllocates, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, conversationFilter, -1, 0)

	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if len(*userAllocates) < 1 {
		return response.OK(nil)
	}
	return response.OK((*userAllocates)[0])
}

func (s *AssignConversation) AllocateConversation(ctx context.Context, authUser *model.AuthUser, data *model.AssignConversation) (int, any) {
	filter := model.ConversationFilter{
		ConversationId: []string{data.ConversationId},
		TenantId:       authUser.TenantId,
	}
	_, conversations, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(*conversations) < 1 {
		log.Errorf("conversation not found")
		return response.ServiceUnavailableMsg("conversation not found")
	}

	allocateFilter := model.UserAllocateFilter{
		ConversationId: (*conversations)[0].ConversationId,
		MainAllocate:   data.Status,
	}
	_, userAllocates, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, allocateFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if len(*userAllocates) < 1 {
		userAllocate := model.UserAllocate{
			Base:               model.InitBase(),
			TenantId:           (*conversations)[0].TenantId,
			AppId:              (*conversations)[0].AppId,
			OaId:               (*conversations)[0].OaId,
			UserId:             data.UserId,
			QueueId:            data.QueueId,
			AllocatedTimestamp: time.Now().Unix(),
			MainAllocate:       "active",
			ConnectionId:       (*conversations)[0].ConversationId,
		}
		log.Infof("conversation %s allocated to user %s", (*conversations)[0].ConversationId, data.UserId)
		if err := repository.UserAllocateRepo.Insert(ctx, repository.DBConn, userAllocate); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg("can not insert to user allocate")
		}
	}

	userIdAssigned := (*userAllocates)[0].UserId
	(*userAllocates)[0].UserId = data.UserId
	if err := repository.UserAllocateRepo.Update(ctx, repository.DBConn, (*userAllocates)[0]); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg("can not update to user allocate")
	}

	// TODO: clear cache
	userAllocateCache := cache.RCache.Get(USER_ALLOCATE + "_" + GenerateConversationId((*userAllocates)[0].AppId, (*userAllocates)[0].OaId, (*userAllocates)[0].ConversationId))
	if userAllocateCache != nil {
		if err = cache.RCache.Del([]string{USER_ALLOCATE + "_" + GenerateConversationId((*userAllocates)[0].AppId, (*userAllocates)[0].OaId, (*userAllocates)[0].ConversationId)}); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
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
		var conversationEvent model.Conversation
		if err := util.ParseAnyToAny((*conversations)[0], &conversationEvent); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		go PublishConversationToOneUser(variables.EVENT_CHAT["conversation_unassigned"], userIdAssigned, subscribers, true, &conversationEvent)
	}

	// Event user_assigned
	userUuids := []string{}
	manageQueueUser, err := GetManageQueueUser(ctx, (*userAllocates)[0].QueueId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err)
	} else if len(manageQueueUser.Id) < 1 {
		log.Error("queue " + (*userAllocates)[0].QueueId + " not found")
		return response.ServiceUnavailableMsg("queue " + (*userAllocates)[0].QueueId + " not found")
	}
	for s := range WsSubscribers.Subscribers {
		if s.TenantId == manageQueueUser.TenantId && s.Level == "admin" {
			userUuids = append(userUuids, s.Id)
		}
		if s.TenantId == manageQueueUser.TenantId && manageQueueUser.ManageId == s.Id && s.Level == "manager" && s.Id != authUser.UserId {
			userUuids = append(userUuids, s.Id)
		}
		if s.TenantId == manageQueueUser.TenantId && s.Id == data.UserId {
			userUuids = append(userUuids, s.Id)
		}
	}

	if len(userUuids) > 0 {
		conversationEvent := model.Conversation{}
		if err := util.ParseAnyToAny((*conversations)[0], &conversationEvent); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}

		go PublishConversationToManyUser(variables.EVENT_CHAT["conversation_assigned"], userUuids, true, &conversationEvent)
	}

	return response.OKResponse()
}
