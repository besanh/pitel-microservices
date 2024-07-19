package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IChatConnectionPipeline interface {
		AttachConnectionQueueToApp(ctx context.Context, authUser *model.AuthUser, data model.AttachConnectionQueueToConnectionAppRequest) (string, error)
		UpsertConnectionQueueInApp(ctx context.Context, authUser *model.AuthUser, id string, data model.EditConnectionQueueInConnectionAppRequest) error
		UpdateConnectionAppStatus(ctx context.Context, authUser *model.AuthUser, id string, status string) error
	}
	ChatConnectionPipeline struct{}
)

var ChatConnectionPipelineService IChatConnectionPipeline

func NewChatConnectionPipeline() IChatConnectionPipeline {
	return &ChatConnectionPipeline{}
}

func (s *ChatConnectionPipeline) AttachConnectionQueueToApp(ctx context.Context, authUser *model.AuthUser, data model.AttachConnectionQueueToConnectionAppRequest) (string, error) {
	id := data.ConnectionAppRequest.Id
	if len(id) < 1 {
		id = uuid.NewString()
	}
	connectionApp := model.ChatConnectionApp{
		Id:             id,
		TenantId:       authUser.TenantId,
		CreatedAt:      time.Now(),
		ConnectionName: data.ConnectionAppRequest.ConnectionName,
		ConnectionType: data.ConnectionAppRequest.ConnectionType,
		Status:         data.ConnectionAppRequest.Status,
	}

	tx, err := repository.ChatConnectionPipelineRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}
	defer tx.Rollback()

	apps, err := ChatAppService.GetChatAppAssign(ctx, authUser)
	if err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}
	var app *model.ChatApp
	for _, item := range apps {
		if data.ConnectionAppRequest.ChatAppId == item.Id {
			app = item
		}
	}
	if app == nil {
		log.Error("app with type " + data.ConnectionAppRequest.ConnectionType + " not found")
		return connectionApp.Id, errors.New("app not found")
	}
	connectionApp.ChatAppId = app.GetId()

	if data.ConnectionAppRequest.OaInfo.Facebook != nil || data.ConnectionAppRequest.OaInfo.Zalo != nil {
		connectionApp.OaInfo = *data.ConnectionAppRequest.OaInfo
	}

	if data.ConnectionAppRequest.ConnectionType == "facebook" && len(connectionApp.OaInfo.Facebook) > 0 {
		connectionApp.OaInfo.Facebook[0].AppId = app.InfoApp.Facebook.AppId
		connectionApp.OaInfo.Facebook[0].CreatedTimestamp = time.Now().Unix()
	} else if data.ConnectionAppRequest.ConnectionType == "zalo" && len(connectionApp.OaInfo.Zalo) > 0 {
		connectionApp.OaInfo.Zalo[0].AppId = app.InfoApp.Zalo.AppId
		connectionApp.OaInfo.Zalo[0].CreatedTimestamp = time.Now().Unix()
	}
	connectionApp.Status = data.ConnectionAppRequest.Status

	if err := repository.ChatConnectionPipelineRepo.InsertConnectionApp(ctx, tx, connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	// Step belows apply when app is available
	// Call ott, if fail => roll back
	if err := common.PostOttAccount(OTT_URL, OTT_VERSION, *app, connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	// Insert share info
	if connectionApp.ConnectionType == "zalo" {
		shareInfo := model.ShareInfoForm{
			Base:         model.InitBase(),
			TenantId:     authUser.TenantId,
			ConnectionId: connectionApp.Id,
			ShareType:    "zalo",
			ShareForm: model.ShareForm{
				Zalo: struct {
					AppId     string "json:\"app_id\""
					OaId      string "json:\"oa_id\""
					ImageName string "json:\"image_name\""
					ImageUrl  string "json:\"image_url\""
					Title     string "json:\"title\""
					Subtitle  string "json:\"subtitle\""
				}{
					AppId:     app.InfoApp.Zalo.AppId,
					OaId:      connectionApp.OaInfo.Zalo[0].OaId,
					ImageName: "oa_zalo.png",
					ImageUrl:  API_DOC + "/bss-image/oa_zalo.png",
					Title:     connectionApp.ConnectionName,
					Subtitle:  ZALO_SHARE_INFO_SUBTITLE,
				},
			},
		}
		if err = repository.ShareInfoRepo.TxInsert(ctx, tx, shareInfo); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
	}

	// select existed queue
	if len(data.ConnectionQueueId) > 0 {
		chatQueueExist, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, data.ConnectionQueueId)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		} else if chatQueueExist == nil {
			err = errors.New("selected chat queue not found")
			log.Error(err)
			return connectionApp.Id, err
		}

		filterConnectionQueue := model.ConnectionQueueFilter{
			TenantId:     authUser.TenantId,
			ConnectionId: connectionApp.Id,
			QueueId:      data.ConnectionQueueId,
		}
		_, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filterConnectionQueue, 1, 0)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
		if len(*connectionQueues) < 1 || len(connectionApp.ConnectionQueueId) < 1 {
			// TODO: delete connection queue with connectionId
			filter := model.ConnectionQueueFilter{
				TenantId:     authUser.TenantId,
				ConnectionId: connectionApp.Id,
			}
			_, connectionQueueExists, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, -1, 0)
			if err != nil {
				log.Error(err)
				return connectionApp.Id, err
			}
			if len(*connectionQueueExists) > 0 {
				if err = repository.ConnectionQueueRepo.TxBulkDelete(ctx, tx, *connectionQueueExists); err != nil {
					log.Error(err)
					return connectionApp.Id, err
				}
			}
			// TODO: insert connection queue
			connectionQueue := model.ConnectionQueue{
				Base:         model.InitBase(),
				TenantId:     authUser.TenantId,
				ConnectionId: connectionApp.Id,
				QueueId:      data.ConnectionQueueId,
			}
			if err = repository.ConnectionQueueRepo.TxInsert(ctx, tx, connectionQueue); err != nil {
				log.Error(err)
				return connectionApp.Id, err
			}

			connectionApp.ConnectionQueueId = connectionQueue.Id
		}
	} else if len(data.ConnectionQueueId) < 1 {
		// create new queue and update it to c.app
		chatQueue := model.ChatQueue{
			Base: model.InitBase(),
		}

		routingExist, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, data.ChatQueue.ChatRoutingId)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		} else if routingExist == nil {
			err = errors.New("chat routing not found")
			return connectionApp.Id, err
		}

		connectionQueue := model.ConnectionQueue{
			Base:         model.InitBase(),
			TenantId:     authUser.TenantId,
			ConnectionId: connectionApp.Id,
			QueueId:      chatQueue.Base.GetId(),
		}
		if err = repository.ConnectionQueueRepo.TxInsert(ctx, tx, connectionQueue); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}

		chatQueue.TenantId = authUser.TenantId
		chatQueue.QueueName = data.ChatQueue.QueueName
		chatQueue.Description = data.ChatQueue.Description
		chatQueue.ChatRoutingId = data.ChatQueue.ChatRoutingId
		chatQueue.Status = data.ChatQueue.Status
		if err = repository.ChatQueueRepo.TxInsert(ctx, tx, chatQueue); err != nil {
			log.Error(err)
			return connectionApp.Id, err
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
		err = repository.ChatQueueUserRepo.TxBulkInsert(ctx, tx, chatQueueUsers)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}

		// insert manage queue user
		manageQueue := model.ChatManageQueueUser{
			Base: model.InitBase(),
		}
		manageQueue.TenantId = authUser.TenantId
		manageQueue.ConnectionId = connectionApp.Id
		manageQueue.QueueId = chatQueue.GetId()
		manageQueue.UserId = data.ChatManageQueueUser.UserId

		chatQueue.ManageQueueId = manageQueue.GetId()

		if err = repository.ManageQueueRepo.TxInsert(ctx, tx, manageQueue); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
		if err = repository.ChatQueueRepo.TxUpdate(ctx, tx, chatQueue); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}

		connectionApp.ConnectionQueueId = connectionQueue.Id
	}

	if len(connectionApp.ConnectionQueueId) > 0 {
		if err = repository.ChatConnectionPipelineRepo.UpdateConnectionApp(ctx, tx, connectionApp); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
	}

	if err = repository.ChatConnectionPipelineRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	return connectionApp.Id, err
}

