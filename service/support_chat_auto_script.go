package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func GenerateChatAutoScriptId(tenantId, channel, appId, oaId, triggerEvent string) string {
	return CHAT_AUTO_SCRIPT + "_" + tenantId + "_" + channel + "_" + appId + "_" + oaId + "_" + triggerEvent
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

/*
 * chat auto scripts after being fetched from db need to be aggregated from 3 fields (send_message, script_link, label_link)
 * into one single field for easier managing
 */
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
		if action.ChatScript == nil {
			log.Error("not found this chat script's info, id: ", action.ChatScriptId)
			continue
		}
		chatAutoScript.ActionScript.Actions = append(chatAutoScript.ActionScript.Actions, model.ActionScriptActionType{
			Type:         string(model.MoveToExistedScript),
			ChatScriptId: action.ChatScriptId,
			Order:        action.Order,
		})
	}

	addLabels := make(map[int][]string)
	removeLabels := make(map[int][]string)
	for _, action := range chatAutoScript.ChatLabelLink {
		if action.ChatLabel == nil {
			log.Error("not found this chat label's info, id: ", action.ChatLabelId)
			continue
		}
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

func processLabels(ctx context.Context, dbCon sqlclient.ISqlClientConn, labels []model.ChatAutoScriptToChatLabel, labelIds []string,
	chatLabelAction model.ChatLabelAction) ([]model.ChatAutoScriptToChatLabel, error) {
	for _, labelId := range labelIds {
		label, err := repository.ChatLabelRepo.GetById(ctx, dbCon, labelId)
		if err != nil || label == nil {
			return labels, fmt.Errorf("not found label id: %v", err)
		}

		labels = append(labels, model.ChatAutoScriptToChatLabel{
			ChatLabelId:      labelId,
			ChatAutoScriptId: chatLabelAction.ChatAutoScriptId,
			ActionType:       chatLabelAction.ActionType,
			Order:            chatLabelAction.Order,
			CreatedAt:        chatLabelAction.CreatedAt,
		})
	}
	return labels, nil
}

/*
 * Handle execute main chat auto script's logics (detecting keywords, offline agents)
 */
func ExecutePlannedAutoScript(ctx context.Context, user model.User, message model.Message, conversation *model.ConversationView) error {
	if err := DetectKeywordsAndExecutePlannedAutoScript(ctx, user, message, conversation); err != nil {
		return err
	}
	if err := ExecutePlannedAutoScriptWhenAgentsOffline(ctx, user, message, conversation); err != nil {
		return err
	}
	return nil
}

/*
 * Handle detect keywords in message's content then executing the first matching script
 */
func DetectKeywordsAndExecutePlannedAutoScript(ctx context.Context, user model.User, message model.Message, conversation *model.ConversationView) error {
	if conversation == nil {
		return errors.New("not found conversation")
	}
	if user.AuthUser == nil {
		log.Error("not found auth user info")
		return nil
	}
	filter := model.ChatAutoScriptFilter{
		TenantId:     message.TenantId,
		Channel:      message.MessageType,
		OaId:         message.OaId,
		TriggerEvent: "keyword",
		Status:       sql.NullBool{Valid: true, Bool: true},
	}

	var chatAutoScripts *[]model.ChatAutoScriptView
	key := GenerateChatAutoScriptId(filter.TenantId, filter.Channel, conversation.AppId, filter.OaId, filter.TriggerEvent)
	chatAutoScriptsCache := cache.RCache.Get(key)
	if chatAutoScriptsCache != nil {
		if err := json.Unmarshal([]byte(chatAutoScriptsCache.(string)), &chatAutoScripts); err != nil {
			log.Error(err)
			return err
		}
	} else {
		total, scripts, err := repository.ChatAutoScriptRepo.GetChatAutoScripts(ctx, repository.DBConn, filter, 0, 0)
		if err != nil {
			log.Error(err)
			return err
		}
		if total == 0 {
			log.Info("not found any auto scripts")
			return nil
		}

		if err = cache.RCache.Set(key, scripts, CHAT_AUTO_SCRIPT_EXPIRE); err != nil {
			log.Error(err)
			return err
		}
		chatAutoScripts = scripts
	}

	chatAutoScripts = mergeActionScripts(chatAutoScripts)
	// try to execute the first script
	var script *model.ChatAutoScriptView
	for _, scriptView := range *chatAutoScripts {
		if util.ContainKeywords(message.Content, scriptView.TriggerKeywords.Keywords) {
			script = &scriptView
			break
		}
	}
	if script == nil {
		// not matching any keywords
		log.Info("not found matching keyword")
		return nil
	}

	if err := executeScriptActions(ctx, user, message, conversation, *script); err != nil {
		return err
	}
	return nil
}

/*
 * Handle detect agents online status then executing the first matching script
 */
func ExecutePlannedAutoScriptWhenAgentsOffline(ctx context.Context, user model.User, message model.Message, conversation *model.ConversationView) error {
	if conversation == nil {
		return errors.New("not found conversation")
	}
	if user.AuthUser == nil {
		log.Error("not found auth user info")
		return nil
	}

	// check if subscribers' level is agent/user
	userLive := false
	for s := range WsSubscribers.Subscribers {
		if s.Level == "user" || s.Level == "agent" {
			userLive = true
			break
		}
	}
	if userLive {
		// has online agents -> do nothing
		log.Info("not executed offline auto script because agents are online")
		return nil
	}

	filter := model.ChatAutoScriptFilter{
		TenantId:     message.TenantId,
		Channel:      message.MessageType,
		OaId:         message.OaId,
		TriggerEvent: "offline",
		Status:       sql.NullBool{Valid: true, Bool: true},
	}

	var chatAutoScripts *[]model.ChatAutoScriptView
	key := GenerateChatAutoScriptId(filter.TenantId, filter.Channel, conversation.AppId, filter.OaId, filter.TriggerEvent)
	chatAutoScriptsCache := cache.RCache.Get(key)
	if chatAutoScriptsCache != nil {
		if err := json.Unmarshal([]byte(chatAutoScriptsCache.(string)), &chatAutoScripts); err != nil {
			log.Error(err)
			return err
		}
	} else {
		total, scripts, err := repository.ChatAutoScriptRepo.GetChatAutoScripts(ctx, repository.DBConn, filter, 0, 0)
		if err != nil {
			log.Error(err)
			return err
		}
		if total == 0 {
			log.Info("not found any auto scripts")
			return nil
		}

		if err = cache.RCache.Set(key, scripts, CHAT_AUTO_SCRIPT_EXPIRE); err != nil {
			log.Error(err)
			return err
		}
		chatAutoScripts = scripts
	}

	chatAutoScripts = mergeActionScripts(chatAutoScripts)
	if len(*chatAutoScripts) == 0 {
		return nil
	}
	// try to execute the first script
	script := (*chatAutoScripts)[0]

	if err := executeScriptActions(ctx, user, message, conversation, script); err != nil {
		return err
	}
	return nil
}

/*
 * Handle chat auto script's logics
 */
func executeScriptActions(ctx context.Context, user model.User, message model.Message, conversation *model.ConversationView, script model.ChatAutoScriptView) error {
	if conversation == nil {
		return errors.New("not found conversation")
	}

	subscribers := make([]string, 0)
	for s := range WsSubscribers.Subscribers {
		if (user.AuthUser != nil && s.TenantId == user.AuthUser.TenantId) || (conversation.TenantId == s.TenantId) {
			subscribers = append(subscribers, s.Id)
		}
	}

	timestamp := time.Now().UnixMilli()
	for _, action := range script.ActionScript.Actions {
		switch action.Type {
		case string(model.MoveToExistedScript):
			if err := executeScript(ctx, user, message, conversation, action.ChatScriptId, 3); err != nil {
				return err
			}
		case string(model.SendMessage):
			if err := executeSendScriptedMessage(ctx, user, conversation, timestamp, "text", action.Content, nil); err != nil {
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
					AppId:           conversation.AppId,
					OaId:            conversation.OaId,
					LabelName:       label.LabelName,
					LabelId:         labelId,
					ExternalLabelId: label.ExternalLabelId,
					ExternalUserId:  conversation.ExternalUserId,
					ConversationId:  conversation.ConversationId,
					Action:          "",
				}

				if conversation.ConversationType == "zalo" {
					request.Action = "create"
					if err = handleLabelZalo(ctx, conversation.ConversationType, request); err != nil {
						return err
					}
					// switch to update because this label already exists
					request.Action = "update"
					if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
						return err
					}

					newLabels, err := UpdateConversationLabelList(conversation.Label, conversation.ConversationType, request.Action, labelId)
					if err != nil {
						return err
					}
					conversation.Label = newLabels
				} else if conversation.ConversationType == "facebook" {
					if len(label.ExternalLabelId) < 1 {
						// request fb to create new external label id
						request.Action = "create"
						externalLabelId, err := handleLabelFacebook(ctx, repository.DBConn, conversation.ConversationType, *label, request)
						if err != nil {
							return err
						}
						// update external label id
						label.ExternalLabelId = externalLabelId
						label.UpdatedBy = user.AuthUser.UserId
						if err := repository.ChatLabelRepo.Update(ctx, repository.DBConn, *label); err != nil {
							return err
						}
						// update external label id to request
						request.ExternalLabelId = externalLabelId
					}
					//switch to update
					request.Action = "update"
					if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
						return err
					}

					newLabels, err := UpdateConversationLabelList(conversation.Label, conversation.ConversationType, request.Action, request.ExternalLabelId)
					if err != nil {
						return err
					}
					conversation.Label = newLabels
				}
			}

			tmp := conversation.Label
			if err := GetLabelsInfo(ctx, conversation); err != nil {
				return err
			}
			if len(action.AddLabels) > 0 {
				PublishConversationToManyUser(variables.EVENT_CHAT["conversation_add_labels"], subscribers, true, conversation)
			}
			// wipe out labels detailed info and set it to label ids list
			conversation.Label = tmp
		case string(model.RemoveLabels):
			for _, labelId := range action.RemoveLabels {
				label, err := repository.ChatLabelRepo.GetById(ctx, repository.DBConn, labelId)
				if err != nil {
					return err
				}
				if label == nil {
					return errors.New("not found label")
				}

				request := model.ConversationLabelRequest{
					AppId:           conversation.AppId,
					OaId:            conversation.OaId,
					LabelName:       label.LabelName,
					LabelId:         labelId,
					ExternalLabelId: label.ExternalLabelId,
					ExternalUserId:  conversation.ExternalUserId,
					ConversationId:  conversation.ConversationId,
					Action:          "",
				}

				if conversation.ConversationType == "zalo" {
					request.Action = "delete"
					if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
						return err
					}

					newLabels, err := UpdateConversationLabelList(conversation.Label, conversation.ConversationType, request.Action, labelId)
					if err != nil {
						return err
					}
					conversation.Label = newLabels
				} else if conversation.ConversationType == "facebook" {
					if len(label.ExternalLabelId) > 0 {
						request.Action = "delete"
						if _, err := PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
							return err
						}

						newLabels, err := UpdateConversationLabelList(conversation.Label, conversation.ConversationType, request.Action, label.ExternalLabelId)
						if err != nil {
							return err
						}
						conversation.Label = newLabels
					} else {
						// do nothing
					}
				}
			}

			tmp := conversation.Label
			if err := GetLabelsInfo(ctx, conversation); err != nil {
				return err
			}
			if len(action.RemoveLabels) > 0 {
				PublishConversationToManyUser(variables.EVENT_CHAT["conversation_remove_labels"], subscribers, true, conversation)
			}
			// wipe out labels detailed info and set it to label ids list
			conversation.Label = tmp
		default:
			return errors.New("invalid action type")
		}
	}
	return nil
}

