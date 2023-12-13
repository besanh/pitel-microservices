package service

import (
	"context"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IBalanceConfig interface {
		InsertBalanceConfig(ctx context.Context, authUser *model.AuthUser, data model.BalanceConfigBodyRequest) error
		GetBalanceConfigs(ctx context.Context, authUser *model.AuthUser, filter model.BalanceConfigFilter, limit, offset int) (total int, result *[]model.BalanceConfigView, err error)
		GetBalanceConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.BalanceConfig, err error)
		PutBalanceConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.BalanceConfigPutBodyRequest) (err error)
		DeleteBalanceConfigById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	BalanceConfig struct{}
)

func NewBalanceConfig() IBalanceConfig {
	return &BalanceConfig{}
}

func (s *BalanceConfig) InsertBalanceConfig(ctx context.Context, authUser *model.AuthUser, data model.BalanceConfigBodyRequest) error {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return err
	}

	filter := model.PluginConfigFilter{
		PluginName: []string{data.Provider},
		PluginType: []string{data.BalanceType},
	}
	if !slices.Contains[[]string](constants.CHANNEL, data.BalanceType) {
		return fmt.Errorf("balance type %s not support", data.BalanceType)
	}
	total, _, err := repository.PluginConfigRepo.GetPluginConfigs(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total < 0 {
		return fmt.Errorf("provider %s not found", data.Provider)
	}

	filterBalanceExist := model.BalanceConfigFilter{
		Provider: []string{data.Provider},
		BalanceType: []string{
			data.BalanceType,
		},
	}
	total, _, err = repository.BalanceConfigRepo.GetBalanceConfigs(ctx, dbCon, filterBalanceExist, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 0 {
		return fmt.Errorf("weight %s with type %s, provider %s is exist", data.Weight, data.BalanceType, data.Provider)
	}

	// TODO: check total 100% with provider and type

	balanceConfig := model.BalanceConfig{
		Base:        model.InitBase(),
		Weight:      data.Weight,
		BalanceType: data.BalanceType,
		Priority:    data.Priority,
		Provider:    data.Provider,
		Status:      data.Status,
		CreatedBy:   authUser.Username,
	}

	if err := repository.BalanceConfigRepo.Insert(ctx, dbCon, balanceConfig); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *BalanceConfig) GetBalanceConfigs(ctx context.Context, authUser *model.AuthUser, filter model.BalanceConfigFilter, limit, offset int) (total int, result *[]model.BalanceConfigView, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return 0, nil, err
	}

	total, result, err = repository.BalanceConfigRepo.GetBalanceConfigs(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}

	return total, result, nil
}

func (s *BalanceConfig) GetBalanceConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.BalanceConfig, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return nil, err
	}

	data, err := repository.BalanceConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *BalanceConfig) PutBalanceConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.BalanceConfigPutBodyRequest) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return err
	}

	balancaConfigExist, err := repository.BalanceConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	filter := model.PluginConfigFilter{
		PluginName: []string{data.Provider},
		PluginType: []string{data.BalanceType},
	}
	if !slices.Contains[[]string](constants.CHANNEL, data.BalanceType) {
		return fmt.Errorf("balance type %s not support", data.BalanceType)
	}
	total, _, err := repository.PluginConfigRepo.GetPluginConfigs(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return err
	}
	if total > 1 {
		return fmt.Errorf("provider %s is exist", data.Provider)
	}
	if total == 0 {
		return fmt.Errorf("provider %s not found", data.Provider)
	}

	// TODO: check total 100% with provider and type

	balancaConfigExist.Weight = data.Weight
	balancaConfigExist.BalanceType = data.BalanceType
	balancaConfigExist.Priority = data.Priority
	balancaConfigExist.Status = data.Status
	balancaConfigExist.Provider = data.Provider
	balancaConfigExist.UpdatedAt = time.Now()

	if err := repository.BalanceConfigRepo.Update(ctx, dbCon, *balancaConfigExist); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *BalanceConfig) DeleteBalanceConfigById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err != nil {
		return err
	}

	_, err = repository.BalanceConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	if err = repository.BalanceConfigRepo.Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
