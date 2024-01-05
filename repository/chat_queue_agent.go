package repository

import "github.com/tel4vn/fins-microservices/model"

type (
	IChatQueueAgent interface {
		IRepo[model.ChatQueueAgent]
	}
	ChatQueueAgent struct {
		Repo[model.ChatQueueAgent]
	}
)

var ChatQueueAgentRepo IChatQueueAgent

func NewChatQueueAgent() IChatQueueAgent {
	return &ChatQueueAgent{}
}
