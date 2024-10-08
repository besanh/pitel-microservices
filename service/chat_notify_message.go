package service

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatNotifyMessage interface {
		GetChatNotifyMessages(ctx context.Context, authUser *model.AuthUser, filter model.ChatNotifyMessageFilter, limit, offset int) (int, *[]model.ChatNotifyMessage, error)
		InsertChatNotifyMessage(ctx context.Context, authUser *model.AuthUser, request model.ChatNotifyMessageRequest) (string, error)
		GetChatNotifyMessageById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatNotifyMessage, error)
		UpdateChatNotifyMessageById(ctx context.Context, authUser *model.AuthUser, id string, request model.ChatNotifyMessageRequest) error
		DeleteChatNotifyMessageById(ctx context.Context, authUser *model.AuthUser, id string) error
	}
	ChatNotifyMessage struct{}
)

var ChatNotifyMessageService IChatNotifyMessage

func NewChatNotifyMessage() IChatNotifyMessage {
	return &ChatNotifyMessage{}
}

func (s *ChatNotifyMessage) GetChatNotifyMessages(ctx context.Context, authUser *model.AuthUser, filter model.ChatNotifyMessageFilter, limit, offset int) (total int, result *[]model.ChatNotifyMessage, err error) {
	filter.TenantId = authUser.TenantId
	total, result, err = repository.ChatNotifyMessageRepo.GetChatNotifyMessages(ctx, repository.DBConn, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatNotifyMessage) InsertChatNotifyMessage(ctx context.Context, authUser *model.AuthUser, request model.ChatNotifyMessageRequest) (id string, err error) {
	chatNotifyMessage := model.ChatNotifyMessage{
		Base: model.InitBase(),
	}
	id = chatNotifyMessage.Id

	// check if connectionApp id exists
	connectionApp, err := repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, request.ConnectionId)
	if err != nil {
		log.Error(err)
		return
	}
	if connectionApp == nil {
		err = errors.New("not found connection id")
		log.Error(err)
		return
	}

	filter := model.ChatNotifyMessageFilter{
		TenantId:     authUser.TenantId,
		NotifyType:   []string{request.NotifyType},
		ConnectionId: request.ConnectionId,
	}
	_, chatNotifyMessages, err := repository.ChatNotifyMessageRepo.GetChatNotifyMessages(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*chatNotifyMessages) > 0 {
		err = errors.New("chat notify message connection id " + request.ConnectionId + " has already existed")
		log.Error(err)
		return
	}
	chatNotifyMessage.TenantId = authUser.TenantId
	chatNotifyMessage.ConnectionId = request.ConnectionId
	chatNotifyMessage.MessageNotifyAfter = request.MessageNotifyAfter
	chatNotifyMessage.NotifyType = request.NotifyType
	chatNotifyMessage.ReceiverType = request.ReceiverType
	if connectionApp.ConnectionType == "facebook" {
		chatNotifyMessage.OaId = connectionApp.OaInfo.Facebook[0].OaId
	} else if connectionApp.ConnectionType == "zalo" {
		chatNotifyMessage.OaId = connectionApp.OaInfo.Zalo[0].OaId
	}

	if err = repository.ChatNotifyMessageRepo.Insert(ctx, repository.DBConn, chatNotifyMessage); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatNotifyMessage) GetChatNotifyMessageById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ChatNotifyMessage, err error) {
	result, err = repository.ChatNotifyMessageRepo.GetChatNotifyMessageById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatNotifyMessage) UpdateChatNotifyMessageById(ctx context.Context, authUser *model.AuthUser, id string, request model.ChatNotifyMessageRequest) (err error) {
	// check if connectionApp id exists
	connectionApp, err := repository.ChatConnectionAppRepo.GetById(ctx, repository.DBConn, request.ConnectionId)
	if err != nil {
		log.Error(err)
		return
	}
	if connectionApp == nil {
		err = errors.New("not found connection id")
		log.Error(err)
		return
	}

	chatNotifyMessageExist, err := repository.ChatNotifyMessageRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if chatNotifyMessageExist == nil {
		err = errors.New("chat notify message " + id + " not found")
		log.Error(err)
		return
	}
	filter := model.ChatNotifyMessageFilter{
		TenantId:     authUser.TenantId,
		NotifyType:   []string{request.NotifyType},
		ConnectionId: request.ConnectionId,
	}
	_, chatNotifyMessages, err := repository.ChatNotifyMessageRepo.GetChatNotifyMessages(ctx, repository.DBConn, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*chatNotifyMessages) > 0 && chatNotifyMessageExist.NotifyType != (*chatNotifyMessages)[0].NotifyType {
		err = errors.New("chat notify message connection id " + request.ConnectionId + " has already existed")
		log.Error(err)
		return
	}
	chatNotifyMessageExist.ConnectionId = request.ConnectionId
	chatNotifyMessageExist.MessageNotifyAfter = request.MessageNotifyAfter
	chatNotifyMessageExist.NotifyType = request.NotifyType
	chatNotifyMessageExist.ReceiverType = request.ReceiverType
	if connectionApp.ConnectionType == "facebook" {
		chatNotifyMessageExist.OaId = connectionApp.OaInfo.Facebook[0].OaId
	} else if connectionApp.ConnectionType == "zalo" {
		chatNotifyMessageExist.OaId = connectionApp.OaInfo.Zalo[0].OaId
	}

	if err = repository.ChatNotifyMessageRepo.Update(ctx, repository.DBConn, *chatNotifyMessageExist); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *ChatNotifyMessage) DeleteChatNotifyMessageById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	chatNotifyMessageExist, err := repository.ChatNotifyMessageRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		log.Error(err)
		return
	} else if chatNotifyMessageExist == nil {
		err = errors.New("chat notify message " + id + " not found")
		log.Error(err)
		return
	}
	if err = repository.ChatNotifyMessageRepo.Delete(ctx, repository.DBConn, id); err != nil {
		log.Error(err)
		return
	}
	return
}
