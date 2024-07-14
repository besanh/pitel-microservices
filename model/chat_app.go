package model

import (
	"errors"
	"slices"

	"github.com/uptrace/bun"
)

/**
* app_name outside use for management, helping to track
* app_name inside use for creating account ott
 */
type ChatApp struct {
	*Base
	bun.BaseModel `bun:"table:chat_app,alias:ca"`
	AppName       string   `json:"app_name" bun:"app_name,type:text,notnull"`
	Status        string   `json:"status" bun:"status,notnull"`
	InfoApp       *InfoApp `json:"info_app" bun:"info_app,type:jsonb,notnull"`

	// relations
	Systems []*ChatIntegrateSystem `json:"integrate_systems" bun:"m2m:chat_app_integrate_system,join:ChatApp=ChatIntegrateSystem"`
}

type ChatAppRequest struct {
	AppName   string   `json:"app_name"`
	Status    string   `json:"status"` //active/deactive
	InfoApp   *InfoApp `json:"info_app"`
	SystemIds []string `json:"system_ids"`
}

type InfoApp struct {
	Zalo     *Zalo     `json:"zalo" bun:"zalo"`
	Facebook *Facebook `json:"facebook" bun:"facebook"`
}

type Zalo struct {
	AppId     string `json:"app_id"`
	AppName   string `json:"app_name"`
	AppSecret string `json:"app_secret"`
}

type Facebook struct {
	AppId     string `json:"app_id"`
	AppName   string `json:"app_name"`
	AppSecret string `json:"app_secret"`
}

func (m *ChatAppRequest) Validate() error {
	if len(m.AppName) < 1 {
		return errors.New("app name is required")
	}

	var countOk int

	if m.InfoApp.Zalo != nil && m.InfoApp.Facebook != nil {
		return errors.New("only one info app is allowed")
	}
	if m.InfoApp.Zalo == nil && m.InfoApp.Facebook == nil {
		return errors.New("info app is required")
	}

	if m.InfoApp.Zalo != nil {
		if len(m.InfoApp.Zalo.AppId) < 1 {
			return errors.New("app id is required")
		}
		if len(m.InfoApp.Zalo.AppName) < 1 {
			return errors.New("app name is required")
		}
		if len(m.InfoApp.Zalo.AppSecret) < 1 {
			return errors.New("app secret is required")
		}
		countOk += 1
	}

	if m.InfoApp.Facebook != nil {
		if len(m.InfoApp.Facebook.AppId) < 1 {
			return errors.New("app id is required")
		}
		if len(m.InfoApp.Facebook.AppName) < 1 {
			return errors.New("app name is required")
		}
		if len(m.InfoApp.Facebook.AppSecret) < 1 {
			return errors.New("app secret is required")
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
