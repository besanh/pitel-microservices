package repository

import "github.com/tel4vn/fins-microservices/model"

type (
	IChatAppIntegrateSystem interface {
		IRepo[model.ChatAppIntegrateSystem]
	}
	ChatAppIntegrateSystem struct {
		Repo[model.ChatAppIntegrateSystem]
	}
)

var ChatAppIntegrateSystemRepo IChatAppIntegrateSystem

func NewChatAppIntegrateSystem() IChatAppIntegrateSystem {
	return &ChatAppIntegrateSystem{}
}