/*
 * send scripted message pre-defined from chat auto script or chat script to ott & es
 */
func executeSendScriptedMessage(ctx context.Context, user model.User, conversation *model.ConversationView,
	timestamp int64, eventName, content string, attachments []*model.OttAttachments) error {
	if conversation == nil {
		return errors.New("not found conversation")
	}

	if util.ContainKeywords(content, variables.PERSONALIZATION_KEYWORDS) {
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
	if err := InsertES(ctx, conversation.TenantId, ES_INDEX, scriptedMessage.AppId, docId, scriptedMessage); err != nil {
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
	if err = json.Unmarshal(tmpBytes, &esDoc); err != nil {
		return err
	}
	if err = repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversation.AppId, conversation.ConversationId, esDoc); err != nil {
		log.Error(err)
		return err
	}

	// >send message to manager/admin
	if user.AuthUser.Level != "manager" {
		if len(user.QueueId) > 0 {
			if err := SendEventToManage(ctx, user.AuthUser, scriptedMessage, user.QueueId); err != nil {
				log.Error(err)
				return err
			}
		} else {
			err = errors.New(fmt.Sprintf("queue %s not found in send event to manage", user.QueueId))
			log.Error(err)
			return err
		}
	}
	return nil
}

/*
 * execute chat script based on its script type accordingly
 */
func executeScript(ctx context.Context, user model.User, message model.Message, conversation *model.ConversationView,
	id string, limit int) error {
	if conversation == nil {
		return errors.New("not found conversation")
	}
	if limit < 1 {
		return errors.New("out of limit in executing chat script")
	}
	chatScript, err := repository.ChatScriptRepo.GetById(ctx, repository.DBConn, id)
	if err != nil {
		return errors.New("not found chat script")
	}
	if !chatScript.Status {
		// this script is inactive so just skip it
		return nil
	}

	timestamp := time.Now().UnixMilli()
	switch chatScript.ScriptType {
	case "text":
		content := chatScript.Content
		if err = executeSendScriptedMessage(ctx, user, conversation, timestamp, "text", content, nil); err != nil {
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
		if err = executeSendScriptedMessage(ctx, user, conversation, timestamp, "file", "", attachments); err != nil {
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
		if err = executeSendScriptedMessage(ctx, user, conversation, timestamp, "image", "", attachments); err != nil {
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

func GetLabelsInfo(ctx context.Context, conversation *model.ConversationView) error {
	if !reflect.DeepEqual(conversation.Label, "") {
		var labels []map[string]string
		if err := json.Unmarshal([]byte(conversation.Label), &labels); err != nil {
			log.Error(err)
			return err
		}
		chatLabelIds := []string{}
		if len(labels) > 0 {
			for _, item := range labels {
				chatLabelIds = append(chatLabelIds, item["label_id"])
			}
			if len(chatLabelIds) > 0 {
				_, chatLabelExist, err := repository.ChatLabelRepo.GetChatLabels(ctx, repository.DBConn, model.ChatLabelFilter{
					LabelIds: chatLabelIds,
				}, -1, 0)
				if err != nil {
					log.Error(err)
					return err
				}
				if len(*chatLabelExist) > 0 {
					tmp, err := json.Marshal(*chatLabelExist)
					if err != nil {
						log.Error(err)
						return err
					}
					conversation.Label = tmp
				} else {
					conversation.Label = []byte("[]")
				}
			} else {
				conversation.Label = []byte("[]")
			}
		} else {
			conversation.Label = []byte("[]")
		}
	} else {
		conversation.Label = []byte("[]")
	}
	return nil
}
