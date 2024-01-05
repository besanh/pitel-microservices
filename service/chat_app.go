package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatApp interface {
		InsertChatApp(ctx context.Context, authUser *model.AuthUser, data model.ChatAppRequest) (string, error)
		GetChatApp(ctx context.Context, authUser *model.AuthUser, filter model.AppFilter, limit, offset int) (int, *[]model.ChatApp, error)
	}
	ChatApp struct{}
)

func NewChatApp() IChatApp {
	return &ChatApp{}
}

func (s *ChatApp) InsertChatApp(ctx context.Context, authUser *model.AuthUser, data model.ChatAppRequest) (string, error) {
	app := model.ChatApp{
		Base: model.InitBase(),
	}
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return app.Base.GetId(), err
	}

	total, _, err := repository.ChatAppRepo.GetChatApp(ctx, dbCon, model.AppFilter{
		AppName: data.AppName,
		Status: sql.NullBool{
			Valid: true,
			Bool:  true,
		},
	}, 1, 0)
	if err != nil {
		log.Error(err)
		return app.Base.GetId(), err
	} else if total > 0 {
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

func (s *ChatApp) GetChatApp(ctx context.Context, authUser *model.AuthUser, filter model.AppFilter, limit, offset int) (int, *[]model.ChatApp, error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return 0, nil, err
	}

	total, apps, err := repository.ChatAppRepo.GetChatApp(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, apps, nil
}
