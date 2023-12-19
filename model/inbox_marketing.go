package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type InboxMarketingRequest struct {
	RoutingConfig string `json:"routing_config"`
	PhoneNumber   string `json:"phone_number"`
	Content       string `json:"content"`
	Template      string `json:"template"` // template uuid
}

func (r *InboxMarketingRequest) Validate() error {
	if len(r.RoutingConfig) < 1 {
		return errors.New("routing_config is missing")
	}

	if len(r.Content) < 1 {
		return errors.New("content is missing")
	}

	if len(r.PhoneNumber) < 1 {
		return errors.New("phone_number is missing")
	}

	if len(r.Template) < 1 {
		return errors.New("template is missing")
	}

	return nil
}

type InboxMarketingLog struct {
	bun.BaseModel     `bun:"inbox_marketing_log,alias:iml"`
	Id                string          `json:"id" bun:"id,type: uuid, default: uuid_generate_v4()"`
	TenantId          string          `json:"tenant_id" bun:"tenant_id,type:text,notnull"`
	BusinessUnitId    string          `json:"business_unit_id" bun:"business_unit_id,type:text,notnull"`
	UserId            string          `json:"user_id" bun:"user_id,type:text,notnull"`
	Username          string          `json:"username" bun:"username,type:text,notnull"`
	Services          []string        `json:"services" bun:"services,type:text"`
	RoutingConfigUuid string          `json:"routing_config_uuid" bun:"routing_config_uuid,type:text"`
	FlowType          string          `json:"flow_type" bun:"flow_type,type:text"`
	FlowUuid          string          `json:"flow_uuid" bun:"flow_uuid,type:text"`
	ExternalMessageId string          `json:"external_message_id" bun:"external_message_id,type:text"`
	CampaignUuid      string          `json:"campaign_uuid" bun:"campaign_uuid,type:text"`
	Plugin            string          `json:"plugin" bun:"plugin,type:text,notnull"`
	ChannelHook       string          `json:"channel_hook" bun:"channel_hook,type:text"`
	StatusHook        string          `json:"status_hook" bun:"status_hook,type:text"`
	ErrorCodeHook     string          `json:"error_code_hook" bun:"error_code_hook,type:text"`
	PhoneNumber       string          `json:"phone_number" bun:"phone_number,type:text"`
	Message           string          `json:"message" bun:"message,type:text"`
	ListParam         json.RawMessage `json:"list_param" bun:"list_param,type:text"`
	SendTime          string          `json:"send_time" bun:"send_time,type:text"`
	TemplateCode      string          `json:"template_code" bun:"template_code,type:text"`
	Channel           string          `json:"channel" bun:"channel,type:text"`
	Status            string          `json:"status" bun:"status,type:text"`
	ErrorCode         string          `json:"error_code" bun:"error_code,type:text"`
	Quantity          int             `json:"quantity" bun:"quantity,type:integer"`
	TelcoId           int             `json:"telco_id" bun:"telco_id,type:integer"`
	RouteRule         []string        `json:"route_rule" bun:"route_rule,type:text"`
	ServiceTypeId     int             `json:"service_type_id" bun:"service_type_id,integer,default:0"`
	Ext               string          `json:"ext" bun:"ext,type:text"`
	IsChargedZns      bool            `json:"is_charged_zns" bun:"is_charged_zns,type:boolean,notnull,default:false"`
	IsCheck           bool            `json:"is_check" bun:"is_check,type:boolean,notnull,default:false"`
	Log               []string        `json:"log" bun:"log,type:text"`
	Code              int             `json:"code" bun:"code,type:integer"`
	CountAction       int             `json:"count_action" bun:"count_action,type:integer,default:0"`
	CreatedBy         string          `json:"created_by" bun:"created_by,type:text"`
	UpdatedBy         string          `json:"updated_by" bun:"updated_by,type:text"`
	CreatedAt         time.Time       `json:"created_at" bun:"created_at,type: timestamp,notnull,nullzero,default:current_timestamp"`
	UpdatedAt         time.Time       `json:"updated_at" bun:"updated_at,type: timestamp,nullzero"`
}

type InboxMarketingLogInfo struct {
	*Base
	bun.BaseModel     `bun:"inbox_marketing_log,alias:iml"`
	Id                string          `json:"id" bun:"id"`
	TenantId          string          `json:"tenant_id" bun:"tenant_id"`
	BusinessUnitId    string          `json:"business_unit_id" bun:"business_unit_id"`
	UserId            string          `json:"user_id" bun:"user_id"`
	Username          string          `json:"username" bun:"username"`
	Services          []string        `json:"services" bun:"services"`
	RoutingConfigUuid string          `json:"routing_config_uuid" bun:"routing_config_uuid"`
	FlowType          string          `json:"flow_type" bun:"flow_type"`
	FlowUuid          string          `json:"flow_uuid" bun:"flow_uuid"`
	ExternalMessageId string          `json:"external_message_id" bun:"external_message_id"`
	CampaignUuid      string          `json:"campaign_uuid" bun:"campaign_uuid"`
	Plugin            string          `json:"plugin" bun:"plugin"`
	ChannelHook       string          `json:"channel_hook" bun:"channel_hook"`
	StatusHook        string          `json:"status_hook" bun:"status_hook"`
	ErrorCodeHook     string          `json:"error_code_hook" bun:"error_code_hook"`
	PhoneNumber       string          `json:"phone_number" bun:"phone_number"`
	Message           string          `json:"message" bun:"message"`
	ListParam         json.RawMessage `json:"list_param" bun:"list_param"`
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

type InboxMarketingLogReport struct {
	Id                string          `json:"id"`
	TenantId          string          `json:"tenant_id"`
	BusinessUnitId    string          `json:"business_unit_id"`
	UserId            string          `json:"user_id"`
	Username          string          `json:"username"`
	Services          []string        `json:"services"`
	RoutingConfigUuid string          `json:"routing_config_uuid"`
	FlowType          string          `json:"flow_type"`
	FlowUuid          string          `json:"flow_uuid"`
	ExternalMessageId string          `json:"external_message_id"`
	PhoneNumber       string          `json:"phone_number"`
	Message           string          `json:"message"`
	ListParam         json.RawMessage `json:"list_param"`
	TemplateCode      string          `json:"template_code"`
	Channel           string          `json:"channel"`
	Status            string          `json:"status"`
	ErrorCode         string          `json:"error_code"`
	Quantity          int             `json:"quantity"`
	TelcoId           int             `json:"telco_id"`
	RouteRule         []string        `json:"route_rule"`
	ServiceTypeId     int             `json:"service_type_id"`
	Ext               string          `json:"ext"`
	IsChargedZns      bool            `json:"is_charged_zns"`
	CreatedBy         string          `json:"created_by"`
	UpdatedBy         string          `json:"updated_by"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type ResponseInboxMarketing struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
	// Quantity int `json:"quantity"`
}

type InboxMarketingBasic struct {
	Id                string   `json:"id"`
	TenantId          string   `json:"tenant_id"`
	BusinessUnitId    string   `json:"business_unit_id"`
	UserId            string   `json:"user_id"`
	Username          string   `json:"username"`
	Services          []string `json:"services"`
	RoutingConfigUuid string   `json:"routing_config_uuid"`
	ExternalMessageId string   `json:"external_message_id"`
	DocId             string   `json:"doc_id"`
	Index             string   `json:"index"`
	UpdatedBy         string   `json:"updated_by"`
}
