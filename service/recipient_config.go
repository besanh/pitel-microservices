package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IRecipientConfig interface {
		InsertRecipientConfig(ctx context.Context, authUser *model.AuthUser, data model.RecipientConfigRequest) (string, error)
		GetRecipientConfigs(ctx context.Context, authUser *model.AuthUser, filter model.RecipientConfigFilter, limit, offset int) (total int, result *[]model.RecipientConfigView, err error)
		GetRecipientConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.RecipientConfig, err error)
		PutRecipientConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.RecipientConfigPutRequest) (err error)
		DeleteRecipientConfigById(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	RecipientConfig struct{}
)

func NewRecipientConfig() IRecipientConfig {
	return &RecipientConfig{}
}

func (s *RecipientConfig) InsertRecipientConfig(ctx context.Context, authUser *model.AuthUser, data model.RecipientConfigRequest) (string, error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return "", errors.New(response.ERR_EMPTY_CONN)
	}
	filter := model.PluginConfigFilter{
		PluginName: []string{data.Provider},
		PluginType: []string{data.RecipientType},
	}
	if !slices.Contains[[]string](constants.CHANNEL, data.RecipientType) {
		return "", fmt.Errorf("recipient type %s not support", data.RecipientType)
	}
	total, _, err := repository.PluginConfigRepo.GetPluginConfigs(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if total < 0 {
		return "", fmt.Errorf("provider %s not found", data.Provider)
	}
	if data.Provider == "abenla" && !slices.Contains[[]string](constants.ROLE_ABELA, data.Recipient) {
		return "", fmt.Errorf("recipient %s with provider %s not support", data.Recipient, data.Provider)
	} else if data.Provider == "incom" && !slices.Contains[[]string](constants.ROLE_INCOM, data.Recipient) {
		return "", fmt.Errorf("recipient %s with provider %s not support", data.Recipient, data.Provider)
	} else if data.Provider == "zalo" && !slices.Contains[[]string](constants.ROLE_FPT, data.Recipient) {
		return "", fmt.Errorf("recipient %s with provider %s not support", data.Recipient, data.Provider)
	}
	filterTmp := model.RecipientConfigFilter{
		Provider:      []string{data.Provider},
		RecipientType: []string{data.RecipientType},
		Recipient:     []string{data.Recipient},
	}
	total, _, err = repository.RecipientConfigRepo.GetRecipientConfigs(ctx, dbCon, filterTmp, 1, 0)
	if err != nil {
		log.Error(err)
		return "", err
	}
	if total > 0 {
		return "", fmt.Errorf("recipient %s with provider %s already exist", data.Recipient, data.Provider)
	}
	recipientConfig := model.RecipientConfig{
		Base:          model.InitBase(),
		Recipient:     data.Recipient,
		RecipientType: data.RecipientType,
		Priority:      data.Priority,
		Provider:      data.Provider,
		Status:        data.Status,
		CreatedBy:     authUser.UserId,
	}

	if err := repository.RecipientConfigRepo.Insert(ctx, dbCon, recipientConfig); err != nil {
		log.Error(err)
		return "", err
	}

	return recipientConfig.Base.GetId(), nil
}

func (s *RecipientConfig) GetRecipientConfigs(ctx context.Context, authUser *model.AuthUser, filter model.RecipientConfigFilter, limit, offset int) (total int, result *[]model.RecipientConfigView, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return 0, nil, errors.New(response.ERR_EMPTY_CONN)
	}
	total, result, err = repository.RecipientConfigRepo.GetRecipientConfigs(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return 0, result, errors.New(response.ERR_GET_FAILED)
	}
	return total, result, nil
}

func (s *RecipientConfig) GetRecipientConfigById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.RecipientConfig, err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return nil, errors.New(response.ERR_EMPTY_CONN)
	}
	data, err := repository.RecipientConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if data == nil {
		return nil, errors.New("recipient config is not exist")
	}

	result = data
	return result, nil
}

func (s *RecipientConfig) PutRecipientConfigById(ctx context.Context, authUser *model.AuthUser, id string, data model.RecipientConfigPutRequest) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}

	dataExist, err := repository.RecipientConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if len(dataExist.Id) < 1 {
		return errors.New("recipient config is not exist")
	}

	filter := model.PluginConfigFilter{
		PluginName: []string{data.Provider},
		PluginType: []string{data.RecipientType},
	}
	if !slices.Contains[[]string](constants.CHANNEL, data.RecipientType) {
		return fmt.Errorf("recipient type %s not support", data.RecipientType)
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

	dataExist.Recipient = data.Recipient
	dataExist.RecipientType = data.RecipientType
	dataExist.Priority = data.Priority
	dataExist.Provider = data.Provider
	dataExist.Status = data.Status
	dataExist.UpdatedAt = time.Now()
	if err := repository.RecipientConfigRepo.Update(ctx, dbCon, *dataExist); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *RecipientConfig) DeleteRecipientConfigById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := GetDBConnOfUser(*authUser)
	if err == ERR_EMPTY_CONN {
		return errors.New(response.ERR_EMPTY_CONN)
	}
	dataExist, err := repository.RecipientConfigRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	} else if len(dataExist.Id) < 1 {
		return errors.New("recipient config is not exist")
	}

	if err := repository.RecipientConfigRepo.Delete(ctx, dbCon, id); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
