package model

import "time"

type Conversation struct {
	TenantId         string     `json:"tenant_id"`
	ConversationId   string     `json:"conversation_id"`
	ConversationType string     `json:"conversation_type"`
	AppId            string     `json:"app_id"`
	OaId             string     `json:"oa_id"`
	OaName           string     `json:"oa_name"`
	OaAvatar         string     `json:"oa_avatar"`
	ShareInfo        *ShareInfo `json:"share_info"`
	ExternalUserId   string     `json:"external_user_id"`
	Username         string     `json:"username"`
	Avatar           string     `json:"avatar"`
	IsDone           bool       `json:"is_done"`
	IsDoneAt         time.Time  `json:"is_done_at"`
	IsDoneBy         string     `json:"is_done_by"`
	CreatedAt        string     `json:"created_at"`
	UpdatedAt        string     `json:"updated_at"`
}

type ConversationView struct {
	TenantId               string     `json:"tenant_id"`
	ConversationId         string     `json:"conversation_id"`
	ConversationType       string     `json:"conversation_type"`
	AppId                  string     `json:"app_id"`
	OaId                   string     `json:"oa_id"`
	OaName                 string     `json:"oa_name"`
	OaAvatar               string     `json:"oa_avatar"`
	ShareInfo              *ShareInfo `json:"share_info"`
	ExternalUserId         string     `json:"external_user_id"`
	Username               string     `json:"username"`
	Avatar                 string     `json:"avatar"`
	IsDone                 bool       `json:"is_done"`
	IsDoneAt               string     `json:"is_done_at"`
	IsDoneBy               string     `json:"is_done_by"`
	CreatedAt              string     `json:"created_at"`
	UpdatedAt              string     `json:"updated_at"`
	TotalUnRead            int64      `json:"total_unread"`
	LatestMessageContent   string     `json:"latest_message_content"`
	LatestMessageDirection string     `json:"latest_message_direction"`
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
