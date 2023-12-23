package model

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type InboxMarketingLogInfo struct {
	*Base
	bun.BaseModel  `bun:"inbox_marketing_log,alias:iml"`
	Id             string `json:"id" bun:"id"`
	TenantId       string `json:"tenant_id" bun:"tenant_id"`
	BusinessUnitId string `json:"business_unit_id" bun:"business_unit_id"`
	UserId         string `json:"user_id" bun:"user_id"`
	Username       string `json:"username" bun:"username"`
	// Services          []string        `json:"services" bun:"services"`
	RoutingConfigUuid string          `json:"routing_config_uuid" bun:"routing_config_uuid"`
	FlowType          string          `json:"flow_type" bun:"flow_type"`
	FlowUuid          string          `json:"flow_uuid" bun:"flow_uuid"`
	ExternalMessageId string          `json:"external_message_id" bun:"external_message_id"`
	CampaignUuid      string          `json:"campaign_uuid" bun:"campaign_uuid"`
	CampaignName      string          `json:"campaign_name" bun:"campaign_name"`
	Plugin            string          `json:"plugin" bun:"plugin"`
	ChannelHook       string          `json:"channel_hook" bun:"channel_hook"`
	StatusHook        string          `json:"status_hook" bun:"status_hook"`
	ErrorCodeHook     string          `json:"error_code_hook" bun:"error_code_hook"`
	PhoneNumber       string          `json:"phone_number" bun:"phone_number"`
	Message           string          `json:"message" bun:"message"`
	ListParam         json.RawMessage `json:"list_param" bun:"list_param"`
	TemplateUuid      string          `json:"template_uuid" bun:"template_uuid"`
	TemplateCode      string          `json:"template_code" bun:"template_code"`
	Channel           string          `json:"channel" bun:"channel"`
	Status            string          `json:"status" bun:"status"`
	ErrorCode         string          `json:"error_code" bun:"error_code"`
	Quantity          int             `json:"quantity" bun:"quantity"`
	TelcoId           int             `json:"telco_id" bun:"telco_id"`
	RouteRule         []string        `json:"route_rule" bun:"route_rule"`
	ServiceTypeId     int             `json:"service_type_id" bun:"service_type_id"`
	SendTime          string          `json:"send_time" bun:"send_time"`
	Ext               string          `json:"ext" bun:"ext"`
	IsChargedZns      bool            `json:"is_charged_zns" bun:"is_charged_zns"`
	IsCheck           bool            `json:"is_check" bun:"is_check"`
	Code              int             `json:"code" bun:"code"`
	Log               []string        `json:"log" bun:"log"`
	CountAction       int             `json:"count_action" bun:"count_action"`
	CreatedBy         string          `json:"created_by" bun:"created_by"`
	UpdatedBy         string          `json:"updated_by" bun:"updated_by"`
	CreatedAt         time.Time       `json:"created_at" bun:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" bun:"updated_at"`
}
