package repository

import "github.com/tel4vn/fins-microservices/model"

type (
	IChatTenantIntegrateSystem interface {
		IRepo[model.ChatTenantIntegrateSystem]
	}
	ChatTenantIntegrateSystem struct {
		Repo[model.ChatTenantIntegrateSystem]
	}
)

func NewChatTenantIntegrateSystem() IChatTenantIntegrateSystem {
	repo := &ChatTenantIntegrateSystem{}
	return repo
}
