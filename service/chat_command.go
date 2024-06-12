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
	IChatCommand interface {
		GetChatCommands(ctx context.Context, authUser *model.AuthUser, limit int, offset int) (int, []model.ChatCommandView, error)
		GetChatCommandById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatCommand, error)
		InsertChatCommand(ctx context.Context, authUser *model.AuthUser, cmd model.ChatCommandRequest, file *multipart.FileHeader) (string, error)
		UpdateChatCommandById(ctx context.Context, authUser *model.AuthUser, id string, cmd model.ChatCommandRequest, file *multipart.FileHeader) error
		DeleteChatCommandById(ctx context.Context, authUser *model.AuthUser, id string) error
	}

	ChatCommand struct{}
)

func NewChatCommand() IChatCommand {
	return &ChatCommand{}
}

func (s *ChatCommand) GetChatCommands(ctx context.Context, authUser *model.AuthUser, limit int, offset int) (total int, commands []model.ChatCommandView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, commands, err = repository.ChatCommandRepo.GetChatCommands(ctx, dbCon, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatCommand) GetChatCommandById(ctx context.Context, authUser *model.AuthUser, id string) (rs *model.ChatCommand, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	rs, err = repository.ChatCommandRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatCommand) InsertChatCommand(ctx context.Context, authUser *model.AuthUser, cmd model.ChatCommandRequest, file *multipart.FileHeader) (string, error) {
	chatCommand := model.ChatCommand{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatCommand.Id, err
	}

	// check if page id exists
	page, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, cmd.PageId)
	if err != nil {
		log.Error(err)
		return chatCommand.Id, err
	}
	if page == nil {
		err = errors.New("not found page id")
		log.Error(err)
		return chatCommand.Id, err
	}

	var imageUrl string
	if file != nil {
		imageUrl, err = uploadImageToStorageChatCommand(ctx, file)
		if err != nil {
			log.Error(err)
			return chatCommand.Id, err
		}

	}

	chatCommand.CreatorId = authUser.UserId
	chatCommand.Channel = cmd.Channel
	chatCommand.PageId = cmd.PageId
	chatCommand.Content = cmd.Content
	chatCommand.Keyword = cmd.Keyword
	chatCommand.Theme = cmd.Theme
	chatCommand.ImageUrl = imageUrl
	chatCommand.CreatedAt = time.Now()

	err = repository.ChatCommandRepo.Insert(ctx, dbCon, chatCommand)
	if err != nil {
		log.Error(err)
		return chatCommand.Id, err
	}

	return chatCommand.Id, nil
}

func (s *ChatCommand) UpdateChatCommandById(ctx context.Context, authUser *model.AuthUser, id string, cmd model.ChatCommandRequest, file *multipart.FileHeader) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatCommand, err := repository.ChatCommandRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatCommand == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	var imageUrl string
	if file != nil {
		imageUrl, err = uploadImageToStorageChatCommand(ctx, file)
		if err != nil {
			log.Error(err)
			return err
		}
		err = removeImageFromStorageChatCommand(ctx, chatCommand.ImageUrl)
		if err != nil {
			log.Error(err)
			//remove image just uploaded
			if err = removeImageFromStorageChatCommand(ctx, imageUrl); err != nil {
				log.Error(err)
			}
			return err
		}
	}

	chatCommand.Keyword = cmd.Keyword
	chatCommand.Theme = cmd.Theme
	chatCommand.Content = cmd.Content
	chatCommand.UpdatedAt = time.Now()
	chatCommand.ImageUrl = imageUrl
	err = repository.ChatCommandRepo.Update(ctx, dbCon, *chatCommand)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatCommand) DeleteChatCommandById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	chatCommand, err := repository.ChatCommandRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatCommand == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	err = removeImageFromStorageChatCommand(ctx, chatCommand.ImageUrl)
	if err != nil {
		log.Error(err)
		return err
	}

	err = repository.ChatCommandRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func uploadImageToStorageChatCommand(c context.Context, file *multipart.FileHeader) (url string, err error) {
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

	url = API_DOC + "/bss-message/v1/chat-command/image/" + input.Path

	return
}

func removeImageFromStorageChatCommand(c context.Context, fileName string) error {
	input := storage.NewRetrieveInput(fileName)
	return storage.Instance.RemoveFile(c, *input)
}
