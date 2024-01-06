package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IConnectionApp interface {
		InsertConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ConnectionAppRequest) (string, error)
		GetConnectionApp(ctx context.Context, authUser *model.AuthUser, filter model.ConnectionAppFilter, limit, offset int) (int, *[]model.ConnectionApp, error)
	}
	ConnectionApp struct{}
)

func NewConnectionApp() IConnectionApp {
	return &ConnectionApp{}
}

func (s *ConnectionApp) InsertConnectionApp(ctx context.Context, authUser *model.AuthUser, data model.ConnectionAppRequest) (string, error) {
	connectionApp := model.ConnectionApp{
		Base:           model.InitBase(),
		TenantId:       authUser.TenantId,
		BusinessUnitId: authUser.BusinessUnitId,
		ConnectionName: data.ConnectionName,
		ConnectionType: data.ConnectionType,
		Status:         data.Status,
	}
	db, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return connectionApp.Base.GetId(), err
	}

	_, err = repository.ChatAppRepo.GetById(ctx, db, data.AppId)
	if err != nil {
		log.Error(err)
		return connectionApp.Base.GetId(), err
	}

	_, err = repository.ChatQueueAgentRepo.GetById(ctx, db, data.QueueId)
	if err != nil {
		log.Error(err)
		return connectionApp.Base.GetId(), err
	}

	connectionApp.AppId = data.AppId
	connectionApp.QueueId = data.QueueId

	if err := repository.ConnectionAppRepo.Insert(ctx, db, connectionApp); err != nil {
		log.Error(err)
		return connectionApp.Base.GetId(), err
	}
	return connectionApp.Base.GetId(), nil
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
