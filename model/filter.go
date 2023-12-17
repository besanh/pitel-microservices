package model

import "database/sql"

type BalanceConfigFilter struct {
	Weight      []string     `json:"weight"`
	BalanceType []string     `json:"balance_type"`
	Provider    []string     `json:"provider"`
	Priority    []string     `json:"priority"`
	Status      sql.NullBool `json:"status"`
}

type PluginConfigFilter struct {
	PluginName []string     `json:"plugin_name"`
	PluginType []string     `json:"plugin_type"`
	Status     sql.NullBool `json:"status"`
}

type RecipientConfigFilter struct {
	Recipient     []string     `json:"recipient"`
	RecipientType []string     `json:"recipient_type"`
	Priority      []string     `json:"priority"`
	Provider      []string     `json:"provider"`
	Status        sql.NullBool `json:"status"`
}

type TemplateBssFilter struct {
	TemplateName string       `json:"template_name"`
	TemplateCode []string     `json:"template_code"`
	TemplateType []string     `json:"template_type"`
	Content      string       `json:"content"`
	Status       sql.NullBool `json:"status"`
}

type RoutingConfigFilter struct {
	RoutingName string       `json:"routing_name"`
	RoutingType []string     `json:"routing_type"`
	Brandname   string       `json:"brand_name"`
	Status      sql.NullBool `json:"status"`
}

type InboxMarketingFilter struct {
	Id                string       `json:"id"`
	TenantId          string       `json:"tenant_id"`
	BusinessUnitId    string       `json:"business_unit_id"`
	UserId            string       `json:"user_id"`
	Username          string       `json:"username"`
	Services          []string     `json:"services"`
	RoutingConfigUuid string       `json:"routing_config_uuid"`
	Plugin            []string     `json:"plugin" bun:"plugin,array"`
	PhoneNumber       string       `json:"phone_number"`
	Message           string       `json:"message"`
	TemplateCode      string       `json:"template_code"`
	Channel           []string     `json:"channel"`
	Status            []string     `json:"status"`
	ErrorCode         []string     `json:"error_code"`
	Quantity          string       `json:"quantity"`
	TelcoId           []int        `json:"telco_id"`
	RouteRule         []string     `json:"route_rule"`
	ServiceTypeId     []int        `json:"service_type_id"`
	SendTime          string       `json:"send_time"`
	Ext               string       `json:"ext"`
	IsChargedZns      sql.NullBool `json:"is_charged_zns"`
	Code              string       `json:"code"`
	CountAction       string       `json:"count_action"`
	CampaignUuid      []string     `json:"campaign_uuid"`
	StartTime         string       `json:"start_time"`
	EndTime           string       `json:"end_time"`
	Limit             int          `json:"limit"`
	Offset            int          `json:"offset"`
}
