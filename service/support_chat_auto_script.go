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
	"github.com/redis/go-redis/v9"
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
func ExecutePlannedAutoScript(ctx context.Context, user model.User, message model.Message, conversationView *model.ConversationView) error {
	if user.AuthUser == nil {
		log.Error("not found auth user info")
		return nil
	}
	if message.IsEcho {
		return nil
	}

	if err := DetectKeywordsAndExecutePlannedAutoScript(ctx, user, message, conversationView); err != nil {
		return err
	}
	if err := ExecutePlannedAutoScriptWhenAgentsOffline(ctx, user, message, conversationView); err != nil {
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
func ExecutePlannedAutoScriptWhenAgentsOffline(ctx context.Context, user model.User, message model.Message, conversation *model.ConversationView) (err error) {
	if conversation == nil {
		return errors.New("not found conversation")
	}
	if user.AuthUser == nil {
		log.Error("not found auth user info")
		return nil
	}

	queueUserExist := make(map[string]struct{})
	queueUserExistCache, errTmp := cache.RCache.HGet(CHAT_QUEUE_USER+"_"+user.AuthUser.TenantId, conversation.ExternalConversationId)
	if errTmp != nil && !errors.Is(errTmp, redis.Nil) {
		err = errTmp
		log.Error(err)
		return
	}
	if len(queueUserExistCache) > 0 {
		if err = json.Unmarshal([]byte(queueUserExistCache), &queueUserExist); err != nil {
			log.Error("failed to unmarshal: ", err)
			return
		}
	} else {
		queueUserFilter := model.ChatQueueUserFilter{
			TenantId: user.AuthUser.TenantId,
			QueueId:  []string{user.QueueId},
		}
		_, queueUsers, err := repository.ChatQueueUserRepo.GetChatQueueUsers(ctx, repository.DBConn, queueUserFilter, -1, 0)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(*queueUsers) < 1 {
			err = errors.New("not found any users in chat queue user")
			log.Error(err)
			return err
		}
		_, manageQueueUser, err := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, model.ChatManageQueueUserFilter{
			TenantId: user.AuthUser.UserId,
			QueueId:  user.QueueId,
		}, 1, 0)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(*manageQueueUser) < 1 {
			err = errors.New("not found manager of queue " + user.QueueId)
			log.Error(err)
			return err
		}
		// add manager's id to the map
		queueUserExist[(*manageQueueUser)[0].UserId] = struct{}{}

		for _, item := range *queueUsers {
			queueUserExist[item.UserId] = struct{}{}
		}

		// set to cache
		jsonByte, err := json.Marshal(&queueUserExist)
		if err != nil {
			log.Error(err)
			return err
		}
		if err = cache.RCache.HSetRaw(ctx, CHAT_QUEUE_USER+"_"+user.AuthUser.TenantId, conversation.ExternalConversationId, string(jsonByte)); err != nil {
			log.Error(err)
			return err
		}
	}
	if len(queueUserExist) < 1 {
		err = errors.New("not found any users in chat queue user")
		log.Error(err)
		return
	}

	// check if subscribers' level is agent/user
	userLive := false
	for s := range WsSubscribers.Subscribers {
		if s.TenantId == user.AuthUser.TenantId && (s.Level == "user" || s.Level == "agent") {
			if _, ok := queueUserExist[s.UserId]; !ok {
				continue
			}
			userLive = true
			break
		}
	}
	if userLive {
		// has online agents -> do nothing
		log.Info("not executed offline auto script because agents are online")
		return
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
		if err = json.Unmarshal([]byte(chatAutoScriptsCache.(string)), &chatAutoScripts); err != nil {
			log.Error(err)
			return
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
		return
	}
	// try to execute the first script
	script := (*chatAutoScripts)[0]

	if err = executeScriptActions(ctx, user, message, conversation, script); err != nil {
		return
	}
	return
}

/*
 * Handle chat auto script's logics
 */
func executeScriptActions(ctx context.Context, user model.User, message model.Message, conversation *model.ConversationView, script model.ChatAutoScriptView) (err error) {
	if conversation == nil {
		err = errors.New("conversation not found")
		return
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
			if err = executeScript(ctx, user, message, conversation, action.ChatScriptId, 3); err != nil {
				return
			}
		case string(model.SendMessage):
			if err = executeSendScriptedMessage(ctx, user, conversation, timestamp, "text", action.Content, nil); err != nil {
				return
			}
		case string(model.AddLabels):
			for _, labelId := range action.AddLabels {
				label, errTmp := repository.ChatLabelRepo.GetById(ctx, repository.DBConn, labelId)
				if errTmp != nil {
					err = errTmp
					return
				}
				if label == nil {
					err = errors.New("label " + labelId + " not found")
					return
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

					newLabels, errTmp := UpdateConversationLabelList(conversation.Labels, conversation.ConversationType, request.Action, labelId)
					if errTmp != nil {
						err = errTmp
						return
					}
					conversation.Labels = newLabels
				} else if conversation.ConversationType == "facebook" {
					if len(label.ExternalLabelId) < 1 {
						// request fb to create new external label id
						request.Action = "create"
						externalLabelId, errTmp := handleLabelFacebook(ctx, conversation.ConversationType, *label, request)
						if errTmp != nil {
							err = errTmp
							return
						}
						// update external label id
						label.ExternalLabelId = externalLabelId
						label.UpdatedBy = user.AuthUser.UserId
						if err = repository.ChatLabelRepo.Update(ctx, repository.DBConn, *label); err != nil {
							return
						}
						// update external label id to request
						request.ExternalLabelId = externalLabelId
					}
					//switch to update
					request.Action = "update"
					if _, err = PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
						return
					}

					newLabels, errTmp := UpdateConversationLabelList(conversation.Labels, conversation.ConversationType, request.Action, request.ExternalLabelId)
					if errTmp != nil {
						err = errTmp
						return
					}
					conversation.Labels = newLabels
				}
			}

			tmp := conversation.Labels
			if err = GetLabelsInfo(ctx, conversation); err != nil {
				return
			}
			if len(action.AddLabels) > 0 {
				PublishConversationToManyUser(variables.EVENT_CHAT["conversation_add_labels"], subscribers, true, conversation)
			}
			// wipe out labels detailed info and set it to label ids list
			conversation.Labels = tmp
		case string(model.RemoveLabels):
			for _, labelId := range action.RemoveLabels {
				label, errTmp := repository.ChatLabelRepo.GetById(ctx, repository.DBConn, labelId)
				if errTmp != nil {
					err = errTmp
					return
				}
				if label == nil {
					err = errors.New("label " + labelId + " not found")
					return
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
					if _, err = PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
						return
					}

					newLabels, errTmp := UpdateConversationLabelList(conversation.Labels, conversation.ConversationType, request.Action, labelId)
					if errTmp != nil {
						err = errTmp
						return
					}
					conversation.Labels = newLabels
				} else if conversation.ConversationType == "facebook" {
					if len(label.ExternalLabelId) > 0 {
						request.Action = "delete"
						if _, err = PutLabelToConversation(ctx, user.AuthUser, message.MessageType, request); err != nil {
							return
						}

						newLabels, errTmp := UpdateConversationLabelList(conversation.Labels, conversation.ConversationType, request.Action, label.ExternalLabelId)
						if errTmp != nil {
							err = errTmp
							return
						}
						conversation.Labels = newLabels
					} else {
						// do nothing
					}
				}
			}

			tmp := conversation.Labels
			if err = GetLabelsInfo(ctx, conversation); err != nil {
				return
			}
			if len(action.RemoveLabels) > 0 {
				PublishConversationToManyUser(variables.EVENT_CHAT["conversation_remove_labels"], subscribers, true, conversation)
			}
			// wipe out labels detailed info and set it to label ids list
			conversation.Labels = tmp
		default:
			err = errors.New("invalid action type")
			return
		}
	}
	return
}

