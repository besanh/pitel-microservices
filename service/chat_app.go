package service

import (
	"context"
	"errors"
	"time"

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
	systems := make([]model.ChatAppIntegrateSystem, 0)
	for _, id := range data.SystemIds {
		systems = append(systems, model.ChatAppIntegrateSystem{
			ChatAppId:             app.GetId(),
			ChatIntegrateSystemId: id,
			CreatedAt:             time.Now(),
		})
	}

	if err := repository.ChatAppRepo.InsertChatApp(ctx, dbCon, app, systems); err != nil {
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

	chatApp, err := repository.ChatAppRepo.GetChatAppById(ctx, dbCon, id)
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
	systems := make([]model.ChatAppIntegrateSystem, 0)
	for _, systemId := range data.SystemIds {
		systems = append(systems, model.ChatAppIntegrateSystem{
			ChatAppId:             chatAppExist.GetId(),
			ChatIntegrateSystemId: systemId,
			CreatedAt:             time.Now(),
		})
	}

	err = repository.ChatAppRepo.UpdateChatAppById(ctx, dbCon, *chatAppExist, systems)
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
	err = repository.ChatAppRepo.DeleteChatAppById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
