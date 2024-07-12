package service

import (
	"context"
	"crypto/md5"
	"fmt"

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
		UpdateChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, id string, data *model.ChatIntegrateSystemRequest) error
		DeleteChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, id string) error
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

	if len(data.TenantDefaultId) < 1 {
		chatTenant := model.ChatTenant{
			Base:              model.InitBase(),
			TenantName:        data.SystemName,
			IntegrateSystemId: chatIntegrateSystem.Base.GetId(),
			Status:            true,
		}
		if err = repository.ChatTenantRepo.Insert(ctx, repository.DBConn, chatTenant); err != nil {
			log.Error(err)
			return
		}

		chatIntegrateSystem.TenantDefaultId = chatTenant.Base.GetId()
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
		AuthType:      data.AuthType,
		Username:      data.Username,
		Password:      data.Password,
		Token:         data.Token,
		WebsocketUrl:  data.WebsocketUrl,
		ApiUrl:        data.ApiUrl,
		ApiGetUserUrl: data.ApiGetUserUrl,
	}

	if err = repository.ChatIntegrateSystemRepo.Insert(ctx, repository.DBConn, chatIntegrateSystem); err != nil {
		log.Error(err)
		if err = repository.ChatTenantRepo.Delete(ctx, repository.DBConn, chatIntegrateSystem.TenantDefaultId); err != nil {
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
	result, err = repository.ChatIntegrateSystemRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if result == nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatIntegrateSystem) UpdateChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, id string, data *model.ChatIntegrateSystemRequest) error {
	chatIntegrateSystemExist, err := repository.ChatIntegrateSystemRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	} else if chatIntegrateSystemExist == nil {
		log.Error(err)
		return err
	}
	chatIntegrateSystemExist.SystemName = data.SystemName
	chatIntegrateSystemExist.VendorId = data.VendorId
	chatIntegrateSystemExist.Status = data.Status
	chatIntegrateSystemExist.InfoSystem = &model.InfoSystem{
		AuthType:      data.AuthType,
		Username:      data.Username,
		Password:      data.Password,
		Token:         data.Token,
		WebsocketUrl:  data.WebsocketUrl,
		ApiUrl:        data.ApiUrl,
		ApiGetUserUrl: data.ApiGetUserUrl,
	}
	err = repository.ChatIntegrateSystemRepo.Update(ctx, repository.DBConn, *chatIntegrateSystemExist)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *ChatIntegrateSystem) DeleteChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, id string) error {
	_, err := repository.ChatIntegrateSystemRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	}
	err = repository.ChatIntegrateSystemRepo.Delete(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
