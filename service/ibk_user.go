package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/repository/sql_builder"
)

type (
	IIBKUser interface {
		GetUsers(c context.Context, authUser *model.AuthUser, request model.IBKUserQueryParam, limit, offset int) (total int, result []*model.IBKUserInfo, err error)
		GetUserById(ctx context.Context, authUser *model.AuthUser, id string) (entry *model.IBKUserInfo, err error)
		CreateUser(ctx context.Context, authUser *model.AuthUser, body *model.IBKUserBody) (id string, err error)
		PutUser(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKUserBody) (err error)
		DeleteUser(ctx context.Context, authUser *model.AuthUser, id string) (err error)
		GetUserProfile(ctx context.Context, authUser *model.AuthUser) (entry *model.IBKUserInfo, err error)
		PostQueryUsers(ctx context.Context, authUser *model.AuthUser, filter *model.GenericFilter, limit int, offset int) (entries []*model.IBKUser, total int, err error)
	}
	IBKUser struct{}
)

var IBKUserService IIBKUser

func NewIBKUser() IIBKUser {
	return &IBKUser{}
}

func (s *IBKUser) GetUsers(c context.Context, authUser *model.AuthUser, request model.IBKUserQueryParam, limit, offset int) (total int, result []*model.IBKUserInfo, err error) {
	total, result, err = repository.IBKUserRepo.GetInfoByQuery(c, request, limit, offset, "")
	if err != nil {
		log.ErrorContext(c, err)
		return
	}
	return
}

func (s *IBKUser) GetUserById(ctx context.Context, authUser *model.AuthUser, id string) (entry *model.IBKUserInfo, err error) {
	entry, err = repository.IBKUserRepo.GetInfoById(ctx, id, model.IBKUserQueryParam{})
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	} else if entry == nil {
		err = variables.NewError(constants.ERR_USER_NOTFOUND)
		return
	}
	return
}

func (s *IBKUser) CreateUser(ctx context.Context, authUser *model.AuthUser, body *model.IBKUserBody) (id string, err error) {
	salt := util.GenerateRandomString(8, nil)
	body.Password = fmt.Sprintf("%s$%s", salt, hashSalt(salt, body.Password))

	user := &model.IBKUser{
		Base:           model.InitBaseModel(authUser.GetID(), authUser.GetID()),
		Username:       body.Username,
		Password:       body.Password,
		Level:          body.Level,
		BusinessUnitId: body.BusinessUnitId,
		Fullname:       body.Fullname,
		Email:          body.Email,
		IsActivated:    body.IsActivated,
		IsLocked:       body.IsLocked,
		IsSentEmail:    body.IsSentEmail,
	}
	if authUser.GetLevel() == constants.ADMIN {
		user.Level = constants.USER
	}
	if len(user.BusinessUnitId) < 1 {
		user.BusinessUnitId = authUser.BusinessUnitId
	}
	var userExisted *model.IBKUser
	userExisted, err = repository.IBKUserRepo.GetUserByUsername(ctx, body.Username)
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	} else if userExisted != nil {
		err = errors.New("username is existed")
		log.ErrorContext(ctx, err)
		return
	}

	if len(body.RoleId) > 0 {
		// check role in body is existed
		var role *model.IBKRole
		role, err = repository.IBKRoleRepo.GetById(ctx, body.RoleId)
		if err != nil {
			log.ErrorContext(ctx, err)
			return
		} else if role == nil {
			err = variables.NewError(constants.ERR_ROLE_NOTFOUND)
			return
		}
		user.RoleId = role.Id
	} else if len(body.Scopes) > 0 {
		user.Scopes = body.Scopes
	}

	if err = repository.IBKUserRepo.Insert(ctx, *user); err != nil {
		log.ErrorContext(ctx, err)
		return
	}

	id = user.Id
	return
}

