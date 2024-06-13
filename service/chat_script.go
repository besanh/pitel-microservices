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
	IChatScript interface {
		GetChatScripts(ctx context.Context, authUser *model.AuthUser, limit int, offset int) (int, *[]model.ChatScriptView, error)
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

func (s *ChatScript) GetChatScripts(ctx context.Context, authUser *model.AuthUser, limit int, offset int) (total int, chatScripts *[]model.ChatScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, chatScripts, err = repository.ChatScriptRepo.GetChatScripts(ctx, dbCon, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	if chatScripts == nil {
		return
	}

	return
}

func (s *ChatScript) GetChatScriptById(ctx context.Context, authUser *model.AuthUser, id string) (rs *model.ChatScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	rs, err = repository.ChatScriptRepo.GetChatScriptById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}
	if rs == nil {
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
	case "TEXT":
		chatScript.Content = csr.Content
	case "IMAGE", "FILE":
		var fileUrl string
		if file != nil && len(file.Filename) > 0 {
			fileUrl, err = uploadFileToStorageChatScript(ctx, file)
			if err != nil {
				log.Error(err)
				return chatScript.Id, err
			}
		}
		chatScript.FileUrl = fileUrl
	case "OTHER":
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
	chatScript.UpdatedBy = authUser.UserId
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
	case "TEXT":
		chatScript.Content = csr.Content
	case "IMAGE", "FILE":
		var fileUrl string
		if file != nil && len(file.Filename) > 0 {
			fileUrl, err = uploadFileToStorageChatScript(ctx, file)
			if err != nil {
				log.Error(err)
				return err
			}

			if len(chatScript.FileUrl) > 0 {
				err = removeImageFromStorageChatScript(ctx, chatScript.FileUrl)
				if err != nil {
					log.Error(err)
					//remove image just uploaded
					if err = removeImageFromStorageChatScript(ctx, fileUrl); err != nil {
						log.Error(err)
					}
					return err
				}
			}
		}
		chatScript.FileUrl = fileUrl
	case "OTHER":
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
		err = removeImageFromStorageChatScript(ctx, chatScript.FileUrl)
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

func uploadFileToStorageChatScript(c context.Context, file *multipart.FileHeader) (url string, err error) {
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

	url = API_DOC + "/bss-message/v1/chat-script/image/" + input.Path

	return
}

func removeImageFromStorageChatScript(c context.Context, fileName string) error {
	input := storage.NewRetrieveInput(fileName)
	return storage.Instance.RemoveFile(c, *input)
}
