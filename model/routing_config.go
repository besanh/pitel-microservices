package model

import (
	"github.com/uptrace/bun"
)

type RoutingConfig struct {
	*Base
	bun.BaseModel `bun:"table:routing_config,alias:brc"`
	RoutingName   string        `json:"routing_name" bun:"routing_name,type:text,notnull"`
	RoutingType   string        `json:"routing_type" bun:"routing_type,type:text,notnull"` // sms,zns,...
	RoutingFlow   RoutingFlow   `json:"routing_flow" bun:"routing_flow,type:text,notnull"`
	RoutingOption RoutingOption `json:"routing_option" bun:"routing_option,type:text,notnull"`
	Status        bool          `json:"status" bun:"status,type:boolean"`
}

type RoutingConfigView struct {
	bun.BaseModel `bun:"table:routing_config,alias:brc"`
	RoutingName   string        `json:"routing_name" bun:"routing_name,type:text,notnull"`
	RoutingType   string        `json:"routing_type" bun:"routing_type,type:text,notnull"` // sms,zns,...
	RoutingFlow   RoutingFlow   `json:"routing_flow" bun:"routing_flow,type:text,notnull"`
	RoutingOption RoutingOption `json:"routing_option" bun:"routing_option,type:text,notnull"`
	Status        bool          `json:"status" bun:"status,type:boolean"`
}

// Link to table recipient_Routing or balance Routing to control flow send data
type RoutingFlow struct {
	FlowType string `json:"flow_type"` // table recipient or balance
	FlowUuid string `json:"flow_uuid"`
}

// Include account info connected with external plugin
type RoutingOption struct {
	Incom  Incom  `json:"incom" bun:"incom,type:text"`
	Abenla Abenla `json:"abenla" bun:"abenla,type:text"`
	Fpt    Fpt    `json:"fpt" bun:"fpt,type:text"`
}

type Incom struct {
	Username          string `json:"username" bun:"username"`
	Password          string `json:"password" bun:"password"`
	ApiAuthUrl        string `json:"api_auth_url" bun:"api_auth_url,type:text"`
	ApiSendMessageUrl string `json:"api_send_message_url" bun:"api_send_message_url,type:text"`
	WebhookUrl        string `json:"webhook_url" bun:"webhook_url,type:text"`
	MaxAttempts       int    `json:"max_attempts" bun:"max_attempts,type:text"`
	Signature         string `json:"signature" bun:"signature,type:text"`
	Status            bool   `json:"status" bun:"status"`
}

type Abenla struct {
	Username          string `json:"username" bun:"username"`
	Password          string `json:"password" bun:"password"`
	ApiAuthUrl        string `json:"api_auth_url" bun:"api_auth_url,type:text"`
	ApiSendMessageUrl string `json:"api_send_message_url" bun:"api_send_message_url,type:text"`
	ServiceTypeId     string `json:"service_type_id" bun:"service_type_id"`
	WebhookUrl        string `json:"webhook_url" bun:"webhook_url,type:text"`
	MaxAttempts       int    `json:"max_attempts" bun:"max_attempts,type:text"`
	Signature         string `json:"signature" bun:"signature,type:text"`
	Brandname         string `json:"brand_name" bun:"brand_name,type:text"`
	Status            bool   `json:"status" bun:"status"`
}

type Fpt struct {
	ClientId          string `json:"client_id" bun:"client_id,type:text"`
	ClientSecret      string `json:"client_secret" bun:"client_secret,type:text"`
	Scope             string `json:"scope" bun:"scope,type:text"`
	GrantType         string `json:"grant_type" bun:"grant_type,type:text"`
	ApiAuthUrl        string `json:"api_auth_url" bun:"api_auth_url,type:text"`
	ApiSendMessageUrl string `json:"api_send_message_url" bun:"api_send_message_url,type:text"`
	WebhookUrl        string `json:"webhook_url" bun:"webhook_url,type:text"`
	MaxAttempts       int    `json:"max_attempts" bun:"max_attempts,type:text"`
	Signature         string `json:"signature" bun:"signature,type:text"`
	BrandName         string `json:"brand_name" bun:"brand_name,type:text"`
	Status            bool   `json:"status" bun:"status"`
}
