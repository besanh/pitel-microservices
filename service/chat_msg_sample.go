package service

import (
	"context"
	"errors"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/storage"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"io"
	"mime/multipart"
	"time"
)

type (
	IChatMsgSample interface {
		GetChatMsgSamples(ctx context.Context, authUser *model.AuthUser, limit int, offset int) (int, *[]model.ChatMsgSampleView, error)
		GetChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatMsgSampleView, error)
		InsertChatMsgSample(ctx context.Context, authUser *model.AuthUser, cms model.ChatMsgSampleRequest, file *multipart.FileHeader) (string, error)
		UpdateChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string, cms model.ChatMsgSampleRequest, file *multipart.FileHeader) error
		DeleteChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string) error
	}

	ChatMsgSample struct{}
)

func NewChatMsgSample() IChatMsgSample {
	return &ChatMsgSample{}
}

func (s *ChatMsgSample) GetChatMsgSamples(ctx context.Context, authUser *model.AuthUser, limit int, offset int) (total int, msgSamples *[]model.ChatMsgSampleView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, msgSamples, err = repository.ChatMsgSampleRepo.GetChatMsgSamples(ctx, dbCon, limit, offset)
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
		Base: model.InitBase(),
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

	var imageUrl string
	if file != nil && len(file.Filename) > 0 {
		imageUrl, err = uploadImageToStorageChatMsgSample(ctx, file)
		if err != nil {
			log.Error(err)
			return chatMsgSample.Id, err
		}
	}

	chatMsgSample.CreatedBy = authUser.UserId
	chatMsgSample.UpdatedBy = authUser.UserId
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

func (s *ChatMsgSample) UpdateChatMsgSampleById(ctx context.Context, authUser *model.AuthUser, id string, cms model.ChatMsgSampleRequest, file *multipart.FileHeader) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
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

	var imageUrl string
	if file != nil {
		imageUrl, err = uploadImageToStorageChatMsgSample(ctx, file)
		if err != nil {
			log.Error(err)
			return err
		}

		if len(chatMsgSample.ImageUrl) > 0 {
			err = removeImageFromStorageChatMsgSample(ctx, chatMsgSample.ImageUrl)
			if err != nil {
				log.Error(err)
				//remove image just uploaded
				if err = removeImageFromStorageChatMsgSample(ctx, imageUrl); err != nil {
					log.Error(err)
				}
				return err
			}
		}
	}

	if len(cms.Keyword) > 0 {
		chatMsgSample.Keyword = cms.Keyword
	}
	if len(cms.Theme) > 0 {
		chatMsgSample.Theme = cms.Theme
	}
	if len(cms.Content) > 0 {
		chatMsgSample.Content = cms.Content
	}
	if len(imageUrl) > 0 {
		chatMsgSample.ImageUrl = imageUrl
	}
	chatMsgSample.UpdatedBy = authUser.UserId
	chatMsgSample.UpdatedAt = time.Now()
	err = repository.ChatMsgSampleRepo.Update(ctx, dbCon, *chatMsgSample)
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

	if len(chatMsgSample.ImageUrl) > 0 {
		err = removeImageFromStorageChatMsgSample(ctx, chatMsgSample.ImageUrl)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	err = repository.ChatMsgSampleRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func uploadImageToStorageChatMsgSample(c context.Context, file *multipart.FileHeader) (url string, err error) {
	f, err := file.Open()
	if err != nil {
		log.Error(err)
		return
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		log.Error(err)
		return
	}
	metaData := storage.NewStoreInput(fileBytes, file.Filename)
	isSuccess, err := storage.Instance.Store(c, *metaData)
	if err != nil || !isSuccess {
		log.Error(err)
		return
	}

	input := storage.NewRetrieveInput(file.Filename)
	_, err = storage.Instance.Retrieve(c, *input)
	if err != nil {
		log.Error(err)
		return
	}

	url = API_DOC + "/bss-message/v1/chat-msg-sample/image/" + input.Path

	return
}

func removeImageFromStorageChatMsgSample(c context.Context, fileName string) error {
	input := storage.NewRetrieveInput(fileName)
	return storage.Instance.RemoveFile(c, *input)
}
