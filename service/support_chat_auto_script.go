package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"sort"
	"strings"
	"time"
)

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

	addLabels := make(map[int][]string)
	removeLabels := make(map[int][]string)
	for _, action := range chatAutoScript.ChatLabelLink {
		if action.ActionType == string(model.AddLabels) {
			if addLabels[action.Order] == nil || len(addLabels[action.Order]) == 0 {
				addLabels[action.Order] = make([]string, 0)
			}
			addLabels[action.Order] = append(addLabels[action.Order], action.ChatLabelId)
		} else if action.ActionType == string(model.RemoveLabels) {
			if removeLabels[action.Order] == nil || len(removeLabels[action.Order]) == 0 {
				removeLabels[action.Order] = make([]string, 0)
			}
			removeLabels[action.Order] = append(removeLabels[action.Order], action.ChatLabelId)
		}
	}

	for order, labels := range addLabels {
		chatAutoScript.ActionScript.Actions = append(chatAutoScript.ActionScript.Actions, model.ActionScriptActionType{
			Type:      string(model.AddLabels),
			AddLabels: labels,
			Order:     order,
		})
	}

	for order, labels := range removeLabels {
		chatAutoScript.ActionScript.Actions = append(chatAutoScript.ActionScript.Actions, model.ActionScriptActionType{
			Type:         string(model.RemoveLabels),
			RemoveLabels: labels,
			Order:        order,
		})
	}

	sort.Slice(chatAutoScript.ActionScript.Actions, func(i, j int) bool {
		return chatAutoScript.ActionScript.Actions[i].Order < chatAutoScript.ActionScript.Actions[j].Order
	})
	return chatAutoScript
}

func containsKeyword(message string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(message, keyword) {
			return true
		}
	}
	return false
}

func DetectAndExecutePlannedAutoScript(ctx context.Context, user model.User, message model.Message, conversation model.Conversation) error {
	if user.AuthUser == nil {
		return nil
	}
	filter := model.ChatAutoScriptFilter{
		TenantId:     message.TenantId,
		Channel:      message.MessageType,
		OaId:         message.OaId,
		TriggerEvent: "keyword",
		Status:       sql.NullBool{Valid: true, Bool: true},
	}
	total, chatAutoScripts, err := repository.ChatAutoScriptRepo.GetChatAutoScripts(ctx, repository.DBConn, filter, 0, 0)
	if err != nil {
		return err
	}
	if total == 0 || chatAutoScripts == nil {
		return nil
	}
	if len(*chatAutoScripts) == 0 {
		return nil
	}

	chatAutoScripts = mergeActionScripts(chatAutoScripts)
	// try to execute the first script
	var script *model.ChatAutoScriptView
	for _, scriptView := range *chatAutoScripts {
		if containsKeyword(message.Content, scriptView.TriggerKeywords.Keywords) {
			script = &scriptView
			break
		}
	}
	if script == nil {
		// not matching any keywords
		return nil
	}

	timestamp := time.Now().UnixMilli()
	for _, action := range script.ActionScript.Actions {
		switch action.Type {
		case string(model.MoveToExistedScript):
			if err = executeScript(ctx, user, message, conversation, action.ChatScriptId, 3); err != nil {
				return err
			}
		case string(model.SendMessage):
			if err = executeSendScriptedMessage(ctx, user, message, conversation, timestamp, "text", action.Content, nil); err != nil {
				return err
			}
		case string(model.AddLabels):
			for _, labelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, repository.DBConn, labelId)
				if err != nil {
					return err
				}
				if label == nil {
					return errors.New("not found label")
				}

				request := model.ConversationLabelRequest{
					AppId:          conversation.AppId,
					OaId:           conversation.OaId,
					LabelName:      label.LabelName,
					LabelId:        labelId,
					ExternalUserId: conversation.ExternalUserId,
					ConversationId: conversation.ConversationId,
					Action:         "create",
				}
				if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
					return err
				}
			}
		case string(model.RemoveLabels):
			for _, labelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, repository.DBConn, labelId)
				if err != nil {
					return err
				}
				if label == nil {
					return errors.New("not found label")
				}
				request := model.ConversationLabelRequest{
					AppId:          conversation.AppId,
					OaId:           conversation.OaId,
					LabelId:        labelId,
					ExternalUserId: conversation.ExternalUserId,
					LabelName:      label.LabelName,
					ConversationId: conversation.ConversationId,
					Action:         "delete",
				}
				if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
					return err
				}
			}
		default:
			return errors.New("invalid action type")
		}
	}
	return nil
}

