package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatIntegrateSystem interface {
		InsertChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, data *model.ChatIntegrateSystemRequest) (string, error)
		GetChatIntegrateSystems(ctx context.Context, authUser *model.AuthUser, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error)
		GetChatIntegrateSystemById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ChatIntegrateSystem, err error)
		UpdateChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, id string, data *model.ChatIntegrateSystemRequest) error
		DeleteChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, id string) error
	}
	ChatIntegrateSystem struct{}
)

var ChatIntegrateSystemService IChatIntegrateSystem

func NewChatTenantIntegrateSystem() IChatIntegrateSystem {
	return &ChatIntegrateSystem{}
}

func (s *ChatIntegrateSystem) InsertChatIntegrateSystem(ctx context.Context, authUser *model.AuthUser, data *model.ChatIntegrateSystemRequest) (string, error) {
	chatIntegrateSystem := model.ChatIntegrateSystem{
		Base: model.InitBase(),
	}

	_, err := repository.VendorRepo.GetById(ctx, repository.DBConn, data.VendorId)
	if err != nil {
		log.Error(err)
		return chatIntegrateSystem.Base.GetId(), err
	}

	chatIntegrateSystem.SystemName = data.SystemName
	chatIntegrateSystem.VendorId = data.VendorId
	chatIntegrateSystem.Status = data.Status
	chatIntegrateSystem.InfoSystem = data.InfoSystem

	if err := repository.ChatIntegrateSystemRepo.Insert(ctx, repository.DBConn, chatIntegrateSystem); err != nil {
		return chatIntegrateSystem.Base.GetId(), err
	}
	return chatIntegrateSystem.Base.GetId(), nil
}

func (s *ChatIntegrateSystem) GetChatIntegrateSystems(ctx context.Context, authUser *model.AuthUser, filter model.ChatIntegrateSystemFilter, limit, offset int) (total int, result *[]model.ChatIntegrateSystem, err error) {
	total, result, err = repository.ChatIntegrateSystemRepo.GetIntegrateSystem(ctx, repository.DBConn, filter, limit, offset)
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
	chatIntegrateSystemExist.InfoSystem = data.InfoSystem
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
