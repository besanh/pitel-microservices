package service

import (
	"context"
	"errors"
	"strconv"
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

	filter := model.AppFilter{
		AppType:    data.ConnectionType,
		DefaultApp: "active",
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
	connectionApp.QueueId = data.QueueId
	connectionApp.OaInfo = *data.OaInfo
	if data.ConnectionType == "facebook" {
		connectionApp.OaInfo.Facebook[0].AppId = (*app)[0].InfoApp.Facebook.AppId
		connectionApp.OaInfo.Facebook[0].CreatedTimestamp = time.Now().Unix()
	} else if data.ConnectionType == "zalo" {
		connectionApp.OaInfo.Zalo[0].AppId = (*app)[0].InfoApp.Zalo.AppId
		connectionApp.OaInfo.Zalo[0].CreatedTimestamp = time.Now().Unix()
	}
	connectionApp.Status = data.Status

	if err := repository.ChatConnectionAppRepo.Insert(ctx, dbCon, connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	if len(data.QueueId) > 0 {
		connectionQueue := model.ConnectionQueue{
			Base:         model.InitBase(),
			TenantId:     authUser.TenantId,
			ConnectionId: connectionApp.Id,
			QueueId:      data.QueueId,
		}
		if err = repository.ConnectionQueueRepo.Insert(ctx, repository.DBConn, connectionQueue); err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
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
			Base:      model.InitBase(),
			TenantId:  authUser.TenantId,
			ShareType: "zalo",
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
					OaId:      (*app)[0].InfoApp.Zalo.OaId,
					ImageName: "oa_zalo.png",
					ImageUrl:  API_SHARE_INFO_HOST + "/images/oa_zalo.png",
					Title:     (*app)[0].InfoApp.Zalo.OaName,
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
		chatConnectionAppExist.QueueId = data.QueueId
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
			chatConnectionAppExist.OaInfo.Zalo[0].AccessToken = data.AccessToken
			expire, _ := strconv.ParseInt(data.Expire, 10, 64)
			chatConnectionAppExist.OaInfo.Zalo[0].Expire = expire
			chatConnectionAppExist.OaInfo.Zalo[0].UpdatedTimestamp = time.Now().Unix()
		}
	} else if chatConnectionAppExist.ConnectionType == "facebook" && len(data.OaId) > 0 {
		chatConnectionAppExist.OaInfo.Facebook[0].OaId = data.OaId
	}
	chatConnectionAppExist.UpdatedAt = time.Now()

	if len(data.OaId) < 1 {
		if err = repository.ConnectionQueueRepo.DeleteConnectionQueue(ctx, repository.DBConn, "", chatConnectionAppExist.QueueId); err != nil {
			log.Error(err)
			return err
		}
	}

	if err = repository.ChatConnectionAppRepo.Update(ctx, dbCon, *chatConnectionAppExist); err != nil {
		log.Error(err)
		return err
	}

	// Update share form
	if isUpdateFromOtt {
		if chatConnectionAppExist.ConnectionType == "zalo" {
			filter := model.ShareInfoFormFilter{
				AppId: data.AppId,
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

			if err = repository.ShareInfoRepo.Update(ctx, repository.DBConn, (*shareInfo)[0]); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	if len(data.OaId) < 1 {
		connectionQueue := model.ConnectionQueue{
			Base:         model.InitBase(),
			TenantId:     chatConnectionAppExist.TenantId,
			ConnectionId: chatConnectionAppExist.Id,
			QueueId:      data.QueueId,
		}
		if err = repository.ConnectionQueueRepo.Insert(ctx, repository.DBConn, connectionQueue); err != nil {
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
