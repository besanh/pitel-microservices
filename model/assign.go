package model

import "database/sql"

type AssignConversation struct {
	ConversationId string       `json:"conversation_id"`
	Status         sql.NullBool `json:"status"`
	UserId         string       `json:"user_id"`
}
