package service

import (
	"context"
	"sync"
	"time"

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
	var queueId string
	filter := model.ChatConnectionAppFilter{
		AppId:          data.AppId,
		OaId:           data.OaId,
		ConnectionType: data.ConversationType,
	}
	total, connection, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if total < 1 {
		log.Errorf("connection not found")
		return response.ServiceUnavailableMsg("connection not found")
	}

	filterChatManageQueueUser := model.ChatManageQueueUserFilter{
		QueueId: (*connection)[0].QueueId,
	}
	totalManageQueueUser, manageQueueUser, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, filterChatManageQueueUser, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if totalManageQueueUser > 0 {
		queueId = (*manageQueueUser)[0].QueueId
	}

	filterUserInQueue := model.ChatQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  []string{queueId},
	}

	if len(queueId) > 0 {
		filterChatManageQueueUser.QueueId = queueId
	}

	totalUserInQueue, userInQueue, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, filterUserInQueue, -1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	result := []model.ChatQueueUserView{}
	if totalUserInQueue > 0 {
		for _, item := range *userInQueue {
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
			total, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, conversationFilter, -1, 0)

			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}

			if total < 1 {
				log.Errorf("conversation not found")
				return response.ServiceUnavailableMsg("conversation not found")
			}
			manageQueueUser, err := GetManageQueueUser(ctx, (*userAllocations)[0].QueueId)
			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err)
			} else if len(manageQueueUser.Id) < 1 {
				log.Error("queue " + (*userAllocations)[0].QueueId + " not found")
			}

			result = append(result, model.ChatQueueUserView{
				TenantId: manageQueueUser.TenantId,
				QueueId:  manageQueueUser.QueueId,
				UserId:   manageQueueUser.ManageId,
			})

		}
	}

	return response.OK(result)
}

func (s *AssignConversation) GetUserAssigned(ctx context.Context, authUser *model.AuthUser, conversationId string, status string) (int, any) {
	filter := model.ConversationFilter{
		ConversationId: []string{conversationId},
		TenantId:       authUser.TenantId,
	}
	totalConversation, conversation, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if totalConversation < 1 {
		log.Errorf("conversation %s not found", (*conversation)[0].ConversationId)
		return response.ServiceUnavailableMsg("conversation " + (*conversation)[0].ConversationId + " not found")
	}

	conversationFilter := model.UserAllocateFilter{
		ConversationId: (*conversation)[0].ConversationId,
		MainAllocate:   status,
	}
	total, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, conversationFilter, -1, 0)

	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if total < 1 {
		return response.OK(nil)
	}
	return response.OK((*userAllocations)[0])
}

func (s *AssignConversation) AllocateConversation(ctx context.Context, authUser *model.AuthUser, data *model.AssignConversation) (int, any) {
	filter := model.ConversationFilter{
		ConversationId: []string{data.ConversationId},
		TenantId:       authUser.TenantId,
	}
	totalConversation, conversation, err := repository.ConversationESRepo.GetConversations(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if totalConversation < 1 {
		log.Errorf("conversation %s not found", (*conversation)[0].ConversationId)
		return response.ServiceUnavailableMsg("conversation " + (*conversation)[0].ConversationId + " not found")
	}

	allocateFilter := model.UserAllocateFilter{
		ConversationId: (*conversation)[0].ConversationId,
		MainAllocate:   data.Status,
	}

	total, userAllocations, err := repository.UserAllocateRepo.GetUserAllocates(ctx, repository.DBConn, allocateFilter, -1, 0)

	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	if total < 1 {
		userAllocation := model.UserAllocate{
			Base:               model.InitBase(),
			TenantId:           (*conversation)[0].TenantId,
			AppId:              (*conversation)[0].AppId,
			UserId:             data.UserId,
			QueueId:            data.QueueId,
			AllocatedTimestamp: time.Now().Unix(),
			MainAllocate:       "active",
			ConnectionId:       (*conversation)[0].ConversationId,
		}
		log.Infof("conversation %s allocated to user %s", (*conversation)[0].ConversationId, data.UserId)
		if err := repository.UserAllocateRepo.Insert(ctx, repository.DBConn, userAllocation); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg("Can not insert to user allocation")
		}
	}

	if userAllocations != nil {
		userIdAssigned := (*userAllocations)[0].UserId
		(*userAllocations)[0].UserId = data.UserId
		if err := repository.UserAllocateRepo.Update(ctx, repository.DBConn, (*userAllocations)[0]); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg("Can not update to user allocation")
		}

		if authUser.UserId != userIdAssigned {
			var subscribers []*Subscriber
			for s := range WsSubscribers.Subscribers {
				if s.TenantId == authUser.TenantId && s.Id == userIdAssigned {
					subscribers = append(subscribers, s)
				}
			}
			var conversationEvent model.Conversation
			if err := util.ParseAnyToAny((*conversation)[0], &conversationEvent); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
			var mu sync.Mutex
			mu.Lock()
			go PublishConversationToOneUser(variables.EVENT_CHAT["conversation_unassigned"], userIdAssigned, subscribers, true, &conversationEvent)
			mu.Unlock()
		}
	}

	// Event user_assigned
	userUuids := []string{}
	manageQueueUser, err := GetManageQueueUser(ctx, (*userAllocations)[0].QueueId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err)
	} else if len(manageQueueUser.Id) < 1 {
		log.Error("queue " + (*userAllocations)[0].QueueId + " not found")
		return response.ServiceUnavailableMsg("queue " + (*userAllocations)[0].QueueId + " not found")
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
		if err := util.ParseAnyToAny((*conversation)[0], &conversationEvent); err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		var mu sync.Mutex
		mu.Lock()
		go PublishConversationToManyUser(variables.EVENT_CHAT["conversation_assigned"], userUuids, true, &conversationEvent)
		mu.Unlock()
	}

	return response.OKResponse()
}
