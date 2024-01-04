package model

import "github.com/uptrace/bun"

// Round robin
type ChatRouting struct {
	*Base
	bun.BaseModel `bun:"table:chat_routing,alias:cr"`
	RoutingName   string `bun:"routing_name,type:text,notnull"`
}
