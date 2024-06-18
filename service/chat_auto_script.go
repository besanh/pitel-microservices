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

func (s *ChatAutoScript) GetChatAutoScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatAutoScriptFilter, limit int, offset int) (total int, chatAutoScripts *[]model.ChatAutoScriptView, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	filter.TenantId = authUser.TenantId
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
		TenantId:           authUser.TenantId,
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
	labels := make([]model.ChatAutoScriptToChatLabel, 0)
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
		case model.AddLabels:
			for _, addingLabelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, dbCon, addingLabelId)
				if err != nil {
					log.Error(err)
					return chatAutoScript.Id, err
				}
				if label == nil {
					err = errors.New("not found label id")
					log.Error(err)
					return chatAutoScript.Id, err
				}

				currentTime := time.Now()
				labels = append(labels, model.ChatAutoScriptToChatLabel{
					ChatAutoScriptId: chatAutoScript.Id,
					ChatLabelId:      addingLabelId,
					ActionType:       string(model.AddLabels),
					Order:            i,
					CreatedAt:        currentTime,
					UpdatedAt:        currentTime,
				})
			}
		case model.RemoveLabels:
			for _, removingLabelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, dbCon, removingLabelId)
				if err != nil {
					log.Error(err)
					return chatAutoScript.Id, err
				}
				if label == nil {
					err = errors.New("not found label id")
					log.Error(err)
					return chatAutoScript.Id, err
				}

				currentTime := time.Now()
				labels = append(labels, model.ChatAutoScriptToChatLabel{
					ChatAutoScriptId: chatAutoScript.Id,
					ChatLabelId:      removingLabelId,
					ActionType:       string(model.RemoveLabels),
					Order:            i,
					CreatedAt:        currentTime,
					UpdatedAt:        currentTime,
				})
			}
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

	err = repository.ChatAutoScriptRepo.InsertChatAutoScript(ctx, dbCon, chatAutoScript, scripts, labels)
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

	actionTypes := make(map[model.ScriptActionType]bool)
	newScripts := make([]model.ChatAutoScriptToChatScript, 0)
	newLabels := make([]model.ChatAutoScriptToChatLabel, 0)
	// handle actions' content
	for i, action := range chatAutoScriptRequest.ActionScript.Actions {
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

			currentTime := time.Now()
			newScripts = append(newScripts, model.ChatAutoScriptToChatScript{
				ChatAutoScriptId: chatAutoScript.Id,
				ChatScriptId:     action.ChatScriptId,
				Order:            i,
				CreatedAt:        currentTime,
				UpdatedAt:        currentTime,
			})
		case model.SendMessage:
			if _, ok := actionTypes[model.SendMessage]; !ok {
				actionTypes[model.SendMessage] = true
				//create new send message script
				chatAutoScript.SendMessageActions = model.AutoScriptSendMessage{Actions: make([]model.AutoScriptSendMessageType, 0)}
			}
			chatAutoScript.SendMessageActions.Actions = append(chatAutoScript.SendMessageActions.Actions,
				model.AutoScriptSendMessageType{
					Content: action.Content,
					Order:   i,
				})
		case model.AddLabels:
			for _, addingLabelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, dbCon, addingLabelId)
				if err != nil {
					log.Error(err)
					return err
				}
				if label == nil {
					err = errors.New("not found label id")
					log.Error(err)
					return err
				}

				currentTime := time.Now()
				newLabels = append(newLabels, model.ChatAutoScriptToChatLabel{
					ChatAutoScriptId: chatAutoScript.Id,
					ChatLabelId:      addingLabelId,
					ActionType:       string(model.AddLabels),
					Order:            i,
					CreatedAt:        currentTime,
					UpdatedAt:        currentTime,
				})
			}
		case model.RemoveLabels:
			for _, removingLabelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, dbCon, removingLabelId)
				if err != nil {
					log.Error(err)
					return err
				}
				if label == nil {
					err = errors.New("not found label id")
					log.Error(err)
					return err
				}

				currentTime := time.Now()
				newLabels = append(newLabels, model.ChatAutoScriptToChatLabel{
					ChatAutoScriptId: chatAutoScript.Id,
					ChatLabelId:      removingLabelId,
					ActionType:       string(model.RemoveLabels),
					Order:            i,
					CreatedAt:        currentTime,
					UpdatedAt:        currentTime,
				})
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
	err = repository.ChatAutoScriptRepo.UpdateChatAutoScriptById(ctx, dbCon, *chatAutoScript, newScripts, newLabels)
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
