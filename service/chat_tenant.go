package service

import (
	"context"
	"time"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/internal/sqlclient"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
)

type (
	IChatTenant interface {
		InsertChatTenant(ctx context.Context, payload *model.ChatTenantRequest) (id string, err error)
		GetChatTenants(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatTenantFilter, limit, offset int) (total int, result *[]model.ChatTenant, err error)
		UpdateChatTenantById(ctx context.Context, id string, payload *model.ChatTenantRequest) (err error)
	}
	ChatTenant struct {
	}
)

var ChatTenantService IChatTenant

func NewChatTenant() IChatTenant {
	return &ChatTenant{}
}

func (s *ChatTenant) InsertChatTenant(ctx context.Context, payload *model.ChatTenantRequest) (id string, err error) {
	chatTenant := model.ChatTenant{
		Base:       model.InitBase(),
		TenantName: payload.TenantName,
		Status:     payload.Status,
	}

	if err := repository.ChatTenantRepo.Insert(ctx, repository.DBConn, chatTenant); err != nil {
		log.Error(err)
		return chatTenant.Base.GetId(), err
	}
	return chatTenant.Base.GetId(), nil
}

func (s *ChatTenant) GetChatTenants(ctx context.Context, db sqlclient.ISqlClientConn, filter model.ChatTenantFilter, limit, offset int) (total int, result *[]model.ChatTenant, err error) {
	total, result, err = repository.ChatTenantRepo.GetChatTenants(ctx, db, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	return total, result, nil
}

func (s *ChatTenant) UpdateChatTenantById(ctx context.Context, id string, payload *model.ChatTenantRequest) (err error) {
	chatTenantExist, err := repository.ChatTenantRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if chatTenantExist == nil {
		log.Error(err)
		return
	}

	integrateSystemExist, err := repository.ChatIntegrateSystemRepo.GetById(ctx, repository.DBConn, payload.IntegrateSystemId)
	if err != nil {
		log.Error(err)
		return
	} else if integrateSystemExist == nil {
		log.Error(err)
		return
	}

	chatTenantExist.IntegrateSystemId = payload.IntegrateSystemId
	chatTenantExist.TenantName = payload.TenantName
	chatTenantExist.Status = payload.Status
	chatTenantExist.UpdatedAt = time.Now()
	err = repository.ChatTenantRepo.Update(ctx, repository.DBConn, *chatTenantExist)
	if err != nil {
		log.Error(err)
		return
	}
	return nil

}
