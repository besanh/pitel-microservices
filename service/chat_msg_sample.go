package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatMsgSample interface {
		GetChatMsgSamples(ctx context.Context, authUser *model.AuthUser, filter model.ChatMsgSampleFilter, limit int, offset int) (int, *[]model.ChatMsgSampleView, error)
		GetChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatMsgSampleView, error)
		InsertChatMsgSample(ctx context.Context, authUser *model.AuthUser, cms model.ChatMsgSampleRequest, file *multipart.FileHeader) (string, error)
		UpdateChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string, cms model.ChatMsgSampleRequest, file *multipart.FileHeader) error
		DeleteChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string) error
	}

	ChatMsgSample struct{}
)

var ChatMessageSampleService IChatMsgSample

func NewChatMsgSample() IChatMsgSample {
	return &ChatMsgSample{}
}

func (s *ChatMsgSample) GetChatMsgSamples(ctx context.Context, authUser *model.AuthUser, filter model.ChatMsgSampleFilter, limit int, offset int) (total int, msgSamples *[]model.ChatMsgSampleView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, msgSamples, err = repository.ChatMsgSampleRepo.GetChatMsgSamples(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatMsgSample) GetChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string) (rs *model.ChatMsgSampleView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	rs, err = repository.ChatMsgSampleRepo.GetChatMsgSampleById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}
	if rs == nil {
		log.Error(errors.New("not found chat msg sample"))
		return
	}
	return
}

func (s *ChatMsgSample) InsertChatMsgSample(ctx context.Context, authUser *model.AuthUser, cms model.ChatMsgSampleRequest, file *multipart.FileHeader) (string, error) {
	chatMsgSample := model.ChatMsgSample{
		Base:     model.InitBase(),
		TenantId: authUser.TenantId,
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatMsgSample.Id, err
	}

	// check if connectionApp id exists
	connectionApp, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, cms.ConnectionId)
	if err != nil {
		log.Error(err)
		return chatMsgSample.Id, err
	}
	if connectionApp == nil {
		err = errors.New("not found connection id")
		log.Error(err)
		return chatMsgSample.Id, err
	}
	var oaId string
	if connectionApp.ConnectionType == "zalo" && len(connectionApp.OaInfo.Zalo) > 0 {
		oaId = connectionApp.OaInfo.Zalo[0].OaId
	} else if connectionApp.ConnectionType == "facebook" && len(connectionApp.OaInfo.Facebook) > 0 {
		oaId = connectionApp.OaInfo.Facebook[0].OaId
	}

	var imageUrl string
	if file != nil && len(file.Filename) > 0 {
		imageUrl, err = UploadDoc(ctx, connectionApp.AppId, oaId, file)
		if err != nil {
			log.Error(err)
			return chatMsgSample.Id, err
		}
	}

	chatMsgSample.CreatedBy = authUser.UserId
	chatMsgSample.Channel = cms.Channel
	chatMsgSample.ConnectionId = cms.ConnectionId
	chatMsgSample.Content = cms.Content
	chatMsgSample.Keyword = cms.Keyword
	chatMsgSample.Theme = cms.Theme
	chatMsgSample.ImageUrl = imageUrl
	chatMsgSample.CreatedAt = time.Now()

	err = repository.ChatMsgSampleRepo.Insert(ctx, dbCon, chatMsgSample)
	if err != nil {
		log.Error(err)
		return chatMsgSample.Id, err
	}

	return chatMsgSample.Id, nil
}

func (s *ChatMsgSample) UpdateChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string, chatMsgSampleRequest model.ChatMsgSampleRequest, file *multipart.FileHeader) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatMsgSample, err := repository.ChatMsgSampleRepo.GetChatMsgSampleById(ctx, dbCon, id)
	if err != nil || chatMsgSample == nil {
		err = fmt.Errorf("not found id, err=%v", err)
		log.Error(err)
		return err
	}
	if chatMsgSample.ConnectionApp == nil {
		err = errors.New("not found connection id")
		log.Error(err)
		return err
	}
	var oaId string
	if chatMsgSample.ConnectionApp.ConnectionType == "zalo" && len(chatMsgSample.ConnectionApp.OaInfo.Zalo) > 0 {
		oaId = chatMsgSample.ConnectionApp.OaInfo.Zalo[0].OaId
	} else if chatMsgSample.ConnectionApp.ConnectionType == "facebook" && len(chatMsgSample.ConnectionApp.OaInfo.Facebook) > 0 {
		oaId = chatMsgSample.ConnectionApp.OaInfo.Facebook[0].OaId
	}

	var imageUrl string
	if file != nil && len(file.Filename) > 0 {
		imageUrl, err = UploadDoc(ctx, chatMsgSample.ConnectionApp.AppId, oaId, file)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	chatMsgSample.Keyword = chatMsgSampleRequest.Keyword
	chatMsgSample.Theme = chatMsgSampleRequest.Theme
	chatMsgSample.Content = chatMsgSampleRequest.Content
	chatMsgSample.Channel = chatMsgSampleRequest.Channel
	chatMsgSample.ConnectionId = chatMsgSampleRequest.ConnectionId
	if len(imageUrl) > 0 {
		chatMsgSample.ImageUrl = imageUrl
	}
	chatMsgSample.UpdatedBy = authUser.UserId
	chatMsgSample.UpdatedAt = time.Now()
	err = repository.ChatMsgSampleRepo.UpdateById(ctx, dbCon, *chatMsgSample)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatMsgSample) DeleteChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	chatMsgSample, err := repository.ChatMsgSampleRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatMsgSample == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	err = repository.ChatMsgSampleRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}
