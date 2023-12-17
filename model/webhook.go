package model

type WebhookSendData struct {
	Id           string `json:"id"`
	Status       string `json:"status"`
	Channel      string `json:"channel"`
	ErrorCode    string `json:"error_code"`
	Quantity     int    `json:"quantity"`
	TelcoId      int    `json:"telco_id"`
	IsChargedZns bool   `json:"is_charged_zns"`
}

type WebhookPlugin struct {
	MethodAction string `json:"method_action"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	Signature    string `json:"signature"`
	Url          string `json:"url"`
}

type WebhookIncom struct {
	IdOmniMess   string `json:"id_omni_mess" bun:"id_omni_mess"`
	Status       string `json:"status" bun:"status"`
	Channel      string `json:"channel" bun:"channel"`
	ErrorCode    string `json:"error_code" bun:"error_code"`
	Quantity     int    `json:"quantity" bun:"quantity"`
	TelcoId      int    `json:"telco_id" bun:"telco_id"`
	IsChargedZns bool   `json:"is_charged_zns" bun:"is_charged_zns"`
}

type WebhookReceiveSmsStatus struct {
	SmsGuid     string `json:"sms_guid"`
	Status      int    `json:"status"`
	SercretSign string `json:"sercret_sign"`
}
