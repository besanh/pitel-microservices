package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatLabel struct {
	*Base
	bun.BaseModel   `bun:"table:chat_label,alias:cl"`
	TenantId        string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	AppId           string `json:"app_id" bun:"app_id,type:text,notnull"`
	OaId            string `json:"oa_id" bun:"oa_id,type:text,notnull"`
	LabelName       string `json:"label_name" bun:"label_name,type:text,notnull"`
	LabelType       string `json:"label_type" bun:"label_type,type:text,notnull"`
	ExternalLabelId string `json:"external_label_id" bun:"external_label_id,type:text"`
	LabelColor      string `json:"label_color" bun:"label_color,type:text,notnull"`
	LabelStatus     bool   `json:"label_status" bun:"label_status,notnull"`
	CreatedBy       string `json:"created_by" bun:"created_by,type:uuid,notnull"`
	UpdatedBy       string `json:"updated_by" bun:"updated_by,type:uuid,nullzero,default:null"`
}

// Use internal
type ChatLabelRequest struct {
	AppId      string `json:"app_id"`
	OaId       string `json:"oa_id"`
	LabelName  string `json:"label_name"`
	LabelColor string `json:"label_color"`
}

/* Use for external system such as zalo, facebook
* Zalo: labeling customer, removing label customer, delete label, getting labels
* Facebook: creating label, associating label, removing label, retrieving label, retrieving label details, retrieving a list of all labels
 */
type ChatExternalLabelRequest struct {
	AppId          string `json:"app_id"`
	OaId           string `json:"oa_id"`
	ExternalUserId string `json:"external_user_id"`
	LabelId        string `json:"label_id"`
	TagName        string `json:"tag_name"`
	Action         string `json:"action"`
}

type ChatExternalLabelConvertedRequest struct {
	AppId   string `json:"app_id"`
	OaId    string `json:"oa_id"`
	UserId  string `json:"user_id"`
	LabelId string `json:"label_id"`
	TagName string `json:"tag_name"`
	Action  string `json:"action"`
}

type ChatExternalLabelResponse struct {
	Code    string `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Id      string `json:"id"`
	Data    any    `json:"data"`
}

func (m *ChatLabelRequest) Validate() error {
	if len(m.AppId) < 1 {
		return errors.New("app id is required")
	}
	if len(m.OaId) < 1 {
		return errors.New("oa id is required")
	}
	if len(m.LabelName) < 1 {
		return errors.New("label name is required")
	}

	if len(m.LabelColor) < 1 {
		return errors.New("label color is required")
	}

	return nil
}