func (s *ChatConnectionPipeline) UpsertConnectionQueueInApp(ctx context.Context, authUser *model.AuthUser, id string, data model.EditConnectionQueueInConnectionAppRequest) (err error) {
	connectionAppExist, err := repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
	} else if connectionAppExist == nil {
		err = errors.New("connection app not exist")
		log.Error(err)
		return
	}

	tx, err := repository.ChatConnectionPipelineRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer tx.Rollback()

	// remove old relation data
	filter := model.ConnectionQueueFilter{
		TenantId:     authUser.TenantId,
		ConnectionId: connectionAppExist.Id,
	}
	_, connectionQueueExists, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*connectionQueueExists) > 0 {
		if err = repository.ConnectionQueueRepo.TxBulkDelete(ctx, tx, *connectionQueueExists); err != nil {
			log.Error(err)
			return
		}
	}

	// select existed queue
	if len(data.ConnectionQueueId) > 0 {
		chatQueueExist, err := repository.ChatQueueRepo.GetById(ctx, repository.DBConn, data.ConnectionQueueId)
		if err != nil {
			log.Error(err)
			return err
		} else if chatQueueExist == nil {
			err = errors.New("selected chat queue not found")
			log.Error(err)
			return err
		}

		// insert connection queue
		connectionQueue := model.ConnectionQueue{
			Base:         model.InitBase(),
			TenantId:     authUser.TenantId,
			ConnectionId: connectionAppExist.Id,
			QueueId:      data.ConnectionQueueId,
		}
		if err = repository.ConnectionQueueRepo.TxInsert(ctx, tx, connectionQueue); err != nil {
			log.Error(err)
			return err
		}

		connectionAppExist.ConnectionQueueId = connectionQueue.Id
	} else if len(data.ConnectionQueueId) < 1 {
		// create new queue and update it to c.app
		chatQueue := model.ChatQueue{
			Base: model.InitBase(),
		}

		routingExist, err := repository.ChatRoutingRepo.GetById(ctx, repository.DBConn, data.ChatQueue.ChatRoutingId)
		if err != nil {
			log.Error(err)
			return err
		} else if routingExist == nil {
			err = errors.New("chat routing not found")
			return err
		}

		connectionQueue := model.ConnectionQueue{
			Base:         model.InitBase(),
			TenantId:     authUser.TenantId,
			ConnectionId: connectionAppExist.Id,
			QueueId:      chatQueue.Base.GetId(),
		}
		if err = repository.ConnectionQueueRepo.TxInsert(ctx, tx, connectionQueue); err != nil {
			log.Error(err)
			return err
		}

		chatQueue.TenantId = authUser.TenantId
		chatQueue.QueueName = data.ChatQueue.QueueName
		chatQueue.Description = data.ChatQueue.Description
		chatQueue.ChatRoutingId = data.ChatQueue.ChatRoutingId
		chatQueue.Status = data.ChatQueue.Status
		if err = repository.ChatQueueRepo.TxInsert(ctx, tx, chatQueue); err != nil {
			log.Error(err)
			return err
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
		err = repository.ChatQueueUserRepo.TxBulkInsert(ctx, tx, chatQueueUsers)
		if err != nil {
			log.Error(err)
			return err
		}

		// insert manage queue user
		manageQueue := model.ChatManageQueueUser{
			Base: model.InitBase(),
		}
		manageQueue.TenantId = authUser.TenantId
		manageQueue.ConnectionId = connectionAppExist.Id
		manageQueue.QueueId = chatQueue.GetId()
		manageQueue.UserId = data.ChatManageQueueUser.UserId

		chatQueue.ManageQueueId = manageQueue.GetId()

		if err = repository.ManageQueueRepo.TxInsert(ctx, tx, manageQueue); err != nil {
			log.Error(err)
			return err
		}
		if err = repository.ChatQueueRepo.TxUpdate(ctx, tx, chatQueue); err != nil {
			log.Error(err)
			return err
		}

		connectionAppExist.ConnectionQueueId = connectionQueue.Id
	}

	if len(connectionAppExist.ConnectionQueueId) > 0 {
		if err = repository.ChatConnectionPipelineRepo.UpdateConnectionApp(ctx, tx, *connectionAppExist); err != nil {
			log.Error(err)
			return
		}
	}

	if err = repository.ChatConnectionPipelineRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatConnectionPipeline) UpdateConnectionAppStatus(ctx context.Context, authUser *model.AuthUser, id string, status string) (err error) {
	connectionAppExist, err := repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
	} else if connectionAppExist == nil {
		err = errors.New("connection app not exist")
		log.Error(err)
		return
	}

	connectionAppExist.Status = status
	err = repository.ChatConnectionPipelineRepo.UpdateConnectionAppStatus(ctx, repository.DBConn, *connectionAppExist)
	if err != nil {
		log.Error(err)
		return
	}

	return nil
}
