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
		InsertChatConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ChatConnectionAppRequest) (string, error)
		AttachConnectionQueueToApp(ctx context.Context, authUser *model.AuthUser, data model.AttachConnectionQueueToConnectionAppRequest) error
	}
	ChatConnectionPipeline struct{}
)

var ChatConnectionPipelineService IChatConnectionPipeline

func NewChatConnectionPipeline() IChatConnectionPipeline {
	return &ChatConnectionPipeline{}
}

func (s *ChatConnectionPipeline) InsertChatConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ChatConnectionAppRequest) (string, error) {
	id := data.Id
	if len(id) < 1 {
		id = uuid.NewString()
	}
	connectionApp := model.ChatConnectionApp{
		Id:             id,
		TenantId:       authUser.TenantId,
		CreatedAt:      time.Now(),
		ConnectionName: data.ConnectionName,
		ConnectionType: data.ConnectionType,
		Status:         data.Status,
	}

	if len(data.QueueId) > 0 {
		_, err := repository.ChatQueueUserRepo.GetById(ctx, repository.DBConn, data.QueueId)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
	}

	apps, err := ChatAppService.GetChatAppAssign(ctx, authUser)
	if err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}
	var app *model.ChatApp
	for _, item := range apps {
		if data.ChatAppId == item.Id {
			app = item
		}
	}
	if app == nil {
		log.Error("app with type " + data.ConnectionType + " not found")
		return connectionApp.Id, errors.New("app not found")
	}
	connectionApp.ChatAppId = app.GetId()

	tx, err := repository.ChatConnectionPipelineRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}
	defer tx.Rollback()
	connectionQueue := model.ConnectionQueue{
		Base:         model.InitBase(),
		TenantId:     authUser.TenantId,
		ConnectionId: connectionApp.Id,
		QueueId:      data.QueueId,
	}
	// TODO: init connection_queue and add to connection
	if len(data.QueueId) > 0 {
		// TODO: remove on duplicate connection_queue
		if err := repository.ChatConnectionPipelineRepo.DeleteConnectionQueue(ctx, tx, connectionApp.Id, ""); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}

		if err = repository.ConnectionQueueRepo.TxInsert(ctx, tx, connectionQueue); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}

		connectionApp.ConnectionQueueId = connectionQueue.GetId()
	}

	if data.OaInfo.Facebook != nil || data.OaInfo.Zalo != nil {
		connectionApp.OaInfo = *data.OaInfo
	}

	if data.ConnectionType == "facebook" && len(connectionApp.OaInfo.Facebook) > 0 {
		connectionApp.OaInfo.Facebook[0].AppId = app.InfoApp.Facebook.AppId
		connectionApp.OaInfo.Facebook[0].CreatedTimestamp = time.Now().Unix()
	} else if data.ConnectionType == "zalo" && len(connectionApp.OaInfo.Zalo) > 0 {
		connectionApp.OaInfo.Zalo[0].AppId = app.InfoApp.Zalo.AppId
		connectionApp.OaInfo.Zalo[0].CreatedTimestamp = time.Now().Unix()
	}
	connectionApp.Status = data.Status

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
	if err = repository.ChatConnectionPipelineRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	return connectionApp.Id, nil
}

func (s *ChatConnectionPipeline) AttachConnectionQueueToApp(ctx context.Context, authUser *model.AuthUser, data model.AttachConnectionQueueToConnectionAppRequest) (err error) {
	var chatConnectionAppExist *model.ChatConnectionApp
	if data.IsAttachingApp {
		chatConnectionAppExist, err = repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, data.ConnectionId)
		if err != nil {
			log.Error(err)
			return err
		} else if chatConnectionAppExist == nil {
			log.Error("connection app " + data.ConnectionId + " not found")
			return errors.New("connection app " + data.ConnectionId + " not found")
		}
	}

	tx, err := repository.ChatConnectionPipelineRepo.BeginTx(ctx, repository.DBConn, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback()

	// select existed queue
	if data.IsAttachingApp && len(data.ConnectionQueueId) > 0 && chatConnectionAppExist != nil {
		filterConnectionQueue := model.ConnectionQueueFilter{
			TenantId:     authUser.TenantId,
			ConnectionId: chatConnectionAppExist.Id,
			QueueId:      data.ConnectionQueueId,
		}
		_, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filterConnectionQueue, 1, 0)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(*connectionQueues) < 1 || len(chatConnectionAppExist.ConnectionQueueId) < 1 {
			// TODO: delete connection queue with connectionId
			filter := model.ConnectionQueueFilter{
				TenantId:     authUser.TenantId,
				ConnectionId: chatConnectionAppExist.Id,
			}
			_, connectionQueueExists, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, repository.DBConn, filter, -1, 0)
			if err != nil {
				log.Error(err)
				return err
			}
			if len(*connectionQueueExists) > 0 {
				if err = repository.ConnectionQueueRepo.TxBulkDelete(ctx, tx, *connectionQueueExists); err != nil {
					log.Error(err)
					return err
				}
			}
			// TODO: insert connection queue
			connectionQueue := model.ConnectionQueue{
				Base:         model.InitBase(),
				TenantId:     authUser.TenantId,
				ConnectionId: chatConnectionAppExist.Id,
				QueueId:      data.ConnectionQueueId,
			}
			if err = repository.ConnectionQueueRepo.TxInsert(ctx, tx, connectionQueue); err != nil {
				log.Error(err)
				return err
			}

			chatConnectionAppExist.ConnectionQueueId = connectionQueue.Id
		}
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
			ConnectionId: data.ConnectionId,
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
		manageQueue.ConnectionId = data.ConnectionId
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

		if data.IsAttachingApp {
			chatConnectionAppExist.ConnectionQueueId = connectionQueue.Id
		}
	}

	if chatConnectionAppExist != nil && len(chatConnectionAppExist.ConnectionQueueId) > 0 {
		if err = repository.ChatConnectionPipelineRepo.UpdateConnectionApp(ctx, tx, *chatConnectionAppExist); err != nil {
			log.Error(err)
			return err
		}
	}

	if err = repository.ChatConnectionPipelineRepo.CommitTx(ctx, tx); err != nil {
		log.Error(err)
		return err
	}

	return
}
