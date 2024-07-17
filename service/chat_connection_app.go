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
	IChatConnectionApp interface {
		InsertChatConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ChatConnectionAppRequest) (string, error)
		GetChatConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionAppView, error)
		GetChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (model.ChatConnectionApp, error)
		UpdateChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatConnectionAppRequest, isUpdateFromOtt bool) (err error)
		DeleteChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ChatConnectionApp struct{}
)

func NewChatConnectionApp() IChatConnectionApp {
	return &ChatConnectionApp{}
}

func (s *ChatConnectionApp) InsertChatConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ChatConnectionAppRequest) (string, error) {
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
	if len(*app) < 1 {
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

func (s *ChatConnectionApp) GetChatConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ChatConnectionAppFilter, limit, offset int) (total int, connectionApps *[]model.ChatConnectionAppView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}
	filter.TenantId = authUser.TenantId

	total, apps, err := repository.ChatConnectionAppRepo.GetChatConnectionAppCustom(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return total, apps, nil
}

func (s *ChatConnectionApp) GetChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (chatApp model.ChatConnectionApp, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	chatConnectionApp, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	} else if chatConnectionApp == nil {
		log.Error("connection app not found")
		return
	}

	return *chatConnectionApp, nil
}

func (s *ChatConnectionApp) UpdateChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatConnectionAppRequest, isUpdateFromOtt bool) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatConnectionAppExist, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if chatConnectionAppExist == nil {
		log.Error("connection app " + id + " not found")
		return errors.New("connection app " + id + " not found")
	}

	if len(data.ConnectionName) > 0 {
		chatConnectionAppExist.ConnectionName = data.ConnectionName
	}

	if len(data.ConnectionType) > 0 {
		chatConnectionAppExist.ConnectionType = data.ConnectionType
	}

	if len(data.QueueId) > 0 {
		filterConnectionQueue := model.ConnectionQueueFilter{
			TenantId:     authUser.TenantId,
			ConnectionId: chatConnectionAppExist.Id,
			QueueId:      data.QueueId,
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
				QueueId:      data.QueueId,
			}
			if err = repository.ConnectionQueueRepo.Insert(ctx, repository.DBConn, connectionQueue); err != nil {
				log.Error(err)
				return err
			}

			chatConnectionAppExist.ConnectionQueueId = connectionQueue.Id
		} else {
			// TODO: if connection have not connection_queue, we need to create and update it
			// TODO: update queue in connection queue
			connectionQueueExist, err := repository.ConnectionQueueRepo.GetById(ctx, repository.DBConn, chatConnectionAppExist.ConnectionQueueId)
			if err != nil {
				log.Error(err)
				return err
			} else if connectionQueueExist == nil {
				log.Error("connection queue " + chatConnectionAppExist.ConnectionQueueId + " not found")
				return errors.New("connection queue " + chatConnectionAppExist.ConnectionQueueId + " not found")
			}

			// TODO: delete connection queue
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
				for _, item := range *connectionQueueExists {
					if item.Id != connectionQueueExist.Id {
						if err = repository.ConnectionQueueRepo.Delete(ctx, repository.DBConn, item.Id); err != nil {
							log.Error(err)
							return err
						}
					}
				}
			}

			connectionQueueExist.QueueId = data.QueueId
			if err = repository.ConnectionQueueRepo.Update(ctx, repository.DBConn, *connectionQueueExist); err != nil {
				log.Error(err)
				return err
			}

			chatConnectionAppExist.ConnectionQueueId = connectionQueueExist.Id
		}
	}

	if data.OaInfo != nil {
		chatConnectionAppExist.OaInfo.Zalo = data.OaInfo.Zalo
	}
	if len(data.Status) > 0 {
		chatConnectionAppExist.Status = data.Status
	}

	if chatConnectionAppExist.ConnectionType == "zalo" && len(data.OaId) > 0 {
		chatConnectionAppExist.OaInfo.Zalo[0].OaId = data.OaId
		chatConnectionAppExist.OaInfo.Zalo[0].AppId = data.AppId
		chatConnectionAppExist.OaInfo.Zalo[0].OaName = data.OaName
		chatConnectionAppExist.OaInfo.Zalo[0].Avatar = data.Avatar
		chatConnectionAppExist.OaInfo.Zalo[0].Cover = data.Cover
		chatConnectionAppExist.OaInfo.Zalo[0].CateName = data.CateName
		chatConnectionAppExist.OaInfo.Zalo[0].Status = data.Status
		if isUpdateFromOtt {
			createdTimestamp, _ := time.Parse("2006-01-02T15:04:05.999999999Z", data.TokenCreatedAt)
			chatConnectionAppExist.OaInfo.Zalo[0].CreatedTimestamp = createdTimestamp.Unix()
			chatConnectionAppExist.OaInfo.Zalo[0].Expire = data.TokenExpiresIn
			chatConnectionAppExist.OaInfo.Zalo[0].TokenTimeRemainning = data.TokenTimeRemainning
			chatConnectionAppExist.OaInfo.Zalo[0].UpdatedTimestamp = time.Now().Unix()
		}
	} else if chatConnectionAppExist.ConnectionType == "facebook" && len(data.OaId) > 0 {
		chatConnectionAppExist.OaInfo.Facebook[0].OaId = data.OaId
	}
	chatConnectionAppExist.UpdatedAt = time.Now()

	// if len(data.OaId) < 1 && chatConnectionAppExist.ConnectionType == "zalo" {
	// 	if err = repository.ConnectionQueueRepo.DeleteConnectionQueue(ctx, repository.DBConn, "", chatConnectionAppExist.ConnectionQueueId); err != nil {
	// 		log.Error(err)
	// 		return err
	// 	}
	// }

	if len(chatConnectionAppExist.ConnectionQueueId) > 0 {
		if err = repository.ChatConnectionAppRepo.Update(ctx, repository.DBConn, *chatConnectionAppExist); err != nil {
			log.Error(err)
			return err
		}
	} else {
		if err = repository.ChatConnectionAppRepo.UpdateSpecifColumnyById(ctx, dbCon, *chatConnectionAppExist); err != nil {
			log.Error(err)
			return err
		}
	}

	// Update share form
	if isUpdateFromOtt {
		if chatConnectionAppExist.ConnectionType == "zalo" {
			filter := model.ShareInfoFormFilter{
				AppId:        data.AppId,
				ShareType:    "zalo",
				ConnectionId: chatConnectionAppExist.Id,
			}
			_, shareInfo, err := repository.ShareInfoRepo.GetShareInfos(ctx, repository.DBConn, filter, 1, 0)
			if err != nil {
				log.Error(err)
				return err
			}
			if len(*shareInfo) < 1 {
				log.Error("share config app_id " + data.AppId + " not exist")
				err = errors.New("share config app_id " + data.AppId + " not exist")
				return err
			}

			(*shareInfo)[0].ShareForm = model.ShareForm{
				Zalo: struct {
					AppId     string "json:\"app_id\""
					OaId      string "json:\"oa_id\""
					ImageName string "json:\"image_name\""
					ImageUrl  string "json:\"image_url\""
					Title     string "json:\"title\""
					Subtitle  string "json:\"subtitle\""
				}{
					AppId:     (*shareInfo)[0].ShareForm.Zalo.AppId,
					OaId:      data.OaId,
					ImageName: (*shareInfo)[0].ShareForm.Zalo.ImageName,
					ImageUrl:  (*shareInfo)[0].ShareForm.Zalo.ImageUrl,
					Title:     (*shareInfo)[0].ShareForm.Zalo.Title,
					Subtitle:  (*shareInfo)[0].ShareForm.Zalo.Subtitle,
				},
			}

			(*shareInfo)[0].ConnectionId = chatConnectionAppExist.Id
			(*shareInfo)[0].UpdatedAt = time.Now()

			if err = repository.ShareInfoRepo.Update(ctx, repository.DBConn, (*shareInfo)[0]); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	if len(data.OaId) < 1 && len(data.QueueId) < 1 {
		connectionQueue := model.ConnectionQueue{
			Base:         model.InitBase(),
			TenantId:     chatConnectionAppExist.TenantId,
			ConnectionId: chatConnectionAppExist.Id,
			QueueId:      data.QueueId,
		}
		if err := repository.ConnectionQueueRepo.DeleteConnectionQueue(ctx, dbCon, chatConnectionAppExist.Id, ""); err != nil {
			log.Error(err)
			return err
		}
		if err = repository.ConnectionQueueRepo.Insert(ctx, repository.DBConn, connectionQueue); err != nil {
			log.Error(err)
			return err
		}

		// Update connection
		chatConnectionAppExist.ConnectionQueueId = connectionQueue.Id
		chatConnectionAppExist.UpdatedAt = time.Now()
		if err = repository.ChatConnectionAppRepo.Update(ctx, dbCon, *chatConnectionAppExist); err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func (s *ChatConnectionApp) DeleteChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	connectionAppExist, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if connectionAppExist == nil {
		log.Error("connection app not found")
		return errors.New("connection app not found")
	}
	if err = repository.ChatConnectionAppRepo.Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}

	if err = repository.ConnectionQueueRepo.DeleteConnectionQueue(ctx, repository.DBConn, connectionAppExist.Id, ""); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
