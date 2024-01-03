package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IConnectionApp interface {
		InsertConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ConnectionApp) (string, error)
		GetConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ConnectionAppFilter, limit, offset int) (int, *[]model.ConnectionApp, error)
	}
	ConnectionApp struct{}
)

func NewConnectionApp() IConnectionApp {
	return &ConnectionApp{}
}

func (s *ConnectionApp) InsertConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ConnectionApp) (string, error) {
	connectionApp := model.ChatApp{
		Base: model.InitBase(),
	}
	db, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return connectionApp.Base.GetId(), err
	}

	if err := repository.ConnectionAppRepo.Insert(ctx, db, data); err != nil {
		return connectionApp.Base.GetId(), err
	}
	return data.Base.GetId(), nil
}

func (s *ConnectionApp) GetConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ConnectionAppFilter, limit, offset int) (int, *[]model.ConnectionApp, error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return 0, nil, err
	}

	total, apps, err := repository.ConnectionAppRepo.GetConnectionApp(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, apps, nil
}
