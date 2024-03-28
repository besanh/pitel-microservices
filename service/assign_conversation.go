package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IAssignConversation interface {
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
		AppId: data.AppId,
		OaId:  data.OaId,
	}
	total, connection, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	if total < 1 {
		log.Errorf("connection %s not found", (*connection)[0].Id)
		return response.ServiceUnavailableMsg("connection " + (*connection)[0].Id + " not found")
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

	return response.OK(result)
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

	// err = repository.ChatConnectionAppRepo.UpdateChatConnectionApp(ctx, dbCon, data)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.ServiceUnavailableMsg(err.Error())
	// }
	return response.OKResponse()
}
