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
		GetChatApp(ctx context.Context, authUser *model.AuthUser, filter model.AppFilter, limit, offset int) (int, *[]model.ChatApp, error)
		GetChatAppById(ctx context.Context, authUser *model.AuthUser, id string) (app model.ChatApp, err error)
		UpdateChatAppById(ctx context.Context, authUser *model.AuthUser, id string, data model.ChatAppRequest) (err error)
		DeleteChatAppById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	ChatApp struct {
	}
)

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

	filter := model.AppFilter{
		AppName: data.AppName,
		Status:  data.Status,
	}
	total, _, err := repository.ChatAppRepo.GetChatApp(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return app.Base.GetId(), err
	} else if total > 0 {
		return app.Base.GetId(), errors.New("app already exists")
	}

	app.AppName = data.AppName
	app.InfoApp = data.InfoApp
	app.Status = data.Status
	app.DefaultApp = data.DefaultApp

	if err := repository.ChatAppRepo.Insert(ctx, dbCon, app); err != nil {
		log.Error(err)
		return app.Base.GetId(), err
	}

	return app.Base.GetId(), nil
}

func (s *ChatApp) GetChatApp(ctx context.Context, authUser *model.AuthUser, filter model.AppFilter, limit, offset int) (int, *[]model.ChatApp, error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	total, apps, err := repository.ChatAppRepo.GetChatApp(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, apps, nil
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
	// chatAppTmp := chatAppExist

	chatAppExist.AppName = data.AppName
	chatAppExist.InfoApp = data.InfoApp
	chatAppExist.Status = data.Status
	chatAppExist.DefaultApp = data.DefaultApp
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
