package service

import (
	"context"
	"errors"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"mime/multipart"
	"time"
)

type (
	IChatScript interface {
		GetChatScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatScriptFilter, limit int, offset int) (int, *[]model.ChatScriptView, error)
		GetChatScriptById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatScriptView, error)
		InsertChatScript(ctx context.Context, authUser *model.AuthUser, csr model.ChatScriptRequest, file *multipart.FileHeader) (string, error)
		UpdateChatScriptById(ctx context.Context, authUser *model.AuthUser, id string, csr model.ChatScriptRequest, file *multipart.FileHeader) error
		UpdateChatScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, oldStatus string) error
		DeleteChatScriptById(ctx context.Context, authUser *model.AuthUser, id string) error
	}

	ChatScript struct{}
)

func NewChatScript() IChatScript {
	return &ChatScript{}
}

func (s *ChatScript) GetChatScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatScriptFilter, limit int, offset int) (total int, chatScripts *[]model.ChatScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, chatScripts, err = repository.ChatScriptRepo.GetChatScripts(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatScript) GetChatScriptById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ChatScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	result, err = repository.ChatScriptRepo.GetChatScriptById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}
	if result == nil {
		log.Error(errors.New("not found chat script config"))
		return
	}

	return
}

func (s *ChatScript) InsertChatScript(ctx context.Context, authUser *model.AuthUser, csr model.ChatScriptRequest, file *multipart.FileHeader) (string, error) {
	chatScript := model.ChatScript{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatScript.Id, err
	}

	// check if connectionApp id exists
	connectionApp, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, csr.ConnectionId)
	if err != nil {
		log.Error(err)
		return chatScript.Id, err
	}
	if connectionApp == nil {
		err = errors.New("not found connection id")
		log.Error(err)
		return chatScript.Id, err
	}

	switch csr.ScriptType {
	case "text":
		chatScript.Content = csr.Content
	case "image", "file":
		var fileUrl string
		if file != nil && len(file.Filename) > 0 {
			fileUrl, err = uploadImageToStorageShareInfo(ctx, file)
			if err != nil {
				log.Error(err)
				return chatScript.Id, err
			}
		}
		chatScript.FileUrl = fileUrl
	case "other":
		chatScript.OtherScriptId = csr.OtherScriptId
	default:
		err = errors.New("invalid script type")
		log.Error(err)
		return chatScript.Id, err
	}

	if csr.Status == "true" {
		chatScript.Status = true
	}

	chatScript.ScriptType = csr.ScriptType
	chatScript.ScriptName = csr.ScriptName
	chatScript.CreatedBy = authUser.UserId
	chatScript.Channel = csr.Channel
	chatScript.ConnectionId = csr.ConnectionId
	chatScript.CreatedAt = time.Now()

	err = repository.ChatScriptRepo.Insert(ctx, dbCon, chatScript)
	if err != nil {
		log.Error(err)
		return chatScript.Id, err
	}

	return chatScript.Id, nil
}

func (s *ChatScript) UpdateChatScriptById(ctx context.Context, authUser *model.AuthUser, id string, csr model.ChatScriptRequest, file *multipart.FileHeader) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatScript, err := repository.ChatScriptRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatScript == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	if csr.ScriptType != chatScript.ScriptType {
		err = errors.New("not meet script type")
		log.Error(err)
		return err
	}

	switch csr.ScriptType {
	case "text":
		chatScript.Content = csr.Content
	case "image", "file":
		var fileUrl string
		if file != nil && len(file.Filename) > 0 {
			fileUrl, err = uploadImageToStorageShareInfo(ctx, file)
			if err != nil {
				log.Error(err)
				return err
			}

			if len(chatScript.FileUrl) > 0 {
				err = removeFileFromStorageShareInfo(ctx, chatScript.FileUrl)
				if err != nil {
					log.Error(err)
					//remove image just uploaded
					if err = removeFileFromStorageShareInfo(ctx, fileUrl); err != nil {
						log.Error(err)
					}
					return err
				}
			}
		}
		chatScript.FileUrl = fileUrl
	case "other":
		chatScript.OtherScriptId = csr.OtherScriptId
	default:
		err = errors.New("invalid script type")
		log.Error(err)
		return err
	}

	if len(csr.ScriptName) > 0 {
		chatScript.ScriptName = csr.ScriptName
	}
	chatScript.UpdatedBy = authUser.UserId
	chatScript.UpdatedAt = time.Now()
	err = repository.ChatScriptRepo.Update(ctx, dbCon, *chatScript)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatScript) UpdateChatScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, oldStatus string) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatScript, err := repository.ChatScriptRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatScript == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	var status bool
	if oldStatus == "true" {
		status = true
	}
	chatScript.Status = !status
	chatScript.UpdatedBy = authUser.UserId
	chatScript.UpdatedAt = time.Now()
	err = repository.ChatScriptRepo.Update(ctx, dbCon, *chatScript)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatScript) DeleteChatScriptById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	chatScript, err := repository.ChatScriptRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatScript == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	if len(chatScript.FileUrl) > 0 {
		err = removeFileFromStorageShareInfo(ctx, chatScript.FileUrl)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	err = repository.ChatScriptRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}
