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
		GetChatConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionApp, error)
		GetChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (model.ChatConnectionApp, error)
		UpdateChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatConnectionAppRequest) (err error)
		DeleteChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ChatConnectionApp struct {
		OttDomain string
	}
)

func NewChatConnectionApp(ottDomain string) IChatConnectionApp {
	return &ChatConnectionApp{
		OttDomain: ottDomain,
	}
}

func (s *ChatConnectionApp) InsertChatConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ChatConnectionAppRequest) (string, error) {
	id := data.Id
	if len(id) < 1 {
		id = uuid.NewString()
	}
	connectionApp := model.ChatConnectionApp{
		Id:             id,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		TenantId:       authUser.TenantId,
		BusinessUnitId: authUser.BusinessUnitId,
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
		_, err = repository.ChatQueueAgentRepo.GetById(ctx, dbCon, data.QueueId)
		if err != nil {
			log.Error(err)
			return connectionApp.Id, err
		}
	}

	filter := model.AppFilter{
		AppType: data.ConnectionType,
	}
	total, app, err := repository.ChatAppRepo.GetChatApp(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}
	if total > 0 {
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
	connectionApp.OaInfo = data.OaInfo
	connectionApp.Status = data.Status

	if err := repository.ChatConnectionAppRepo.Insert(ctx, dbCon, connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	// Step belows apply when app is available
	// Call ott, if fail => roll back
	if err := common.PostOttAccount(s.OttDomain, (*app)[0], connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Id, err
	}

	return connectionApp.Id, nil
}

func (s *ChatConnectionApp) GetChatConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ChatConnectionAppFilter, limit, offset int) (total int, connectionApps *[]model.ChatConnectionApp, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, apps, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, dbCon, filter, limit, offset)
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

func (s *ChatConnectionApp) UpdateChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatConnectionAppRequest) (err error) {
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
		log.Error("connection app not found")
		return errors.New("connection app not found")
	}

	// Check if page having many connection in one app => reject
	// if data.ConnectionType == "zalo" {
	// 	filter := model.ChatConnectionAppFilter{
	// 		AppId: data.AppId,
	// 		OaId:  data.OaId,
	// 	}
	// 	total, _, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, dbCon, filter, -1, 0)
	// 	if err != nil {
	// 		log.Error(err)
	// 		return err
	// 	} else if total > 1 {
	// 		log.Error("page having many connection in one app")
	// 		return errors.New("page having many connection in one app")
	// 	}
	// }

	chatConnectionAppExist.ConnectionName = data.ConnectionName
	chatConnectionAppExist.ConnectionType = data.ConnectionType
	chatConnectionAppExist.QueueId = data.QueueId
	chatConnectionAppExist.OaInfo.Zalo = data.OaInfo.Zalo
	chatConnectionAppExist.Status = data.Status
	chatConnectionAppExist.UpdatedAt = time.Now()
	if chatConnectionAppExist.ConnectionType == "zalo" && len(data.OaId) > 0 {
		chatConnectionAppExist.OaInfo.Zalo[0].OaId = data.OaId
	} else if chatConnectionAppExist.ConnectionType == "facebook" && len(data.OaId) > 0 {
		chatConnectionAppExist.OaInfo.Facebook[0].OaId = data.OaId
	}

	if err = repository.ConnectionQueueRepo.DeleteConnectionQueue(ctx, repository.DBConn, "", chatConnectionAppExist.QueueId); err != nil {
		log.Error(err)
		return err
	}

	if err = repository.ChatConnectionAppRepo.Update(ctx, dbCon, *chatConnectionAppExist); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatConnectionApp) DeleteChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = repository.ChatConnectionAppRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ChatConnectionAppRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
