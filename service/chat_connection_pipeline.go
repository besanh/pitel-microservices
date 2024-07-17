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
		UpdatedAt:      time.Now(),
		ConnectionName: data.ConnectionName,
		ConnectionType: data.ConnectionType,
		Status:         data.Status,
	}

	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	if len(data.QueueId) > 0 {
		_, err = repository.ChatQueueUserRepo.GetById(ctx, dbCon, data.QueueId)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
	}

	filter := model.ChatAppFilter{
		AppType: data.ConnectionType,
	}
	_, app, err := repository.ChatAppRepo.GetChatApp(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}
	if len(*app) > 0 {
		if data.ConnectionType == "facebook" {
			connectionApp.AppId = (*app)[0].InfoApp.Facebook.AppId
		} else if data.ConnectionType == "zalo" {
			connectionApp.AppId = (*app)[0].InfoApp.Zalo.AppId
		}
	} else {
		log.Error("app with type " + data.ConnectionType + " not found")
		return connectionApp.Id, errors.New("app not found")
	}

	connectionQueue := model.ConnectionQueue{
		Base:         model.InitBase(),
		TenantId:     authUser.TenantId,
		ConnectionId: connectionApp.Id,
		QueueId:      data.QueueId,
	}
	// TODO: init connection_queue and add to connection
	if len(data.ConnectionQueueId) > 0 {
		connectionQueueExist, err := repository.ConnectionQueueRepo.GetById(ctx, dbCon, data.ConnectionQueueId)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		} else if connectionQueueExist == nil {
			log.Error("connection queue " + data.ConnectionQueueId + " not found")
			return connectionApp.Id, errors.New("connection queue " + data.ConnectionQueueId + " not found")
		}
		connectionQueue = *connectionQueueExist
		connectionApp.ConnectionQueueId = connectionQueue.GetId()
	} else if len(data.QueueId) > 0 {
		// TODO: remove on duplicate connection_queue
		if err := repository.ConnectionQueueRepo.DeleteConnectionQueue(ctx, dbCon, connectionApp.Id, ""); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}

		if err = repository.ConnectionQueueRepo.Insert(ctx, repository.DBConn, connectionQueue); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}

		connectionApp.ConnectionQueueId = connectionQueue.GetId()
	}

	if data.OaInfo.Facebook != nil || data.OaInfo.Zalo != nil {
		connectionApp.OaInfo = *data.OaInfo
	}

	if data.ConnectionType == "facebook" && len(connectionApp.OaInfo.Facebook) > 0 {
		connectionApp.OaInfo.Facebook[0].AppId = (*app)[0].InfoApp.Facebook.AppId
		connectionApp.OaInfo.Facebook[0].CreatedTimestamp = time.Now().Unix()
	} else if data.ConnectionType == "zalo" && len(connectionApp.OaInfo.Zalo) > 0 {
		connectionApp.OaInfo.Zalo[0].AppId = (*app)[0].InfoApp.Zalo.AppId
		connectionApp.OaInfo.Zalo[0].CreatedTimestamp = time.Now().Unix()
	}
	connectionApp.Status = data.Status

	if err := repository.ChatConnectionAppRepo.Insert(ctx, dbCon, connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	// Step belows apply when app is available
	// Call ott, if fail => roll back
	if err := common.PostOttAccount(OTT_URL, OTT_VERSION, (*app)[0], connectionApp); err != nil {
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
					AppId:     (*app)[0].InfoApp.Zalo.AppId,
					OaId:      connectionApp.OaInfo.Zalo[0].OaId,
					ImageName: "oa_zalo.png",
					ImageUrl:  API_DOC + "/bss-image/oa_zalo.png",
					Title:     connectionApp.ConnectionName,
					Subtitle:  ZALO_SHARE_INFO_SUBTITLE,
				},
			},
		}
		if err = repository.ShareInfoRepo.Insert(ctx, dbCon, shareInfo); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
	}

	return connectionApp.Id, nil
}

func (s *ChatConnectionPipeline) AttachConnectionQueueToApp(ctx context.Context, authUser *model.AuthUser, data model.AttachConnectionQueueToConnectionAppRequest) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	chatConnectionAppExist, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, data.ConnectionId)
	if err != nil {
		log.Error(err)
		return err
	} else if chatConnectionAppExist == nil {
		log.Error("connection app " + data.ConnectionId + " not found")
		return errors.New("connection app " + data.ConnectionId + " not found")
	}

	// select existed queue
	if len(data.ConnectionQueueId) > 0 {
		filterConnectionQueue := model.ConnectionQueueFilter{
			TenantId:     authUser.TenantId,
			ConnectionId: chatConnectionAppExist.Id,
			QueueId:      data.ConnectionQueueId,
		}
		_, connectionQueues, err := repository.ConnectionQueueRepo.GetConnectionQueues(ctx, dbCon, filterConnectionQueue, 1, 0)
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
				if err = repository.ConnectionQueueRepo.BulkDeleteConnectionQueue(ctx, repository.DBConn, connectionQueueExists); err != nil {
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
			if err = repository.ConnectionQueueRepo.Insert(ctx, repository.DBConn, connectionQueue); err != nil {
				log.Error(err)
				return err
			}

			chatConnectionAppExist.ConnectionQueueId = connectionQueue.Id
		}
	} else {
		// create new queue and update it to c.app
		chatQueue := model.ChatQueue{
			Base: model.InitBase(),
		}

		routingExist, err := repository.ChatRoutingRepo.GetById(ctx, dbCon, data.ChatQueue.ChatRoutingId)
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
		if err = repository.ConnectionQueueRepo.Insert(ctx, dbCon, connectionQueue); err != nil {
			log.Error(err)
			return err
		}

		chatQueue.TenantId = authUser.TenantId
		chatQueue.QueueName = data.ChatQueue.QueueName
		chatQueue.Description = data.ChatQueue.Description
		chatQueue.ChatRoutingId = data.ChatQueue.ChatRoutingId
		chatQueue.Status = data.ChatQueue.Status
		if err = repository.ChatQueueRepo.Insert(ctx, dbCon, chatQueue); err != nil {
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
		err = repository.ChatQueueUserRepo.BulkInsert(ctx, dbCon, chatQueueUsers)
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

		queueExist, err := repository.ChatQueueRepo.GetById(ctx, dbCon, chatQueue.GetId())
		if err != nil {
			log.Error(err)
			return err
		} else if queueExist == nil {
			log.Error("queue not found")
			return errors.New("queue " + chatQueue.GetId() + " not found")
		}
		queueExist.ManageQueueId = manageQueue.GetId()

		if err = repository.ManageQueueRepo.Insert(ctx, dbCon, manageQueue); err != nil {
			log.Error(err)
			return err
		}
		if err = repository.ChatQueueRepo.Update(ctx, dbCon, *queueExist); err != nil {
			log.Error(err)
			return err
		}

		chatConnectionAppExist.ConnectionQueueId = connectionQueue.Id
	}

	if len(chatConnectionAppExist.ConnectionQueueId) > 0 {
		if err = repository.ChatConnectionAppRepo.Update(ctx, repository.DBConn, *chatConnectionAppExist); err != nil {
			log.Error(err)
			return err
		}
	}

	return
}
