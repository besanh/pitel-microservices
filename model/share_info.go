package model

import (
	"errors"
	"mime/multipart"

	"github.com/uptrace/bun"
)

type ShareInfoForm struct {
	*Base
	bun.BaseModel `bun:"table:chat_share_info,alias:csi"`
	TenantId      string    `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	ShareType     string    `json:"share_type" bun:"share_type,type:text,notnull"` //zalo,fb,....
	ShareForm     ShareForm `json:"share_form" bun:"share_form,type:jsonb,notnull"`
}

type ShareInfoFormRequest struct {
	Id             string                `form:"id"`
	ShareType      string                `form:"share_type" binding:"required"`
	EventName      string                `form:"event_name"`
	AppId          string                `form:"app_id" binding:"required"`
	OaId           string                `form:"oa_id" binding:"required"`
	ExternalUserId string                `form:"external_user_id"`
	ImageUrl       string                `form:"image_url"`
	Title          string                `form:"title" binding:"required"`
	Subtitle       string                `form:"subtitle" binding:"required"`
	Files          *multipart.FileHeader `form:"file" binding:"required"`
}

type ShareInfoFormSubmitRequest struct {
	ShareType      string `json:"share_type"`
	EventName      string `json:"event_name"`
	AppId          string `json:"app_id"`
	ExternalUserId string `json:"external_user_id"`
}

type ShareForm struct {
	Facebook struct{} `json:"facebook"`
	Zalo     struct {
		AppId    string `json:"app_id"`
		OaId     string `json:"oa_id"`
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
	return
}

func (s *ShareInfoFormSubmitRequest) Validate() (err error) {
	if len(s.ShareType) < 1 {
		return errors.New("share type is required")
	}
	if len(s.AppId) < 1 {
		return errors.New("app id is required")
	}
	if len(s.ExternalUserId) < 1 {
		return errors.New("external_user_id is required")
	}
	return
}
