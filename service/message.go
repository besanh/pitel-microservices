package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IMessage interface {
		SendMessageToOTT(ctx context.Context, data model.MessageRequest) error
	}
	Message struct {
	}
)

func NewMessage() IMessage {
	return &Message{}
}

func (s *Message) SendMessageToOTT(ctx context.Context, data model.MessageRequest) error {
	var app model.ChatApp
	appCache := cache.RCache.Get(CHAT_APP)
	if appCache == nil {
		app, err := repository.ChatAppRepo.GetById(ctx, repository.DBConn, data.AppId)
		if err != nil {
			log.Error(err)
			return err
		}
		if err = cache.RCache.Set(CHAT_APP, app, CHAT_APP_EXPIRE); err != nil {
			log.Error(err)
			return err
		}
	} else {
		err := json.Unmarshal([]byte(appCache.(string)), &app)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	// var messageType string
	// if app.InfoApp.Zalo.Status {
	// 	messageType = "zalo"
	// } else if app.InfoApp.Facebook.Status {
	// 	messageType = "facebook"
	// }

	// total, previousMessage, err := repository.GetMessage(ctx, repository.DBConn, model.MessageFilter{

	// })

	// message := model.Message{
	// 	ParentMessageId: data.ParentMessageId,
	// 	ConversationId:  data.ConversationId,
	// 	MessageType:     messageType,
	// 	EventName:       data.EventName,
	// 	AppId:           app.Id,
	// 	// UserIdByApp: ,
	// }
	return nil
}
