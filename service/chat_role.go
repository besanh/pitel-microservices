package service

import (
	"context"
	"time"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
)

type (
	IChatRole interface {
		GetChatRoles(ctx context.Context, filter model.ChatRoleFilter, limit, offset int) (total int, result *[]model.ChatRole, err error)
		GetChatRoleById(ctx context.Context, id string) (chatRole *model.ChatRole, err error)
		InsertChatRole(ctx context.Context, role *model.ChatRoleRequest) (id string, err error)
		UpdateChatRoleById(ctx context.Context, id string, role *model.ChatRole) (err error)
		DeleteChatRoleById(ctx context.Context, id string) (err error)
	}
	ChatRole struct{}
)

var ChatRoleService IChatRole

func NewChatRole() IChatRole {
	return &ChatRole{}
}

func (s *ChatRole) GetChatRoles(ctx context.Context, filter model.ChatRoleFilter, limit, offset int) (total int, result *[]model.ChatRole, err error) {
	total, result, err = repository.ChatRoleRepo.GetChatRoles(ctx, repository.DBConn, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatRole) GetChatRoleById(ctx context.Context, id string) (chatRole *model.ChatRole, err error) {
	chatRole, err = repository.ChatRoleRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return chatRole, nil
}

func (s *ChatRole) InsertChatRole(ctx context.Context, role *model.ChatRoleRequest) (string, error) {
	chatRole := model.ChatRole{
		Base:     model.InitBase(),
		RoleName: role.RoleName,
		Status:   role.Status,
		Setting:  &role.Setting,
	}

	if err := repository.ChatRoleRepo.Insert(ctx, repository.DBConn, chatRole); err != nil {
		log.Error(err)
		return chatRole.Base.GetId(), err
	}
	return chatRole.Base.GetId(), nil
}

func (s *ChatRole) UpdateChatRoleById(ctx context.Context, id string, role *model.ChatRole) (err error) {
	chatRoleExist, err := repository.ChatRoleRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	} else if chatRoleExist == nil {
		log.Error(err)
		return err
	}

	chatRoleExist.RoleName = role.RoleName
	chatRoleExist.Status = role.Status
	chatRoleExist.UpdatedAt = time.Now()
	if err = repository.ChatRoleRepo.Update(ctx, repository.DBConn, *chatRoleExist); err != nil {
		log.Error(err)
		return err
	}
	return
}

func (s *ChatRole) DeleteChatRoleById(ctx context.Context, id string) (err error) {
	chatRoleExist, err := repository.ChatRoleRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	} else if chatRoleExist == nil {
		log.Error(err)
		return err
	}
	err = repository.ChatRoleRepo.Delete(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return
}
