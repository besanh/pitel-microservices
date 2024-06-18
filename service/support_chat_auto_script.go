package service

import (
	"github.com/tel4vn/fins-microservices/model"
	"sort"
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