func (s *IBKUser) PutUser(ctx context.Context, authUser *model.AuthUser, id string, body *model.IBKUserBody) (err error) {
	user, err := repository.IBKUserRepo.GetById(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	} else if user == nil {
		err = variables.NewError(constants.ERR_USER_NOTFOUND)
		return
	}

	{
		var userExisted *model.IBKUser
		if userExisted, err = repository.IBKUserRepo.GetUserByUsername(ctx, body.Username); err != nil {
			log.ErrorContext(ctx, err)
			return
		} else if userExisted != nil && userExisted.Id != id {
			err = errors.New("username is existed")
			return
		}
	}

	if body.Password != constants.DEFAULT_PASSWORD {
		salt := util.GenerateRandomString(8, nil)
		body.Password = fmt.Sprintf("%s$%s", salt, hashSalt(salt, body.Password))
	}

	if len(body.RoleId) > 0 {
		// check role in body is existed
		var role *model.IBKRole
		role, err = repository.IBKRoleRepo.GetById(ctx, body.RoleId)
		if err != nil {
			return
		} else if role == nil {
			err = variables.NewError(constants.ERR_ROLE_NOTFOUND)
			return
		}
		user.RoleId = role.Id
	} else if len(body.Scopes) > 0 {
		user.Scopes = body.Scopes
	}

	user.UpdatedAt = time.Now()
	user.Email = body.Email
	user.Fullname = body.Fullname
	user.IsActivated = body.IsActivated
	user.IsSentEmail = body.IsSentEmail
	user.IsLocked = body.IsLocked

	if err = repository.IBKUserRepo.Update(ctx, *user); err != nil {
		log.ErrorContext(ctx, err)
		return
	}
	return
}

func (s *IBKUser) DeleteUser(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	user, err := repository.IBKUserRepo.GetById(ctx, id)
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	} else if user == nil {
		err = variables.NewError(constants.ERR_USER_NOTFOUND)
		return
	}

	if err = repository.IBKUserRepo.Delete(ctx, id); err != nil {
		log.ErrorContext(ctx, err)
		return
	}
	return
}

func (s *IBKUser) GetUserProfile(ctx context.Context, authUser *model.AuthUser) (entry *model.IBKUserInfo, err error) {
	entry, err = repository.IBKUserRepo.GetInfoById(ctx, authUser.GetID(), model.IBKUserQueryParam{})
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	} else if entry == nil {
		err = variables.NewError(constants.ERR_USER_NOTFOUND)
		return
	}

	// Get tenant info to get meta data
	tenant, err := repository.IBKTenantRepo.GetById(ctx, entry.TenantId)
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	} else if tenant == nil {
		err = variables.NewError(constants.ERR_TENANT_NOTFOUND)
		return
	}

	return
}

func (s *IBKUser) PostQueryUsers(ctx context.Context, authUser *model.AuthUser, body *model.GenericFilter, limit int, offset int) (entries []*model.IBKUser, total int, err error) {
	total, businessUnits, err := repository.IBKBusinessUnitRepo.SelectByQuery(ctx, []sql_builder.QueryCondition{
		sql_builder.EqualQuery("tenant_id", authUser.TenantId),
	}, limit, offset, "")
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	} else if len(businessUnits) < 1 {
		err = variables.NewError(constants.ERR_BUSINESS_UNIT_NOTFOUND)
		return
	}
	businessUnitIds := make([]string, 0)
	for _, businessUnit := range businessUnits {
		businessUnitIds = append(businessUnitIds, businessUnit.Id)
	}
	conditions := make([]sql_builder.QueryCondition, 0)
	conditions = append(conditions, sql_builder.QueryCondition{
		Field:    "business_unit_id",
		Value:    businessUnitIds,
		Operator: "IN",
	})

	condition := sql_builder.QueryCondition{}
	if err = util.ParseAnyToAny(body, &condition); err != nil {
		log.ErrorContext(ctx, err)
		return
	}

	conditions = append(conditions, condition)
	total, entries, err = repository.IBKUserRepo.SelectByQuery(ctx, conditions, limit, offset, "")
	if err != nil {
		log.ErrorContext(ctx, err)
		return
	}
	return
}
