package service

import (
	"context"
	"time"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	repository "github.com/tel4vn/fins-microservices/repository"
)

type (
	IIBKRole interface {
		GetRoles(ctx context.Context, authUser *model.AuthUser, params model.IBKRoleQueryParam, limit int, offset int) (total int, entries []*model.IBKRoleInfo, err error)
		GetRoleById(ctx context.Context, authUser *model.AuthUser, id string) (entry *model.IBKRoleInfo, err error)
		CreateRole(ctx context.Context, authUser *model.AuthUser, body *model.IBKRoleBody) (id string, err error)
		PutRole(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKRoleBody) (err error)
		DeleteRole(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}

	IBKRole struct {
	}
)

var IBKRoleService IIBKRole

func NewIBKRole() IIBKRole {
	return &IBKRole{}
}

func (s *IBKRole) GetRoles(ctx context.Context, authUser *model.AuthUser, params model.IBKRoleQueryParam, limit int, offset int) (total int, entries []*model.IBKRoleInfo, err error) {
	entries, total, err = repository.IBKRoleRepo.SelectInfo(ctx, params, limit, offset)
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	}
	return
}

func (s *IBKRole) GetRoleById(ctx context.Context, authUser *model.AuthUser, id string) (entry *model.IBKRoleInfo, err error) {
	entry, err = repository.IBKRoleRepo.FindInfoById(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	}
	return
}

func (*IBKRole) CreateRole(ctx context.Context, authUser *model.AuthUser, body *model.IBKRoleBody) (id string, err error) {
	role := &model.IBKRole{
		Base:        model.InitBaseModel(authUser.GetID(), authUser.GetID()),
		TenantId:    authUser.TenantId,
		RoleName:    body.RoleName,
		Description: body.Description,
	}

	if err = repository.IBKRoleRepo.Insert(ctx, *role); err != nil {
		log.Error(err)
		return
	}

	id = role.Id
	return
}

func (*IBKRole) PutRole(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKRoleBody) (err error) {
	role, err := repository.IBKRoleRepo.GetById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if role == nil {
		err = variables.NewError(constants.ERR_USER_NOTFOUND)
		return
	}

	role.UpdatedAt = time.Now()
	role.RoleName = body.RoleName
	role.Description = body.Description

	if err = repository.IBKRoleRepo.Update(ctx, *role); err != nil {
		log.Error(err)
		return
	}
	return
}

func (*IBKRole) DeleteRole(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	role, err := repository.IBKRoleRepo.GetById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if role == nil {
		err = variables.NewError(constants.ERR_USER_NOTFOUND)
		return
	}

	if err = repository.IBKRoleRepo.Delete(ctx, id); err != nil {
		log.Error(err)
		return
	}
	return
}
