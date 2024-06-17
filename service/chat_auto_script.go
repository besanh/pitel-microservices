package service

import (
	"context"
	"errors"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"sort"
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

func (s *ChatAutoScript) GetChatAutoScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatAutoScriptFilter, limit int, offset int) (total int, chatAutoScripts *[]model.ChatAutoScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	total, chatAutoScripts, err = repository.ChatAutoScriptRepo.GetChatAutoScripts(ctx, dbCon, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}

	chatAutoScripts = mergeActionScripts(chatAutoScripts)

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
	tmp := mergeSingleActionScript(*result)
	result = &tmp

	return
}

func (s *ChatAutoScript) InsertChatAutoScript(ctx context.Context, authUser *model.AuthUser, chatAutoScriptRequest model.ChatAutoScriptRequest) (string, error) {
	chatAutoScript := model.ChatAutoScript{
		Base:               model.InitBase(),
		SendMessageActions: model.AutoScriptSendMessage{Actions: make([]model.AutoScriptSendMessageType, 0)},
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

	scripts := make([]model.ChatAutoScriptToChatScript, 0)

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

			currentTime := time.Now()
			scripts = append(scripts, model.ChatAutoScriptToChatScript{
				ChatAutoScriptId: chatAutoScript.Id,
				ChatScriptId:     action.ChatScriptId,
				Order:            i,
				CreatedAt:        currentTime,
				UpdatedAt:        currentTime,
			})
		case model.SendMessage:
			chatAutoScript.SendMessageActions.Actions = append(chatAutoScript.SendMessageActions.Actions,
				model.AutoScriptSendMessageType{
					Content: action.Content,
					Order:   i,
				})
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

	err = repository.ChatAutoScriptRepo.InsertChatAutoScript(ctx, dbCon, chatAutoScript, scripts)
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
		switch model.ScriptActionType(action.Type) {
		case model.MoveToExistedScript:
			// do nothing
		case model.SendMessage:
			// update content of message
			for j, _ := range chatAutoScript.SendMessageActions.Actions {
				if chatAutoScript.SendMessageActions.Actions[j].Order == i {
					chatAutoScript.SendMessageActions.Actions[j].Content = action.Content
				}
			}
		default:
			err = errors.New("invalid action type: " + action.Type)
			log.Error(err)
			return err
		}
	}

	var status bool
	if chatAutoScriptRequest.Status == "true" {
		status = true
	}
	chatAutoScript.Status = status
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

	var status bool
	if oldStatus == "true" {
		status = true
	}
	chatAutoScript.Status = !status
	chatAutoScript.UpdatedBy = authUser.UserId
	chatAutoScript.UpdatedAt = time.Now()
	err = repository.ChatAutoScriptRepo.Update(ctx, dbCon, *chatAutoScript)
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

	err = repository.ChatAutoScriptRepo.DeleteChatAutoScriptById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func mergeActionScripts(chatAutoScripts *[]model.ChatAutoScriptView) *[]model.ChatAutoScriptView {
	if chatAutoScripts == nil {
		return nil
	}
	for i, cas := range *chatAutoScripts {
		(*chatAutoScripts)[i] = mergeSingleActionScript(cas)
	}
	return chatAutoScripts
}

func mergeSingleActionScript(chatAutoScript model.ChatAutoScriptView) model.ChatAutoScriptView {
	chatAutoScript.ActionScript = new(model.AutoScriptMergedActions)
	chatAutoScript.ActionScript.Actions = make([]model.ActionScriptActionType, 0)

	for _, action := range chatAutoScript.SendMessageActions.Actions {
		chatAutoScript.ActionScript.Actions = append(chatAutoScript.ActionScript.Actions, model.ActionScriptActionType{
			Type:    string(model.SendMessage),
			Content: action.Content,
			Order:   action.Order,
		})
	}

	for _, action := range chatAutoScript.ChatScriptLink {
		chatAutoScript.ActionScript.Actions = append(chatAutoScript.ActionScript.Actions, model.ActionScriptActionType{
			Type:         string(model.MoveToExistedScript),
			ChatScriptId: action.ChatScriptId,
			Order:        action.Order,
		})
	}

	sort.Slice(chatAutoScript.ActionScript.Actions, func(i, j int) bool {
		return chatAutoScript.ActionScript.Actions[i].Order < chatAutoScript.ActionScript.Actions[j].Order
	})
	return chatAutoScript
}
