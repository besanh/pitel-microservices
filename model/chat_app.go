package model

import (
	"errors"
	"slices"

	"github.com/uptrace/bun"
)

type ChatApp struct {
	*Base
	bun.BaseModel `bun:"table:chat_app,alias:ca"`
	AppName       string   `json:"app_name" bun:"app_name,type:text,notnull"`
	Status        string   `json:"status" bun:"status,notnull"`
	DefaultApp    string   `json:"default_app" bun:"default_app,type:text,notnull"`
	InfoApp       *InfoApp `json:"info_app" bun:"info_app,type:jsonb,notnull"`
}

type ChatAppRequest struct {
	AppName    string   `json:"app_name"`
	Status     string   `json:"status"` //active/deactive
	DefaultApp string   `json:"default_app"`
	InfoApp    *InfoApp `json:"info_app"`
}

type InfoApp struct {
	Zalo     *Zalo     `json:"zalo" bun:"zalo"`
	Facebook *Facebook `json:"facebook" bun:"facebook"`
}

type Zalo struct {
	AppId     string `json:"app_id"`
	AppName   string `json:"app_name"`
	AppSecret string `json:"app_secret"`
	Status    string `json:"status"` //active/deactive
	Active    bool   `json:"active"`
}

type Facebook struct {
	AppId     string `json:"app_id"`
	AppName   string `json:"app_name"`
	AppSecret string `json:"app_secret"`
	Status    string `json:"status"`
}

func (m *ChatAppRequest) Validate() error {
	if len(m.AppName) < 1 {
		return errors.New("app name is required")
	}

	if len(m.DefaultApp) < 1 {
		return errors.New("default app is required")
	}

	var countOk int

	if m.InfoApp.Zalo != nil && m.InfoApp.Zalo.Status == "active" {
		if len(m.InfoApp.Zalo.AppId) < 1 {
			return errors.New("app id is required")
		}
		if len(m.InfoApp.Zalo.AppName) < 1 {
			return errors.New("app name is required")
		}
		if len(m.InfoApp.Zalo.AppSecret) < 1 {
			return errors.New("app secret is required")
		}
		if len(m.InfoApp.Zalo.Status) < 1 {
			return errors.New("status is required")
		}
		if !slices.Contains([]string{"active", "deactive"}, m.InfoApp.Zalo.Status) {
			return errors.New("status zalo " + m.InfoApp.Zalo.Status + " is not supported")
		}
		countOk += 1
	}

	if m.InfoApp.Facebook != nil && m.InfoApp.Facebook.Status == "active" {
		if len(m.InfoApp.Facebook.AppId) < 1 {
			return errors.New("app id is required")
		}
		if len(m.InfoApp.Facebook.AppName) < 1 {
			return errors.New("app name is required")
		}
		if len(m.InfoApp.Facebook.AppSecret) < 1 {
			return errors.New("app secret is required")
		}
		if len(m.InfoApp.Facebook.Status) < 1 {
			return errors.New("status is required")
		}
		if !slices.Contains([]string{"active", "deactive"}, m.InfoApp.Facebook.Status) {
			return errors.New("status facebook " + m.InfoApp.Facebook.Status + " is not supported")
		}
		countOk += 1
	}

	if !slices.Contains([]string{"active", "deactive"}, m.Status) {
		return errors.New("status " + m.Status + " is not supported")
	}

	if countOk > 1 {
		return errors.New("only one app can be active")
	}

	return nil
}
