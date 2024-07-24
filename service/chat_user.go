package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatUser interface {
		InsertChatUser(ctx context.Context, request *model.ChatUserRequest) (id string, err error)
	}
	ChatUser struct {
	}
)

var ChatUserService IChatUser

func NewChatUser() IChatUser {
	return &ChatUser{}
}

func (s *ChatUser) InsertChatUser(ctx context.Context, request *model.ChatUserRequest) (id string, err error) {
	chatUser := model.ChatUser{
		Base: model.InitBase(),
		Salt: uuid.NewString(),
	}

	_, chatUsers, err := repository.ChatUserRepo.GetChatUsers(ctx, repository.DBConn, model.ChatUserFilter{
		Username: request.Username,
	}, 1, 0)
	if err != nil {
		log.Error(err)
		return
	} else if len(*chatUsers) > 0 {
		err = fmt.Errorf("username already exists")
		log.Error(err)
		return
	}

	roleExist, err := repository.ChatRoleRepo.GetById(ctx, repository.DBConn, request.RoleId)
	if err != nil {
		log.Error(err)
		return
	} else if roleExist == nil {
		err = errors.New("role " + request.RoleId + " not found")
		log.Error(err)
		return
	}

	chatUser.Username = request.Username
	chatUser.Fullname = request.Fullname
	dbPassword := []byte(chatUser.Salt + request.Password)
	chatUser.Password = fmt.Sprintf("%x", md5.Sum(dbPassword))
	chatUser.Email = request.Email
	chatUser.Avatar = request.Avatar
	chatUser.Level = request.Level
	chatUser.Status = request.Status
	chatUser.RoleId = request.RoleId
	chatUser.CreatedAt = time.Now()
	if err = repository.ChatUserRepo.Insert(ctx, repository.DBConn, chatUser); err != nil {
		log.Error(err)
		return
	}

	id = chatUser.GetId()
	return
}