func ExecutePlannedAutoScriptWhenAgentsOffline(ctx context.Context, user model.User, message model.Message, conversation model.Conversation) error {
	if user.AuthUser == nil {
		return nil
	}

	filterQueueUser := model.ChatQueueUserFilter{
		QueueId: []string{user.QueueId},
	}
	total, chatQueueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, filterQueueUser, 1, 0)
	if err != nil {
		return err
	}
	if total == 0 {
		return nil
	}
	userLives, err := anyUsersOnline(chatQueueUsers)
	if err != nil {
		return err
	}
	if userLives {
		// has online agents -> do nothing
		return nil
	}

	filter := model.ChatAutoScriptFilter{
		TenantId:     message.TenantId,
		Channel:      message.MessageType,
		OaId:         message.OaId,
		TriggerEvent: "offline",
		Status:       sql.NullBool{Valid: true, Bool: true},
	}
	total, chatAutoScripts, err := repository.ChatAutoScriptRepo.GetChatAutoScripts(ctx, repository.DBConn, filter, 0, 0)
	if err != nil {
		return err
	}
	if total == 0 || chatAutoScripts == nil {
		return nil
	}
	if len(*chatAutoScripts) == 0 {
		return nil
	}

	chatAutoScripts = mergeActionScripts(chatAutoScripts)
	// try to execute the first script
	script := (*chatAutoScripts)[0]

	timestamp := time.Now().UnixMilli()
	for _, action := range script.ActionScript.Actions {
		switch action.Type {
		case string(model.MoveToExistedScript):
			if err = executeScript(ctx, user, message, conversation, action.ChatScriptId, 3); err != nil {
				return err
			}
		case string(model.SendMessage):
			if err = executeSendScriptedMessage(ctx, user, message, conversation, timestamp, "text", action.Content, nil); err != nil {
				return err
			}
		case string(model.AddLabels):
			for _, labelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, repository.DBConn, labelId)
				if err != nil {
					return err
				}
				if label == nil {
					return errors.New("not found label")
				}

				request := model.ConversationLabelRequest{
					AppId:          conversation.AppId,
					OaId:           conversation.OaId,
					LabelName:      label.LabelName,
					LabelId:        labelId,
					ExternalUserId: conversation.ExternalUserId,
					ConversationId: conversation.ConversationId,
					Action:         "create",
				}
				if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
					return err
				}
			}
		case string(model.RemoveLabels):
			for _, labelId := range action.AddLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, repository.DBConn, labelId)
				if err != nil {
					return err
				}
				if label == nil {
					return errors.New("not found label")
				}
				request := model.ConversationLabelRequest{
					AppId:          conversation.AppId,
					OaId:           conversation.OaId,
					LabelId:        labelId,
					ExternalUserId: conversation.ExternalUserId,
					LabelName:      label.LabelName,
					ConversationId: conversation.ConversationId,
					Action:         "delete",
				}
				if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
					return err
				}
			}
		default:
			return errors.New("invalid action type")
		}
	}
	return nil
}

func anyUsersOnline(queueUsers *[]model.ChatQueueUser) (bool, error) {
	userLives := make([]Subscriber, 0)
	subscribers, err := cache.RCache.HGetAll(BSS_SUBSCRIBERS)
	if err != nil {
		log.Error(err)
		return false, err
	}
	for _, item := range subscribers {
		s := Subscriber{}
		if err := json.Unmarshal([]byte(item), &s); err != nil {
			log.Error(err)
			return false, err
		}
		if (s.Level == "user" || s.Level == "agent") && CheckInLive(*queueUsers, s.Id) {
			userLives = append(userLives, s)
		}
	}
	if len(userLives) > 0 {
		return true, nil
	}
	return false, nil
}

