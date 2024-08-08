package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IIBKUser interface {
		GetUsers(c context.Context, authUser *model.AuthUser, request model.IBKUserQueryParam, limit, offset int) (total int, result []*model.IBKUserInfo, err error)
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
