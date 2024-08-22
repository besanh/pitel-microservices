package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatApp interface {
		InsertChatApp(ctx context.Context, authUser *model.AuthUser, data model.ChatAppRequest) (string, error)
		GetChatApp(ctx context.Context, authUser *model.AuthUser, filter model.ChatAppFilter, limit, offset int) (int, *[]model.ChatApp, error)
		GetChatAppById(ctx context.Context, authUser *model.AuthUser, id string) (app model.ChatApp, err error)
		UpdateChatAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatAppRequest) (err error)
		DeleteChatAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
		GetChatAppAssign(ctx context.Context, authUser *model.AuthUser) (result []*model.ChatApp, err error)
	}
	ChatApp struct {
	}
)

var ChatAppService IChatApp

func NewChatApp() IChatApp {
	return &ChatApp{}
}

func (s *ChatApp) InsertChatApp(ctx context.Context, authUser *model.AuthUser, data model.ChatAppRequest) (string, error) {
	app := model.ChatApp{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return app.Base.GetId(), err
	}

	filter := model.ChatAppFilter{
		AppName: data.AppName,
	}
	_, chatApps, err := repository.ChatAppRepo.GetChatApp(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return app.Base.GetId(), err
	} else if len(*chatApps) > 0 {
		return app.Base.GetId(), errors.New("app already exists")
	}

	app.AppName = data.AppName
	app.InfoApp = data.InfoApp
	app.Status = data.Status

	if err := repository.ChatAppRepo.Insert(ctx, dbCon, app); err != nil {
		log.Error(err)
		return app.Base.GetId(), err
	}

	return app.Base.GetId(), nil
}

func (s *ChatApp) GetChatApp(ctx context.Context, authUser *model.AuthUser, filter model.ChatAppFilter, limit, offset int) (total int, apps *[]model.ChatApp, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, apps, err = repository.ChatAppRepo.GetChatApp(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return
}

func (s *ChatApp) GetChatAppById(ctx context.Context, authUser *model.AuthUser, id string) (app model.ChatApp, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return app, err
	}

	chatApp, err := repository.ChatAppRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return app, err
	}

	return *chatApp, nil
}

func (s *ChatApp) UpdateChatAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatAppRequest) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatAppExist, err := repository.ChatAppRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	filter := model.ChatAppFilter{
		AppName: data.AppName,
	}
	_, chatApps, err := repository.ChatAppRepo.GetChatApp(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	} else if len(*chatApps) > 1 {
		log.Errorf("app name %s already exists", data.AppName)
		return errors.New("app name " + data.AppName + " already exists")
	}
	chatAppExist.AppName = data.AppName
	chatAppExist.InfoApp = data.InfoApp
	chatAppExist.Status = data.Status

	err = repository.ChatAppRepo.Update(ctx, dbCon, *chatAppExist)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatApp) DeleteChatAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = repository.ChatAppRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ChatAppRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ChatApp) GetChatAppAssign(ctx context.Context, authUser *model.AuthUser) (result []*model.ChatApp, err error) {
	filter := model.ChatIntegrateSystemFilter{
		SystemId: authUser.SystemId,
	}
	_, chatIntegrateSystem, err := repository.ChatIntegrateSystemRepo.GetIntegrateSystems(ctx, repository.DBConn, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}

	if len(*chatIntegrateSystem) > 0 {
		for _, item := range *chatIntegrateSystem {
			if len(item.ChatAppIntegrateSystems) > 0 {
				for _, item2 := range item.ChatAppIntegrateSystems {
					result = append(result, item2.ChatApp)
				}
			}
		}
	}
	return
}
