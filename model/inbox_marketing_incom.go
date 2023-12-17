package model

import "encoding/json"

type IncomSendMessageResponse struct {
	IdOmniMess string `json:"idOmniMess"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

type IncomBodyStatus struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IdOmniMess string `json:"idOmniMess"`
}

type IncomStatusMess struct {
	IncomUuid      string          `json:"incom_uuid" bun:"incom_uuid"`
	IdOmniMess     string          `json:"id_omni_mess" bun:"id_omni_mess"`
	PhoneNumber    string          `json:"phonenumber" bun:"phonenumber"`
	ListParam      json.RawMessage `json:"list_param" bun:"list_param"`
	CreateDatetime string          `json:"createdatetime" bun:"createdatetime"`
	TemplateCode   string          `json:"templatecode" bun:"templatecode"`
	Status         string          `json:"status" bun:"status"`
	Channel        string          `json:"channel" bun:"channel"`
	ErrorCode      string          `json:"errorcode" bun:"errorcode"`
	MtCount        string          `json:"mtcount" bun:"mtcount"`
	TelcoId        string          `json:"telcoid" bun:"telcoid"`
	Ischarged      string          `json:"ischarged" bun:"ischarged"`
}
