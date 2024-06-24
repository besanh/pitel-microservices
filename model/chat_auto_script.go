package model

import (
	"errors"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type ScriptActionType string

const (
	MoveToExistedScript ScriptActionType = "existed_script"
	SendMessage         ScriptActionType = "send_message"
	AddLabels           ScriptActionType = "add_labels"
	RemoveLabels        ScriptActionType = "remove_labels"
)

type ChatAutoScript struct {
	*Base
	bun.BaseModel      `bun:"table:chat_auto_script,alias:cas"`
	TenantId           string                        `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ScriptName         string                        `json:"script_name" bun:"script_name,type:text,notnull"`
	Channel            string                        `json:"channel" bun:"channel,type:text,notnull"`
	ConnectionId       string                        `json:"connection_id" bun:"connection_id,type:uuid,notnull"`
	ConnectionApp      *ChatConnectionApp            `json:"connection_app" bun:"rel:belongs-to,join:connection_id=id"`
	CreatedBy          string                        `json:"created_by" bun:"created_by,type:uuid,notnull"`
	UpdatedBy          string                        `json:"updated_by" bun:"updated_by,type:uuid,default:null"`
	Status             bool                          `json:"status" bun:"status,type:boolean,notnull"`
	TriggerEvent       string                        `json:"trigger_event" bun:"trigger_event,type:text,notnull"`
	TriggerKeywords    AutoScriptTriggerKeywordsType `json:"trigger_keywords" bun:"trigger_keywords,type:jsonb"`
	ChatScriptLink     []*ChatAutoScriptToChatScript `json:"chat_script_link" bun:"rel:has-many,join:id=chat_auto_script_id"`
	SendMessageActions AutoScriptSendMessage         `json:"send_message_actions" bun:"send_message_actions,type:jsonb"`
	ChatLabelLink      []*ChatAutoScriptToChatLabel  `json:"chat_label_link" bun:"rel:has-many,join:id=chat_auto_script_id"`
}

type AutoScriptTriggerKeywordsType struct {
	Keywords []string `json:"keywords"`
}

type AutoScriptSendMessage struct {
	Actions []AutoScriptSendMessageType `json:"actions"`
}

type AutoScriptSendMessageType struct {
	Content string `json:"content"`
	Order   int    `json:"order"`
}

type AutoScriptMergedActions struct {
	Actions []ActionScriptActionType `json:"actions"`
}

type ActionScriptActionType struct {
	Type         string   `json:"type"`
	ChatScriptId string   `json:"chat_script_id"` // use when user selected using an existed chat script
	Content      string   `json:"content"`        // sending message action
	Order        int      `json:"order"`
	AddLabels    []string `json:"add_labels"`
	RemoveLabels []string `json:"remove_labels"`
}

type ChatAutoScriptRequest struct {
	ScriptName      string                        `json:"script_name" form:"script_name" binding:"required"`
	Channel         string                        `json:"channel" form:"channel" binding:"required"`
	ConnectionId    string                        `json:"connection_id" form:"connection_id" binding:"required"`
	Status          string                        `json:"status" form:"status" binding:"required"`
	TriggerEvent    string                        `json:"trigger_event" form:"trigger_event" binding:"required"`
	TriggerKeywords AutoScriptTriggerKeywordsType `json:"trigger_keywords" form:"trigger_keywords"`
	ActionScript    *AutoScriptMergedActions      `json:"action_script" form:"action_script" binding:"required"`
}

type ChatAutoScriptStatusRequest struct {
	Status string `json:"status" form:"status" binding:"required"`
}

type ChatAutoScriptView struct {
	*Base
	bun.BaseModel      `bun:"table:chat_auto_script,alias:cas"`
	TenantId           string                        `json:"tenant_id" bun:"tenant_id"`
	ScriptName         string                        `json:"script_name" bun:"script_name"`
	Channel            string                        `json:"channel" bun:"channel"`
	ConnectionId       string                        `json:"connection_id" bun:"connection_id"`
	ConnectionApp      *ChatConnectionApp            `json:"connection_app" bun:"rel:belongs-to,join:connection_id=id"`
	CreatedBy          string                        `json:"created_by" bun:"created_by"`
	UpdatedBy          string                        `json:"updated_by" bun:"updated_by"`
	Status             bool                          `json:"status" bun:"status"`
	TriggerEvent       string                        `json:"trigger_event" bun:"trigger_event"`
	TriggerKeywords    AutoScriptTriggerKeywordsType `json:"trigger_keywords" bun:"trigger_keywords"`
	ChatScriptLink     []*ChatAutoScriptToChatScript `json:"chat_script_link" bun:"rel:has-many,join:id=chat_auto_script_id"`
	SendMessageActions AutoScriptSendMessage         `json:"send_message_actions" bun:"send_message_actions"`
	ChatLabelLink      []*ChatAutoScriptToChatLabel  `json:"chat_label_link" bun:"rel:has-many,join:id=chat_auto_script_id"`
	ActionScript       *AutoScriptMergedActions      `json:"action_script" bun:"-"`
}

func (r *ChatAutoScriptRequest) Validate() error {
	if len(r.ScriptName) < 1 {
		return errors.New("script name is required")
	}
	if len(r.ConnectionId) < 1 {
		return errors.New("connection id is required")
	}
	if len(r.TriggerEvent) < 1 {
		return errors.New("trigger event is required")
	}
	if len(r.Channel) < 1 {
		return errors.New("channel is required")
	}
	if !slices.Contains[[]string](variables.CONNECTION_TYPE, r.Channel) {
		return errors.New("connection type " + r.Channel + " is not supported")
	}
	if r.ActionScript == nil || len(r.ActionScript.Actions) < 1 {
		return errors.New("action script is required")
	}
	if len(r.ActionScript.Actions) > 3 {
		return errors.New("amount of actions must not exceed 3")
	}
	if !slices.Contains[[]string](variables.CHAT_AUTO_SCRIPT_EVENT, r.TriggerEvent) {
		return errors.New("trigger event " + r.TriggerEvent + " is not supported")
	}
	if r.TriggerEvent == "keyword" && len(r.TriggerKeywords.Keywords) < 1 {
		return errors.New("keyword is required")
	}

	for _, action := range r.ActionScript.Actions {
		switch ScriptActionType(action.Type) {
		case MoveToExistedScript:
			if len(action.ChatScriptId) < 1 {
				return errors.New("chat script id is required")
			}
		case SendMessage:
			if len(action.Content) < 1 {
				return errors.New("message's content is required")
			}
		case AddLabels:
			if len(action.AddLabels) < 1 {
				return errors.New("label id is required")
			}
		case RemoveLabels:
			if len(action.RemoveLabels) < 1 {
				return errors.New("label id is required")
			}
		default:
			return errors.New("invalid action type")
		}
	}

	return nil
}
