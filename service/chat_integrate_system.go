package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatIntegrateSystem interface {
		InsertChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, data *model.ChatIntegrateSystemRequest) (id, systemId string, err error)
		GetChatIntegrateSystems(ctx context.Context, authUser *model.AuthUser, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error)
		GetChatIntegrateSystemById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ChatIntegrateSystem, err error)
		UpdateChatIntegrateSystemById(ctx context.Context, authUser *model.AuthUser, id string, data *model.ChatIntegrateSystemRequest) error
		DeleteChatIntegrateSystemById(ctx context.Context, authUser *model.AuthUser, id string) error
	}
	ChatIntegrateSystem struct{}
)

var ChatIntegrateSystemService IChatIntegrateSystem

func NewChatIntegrateSystem() IChatIntegrateSystem {
	return &ChatIntegrateSystem{}
}

func (s *ChatIntegrateSystem) InsertChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, data *model.ChatIntegrateSystemRequest) (id, systemId string, err error) {
	chatIntegrateSystem := model.ChatIntegrateSystem{
		Base: model.InitBase(),
	}

	chatTenant := model.ChatTenant{
		Base: model.InitBase(),
	}
	if len(data.TenantDefaultId) < 1 {
		chatTenant.TenantName = data.SystemName
		chatTenant.IntegrateSystemId = chatIntegrateSystem.Id
		tenantId := uuid.NewString()
		chatTenant.TenantId = tenantId
		chatTenant.Status = true

		chatIntegrateSystem.TenantDefaultId = tenantId
	} else {
		// Check tenant exist in system
		filter := model.ChatIntegrateSystemFilter{
			TenantDefaultId: data.TenantDefaultId,
		}
		_, chatIntegrateSystems, errTmp := repository.ChatIntegrateSystemRepo.GetIntegrateSystems(ctx, repository.DBConn, filter, 1, 0)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if len(*chatIntegrateSystems) > 0 {
			log.Error("tenant " + data.TenantDefaultId + " already exist in system")
			err = errors.New("tenant " + data.TenantDefaultId + " already exist in system")
			return
		}
		chatTenant.TenantName = data.SystemName
		chatTenant.IntegrateSystemId = chatIntegrateSystem.Id
		chatTenant.TenantId = data.TenantDefaultId
		chatIntegrateSystem.TenantDefaultId = data.TenantDefaultId
	}

	_, err = repository.VendorRepo.GetById(ctx, repository.DBConn, data.VendorId)
	if err != nil {
		log.Error(err)
		return
	}

	chatIntegrateSystem.Salt = util.GenerateRandomString(50)
	systemIdByte := []byte(chatIntegrateSystem.Salt + chatIntegrateSystem.Id)
	chatIntegrateSystem.SystemId = fmt.Sprintf("%x", md5.Sum(systemIdByte))
	systemId = chatIntegrateSystem.SystemId

	chatIntegrateSystem.SystemName = data.SystemName
	chatIntegrateSystem.VendorId = data.VendorId
	chatIntegrateSystem.Status = data.Status
	chatIntegrateSystem.InfoSystem = &model.InfoSystem{
		AuthType:            data.AuthType,
		Username:            data.Username,
		Password:            data.Password,
		Token:               data.Token,
		WebsocketUrl:        data.WebsocketUrl,
		ApiUrl:              data.ApiUrl,
		ApiGetUserUrl:       data.ApiGetUserUrl,
		ApiAuthUrl:          data.ApiAuthUrl,
		ApiGetUserDetailUrl: data.ApiGetUserDetailUrl,
	}

	chatAppIntegrateSystems := make([]model.ChatAppIntegrateSystem, 0)
	chatApps := []model.ChatApp{}
	if len(data.ChatAppIds) > 0 {
		for _, item := range data.ChatAppIds {
			chatAppExist, err := repository.ChatAppRepo.GetById(ctx, repository.DBConn, item)
			if err != nil {
				log.Error(err)
				continue
			} else if chatAppExist != nil {
				chatAppIntegrateSystems = append(chatAppIntegrateSystems, model.ChatAppIntegrateSystem{
					Base:                  model.InitBase(),
					ChatAppId:             item,
					ChatIntegrateSystemId: chatIntegrateSystem.Id,
				})
			}
		}
	} else if len(data.ChatApps) > 0 {
		for _, item := range data.ChatApps {
			chatApp := model.ChatApp{
				Base:    model.InitBase(),
				Status:  "active",
				AppName: item.AppName,
				InfoApp: item.InfoApp,
			}
			chatApps = append(chatApps, chatApp)
			chatAppIntegrateSystems = append(chatAppIntegrateSystems, model.ChatAppIntegrateSystem{
				ChatAppId:             chatApp.GetId(),
				ChatIntegrateSystemId: chatIntegrateSystem.Id,
			})
		}
	}

	if err = repository.ChatIntegrateSystemRepo.InsertIntegrateSystemTransaction(ctx, repository.DBConn, chatApps, chatIntegrateSystem, chatAppIntegrateSystems, chatTenant); err != nil {
		log.Error(err)
		if errTmp := repository.ChatTenantRepo.DeleteByTenantId(ctx, repository.DBConn, chatIntegrateSystem.TenantDefaultId); errTmp != nil {
			log.Error(err)
			return
		}
		return
	}

	id = chatIntegrateSystem.Base.GetId()
	return
}

