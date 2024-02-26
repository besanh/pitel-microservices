package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ShareInfoForm struct {
	*Base
	bun.BaseModel `bun:"table:share_info,alias:si"`
	TenantId      string    `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ShareType     string    `json:"share_type" bun:"share_type,type:text,notnull"` //zalo,fb,....
	ShareForm     ShareForm `json:"share_form" bun:"share_form,type:jsonb,notnull"`
}

type ShareInfoFormRequest struct {
	ShareType string `json:"share_type"`
	EventName string `json:"event_name"`
	AppId     string `json:"app_id"`
	OaId      string `json:"oa_id"`
	Uid       string `json:"uid"`
	ImageUrl  string `json:"image_url"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
}

type ShareForm struct {
	Facebook struct{} `json:"facebook"`
	Zalo     struct {
		ImageUrl string `json:"image_url"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
	} `json:"zalo"`
}

type OttShareInfoRequest struct {
	Type      string `json:"type"`
	EventName string `json:"event_name"`
	AppId     string `json:"app_id"`
	OaId      string `json:"oa_id"`
	Uid       string `json:"uid"`
	ImageUrl  string `json:"image_url"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
}

func (s *ShareInfoFormRequest) Validate() (err error) {
	if len(s.ShareType) < 1 {
		return errors.New("share type is required")
	}
	if len(s.AppId) < 1 {
		return errors.New("app id is required")
	}
	if len(s.OaId) < 1 {
		return errors.New("oa id is required")
	}
	if len(s.Uid) < 1 {
		return errors.New("uid is required")
	}
	return
}
