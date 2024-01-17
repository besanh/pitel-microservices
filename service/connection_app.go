package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatConnectionApp interface {
		InsertChatConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ChatConnectionAppRequest) (string, error)
		GetChatConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ChatConnectionAppFilter, limit, offset int) (int, *[]model.ChatConnectionApp, error)
		GetChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (model.ChatConnectionApp, error)
		UpdateChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatConnectionAppRequest) (err error)
		DeleteChatConnectionAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ChatConnectionApp struct{}
)

func NewChatConnectionApp() IChatConnectionApp {
	return &ChatConnectionApp{}
}

func (s *ChatConnectionApp) InsertChatConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ChatConnectionAppRequest) (string, error) {
	connectionApp := model.ChatConnectionApp{
		Base:           model.InitBase(),
		TenantId:       authUser.TenantId,
		BusinessUnitId: authUser.BusinessUnitId,
		ConnectionName: data.ConnectionName,
		ConnectionType: data.ConnectionType,
		Status:         data.Status,
	}

	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return connectionApp.Base.GetId(), err
	}

	_, err = repository.ChatAppRepo.GetById(ctx, dbCon, data.AppId)
	if err != nil {
		log.Error(err)
		return connectionApp.Base.GetId(), err
	}

	_, err = repository.ChatQueueAgentRepo.GetById(ctx, dbCon, data.QueueId)
	if err != nil {
		log.Error(err)
		return connectionApp.Base.GetId(), err
	}

	connectionApp.AppId = data.AppId
	connectionApp.QueueId = data.QueueId

	if err := repository.ChatConnectionAppRepo.Insert(ctx, dbCon, connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Base.GetId(), err
	}
	return connectionApp.Base.GetId(), nil
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
	}

	chatConnectionAppExist.ConnectionName = data.ConnectionName
	chatConnectionAppExist.ConnectionType = data.ConnectionType
	chatConnectionAppExist.Status = data.Status
	err = repository.ChatConnectionAppRepo.Update(ctx, dbCon, *chatConnectionAppExist)
	if err != nil {
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