func (s *ChatIntegrateSystem) GetChatIntegrateSystems(ctx context.Context, authUser *model.AuthUser, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error) {
	total, result, err = repository.ChatIntegrateSystemRepo.GetIntegrateSystems(ctx, repository.DBConn, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatIntegrateSystem) GetChatIntegrateSystemById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ChatIntegrateSystem, err error) {
	result, err = repository.ChatIntegrateSystemRepo.GetIntegrateSystemById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if result == nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatIntegrateSystem) UpdateChatIntegrateSystemById(ctx context.Context, authUser *model.AuthUser, id string, data *model.ChatIntegrateSystemRequest) (err error) {
	chatIntegrateSystemExist, err := repository.ChatIntegrateSystemRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if chatIntegrateSystemExist == nil {
		log.Error(err)
		return
	}

	chatTenantExist, err := repository.ChatTenantRepo.GetByTenantId(ctx, repository.DBConn, chatIntegrateSystemExist.TenantDefaultId)
	if err != nil {
		log.Error(err)
		return
	}

	chatTenant := &model.ChatTenant{}
	// Check tenant exist in system
	if len(data.TenantDefaultId) < 1 || chatTenantExist == nil {
		chatTenant.Base = model.InitBase()
		chatTenant.TenantName = "default"
		chatTenant.IntegrateSystemId = uuid.NewString()
		chatTenant.Status = true
		tenantId := uuid.NewString()
		chatTenant.TenantId = tenantId
	} else {
		filter := model.ChatIntegrateSystemFilter{
			TenantDefaultId: data.TenantDefaultId,
		}
		_, chatIntegrateSystems, errTmp := repository.ChatIntegrateSystemRepo.GetIntegrateSystems(ctx, repository.DBConn, filter, 1, 0)
		if errTmp != nil {
			log.Error(errTmp)
			err = errTmp
			return
		} else if len(*chatIntegrateSystems) > 1 {
			log.Error("tenant " + data.TenantDefaultId + " already exist in system")
			err = errors.New("tenant " + data.TenantDefaultId + " already exist in system")
			return
		}
		chatTenant = chatTenantExist
	}

	chatIntegrateSystemExist.TenantDefaultId = chatTenant.TenantId
	chatIntegrateSystemExist.SystemName = data.SystemName
	chatIntegrateSystemExist.VendorId = data.VendorId
	chatIntegrateSystemExist.Status = data.Status
	chatIntegrateSystemExist.InfoSystem = &model.InfoSystem{
		AuthType:            data.AuthType,
		Username:            data.Username,
		Password:            data.Password,
		Token:               data.Token,
		WebsocketUrl:        data.WebsocketUrl,
		ApiUrl:              data.ApiUrl,
		ApiAuthUrl:          data.ApiAuthUrl,
		ApiGetUserUrl:       data.ApiGetUserUrl,
		ApiGetUserDetailUrl: data.ApiGetUserDetailUrl,
	}
	chatAppIntegrateSystems := make([]model.ChatAppIntegrateSystem, 0)
	if len(data.ChatAppIds) > 0 {
		for _, item := range data.ChatAppIds {
			chatAppExist, err := repository.ChatAppRepo.GetById(ctx, repository.DBConn, item)
			if err != nil {
				log.Error(err)
				continue
			} else if chatAppExist != nil {
				chatAppIntegrateSystems = append(chatAppIntegrateSystems, model.ChatAppIntegrateSystem{
					Base:                  model.InitBase(),
					ChatAppId:             item,
					ChatIntegrateSystemId: chatIntegrateSystemExist.GetId(),
				})
			}
		}
	}

	err = repository.ChatIntegrateSystemRepo.UpdateIntegrateSystemTransaction(ctx, repository.DBConn, *chatIntegrateSystemExist, chatAppIntegrateSystems, chatTenantExist, chatTenant)
	if err != nil {
		log.Error(err)
		return
	}

	return nil
}

func (s *ChatIntegrateSystem) DeleteChatIntegrateSystemById(ctx context.Context, authUser *model.AuthUser, id string) error {
	_, err := repository.ChatIntegrateSystemRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ChatIntegrateSystemRepo.DeleteIntegrateSystemById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
