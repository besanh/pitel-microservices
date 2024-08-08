package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	repository "github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/repository/sql_builder"
)

type (
	IIBKTenant interface {
		GetTenants(ctx context.Context, authUser *model.AuthUser, params model.IBKTenantQueryParam, limit, offset int) (total int, tenants []*model.IBKTenantInfo, err error)
		GetTenantById(ctx context.Context, authUser *model.AuthUser, id string) (tenant *model.IBKTenantInfo, err error)
		GetById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (entity *model.IBKTenant, err error)
		PostTenant(ctx context.Context, authUser *model.AuthUser, body *model.IBKTenantBody) (err error)
		PutTenant(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKTenantBody) (err error)
		DeleteTenant(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	IBKTenant struct {
	}
)

var IBKTenantService IIBKTenant

func NewIBKTenant() IIBKTenant {
	return &IBKTenant{}
}

func (s *IBKTenant) GetTenants(ctx context.Context, authUser *model.AuthUser, params model.IBKTenantQueryParam, limit, offset int) (total int, tenants []*model.IBKTenantInfo, err error) {
	params.TenantId_Eq = authUser.TenantId
	total, tenants, err = repository.IBKTenantRepo.GetInfoByQuery(ctx, params, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

// func encryptPassword(key string, password string) (encrypted string, err error) {
// 	passEncrypted, err := hashUtil.AesEncrypt([]byte(key), []byte(password))
// 	if err != nil {
// 		return
// 	}
// 	encrypted = hex.EncodeToString(passEncrypted)
// 	return
// }

func (s *IBKTenant) GetTenantById(ctx context.Context, authUser *model.AuthUser, id string) (tenant *model.IBKTenantInfo, err error) {
	tenant, err = repository.IBKTenantRepo.GetInfoById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if tenant == nil {
		err = variables.NewError(constants.ERR_TENANT_NOTFOUND)
		return
	}
	return
}

func (s *IBKTenant) PostTenant(ctx context.Context, authUser *model.AuthUser, body *model.IBKTenantBody) (err error) {
	tenant := &model.IBKTenant{
		Base:       model.InitBaseModel(authUser.GetID(), authUser.GetID()),
		TenantName: body.TenantName,
	}
	if err = repository.IBKTenantRepo.Insert(ctx, *tenant); err != nil {
		log.Error(err)
		return
	}
	return
}

func (*IBKTenant) PutTenant(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKTenantBody) (err error) {
	var tenant *model.IBKTenant
	tenant, err = repository.IBKTenantRepo.GetById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if tenant == nil {
		err = variables.NewError(constants.ERR_TENANT_NOTFOUND)
		return
	}
	tenant.TenantName = body.TenantName
	if err = repository.IBKTenantRepo.Update(ctx, *tenant); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *IBKTenant) DeleteTenant(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	tenant, err := repository.IBKTenantRepo.GetById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if tenant == nil {
		err = variables.NewError("tenant not found")
		return
	}
	if err = repository.IBKTenantRepo.Delete(ctx, id); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *IBKTenant) GetById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (entity *model.IBKTenant, err error) {
	entity, err = repository.IBKTenantRepo.GetById(ctx, id, conditions...)
	if err != nil {
		log.Error(err)
		return
	} else if entity == nil {
		err = variables.NewError(constants.ERR_TENANT_NOTFOUND)
		return
	}
	return
}
