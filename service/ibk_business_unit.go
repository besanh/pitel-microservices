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
	IIBKBusinessUnit interface {
		GetBusinessUnits(ctx context.Context, authUser *model.AuthUser, params model.IBKBusinessUnitQueryParam, limit, offset int) (total int, result []*model.IBKBusinessUnit, err error)
		GetBusinessUnitById(ctx context.Context, authUser *model.AuthUser, id string) (businessUnit *model.IBKBusinessUnit, err error)
		FindById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (entity *model.IBKBusinessUnit, err error)
		CreateBusinessUnit(ctx context.Context, authUser *model.AuthUser, body *model.IBKBusinessUnitBody) (id string, err error)
		PutBusinessUnit(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKBusinessUnitBody) (err error)
		DeleteBusinessUnit(ctx context.Context, authUser *model.AuthUser, id string) (err error)
	}
	IBKBusinessUnit struct {
	}
)

var IBKBusinessUnitService IIBKBusinessUnit

func NewIBKBusinessUnit() IIBKBusinessUnit {
	return &IBKBusinessUnit{}
}

func (*IBKBusinessUnit) GetBusinessUnits(ctx context.Context, authUser *model.AuthUser, params model.IBKBusinessUnitQueryParam, limit, offset int) (total int, result []*model.IBKBusinessUnit, err error) {
	total, result, err = repository.IBKBusinessUnitRepo.SelectByQuery(ctx, nil, limit, offset, "")
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (*IBKBusinessUnit) GetBusinessUnitById(ctx context.Context, authUser *model.AuthUser, id string) (entry *model.IBKBusinessUnit, err error) {
	entry, err = repository.IBKBusinessUnitRepo.GetById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if entry == nil {
		err = variables.NewError(constants.ERR_TENANT_NOTFOUND)
		return
	}
	return
}

func (*IBKBusinessUnit) CreateBusinessUnit(ctx context.Context, authUser *model.AuthUser, body *model.IBKBusinessUnitBody) (id string, err error) {
	businessUnit := &model.IBKBusinessUnit{
		Base:             model.InitBaseModel(authUser.GetID(), authUser.GetID()),
		BusinessUnitName: body.BusinessUnitName,
		TenantId:         authUser.TenantId,
	}
	if err = repository.IBKBusinessUnitRepo.Insert(ctx, *businessUnit); err != nil {
		log.Error(err)
		return
	}
	id = businessUnit.Id
	return
}

func (*IBKBusinessUnit) PutBusinessUnit(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKBusinessUnitBody) (err error) {
	var businessUnit *model.IBKBusinessUnit
	businessUnit, err = repository.IBKBusinessUnitRepo.GetById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if businessUnit == nil {
		err = variables.NewError(constants.ERR_TENANT_NOTFOUND)
		return
	}
	businessUnit.BusinessUnitName = body.BusinessUnitName
	if err = repository.IBKBusinessUnitRepo.Update(ctx, *businessUnit); err != nil {
		log.Error(err)
		return
	}
	return
}

func (*IBKBusinessUnit) DeleteBusinessUnit(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	businessUnit, err := repository.IBKBusinessUnitRepo.GetById(ctx, id)
	if err != nil {
		log.Error(err)
		return
	} else if businessUnit == nil {
		err = variables.NewError("businessUnit not found")
		return
	}
	if err = repository.IBKBusinessUnitRepo.Delete(ctx, id); err != nil {
		log.Error(err)
		return
	}
	return
}

func (*IBKBusinessUnit) FindById(ctx context.Context, id string, conditions ...sql_builder.QueryCondition) (entity *model.IBKBusinessUnit, err error) {
	entity, err = repository.IBKBusinessUnitRepo.GetById(ctx, id, conditions...)
	if err != nil {
		log.Error(err)
		return
	} else if entity == nil {
		err = variables.NewError(constants.ERR_TENANT_NOTFOUND)
		return
	}
	return
}
