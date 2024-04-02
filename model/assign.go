package model

import (
	"errors"
)

type AssignConversation struct {
	ConversationId string `json:"conversation_id"`
	Status         string `json:"status"`
	UserId         string `json:"user_id"`
	QueueId        string `json:"queue_id"`
}

func (s *AssignConversation) Validate() (err error) {
	if len(s.ConversationId) < 1 {
		return errors.New("conversation_id is required")
	}
	if len(s.UserId) < 1 {
		return errors.New("user_id is required")
	}
	if len(s.QueueId) < 1 {
		return errors.New("queue_id is required")
	}
	if len(s.Status) < 1 {
		return errors.New("status is required")
	}
	return
}
