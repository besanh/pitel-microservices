package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IRoutingConfig interface {
		InsertRoutingConfig(ctx context.Context, authUser *model.AuthUser, data model.RoutingConfig) error
		GetRoutingConfigs(ctx context.Context, authUser *model.AuthUser, filter model.RoutingConfigFilter, limit, offset int) (total int, result *[]model.RoutingConfigView, err error)
		GetRoutingConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.RoutingConfig, err error)
		PutRoutingConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.RoutingConfig) error
		DeleteRoutingConfigById(ctx context.Context, authUser *model.AuthUser, id string) error
	}
	RoutingConfig struct{}
)

func NewRoutingConfig() IRoutingConfig {
	return &RoutingConfig{}
}

func (s *RoutingConfig) InsertRoutingConfig(ctx context.Context, authUser *model.AuthUser, data model.RoutingConfig) error {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}

	routingConfig := model.RoutingConfig{
		Base: model.InitBase(),
	}

	if err := util.ParseAnyToAny(data, &routingConfig); err != nil {
		log.Error(err)
		return err
	}

	if data.RoutingFlow.FlowType == "recipient" {
		for _, val := range data.RoutingFlow.FlowUuid {
			ok, err := common.CheckRecipientExist(ctx, dbCon, val)
			if err != nil {
				log.Error(err)
				return err
			}
			if !ok {
				log.Error(errors.New("recipient config is exist"))
				return errors.New("recipient config is exist")
			}
		}
	}

	if data.RoutingFlow.FlowType == "balance" {
		for _, val := range data.RoutingFlow.FlowUuid {
			ok, err := common.CheckBalanceExist(ctx, dbCon, val)
			if err != nil {
				log.Error(err)
				return err
			}
			if !ok {
				log.Error(errors.New("balance config is exist"))
				return errors.New("balance config is exist")
			}
		}
	}

	if err := repository.RoutingConfigRepo.Insert(ctx, dbCon, routingConfig); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *RoutingConfig) GetRoutingConfigs(ctx context.Context, authUser *model.AuthUser, filter model.RoutingConfigFilter, limit, offset int) (total int, result *[]model.RoutingConfigView, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return 0, nil, errors.New(response.ERR_EMPTY_CONN)
	}

	total, result, err = repository.RoutingConfigRepo.GetRoutingConfigs(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, result, errors.New(response.ERR_GET_FAILED)
	}

	return total, result, nil
}

func (s *RoutingConfig) GetRoutingConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.RoutingConfig, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return nil, errors.New(response.ERR_EMPTY_CONN)
	}

	data, err := repository.RoutingConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *RoutingConfig) PutRoutingConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.RoutingConfig) error {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}

	routingConfigExist, err := repository.RoutingConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if routingConfigExist == nil {
		log.Error(errors.New("routing config is not exist"))
		return errors.New("routing config is not exist")
	}

	if data.RoutingFlow.FlowType == "recipient" {
		for _, val := range data.RoutingFlow.FlowUuid {
			ok, err := common.CheckRecipientExist(ctx, dbCon, val)
			if err != nil {
				log.Error(err)
				return err
			}
			if !ok {
				log.Error(errors.New("recipient config is exist"))
				return errors.New("recipient config is exist")
			}
		}
	}

	if data.RoutingFlow.FlowType == "balance" {
		for _, val := range data.RoutingFlow.FlowUuid {
			ok, err := common.CheckBalanceExist(ctx, dbCon, val)
			if err != nil {
				log.Error(err)
				return err
			}
			if !ok {
				log.Error(errors.New("balance config is exist"))
				return errors.New("balance config is exist")
			}
		}
	}

	if err := util.ParseAnyToAny(data, routingConfigExist); err != nil {
		log.Error(err)
		return err
	}

	if err := repository.RoutingConfigRepo.Update(ctx, dbCon, *routingConfigExist); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *RoutingConfig) DeleteRoutingConfigById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}

	dataExist, err := repository.RoutingConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if len(dataExist.Id) < 1 {
		return errors.New("routing config is not exist")
	}

	if err := repository.RoutingConfigRepo.Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
