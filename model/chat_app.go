package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type ChatApp struct {
	*Base
	bun.BaseModel `bun:"table:chat_app,alias:ca"`
	AppName       string   `json:"app_name" bun:"app_name,type:text,notnull"`
	Status        bool     `json:"status" bun:"status,notnull"`
	InfoApp       *InfoApp `json:"info_app" bun:"info_app,type:jsonb,notnull"`
}

type ChatAppRequest struct {
	AppName string   `json:"app_name"`
	State   string   `json:"state"`
	Status  bool     `json:"status"`
	InfoApp *InfoApp `json:"info_app"`
}

type InfoApp struct {
	Zalo     *Zalo     `json:"zalo" bun:"zalo"`
	Facebook *Facebook `json:"facebook" bun:"facebook"`
}

type Zalo struct {
	AppId         string `json:"app_id"`
	AppName       string `json:"app_name"`
	SecretKey     string `json:"secret_key"`
	OaId          string `json:"oa_id"`
	OaName        string `json:"oa_name"`
	State         string `json:"state"`
	CodeChallenge string `json:"code_challenge"`
	Status        bool   `json:"status"`
}

type Facebook struct {
	AppId    string `json:"app_id"`
	AppName  string `json:"app_name"`
	AppToken string `json:"app_token"`
	Status   bool   `json:"status"`
}

func (m *ChatAppRequest) Validate() error {
	if len(m.AppName) < 1 {
		return errors.New("app name is required")
	}

	// var countOk int
	if m.InfoApp.Zalo.Status {
		if len(m.InfoApp.Zalo.AppId) < 1 {
			return errors.New("app id is required")
		}
		if len(m.InfoApp.Zalo.AppName) < 1 {
			return errors.New("app name is required")
		}
		if len(m.InfoApp.Zalo.SecretKey) < 1 {
			return errors.New("secret key is required")
		}
		if len(m.InfoApp.Zalo.OaId) < 1 {
			return errors.New("oat id is required")
		}
		if len(m.InfoApp.Zalo.OaName) < 1 {
			return errors.New("oa name is required")
		}
		if len(m.InfoApp.Zalo.State) < 1 {
			return errors.New("state is required")
		}
		// countOk += 1
	}
	if m.InfoApp.Facebook.Status {
		if len(m.InfoApp.Facebook.AppId) < 1 {
			return errors.New("app id is required")
		}
		if len(m.InfoApp.Facebook.AppName) < 1 {
			return errors.New("app name is required")
		}
		if len(m.InfoApp.Facebook.AppToken) < 1 {
			return errors.New("token is required")
		}
		// countOk += 1
	}
	// if countOk > 1 {
	// 	return errors.New("only one app can be active")
	// }
	return nil
}