/*
 * send scripted message pre-defined from chat auto script or chat script to ott & es
 */
func executeSendScriptedMessage(ctx context.Context, user model.User, conversationView *model.ConversationView, timestamp int64, eventName, content string, attachments []*model.OttAttachments) error {
	if conversationView == nil {
		return errors.New("not found conversation")
	}
	if ENABLE_CHAT_POLICY_SETTINGS {
		conversationTime := conversationView.UpdatedAt
		if len(conversationTime) < 1 {
			conversationTime = conversationView.CreatedAt
		}
		if err := CheckOutOfChatWindowTime(ctx, conversationView.TenantId, conversationView.ConversationType, conversationTime); err != nil {
			return err
		}
	}

	if util.ContainKeywords(content, variables.PERSONALIZATION_KEYWORDS) {
		pageName := conversationView.OaName
		customerName := conversationView.Username
		content = strings.ReplaceAll(content, "{{page_name}}", pageName)
		content = strings.ReplaceAll(content, "{{customer_name}}", customerName)
	}

	// send message to ott
	ottMessage := model.SendMessageToOtt{
		Type:          conversationView.ConversationType,
		EventName:     eventName,
		AppId:         conversationView.AppId,
		OaId:          conversationView.OaId,
		Uid:           conversationView.ExternalUserId,
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
		TenantId:            conversationView.TenantId,
		ParentExternalMsgId: "",
		MessageId:           docId,
		MessageType:         conversationView.ConversationType,
		ConversationId:      conversationView.ConversationId,
		ExternalMsgId:       resOtt.Data.MsgId,
		EventName:           eventName,
		Direction:           variables.DIRECTION["send"],
		AppId:               conversationView.AppId,
		OaId:                conversationView.OaId,
		Avatar:              conversationView.Avatar,
		SupporterId:         user.AuthUser.UserId,
		SupporterName:       user.AuthUser.Fullname,
		SendTime:            time.Now(),
		SendTimestamp:       timestamp,
		Content:             content,
		Attachments:         attachments,
	}
	log.Info("message to es: ", scriptedMessage)

	// Should to queue
	if err := InsertMessage(ctx, conversationView.TenantId, ES_INDEX_MESSAGE, scriptedMessage.AppId, docId, scriptedMessage); err != nil {
		log.Error(err)
		return err
	}

	// >update conversation doc on ES
	conversationView.UpdatedAt = time.Now().Format(time.RFC3339)
	conversation := model.Conversation{}
	if err := util.ParseAnyToAny(conversationView, &conversation); err != nil {
		log.Error(err)
		return err
	}
	conversationQueue := model.ConversationQueue{
		DocId:        conversation.ConversationId,
		Conversation: conversation,
	}
	if err = PublishPutConversationToChatQueue(ctx, conversationQueue); err != nil {
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
			err = errors.New("queue " + user.QueueId + " not found in send event to manage")
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
		return errors.New("invalid script type")
	}
	return nil
}

func GetLabelsInfo(ctx context.Context, conversation *model.ConversationView) (err error) {
	if !reflect.DeepEqual(conversation.Labels, "") {
		labels := []any{}
		if conversation.Labels != nil {
			if err = json.Unmarshal(conversation.Labels, &labels); err != nil {
				log.Error(err)
				return
			}
		}

		if len(labels) > 0 {
			chatLabelIds := []string{}
			for _, item := range labels {
				itm := map[string]string{}
				if err = util.ParseAnyToAny(item, &itm); err != nil {
					log.Error(err)
					return err
				}
				chatLabelIds = append(chatLabelIds, itm["label_id"])
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
					conversation.Labels = tmp
				} else {
					conversation.Labels = []byte("[]")
				}
			} else {
				conversation.Labels = []byte("[]")
			}
		} else {
			conversation.Labels = []byte("[]")
		}
	} else {
		conversation.Labels = []byte("[]")
	}
	return nil
}
