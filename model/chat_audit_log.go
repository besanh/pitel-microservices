package model

import (
	"encoding/json"
	"time"
)

type ChatAuditLog struct {
	TenantId     string          `json:"tenant_id"`
	AuditLogUuid string          `json:"audit_log_uuid"`
	Entity       string          `json:"entity"`
	Action       string          `json:"action"`
	Status       string          `json:"status"`
	OldData      json.RawMessage `json:"old_data"`
	NewData      json.RawMessage `json:"new_data"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
