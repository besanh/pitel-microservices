package service

import (
	"context"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatQueue interface {
		InsertChatQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequest) (string, error)
		GetChatQueues(ctx context.Context, authUser *model.AuthUser, filter model.QueueFilter, limit, offset int) (int, *[]model.ChatQueue, error)
		GetChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatQueue, error)
		UpdateChatQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatQueueRequest) error
		DeleteChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) error

		//new version
		InsertChatQueueV2(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequestV2) (string, error)
		UpdateChatQueueByIdV2(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatQueueRequestV2) error
		DeleteChatQueueByIdV2(ctx context.Context, authUser *model.AuthUser, id string) error
		UpdateChatQueueStatus(ctx context.Context, authUser *model.AuthUser, id string, status string) error
	}
	ChatQueue struct{}
)

var ChatQueueService IChatQueue

func NewChatQueue() IChatQueue {
	return &ChatQueue{}
}

func (s *ChatQueue) InsertChatQueue(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequest) (string, error) {
	chatQueue := model.ChatQueue{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}
	routingExist, err := repository.ChatRoutingRepo.GetById(ctx, dbCon, data.ChatRoutingId)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	} else if routingExist == nil {
		err = errors.New("chat routing not found")
		return chatQueue.Base.GetId(), err
	}

	connectionUsers := []model.ConnectionQueue{}
	if len(data.ConnectionId) > 0 {
		for _, item := range data.ConnectionId {
			connectionUser := model.ConnectionQueue{
				Base:         model.InitBase(),
				TenantId:     authUser.TenantId,
				ConnectionId: item,
				QueueId:      chatQueue.Base.GetId(),
			}
			connectionUsers = append(connectionUsers, connectionUser)
		}
	}

	if len(connectionUsers) > 0 {
		if err = repository.ConnectionQueueRepo.BulkInsert(ctx, dbCon, connectionUsers); err != nil {
			log.Error(err)
			return chatQueue.Base.GetId(), err
		}
	}

	chatQueue.TenantId = authUser.TenantId
	chatQueue.QueueName = data.QueueName
	chatQueue.Description = data.Description
	chatQueue.ChatRoutingId = data.ChatRoutingId
	chatQueue.Status = data.Status

	if err = repository.ChatQueueRepo.Insert(ctx, dbCon, chatQueue); err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}
	return chatQueue.Base.GetId(), nil
}

func (s *ChatQueue) GetChatQueues(ctx context.Context, authUser *model.AuthUser, filter model.QueueFilter, limit, offset int) (int, *[]model.ChatQueue, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	filter.TenantId = authUser.TenantId

	total, queues, err := repository.ChatQueueRepo.GetQueues(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, queues, nil
}

func (s *ChatQueue) GetChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatQueue, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err == ERR_EMPTY_CONN {
		log.Error(err)
		return nil, err
	}
	queue, err := repository.ChatQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if queue == nil {
		log.Error("queue " + id + " not found")
		return nil, errors.New("queue " + id + " not found")
	}

	return queue, nil
}

