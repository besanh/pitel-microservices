package service

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatScript interface {
		GetChatScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatScriptFilter, limit int, offset int) (int, *[]model.ChatScriptView, error)
		GetChatScriptById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatScriptView, error)
		InsertChatScript(ctx context.Context, authUser *model.AuthUser, chatScriptRequest model.ChatScriptRequest, file *multipart.FileHeader) (string, error)
		UpdateChatScriptById(ctx context.Context, authUser *model.AuthUser, id string, chatScriptRequest model.ChatScriptRequest, file *multipart.FileHeader) error
		UpdateChatScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, status sql.NullBool) error
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

func (s *ChatScript) InsertChatScript(ctx context.Context, authUser *model.AuthUser, chatScriptRequest model.ChatScriptRequest, file *multipart.FileHeader) (string, error) {
	chatScript := model.ChatScript{
		Base:     model.InitBase(),
		TenantId: authUser.TenantId,
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatScript.Id, err
	}

	switch chatScriptRequest.ScriptType {
	case "text":
		chatScript.Content = chatScriptRequest.Content
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
		if len(chatScriptRequest.OtherScriptId) > 0 {
			chatScript.OtherScriptId = chatScriptRequest.OtherScriptId
		}
	default:
		err = errors.New("invalid script type")
		log.Error(err)
		return chatScript.Id, err
	}

	statusTmp := chatScriptRequest.Status
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}
	if status.Valid {
		chatScript.Status = status.Bool
	}

	chatScript.ScriptType = chatScriptRequest.ScriptType
	chatScript.ScriptName = chatScriptRequest.ScriptName
	chatScript.CreatedBy = authUser.UserId
	chatScript.Channel = chatScriptRequest.Channel
	chatScript.CreatedAt = time.Now()

	err = repository.ChatScriptRepo.Insert(ctx, dbCon, chatScript)
	if err != nil {
		log.Error(err)
		return chatScript.Id, err
	}

	return chatScript.Id, nil
}

func (s *ChatScript) UpdateChatScriptById(ctx context.Context, authUser *model.AuthUser, id string, chatScriptRequest model.ChatScriptRequest, file *multipart.FileHeader) error {
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

	// request changes script type
	if chatScript.ScriptType != chatScriptRequest.ScriptType {
		chatScript.Content = ""
		chatScript.OtherScriptId = ""
		if chatScript.ScriptType == "image" || chatScript.ScriptType == "file" {
			if len(chatScript.FileUrl) > 0 {
				if err = removeFileFromStorageShareInfo(ctx, chatScript.FileUrl); err != nil {
					log.Error(err)
					return err
				}
				chatScript.FileUrl = ""
			}
		}
	}

	switch chatScriptRequest.ScriptType {
	case "text":
		if len(chatScriptRequest.Content) > 0 {
			chatScript.Content = chatScriptRequest.Content
		}
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
		if len(fileUrl) > 0 {
			chatScript.FileUrl = fileUrl
		}
	case "other":
		if len(chatScriptRequest.OtherScriptId) > 0 {
			chatScript.OtherScriptId = chatScriptRequest.OtherScriptId
		}
	default:
		err = errors.New("invalid script type")
		log.Error(err)
		return err
	}

	statusTmp := chatScriptRequest.Status
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}
	if status.Valid {
		chatScript.Status = status.Bool
	}
	if len(chatScriptRequest.ScriptName) > 0 {
		chatScript.ScriptName = chatScriptRequest.ScriptName
	}
	chatScript.ScriptType = chatScriptRequest.ScriptType
	chatScript.Channel = chatScriptRequest.Channel
	chatScript.UpdatedBy = authUser.UserId
	chatScript.UpdatedAt = time.Now()
	err = repository.ChatScriptRepo.Update(ctx, dbCon, *chatScript)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatScript) UpdateChatScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, status sql.NullBool) error {
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

	if status.Valid {
		chatScript.Status = status.Bool
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
