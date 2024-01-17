package model

import (
	"encoding/json"
)

type Conversation struct {
	ConversationId       string          `json:"conversation_id"`
	ConversationType     string          `json:"conversation_type"`
	AppId                string          `json:"app_id"`
	OaId                 string          `json:"oa_id"`
	UserIdByApp          string          `json:"user_id_by_app"`
	Uid                  string          `json:"uid"`
	Username             string          `json:"username"`
	PhoneNumber          string          `json:"phone_number"`
	Email                string          `json:"email"`
	Avatar               string          `json:"avatar"`
	IsDone               bool            `json:"is_done"`
	IsDoneAt             int64           `json:"is_done_at"`
	IsDoneBy             json.RawMessage `json:"is_done_by"`
	CreatedAt            string          `json:"created_at"`
	UpdatedAt            string          `json:"updated_at"`
	TotalUnRead          int64           `json:"total_unread"`
	LatestMessageContent string          `json:"latest_message_content"`
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