func (s *ChatQueue) UpdateChatQueueById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatQueueRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	queueExist, err := repository.ChatQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	queueExist.QueueName = data.QueueName
	queueExist.Description = data.Description
	queueExist.ChatRoutingId = data.ChatRoutingId
	queueExist.Status = data.Status
	err = repository.ChatQueueRepo.Update(ctx, dbCon, *queueExist)
	if err != nil {
		log.Error(err)
		return err
	}

	// filter := model.ConnectionQueueFilter{
	// 	QueueId: queueExist.Id,
	// }
	// _, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, dbCon, filter, -1, 0)
	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }

	// if len(*connectionQueues) > 0 {
	// 	for _, item := range *connectionQueues {
	// 		if err := repository.ConnectionQueueRepo.Delete(ctx, dbCon, item.Id); err != nil {
	// 			log.Error(err)
	// 			return err
	// 		}
	// 	}
	// }

	if len(data.ConnectionId) > 0 {
		connectionUsers := []model.ConnectionQueue{}
		for _, item := range data.ConnectionId {
			connectionUser := model.ConnectionQueue{
				Base:         model.InitBase(),
				TenantId:     authUser.TenantId,
				ConnectionId: item,
				QueueId:      queueExist.Id,
			}
			connectionUsers = append(connectionUsers, connectionUser)
		}
		// TODO: clear cache
		chatQueueUserCache := cache.RCache.Get(CHAT_QUEUE + "_" + id)
		if chatQueueUserCache != nil {
			if err = cache.RCache.Del([]string{CHAT_QUEUE + "_" + id}); err != nil {
				log.Error(err)
				return err
			}
		}
		if err = repository.ConnectionQueueRepo.BulkInsert(ctx, dbCon, connectionUsers); err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func (s *ChatQueue) DeleteChatQueueById(ctx context.Context, authUser *model.AuthUser, id string) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = repository.ChatQueueRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ChatQueueRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// Delete queue user
	if err := repository.ChatQueueUserRepo.DeleteChatQueueUsers(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatQueue) InsertChatQueueV2(ctx context.Context, authUser *model.AuthUser, data model.ChatQueueRequestV2) (string, error) {
	chatQueue := model.ChatQueue{
		Base: model.InitBase(),
	}
	routingExist, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, data.ChatQueue.ChatRoutingId)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	} else if routingExist == nil {
		err = errors.New("chat routing not found")
		return chatQueue.Base.GetId(), err
	}

	tx, err := repository.ChatQueueRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}
	defer tx.Rollback()

	chatQueue.TenantId = authUser.TenantId
	chatQueue.QueueName = data.ChatQueue.QueueName
	chatQueue.Description = data.ChatQueue.Description
	chatQueue.ChatRoutingId = data.ChatQueue.ChatRoutingId
	chatQueue.Status = data.ChatQueue.Status

	if err = repository.ChatQueueRepo.TxInsert(ctx, tx, chatQueue); err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}

	// insert queue user
	chatQueueUsers := make([]model.ChatQueueUser, 0)
	for _, item := range data.ChatQueueUser.UserId {
		chatQueueUser := model.ChatQueueUser{
			Base:     model.InitBase(),
			TenantId: authUser.TenantId,
			QueueId:  chatQueue.GetId(),
			UserId:   item,
			Source:   authUser.Source,
		}
		chatQueueUsers = append(chatQueueUsers, chatQueueUser)
	}
	if len(chatQueueUsers) > 0 {
		err = repository.ChatQueueUserRepo.TxBulkInsert(ctx, tx, chatQueueUsers)
		if err != nil {
			log.Error(err)
			return chatQueue.Base.GetId(), err
		}
	}

	// insert manage queue user
	manageQueue := model.ChatManageQueueUser{
		Base: model.InitBase(),
	}
	manageQueue.TenantId = authUser.TenantId
	//manageQueue.ConnectionId = connectionAppExist.Id //no data
	manageQueue.QueueId = chatQueue.GetId()
	manageQueue.UserId = data.ChatManageQueueUser.UserId

	chatQueue.ManageQueueId = manageQueue.GetId()

	if err = repository.ManageQueueRepo.TxInsert(ctx, tx, manageQueue); err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}
	if err = repository.ChatQueueRepo.TxUpdate(ctx, tx, chatQueue); err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}

	if err = repository.ChatQueueRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return chatQueue.Base.GetId(), err
	}

	return chatQueue.Base.GetId(), nil
}

