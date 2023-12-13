package service

import (
	"context"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IPluginConfig interface {
		GetPluginConfigs(ctx context.Context, authUser *model.AuthUser, filter model.PluginConfigFilter, limit, offset int) (total int, result *[]model.PluginConfigView, err error)
		PostPluginConfig(ctx context.Context, authUser *model.AuthUser, data model.PluginConfigRequest) (err error)
		GetPluginConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.PluginConfig, err error)
		PutPluginConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.PluginConfigRequest) (err error)
		DeletePluginConfigById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	PluginConfig struct{}
)

func NewPluginConfig() IPluginConfig {
	return &PluginConfig{}
}

func (s *PluginConfig) PostPluginConfig(ctx context.Context, authUser *model.AuthUser, data model.PluginConfigRequest) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}
	filter := model.PluginConfigFilter{
		PluginName: []string{data.PluginName},
		PluginType: []string{data.PluginType},
	}
	total, _, err := repository.PluginConfigRepo.GetPluginConfigs(ctx, dbCon, filter, -1, 0)
	if err != nil {
		return errors.New(response.ERR_GET_FAILED)
	}
	if total > 0 {
		return errors.New("plugin is exist")
	}
	pluginConfig := model.PluginConfig{
		Base:       model.InitBase(),
		PluginName: data.PluginName,
		PluginType: data.PluginType,
		Status:     data.Status,
	}

	if err := repository.PluginConfigRepo.Insert(ctx, dbCon, pluginConfig); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *PluginConfig) GetPluginConfigs(ctx context.Context, authUser *model.AuthUser, filter model.PluginConfigFilter, limit, offset int) (total int, result *[]model.PluginConfigView, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return 0, nil, errors.New(response.ERR_EMPTY_CONN)
	}
	total, result, err = repository.PluginConfigRepo.GetPluginConfigs(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, result, errors.New(response.ERR_GET_FAILED)
	}
	return total, result, nil
}

func (s *PluginConfig) GetPluginConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.PluginConfig, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return nil, errors.New(response.ERR_EMPTY_CONN)
	}
	data, err := repository.PluginConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if data == nil {
		return nil, errors.New("plugin config is not exist")
	}

	result = data

	return result, nil
}

func (s *PluginConfig) PutPluginConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.PluginConfigRequest) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}
	dataExist, err := repository.PluginConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if len(dataExist.Id) < 1 {
		return errors.New("plugin config is not exist")
	}

	dataExist.PluginName = data.PluginName
	dataExist.Status = data.Status
	dataExist.CreatedAt = time.Now()
	if err := repository.PluginConfigRepo.Update(ctx, dbCon, *dataExist); err != nil {
		log.Error(err)
		return err
	}

	return
}

func (s *PluginConfig) DeletePluginConfigById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}
	dataExist, err := repository.PluginConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if len(dataExist.Id) < 1 {
		return errors.New("plugin config is not exist")
	}

	if err := repository.PluginConfigRepo.Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}

	return
}
