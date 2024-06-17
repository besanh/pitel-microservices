package service

import (
	"context"
	"errors"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"time"
)

type (
	IChatAutoScript interface {
		GetChatAutoScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatAutoScriptFilter, limit int, offset int) (int, *[]model.ChatAutoScriptView, error)
		GetChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatAutoScriptView, error)
		InsertChatAutoScript(ctx context.Context, authUser *model.AuthUser, chatAutoScriptRequest model.ChatAutoScriptRequest) (string, error)
		UpdateChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string, chatAutoScriptRequest model.ChatAutoScriptRequest) error
		UpdateChatAutoScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, oldStatus string) error
		DeleteChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string) error
	}

	ChatAutoScript struct{}
)

func NewChatAutoScript() IChatAutoScript {
	return &ChatAutoScript{}
}

func (s *ChatAutoScript) GetChatAutoScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatAutoScriptFilter, limit int, offset int) (total int, chatScripts *[]model.ChatAutoScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, chatScripts, err = repository.ChatAutoScriptRepo.GetChatAutoScripts(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (s *ChatAutoScript) GetChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string) (result *model.ChatAutoScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	result, err = repository.ChatAutoScriptRepo.GetChatAutoScriptById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}
	if result == nil {
		log.Error(errors.New("not found chat auto script config"))
		return
	}

	return
}

func (s *ChatAutoScript) InsertChatAutoScript(ctx context.Context, authUser *model.AuthUser, chatAutoScriptRequest model.ChatAutoScriptRequest) (string, error) {
	chatAutoScript := model.ChatAutoScript{
		Base: model.InitBase(),
	}
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return chatAutoScript.Id, err
	}

	// check if connectionApp id exists
	connectionApp, err := repository.ChatConnectionAppRepo.GetById(ctx, dbCon, chatAutoScriptRequest.ConnectionId)
	if err != nil {
		log.Error(err)
		return chatAutoScript.Id, err
	}
	if connectionApp == nil {
		err = errors.New("not found connection id")
		log.Error(err)
		return chatAutoScript.Id, err
	}

	// handle actions' content
	for i, action := range chatAutoScriptRequest.ActionScript.Actions {
		switch model.ScriptActionType(action.Type) {
		case model.MoveToExistedScript:
			// check if script id exists
			chatScript, err := repository.ChatScriptRepo.GetById(ctx, dbCon, action.ChatScriptId)
			if err != nil {
				log.Error(err)
				return chatAutoScript.Id, err
			}
			if chatScript == nil {
				err = errors.New("not found chat script id")
				log.Error(err)
				return chatAutoScript.Id, err
			}

			chatAutoScript.ActionScript.Actions[i].ChatScriptId = action.ChatScriptId
		case model.SendMessage:
			chatAutoScript.ActionScript.Actions[i].Content = action.Content
			//TODO: handle label case
		default:
			err = errors.New("invalid action type: " + action.Type)
			log.Error(err)
			return chatAutoScript.Id, err
		}
	}

	if chatAutoScriptRequest.Status == "true" {
		chatAutoScript.Status = true
	}

	chatAutoScript.TriggerEvent = chatAutoScriptRequest.TriggerEvent
	chatAutoScript.ScriptName = chatAutoScriptRequest.ScriptName
	chatAutoScript.CreatedBy = authUser.UserId
	chatAutoScript.Channel = chatAutoScriptRequest.Channel
	chatAutoScript.ConnectionId = chatAutoScriptRequest.ConnectionId
	chatAutoScript.CreatedAt = time.Now()

	err = repository.ChatAutoScriptRepo.Insert(ctx, dbCon, chatAutoScript)
	if err != nil {
		log.Error(err)
		return chatAutoScript.Id, err
	}

	return chatAutoScript.Id, nil
}

func (s *ChatAutoScript) UpdateChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string, chatAutoScriptRequest model.ChatAutoScriptRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatAutoScript, err := repository.ChatAutoScriptRepo.GetById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatAutoScript == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	// handle actions' content
	for i, action := range chatAutoScriptRequest.ActionScript.Actions {
		// TODO: is action type changeable?
		if chatAutoScript.ActionScript.Actions[i].Type != action.Type {
			err = errors.New("cannot change action's type")
			log.Error(err)
			return err
		}
		switch model.ScriptActionType(action.Type) {
		case model.MoveToExistedScript:
			// check if script id exists
			chatScript, err := repository.ChatScriptRepo.GetById(ctx, dbCon, action.ChatScriptId)
			if err != nil {
				log.Error(err)
				return err
			}
			if chatScript == nil {
				err = errors.New("not found chat script id")
				log.Error(err)
				return err
			}

			chatAutoScript.ActionScript.Actions[i].ChatScriptId = action.ChatScriptId
		case model.SendMessage:
			chatAutoScript.ActionScript.Actions[i].Content = action.Content
			//TODO: handle label case
		default:
			err = errors.New("invalid action type: " + action.Type)
			log.Error(err)
			return err
		}
	}

	chatAutoScript.TriggerEvent = chatAutoScriptRequest.TriggerEvent
	chatAutoScript.ScriptName = chatAutoScriptRequest.ScriptName
	chatAutoScript.UpdatedBy = authUser.UserId
	chatAutoScript.UpdatedAt = time.Now()
	err = repository.ChatAutoScriptRepo.Update(ctx, dbCon, *chatAutoScript)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatAutoScript) UpdateChatAutoScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, oldStatus string) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatScript, err := repository.ChatAutoScriptRepo.GetById(ctx, dbCon, id)
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
	err = repository.ChatAutoScriptRepo.Update(ctx, dbCon, *chatScript)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *ChatAutoScript) DeleteChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	chatScript, err := repository.ChatAutoScriptRepo.GetById(ctx, dbCon, id)
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

	err = repository.ChatAutoScriptRepo.Delete(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}
