package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IExternalPluginConnect interface {
		InsertExternalPluginConnect(ctx context.Context, authUser *model.AuthUser, data model.ExternalPluginConnect) error
	}
	ExternalPluginConnect struct{}
)

func NewExternalPluginConnect() IExternalPluginConnect {
	return &ExternalPluginConnect{}
}

func (s *ExternalPluginConnect) InsertExternalPluginConnect(ctx context.Context, authUser *model.AuthUser, data model.ExternalPluginConnect) error {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return err
	}
	dataExist, err := repository.ExternalPluginConnectRepo.GetExternalPluginByType(ctx, dbCon, data.PluginType)
	if err != nil {
		log.Error(err)
		return err
	} else if dataExist != nil {
		return errors.New("external plugin is not exist")
	}
	externalPluginConnect := model.ExternalPluginConnect{
		Base:       model.InitBase(),
		PluginName: data.PluginName,
		PluginType: data.PluginType,
		Config:     data.Config,
	}
	if err := repository.ExternalPluginConnectRepo.Insert(ctx, dbCon, externalPluginConnect); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
