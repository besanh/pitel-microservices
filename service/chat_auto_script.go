package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IChatAutoScript interface {
		GetChatAutoScripts(ctx context.Context, authUser *model.AuthUser, filter model.ChatAutoScriptFilter, limit int, offset int) (int, *[]model.ChatAutoScriptView, error)
		GetChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string) (*model.ChatAutoScriptView, error)
		InsertChatAutoScript(ctx context.Context, authUser *model.AuthUser, chatAutoScriptRequest model.ChatAutoScriptRequest) (string, error)
		UpdateChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string, chatAutoScriptRequest model.ChatAutoScriptRequest) error
		UpdateChatAutoScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, status sql.NullBool) error
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
		TriggerKeywords:    model.AutoScriptTriggerKeywordsType{Keywords: make([]string, 0)},
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

	currentTime := time.Now()
	scripts := make([]model.ChatAutoScriptToChatScript, 0)
	labels := make([]model.ChatAutoScriptToChatLabel, 0)
	// handle actions' content
	for i, action := range chatAutoScriptRequest.ActionScript.Actions {
		switch model.ScriptActionType(action.Type) {
		case model.MoveToExistedScript:
			// check if script id exists
			chatScript, err := repository.ChatScriptRepo.GetById(ctx, dbCon, action.ChatScriptId)
			if err != nil || chatScript == nil {
				err = fmt.Errorf("not found chat script id, err=%v", err)
				log.Error(err)
				return chatAutoScript.Id, err
			}
			if !chatScript.Status {
				err = errors.New("can not refer to inactive chat script")
				log.Error(err)
				return chatAutoScript.Id, err
			}

			scripts = append(scripts, model.ChatAutoScriptToChatScript{
				ChatAutoScriptId: chatAutoScript.Id,
				ChatScriptId:     action.ChatScriptId,
				Order:            i,
				CreatedAt:        currentTime,
			})
		case model.SendMessage:
			chatAutoScript.SendMessageActions.Actions = append(chatAutoScript.SendMessageActions.Actions,
				model.AutoScriptSendMessageType{
					Content: action.Content,
					Order:   i,
				})
		case model.AddLabels:
			labels, err = processLabels(ctx, dbCon, labels, action.AddLabels, model.ChatLabelAction{
				ChatAutoScriptId: chatAutoScript.Id,
				ActionType:       string(model.AddLabels),
				Order:            i,
				CreatedAt:        currentTime,
			})
			if err != nil {
				log.Error(err)
				return chatAutoScript.Id, err
			}
		case model.RemoveLabels:
			labels, err = processLabels(ctx, dbCon, labels, action.RemoveLabels, model.ChatLabelAction{
				ChatAutoScriptId: chatAutoScript.Id,
				ActionType:       string(model.RemoveLabels),
				Order:            i,
				CreatedAt:        currentTime,
			})
			if err != nil {
				log.Error(err)
				return chatAutoScript.Id, err
			}
		default:
			err = errors.New("invalid action type: " + action.Type)
			log.Error(err)
			return chatAutoScript.Id, err
		}
	}

	statusTmp := chatAutoScriptRequest.Status
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}
	if status.Valid {
		chatAutoScript.Status = status.Bool
	}

	if chatAutoScriptRequest.TriggerEvent == "keyword" {
		chatAutoScript.TriggerKeywords.Keywords = chatAutoScriptRequest.TriggerKeywords.Keywords
	}
	chatAutoScript.TriggerEvent = chatAutoScriptRequest.TriggerEvent
	chatAutoScript.ScriptName = chatAutoScriptRequest.ScriptName
	chatAutoScript.CreatedBy = authUser.UserId
	chatAutoScript.Channel = chatAutoScriptRequest.Channel
	chatAutoScript.ConnectionId = chatAutoScriptRequest.ConnectionId
	chatAutoScript.CreatedAt = currentTime

	err = repository.ChatAutoScriptRepo.InsertChatAutoScript(ctx, dbCon, chatAutoScript, scripts, labels)
	if err != nil {
		log.Error(err)
		return chatAutoScript.Id, err
	}

	// clear cache
	var key string
	if chatAutoScript.Channel == "zalo" && len(connectionApp.OaInfo.Zalo) > 0 {
		key = GenerateChatAutoScriptId(chatAutoScript.TenantId, chatAutoScript.Channel, connectionApp.OaInfo.Zalo[0].AppId,
			connectionApp.OaInfo.Zalo[0].OaId, chatAutoScript.TriggerEvent)
	} else if chatAutoScript.Channel == "facebook" && len(connectionApp.OaInfo.Facebook) > 0 {
		key = GenerateChatAutoScriptId(chatAutoScript.TenantId, chatAutoScript.Channel, connectionApp.OaInfo.Facebook[0].AppId,
			connectionApp.OaInfo.Facebook[0].OaId, chatAutoScript.TriggerEvent)
	}
	if len(key) > 0 {
		if err = cache.RCache.Del([]string{key}); err != nil {
			log.Error(err)
			return chatAutoScript.Id, err
		}
	}

	return chatAutoScript.Id, nil
}

