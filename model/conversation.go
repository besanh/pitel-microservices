package model

import (
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/tel4vn/fins-microservices/common/variables"
)

type Conversation struct {
	TenantId         string          `json:"tenant_id"`
	ConversationId   string          `json:"conversation_id"`
	ConversationType string          `json:"conversation_type"`
	AppId            string          `json:"app_id"`
	OaId             string          `json:"oa_id"`
	OaName           string          `json:"oa_name"`
	OaAvatar         string          `json:"oa_avatar"`
	ShareInfo        *ShareInfo      `json:"share_info"`
	ExternalUserId   string          `json:"external_user_id"`
	Username         string          `json:"username"`
	Avatar           string          `json:"avatar"`
	Label            json.RawMessage `json:"label"`
	IsDone           bool            `json:"is_done"`
	IsDoneAt         time.Time       `json:"is_done_at"`
	IsDoneBy         string          `json:"is_done_by"`
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
}

type ConversationView struct {
	TenantId               string          `json:"tenant_id"`
	ConversationId         string          `json:"conversation_id"`
	ConversationType       string          `json:"conversation_type"`
	AppId                  string          `json:"app_id"`
	OaId                   string          `json:"oa_id"`
	OaName                 string          `json:"oa_name"`
	OaAvatar               string          `json:"oa_avatar"`
	ShareInfo              *ShareInfo      `json:"share_info"`
	ExternalUserId         string          `json:"external_user_id"`
	Username               string          `json:"username"`
	Avatar                 string          `json:"avatar"`
	Label                  json.RawMessage `json:"label"`
	IsDone                 bool            `json:"is_done"`
	IsDoneAt               string          `json:"is_done_at"`
	IsDoneBy               string          `json:"is_done_by"`
	CreatedAt              string          `json:"created_at"`
	UpdatedAt              string          `json:"updated_at"`
	TotalUnRead            int64           `json:"total_unread"`
	LatestMessageContent   string          `json:"latest_message_content"`
	LatestMessageDirection string          `json:"latest_message_direction"`
}

type ConversationCustomView struct {
	TenantId               string       `json:"tenant_id"`
	ConversationId         string       `json:"conversation_id"`
	ConversationType       string       `json:"conversation_type"`
	AppId                  string       `json:"app_id"`
	OaId                   string       `json:"oa_id"`
	OaName                 string       `json:"oa_name"`
	OaAvatar               string       `json:"oa_avatar"`
	ShareInfo              *ShareInfo   `json:"share_info"`
	ExternalUserId         string       `json:"external_user_id"`
	Username               string       `json:"username"`
	Avatar                 string       `json:"avatar"`
	Label                  *[]ChatLabel `json:"label"`
	IsDone                 bool         `json:"is_done"`
	IsDoneAt               string       `json:"is_done_at"`
	IsDoneBy               string       `json:"is_done_by"`
	CreatedAt              string       `json:"created_at"`
	UpdatedAt              string       `json:"updated_at"`
	TotalUnRead            int64        `json:"total_unread"`
	LatestMessageContent   string       `json:"latest_message_content"`
	LatestMessageDirection string       `json:"latest_message_direction"`
}

type ElasticsearchChatResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore any `json:"max_score"`
		Hits     []struct {
			Index   string   `json:"_index"`
			Type    string   `json:"_type"`
			ID      string   `json:"_id"`
			Score   any      `json:"_score"`
			Routing string   `json:"_routing"`
			Source  any      `json:"_source"`
			Sort    []string `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

type ResponseData struct {
	Data []map[string]any `json:"data"`
}

type ConversationLabelRequest struct {
	AppId           string `json:"app_id"`
	OaId            string `json:"oa_id"`
	LabelName       string `json:"label_name"`
	LabelId         string `json:"label_id"`
	ExternalLabelId string `json:"external_label_id"`
	ExternalUserId  string `json:"external_user_id"`
	ConversationId  string `json:"conversation_id"`
	Action          string `json:"action"` // create, update, remove label
}

func (m *ConversationLabelRequest) Validate() error {
	if len(m.AppId) < 1 {
		return errors.New("app id is required")
	}

	if len(m.OaId) < 1 {
		return errors.New("oa id is required")
	}

	if len(m.LabelName) < 1 {
		return errors.New("label name is required")
	}

	if len(m.ExternalUserId) < 1 {
		return errors.New("external user id is required")
	}

	if len(m.ConversationId) < 1 {
		return errors.New("conversation id is required")
	}

	if !slices.Contains(variables.CHAT_LABEL_ACTION, m.Action) {
		return errors.New("action " + m.Action + " is not supported")
	}

	return nil
}
