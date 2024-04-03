package model

import "errors"

type ProfileRequest struct {
	AppId          string `json:"app_id"`
	OaId           string `json:"oa_id"`
	UserId         string `json:"user_id"`
	ProfileType    string `json:"profile_type"`
	ConversationId string `json:"conversation_id"`
}

type ProfileResponse struct {
	Data struct {
		Avatar      string           `json:"avatart"`
		DisplayName string           `json:"display_name"`
		ShareInfo   ShareInfoReceive `json:"share_info"`
		UserId      string           `json:"user_id"`
		UserIdByApp string           `json:"user_id_by_app"`
	} `json:"data"`
}

type ShareInfoReceive struct {
	Address  string `json:"addresss"`
	City     string `json:"city"`
	District string `json:"district"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Ward     string `json:"ward"`
}

func (m *ProfileRequest) Validate() (err error) {
	if len(m.AppId) < 1 {
		err = errors.New("app_id is required")
		return
	}
	if len(m.OaId) < 1 {
		err = errors.New("oa_id is required")
		return
	}
	if len(m.UserId) < 1 {
		err = errors.New("user_id is required")
		return
	}
	return
}
