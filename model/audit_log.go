package model

type LogInboxMarketing struct {
	TenantId          string   `json:"tenant_id"`
	BusinessUnitId    string   `json:"business_unit_id"`
	UserId            string   `json:"user_id"`
	Username          string   `json:"username"`
	Services          []string `json:"services"`
	Id                string   `json:"id"`
	Plugin            string   `json:"plugin"`
	RoutingConfigUuid string   `json:"routing_config_uuid"`
	FlowType          string   `json:"flow_type" bun:"flow_type,type:text"`
	FlowUuid          string   `json:"flow_uuid" bun:"flow_uuid,type:text"`
	ExternalMessageId string   `json:"external_message_id"`
	CampaignUuid      string   `json:"campaign_uuid" bun:"campaign_uuid,type:text"`
	Status            string   `json:"status"`
	Channel           string   `json:"channel"`
	ChannelHook       string   `json:"channel_hook"`
	ErrorCode         string   `json:"error_code"`
	ErrorCodeHook     string   `json:"error_code_hook"`
	Quantity          int      `json:"quantity"`
	TelcoId           int      `json:"telco_id"`
	IsChargedZns      bool     `json:"is_charged_zns"`
	IsCheck           bool     `json:"is_check"`
	Code              int      `json:"code"`
	CountAction       int      `json:"count_action"`
	UpdatedBy         string   `json:"updated_by"`
}