func (s *ChatQueue) UpdateChatQueueByIdV2(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatQueueRequestV2) error {
	queueExist, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	} else if queueExist == nil {
		log.Error("chat queue " + id + " not found")
		return errors.New("chat queue " + id + " not found")
	}

	tx, err := repository.ChatQueueRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback()

	queueExist.QueueName = data.ChatQueue.QueueName
	queueExist.Description = data.ChatQueue.Description
	queueExist.ChatRoutingId = data.ChatQueue.ChatRoutingId
	queueExist.Status = data.ChatQueue.Status

	// remove old queue user
	queueUserFilter := model.ChatQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  []string{queueExist.GetId()},
	}
	_, oldQueueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, queueUserFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(*oldQueueUsers) > 0 {
		if err = repository.ChatQueueUserRepo.TxBulkDelete(ctx, tx, *oldQueueUsers); err != nil {
			log.Error(err)
			return err
		}
	}

	// insert new queue user
	chatQueueUsers := make([]model.ChatQueueUser, 0)
	for _, item := range data.ChatQueueUser.UserId {
		chatQueueUser := model.ChatQueueUser{
			Base:     model.InitBase(),
			TenantId: authUser.TenantId,
			QueueId:  queueExist.GetId(),
			UserId:   item,
			Source:   authUser.Source,
		}
		chatQueueUsers = append(chatQueueUsers, chatQueueUser)
	}
	err = repository.ChatQueueUserRepo.TxBulkInsert(ctx, tx, chatQueueUsers)
	if err != nil {
		log.Error(err)
		return err
	}

	// update manager's user id
	manageFilter := model.ChatManageQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  queueExist.GetId(),
	}
	_, chatManagers, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, manageFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if chatManagers != nil {
		isOk := false
		for i, manageQueueUser := range *chatManagers {
			if len(data.ChatManageQueueUser.UserId) > 0 && data.ChatManageQueueUser.UserId != manageQueueUser.UserId {
				(*chatManagers)[i].UserId = data.ChatManageQueueUser.UserId
				isOk = true
			}
		}

		if len(*chatManagers) > 0 {
			queueExist.ManageQueueId = (*chatManagers)[0].Id
			if isOk {
				if err = repository.ManageQueueRepo.TxBulkUpdate(ctx, tx, *chatManagers); err != nil {
					log.Error(err)
					return err
				}
			}
		}
	}
	if chatManagers != nil && len(*chatManagers) < 1 {
		// insert manager for this queue
		manageQueue := model.ChatManageQueueUser{
			Base:     model.InitBase(),
			TenantId: authUser.TenantId,
			QueueId:  queueExist.GetId(),
			UserId:   data.ChatManageQueueUser.UserId,
		}
		queueExist.ManageQueueId = manageQueue.GetId()

		if err = repository.ManageQueueRepo.TxInsert(ctx, tx, manageQueue); err != nil {
			log.Error(err)
			return err
		}
	}

	// clear cache
	if err = cache.RCache.Del([]string{CHAT_QUEUE_USER + "_" + authUser.TenantId}); err != nil {
		log.Error(err)
		return err
	}
	if err = repository.ChatQueueRepo.TxUpdate(ctx, tx, *queueExist); err != nil {
		log.Error(err)
		return err
	}
	if err = repository.ChatQueueRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatQueue) UpdateChatQueueStatus(ctx context.Context, authUser *model.AuthUser, id string, status string) (err error) {
	queueExist, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if queueExist == nil {
		err = errors.New("chat queue " + id + " not found")
		log.Error(err)
		return
	}

	queueExist.Status = status
	err = repository.ChatQueueRepo.UpdateChatQueueStatus(ctx, repository.DBConn, *queueExist)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatQueue) DeleteChatQueueByIdV2(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	queueExist, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	} else if queueExist == nil {
		log.Error("chat queue " + id + " not found")
		return errors.New("chat queue " + id + " not found")
	}

	// do not delete queue if there's at least a conversation active in that queue
	allocateFilter := model.AllocateUserFilter{
		TenantId:     authUser.TenantId,
		QueueId:      []string{id},
		MainAllocate: "active",
	}
	total, _, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, allocateFilter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if total > 0 {
		err = errors.New("unable to delete queue " + id + " because it is currently in use by a conversation")
		log.Error(err)
		return
	}

	tx, err := repository.ChatQueueRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback()

	// remove chat queue users
	queueUserFilter := model.ChatQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  []string{queueExist.GetId()},
	}
	_, oldQueueUser, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, queueUserFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(*oldQueueUser) > 0 {
		if err = repository.ChatQueueUserRepo.TxBulkDelete(ctx, tx, *oldQueueUser); err != nil {
			log.Error(err)
			return err
		}
	}

	// remove old manage queue user
	manageFilter := model.ChatManageQueueUserFilter{
		TenantId: authUser.TenantId,
		QueueId:  queueExist.GetId(),
	}
	_, oldManageQueues, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, manageFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(*oldManageQueues) > 0 {
		if err = repository.ManageQueueRepo.TxBulkDelete(ctx, tx, *oldManageQueues); err != nil {
			log.Error(err)
			return err
		}
	}
	// delete related connection queues
	connectionQueueFilter := model.ConnectionQueueFilter{
		TenantId: authUser.TenantId,
		QueueId:  queueExist.GetId(),
	}
	_, oldConnectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, connectionQueueFilter, -1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	currentTime := time.Now()
	// update field in connection apps
	for _, connectionQueue := range *oldConnectionQueues {
		connectionAppFilter := model.ChatConnectionAppFilter{
			TenantId:          authUser.TenantId,
			ConnectionQueueId: connectionQueue.GetId(),
		}
		_, connectionApps, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, repository.DBConn, connectionAppFilter, -1, 0)
		if err != nil {
			log.Error(err)
			return err
		}
		for i := range *connectionApps {
			(*connectionApps)[i].UpdatedAt = currentTime
		}
		if len(*connectionApps) > 0 {
			if err = repository.ChatConnectionPipelineRepo.BulkUpdateConnectionApp(ctx, tx, *connectionApps, "connection_queue_id", "NULL"); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	if len(*oldConnectionQueues) > 0 {
		if err = repository.ConnectionQueueRepo.TxBulkDelete(ctx, tx, *oldConnectionQueues); err != nil {
			log.Error(err)
			return err
		}
	}

	err = repository.ChatQueueRepo.TxDelete(ctx, tx, *queueExist)
	if err != nil {
		log.Error(err)
		return err
	}
	// clear cache
	if err = cache.RCache.Del([]string{CHAT_QUEUE_USER + "_" + authUser.TenantId}); err != nil {
		log.Error(err)
		return err
	}

	if err = repository.ChatQueueRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
