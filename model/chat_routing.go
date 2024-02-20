package model

import (
	"errors"

	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

// Round robin
type ChatRouting struct {
	*Base
	bun.BaseModel `bun:"table:chat_routing,alias:cr"`
	TenantId      string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	RoutingName   string `json:"routing_name" bun:"routing_name,type:text,notnull"`
	RoutingAlias  string `json:"routing_alias" bun:"routing_alias,type:text,notnull"`
	Status        bool   `json:"status" bun:"status,notnull"`
}

type ChatRoutingRequest struct {
	RoutingName  string `json:"routing_name"`
	RoutingAlias string `json:"routing_alias"`
	Status       bool   `json:"status"`
}

func (m *ChatRoutingRequest) Validate() error {
	if len(m.RoutingName) < 1 {
		return errors.New("routing name is required")
	}
	if !slices.Contains[[]string](variables.CHAT_ROUTING, m.RoutingAlias) {
		return errors.New("chat routing method " + m.RoutingAlias + " is not supported")
	}
	return nil
}