func executeSendScriptedMessage(ctx context.Context, user model.User, message model.Message, conversation model.Conversation,
	timestamp int64, eventName, content string, attachments []*model.OttAttachments) error {
	if containsKeyword(content, variables.PERSONALIZATION_KEYWORD) {
		pageName := conversation.OaName
		customerName := conversation.Username
		content = strings.ReplaceAll(content, "{{page_name}}", pageName)
		content = strings.ReplaceAll(content, "{{customer_name}}", customerName)
	}

	// >send message to ott
	ottMessage := model.SendMessageToOtt{
		Type:          conversation.ConversationType,
		EventName:     eventName,
		AppId:         conversation.AppId,
		OaId:          conversation.OaId,
		Uid:           conversation.ExternalUserId,
		SupporterId:   user.AuthUser.UserId,
		SupporterName: user.AuthUser.Username,
		Timestamp:     fmt.Sprintf("%d", timestamp),
		Text:          content,
	}

	log.Info("message to ott: ", ottMessage)
	resOtt, err := sendMessageToOTT(ottMessage, attachments)
	if err != nil {
		log.Error(err)
		return err
	}
	docId := uuid.NewString()
	// Store ES
	scriptedMessage := model.Message{
		TenantId:            conversation.TenantId,
		ParentExternalMsgId: "",
		Id:                  docId,
		MessageType:         conversation.ConversationType,
		ConversationId:      conversation.ConversationId,
		ExternalMsgId:       resOtt.Data.MsgId,
		EventName:           eventName,
		Direction:           variables.DIRECTION["send"],
		AppId:               conversation.AppId,
		OaId:                conversation.OaId,
		Avatar:              conversation.Avatar,
		SupporterId:         user.AuthUser.UserId,
		SupporterName:       user.AuthUser.Fullname,
		SendTime:            time.Now(),
		SendTimestamp:       timestamp,
		Content:             content,
		Attachments:         attachments,
	}
	log.Info("message to es: ", scriptedMessage)

	// Should to queue
	if err := InsertES(ctx, conversation.TenantId, ES_INDEX, message.AppId, docId, scriptedMessage); err != nil {
		log.Error(err)
		return err
	}

	// >update conversation doc on ES
	conversation.UpdatedAt = time.Now().Format(time.RFC3339)
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		return err
	}
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		return err
	}
	if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversation.AppId, conversation.ConversationId, esDoc); err != nil {
		log.Error(err)
		return err
	}

	// >send message to manager/admin
	if user.AuthUser.Level != "manager" {
		if len(user.QueueId) > 0 {
			if err := SendEventToManage(ctx, user.AuthUser, message, user.QueueId); err != nil {
				log.Error(err)
				return err
			}
		} else {
			log.Errorf("queue %s not found in send event to manage", user.QueueId)
		}
	}
	return nil
}

func executeScript(ctx context.Context, user model.User, message model.Message, conversation model.Conversation,
	id string, limit int) error {
	if limit < 1 {
		return errors.New("out of limit in executing chat script")
	}
	chatScript, err := repository.ChatScriptRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		return errors.New("not found chat script")
	}
	if !chatScript.Status {
		return errors.New("can not execute inactive chat script")
	}

	timestamp := time.Now().UnixMilli()
	switch chatScript.ScriptType {
	case "text":
		content := chatScript.Content
		if err = executeSendScriptedMessage(ctx, user, message, conversation, timestamp, "text", content, nil); err != nil {
			return err
		}
	case "file":
		attachments := make([]*model.OttAttachments, 0)
		attachments = append(attachments, &model.OttAttachments{
			Payload: &model.OttPayloadMedia{
				Url: chatScript.FileUrl,
			},
			AttType: "file",
		})
		if err = executeSendScriptedMessage(ctx, user, message, conversation, timestamp, "file", "", attachments); err != nil {
			return err
		}
	case "image":
		attachments := make([]*model.OttAttachments, 0)
		attachments = append(attachments, &model.OttAttachments{
			Payload: &model.OttPayloadMedia{
				Url: chatScript.FileUrl,
			},
			AttType: "image",
		})
		if err = executeSendScriptedMessage(ctx, user, message, conversation, timestamp, "image", "", attachments); err != nil {
			return err
		}
	case "other":
		if err = executeScript(ctx, user, message, conversation, chatScript.OtherScriptId, limit-1); err != nil {
			return err
		}
	default:
		if err != nil {
			return errors.New("invalid script type")
		}
	}
	return nil
}