func (s *ChatAutoScript) UpdateChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string, chatAutoScriptRequest model.ChatAutoScriptRequest) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatAutoScript, err := repository.ChatAutoScriptRepo.GetChatAutoScriptById(ctx, dbCon, id)
	if err != nil || chatAutoScript == nil {
		err = fmt.Errorf("not found chat script id, err=%v", err)
		log.Error(err)
		return err
	}

	// clear old messages
	chatAutoScript.SendMessageActions.Actions = make([]model.AutoScriptSendMessageType, 0)

	currentTime := time.Now()
	newScripts := make([]model.ChatAutoScriptToChatScript, 0)
	newLabels := make([]model.ChatAutoScriptToChatLabel, 0)
	// handle actions' content
	for i, action := range chatAutoScriptRequest.ActionScript.Actions {
		switch model.ScriptActionType(action.Type) {
		case model.MoveToExistedScript:
			// check if script id exists
			chatScript, err := repository.ChatScriptRepo.GetById(ctx, dbCon, action.ChatScriptId)
			if err != nil || chatScript == nil {
				err = fmt.Errorf("not found chat script id, err=%v", err)
				log.Error(err)
				return err
			}
			if !chatScript.Status {
				err = errors.New("can not refer to inactive chat script")
				log.Error(err)
				return err
			}

			newScripts = append(newScripts, model.ChatAutoScriptToChatScript{
				ChatAutoScriptId: chatAutoScript.Id,
				ChatScriptId:     action.ChatScriptId,
				Order:            i,
				CreatedAt:        currentTime,
			})
		case model.SendMessage:
			chatAutoScript.SendMessageActions.Actions = append(chatAutoScript.SendMessageActions.Actions,
				model.AutoScriptSendMessageType{
					Content: action.Content,
					Order:   i,
				})
		case model.AddLabels:
			newLabels, err = processLabels(ctx, dbCon, newLabels, action.AddLabels, model.ChatLabelAction{
				ChatAutoScriptId: chatAutoScript.Id,
				ActionType:       string(model.AddLabels),
				Order:            i,
				CreatedAt:        currentTime,
			})
			if err != nil {
				log.Error(err)
				return err
			}
		case model.RemoveLabels:
			newLabels, err = processLabels(ctx, dbCon, newLabels, action.RemoveLabels, model.ChatLabelAction{
				ChatAutoScriptId: chatAutoScript.Id,
				ActionType:       string(model.RemoveLabels),
				Order:            i,
				CreatedAt:        currentTime,
			})
			if err != nil {
				log.Error(err)
				return err
			}
		default:
			err = errors.New("invalid action type: " + action.Type)
			log.Error(err)
			return err
		}
	}

	statusTmp := chatAutoScriptRequest.Status
	var status sql.NullBool
	if len(statusTmp) > 0 {
		statusTmp, _ := strconv.ParseBool(statusTmp)
		status.Valid = true
		status.Bool = statusTmp
	}
	if status.Valid {
		chatAutoScript.Status = status.Bool
	}

	if chatAutoScriptRequest.TriggerEvent == "keyword" {
		chatAutoScript.TriggerKeywords.Keywords = chatAutoScriptRequest.TriggerKeywords.Keywords
	}
	chatAutoScript.ConnectionId = chatAutoScriptRequest.ConnectionId
	chatAutoScript.TriggerEvent = chatAutoScriptRequest.TriggerEvent
	chatAutoScript.ScriptName = chatAutoScriptRequest.ScriptName
	chatAutoScript.UpdatedBy = authUser.UserId
	chatAutoScript.UpdatedAt = currentTime
	err = repository.ChatAutoScriptRepo.UpdateChatAutoScriptById(ctx, dbCon, *chatAutoScript, newScripts, newLabels)
	if err != nil {
		log.Error(err)
		return err
	}

	// clear cache
	if chatAutoScript.ConnectionApp != nil {
		var key string
		if chatAutoScript.Channel == "zalo" && len(chatAutoScript.ConnectionApp.OaInfo.Zalo) > 0 {
			key = GenerateChatAutoScriptId(chatAutoScript.TenantId, chatAutoScript.Channel, chatAutoScript.ConnectionApp.OaInfo.Zalo[0].AppId,
				chatAutoScript.ConnectionApp.OaInfo.Zalo[0].OaId, chatAutoScript.TriggerEvent)
		} else if chatAutoScript.Channel == "facebook" && len(chatAutoScript.ConnectionApp.OaInfo.Facebook) > 0 {
			key = GenerateChatAutoScriptId(chatAutoScript.TenantId, chatAutoScript.Channel, chatAutoScript.ConnectionApp.OaInfo.Facebook[0].AppId,
				chatAutoScript.ConnectionApp.OaInfo.Facebook[0].OaId, chatAutoScript.TriggerEvent)
		}
		if len(key) > 0 {
			if err = cache.RCache.Del([]string{key}); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}

func (s *ChatAutoScript) UpdateChatAutoScriptStatusById(ctx context.Context, authUser *model.AuthUser, id string, status sql.NullBool) error {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return err
	}

	chatAutoScriptView, err := repository.ChatAutoScriptRepo.GetChatAutoScriptById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return err
	}

	// check if exists
	if chatAutoScriptView == nil {
		err = errors.New("not found id")
		log.Error(err)
		return err
	}

	if status.Valid {
		chatAutoScriptView.Status = status.Bool
	}
	chatAutoScriptView.UpdatedBy = authUser.UserId

	chatAutoScript := model.ChatAutoScript{
		Base:               chatAutoScriptView.Base,
		BaseModel:          chatAutoScriptView.BaseModel,
		TenantId:           chatAutoScriptView.TenantId,
		ScriptName:         chatAutoScriptView.ScriptName,
		Channel:            chatAutoScriptView.Channel,
		ConnectionId:       chatAutoScriptView.ConnectionId,
		ConnectionApp:      chatAutoScriptView.ConnectionApp,
		CreatedBy:          chatAutoScriptView.CreatedBy,
		UpdatedBy:          chatAutoScriptView.UpdatedBy,
		Status:             chatAutoScriptView.Status,
		TriggerEvent:       chatAutoScriptView.TriggerEvent,
		TriggerKeywords:    chatAutoScriptView.TriggerKeywords,
		ChatScriptLink:     chatAutoScriptView.ChatScriptLink,
		SendMessageActions: chatAutoScriptView.SendMessageActions,
		ChatLabelLink:      chatAutoScriptView.ChatLabelLink,
	}

	err = repository.ChatAutoScriptRepo.Update(ctx, dbCon, chatAutoScript)
	if err != nil {
		log.Error(err)
		return err
	}

	if chatAutoScriptView.ConnectionApp != nil {
		// clear cache
		var key string
		if chatAutoScriptView.Channel == "zalo" && len(chatAutoScriptView.ConnectionApp.OaInfo.Zalo) > 0 {
			key = GenerateChatAutoScriptId(chatAutoScriptView.TenantId, chatAutoScriptView.Channel, chatAutoScriptView.ConnectionApp.OaInfo.Zalo[0].AppId,
				chatAutoScriptView.ConnectionApp.OaInfo.Zalo[0].OaId, chatAutoScriptView.TriggerEvent)
		} else if chatAutoScriptView.Channel == "facebook" && len(chatAutoScriptView.ConnectionApp.OaInfo.Facebook) > 0 {
			key = GenerateChatAutoScriptId(chatAutoScriptView.TenantId, chatAutoScriptView.Channel, chatAutoScriptView.ConnectionApp.OaInfo.Facebook[0].AppId,
				chatAutoScriptView.ConnectionApp.OaInfo.Facebook[0].OaId, chatAutoScriptView.TriggerEvent)
		}
		if len(key) > 0 {
			if err = cache.RCache.Del([]string{key}); err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}

func (s *ChatAutoScript) DeleteChatAutoScriptById(ctx context.Context, authUser *model.AuthUser, id string) (err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}

	chatAutoScript, err := repository.ChatAutoScriptRepo.GetChatAutoScriptById(ctx, dbCon, id)
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

	err = repository.ChatAutoScriptRepo.DeleteChatAutoScriptById(ctx, dbCon, id)
	if err != nil {
		log.Error(err)
		return
	}

	if chatAutoScript.ConnectionApp != nil {
		// clear cache
		var key string
		if chatAutoScript.Channel == "zalo" && len(chatAutoScript.ConnectionApp.OaInfo.Zalo) > 0 {
			key = GenerateChatAutoScriptId(chatAutoScript.TenantId, chatAutoScript.Channel, chatAutoScript.ConnectionApp.OaInfo.Zalo[0].AppId,
				chatAutoScript.ConnectionApp.OaInfo.Zalo[0].OaId, chatAutoScript.TriggerEvent)
		} else if chatAutoScript.Channel == "facebook" && len(chatAutoScript.ConnectionApp.OaInfo.Facebook) > 0 {
			key = GenerateChatAutoScriptId(chatAutoScript.TenantId, chatAutoScript.Channel, chatAutoScript.ConnectionApp.OaInfo.Facebook[0].AppId,
				chatAutoScript.ConnectionApp.OaInfo.Facebook[0].OaId, chatAutoScript.TriggerEvent)
		}
		if len(key) > 0 {
			if err = cache.RCache.Del([]string{key}); err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return
}
