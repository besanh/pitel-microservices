package model

import (
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
)

type ChatApp struct {
	*Base
	bun.BaseModel `bun:"table:chat_app,alias:ca"`
	AppName       string   `json:"app_name" bun:"app_name,type:text,notnull"`
	State         string   `json:"state" bun:"state,type:text,notnull"`
	Status        bool     `json:"status" bun:"status,notnull"`
	InfoApp       *InfoApp `json:"info_app" bun:"info_app,type:text"`
}

type AppRequest struct {
	AppName string       `json:"app_name" bun:"app_name,type:text,notnull"`
	State   string       `json:"state" bun:"state,type:text,notnull"`
	Status  sql.NullBool `json:"status" bun:"status,notnull"`
	InfoApp *InfoApp     `json:"info_app" bun:"info_app,type:text"`
}

type InfoApp struct {
	Zalo     *ZaloApp     `json:"zalo_app" bun:"zalo_app,type:text"`
	Facebook *FacebookApp `json:"facebook_app" bun:"facebook_app,type:text"`
}

type ZaloApp struct {
	AppId       string `json:"app_id"`
	AppName     string `json:"app_name"`
	SecretKey   string `json:"secret_key"`
	AccessToken string `json:"access_token"`
	Status      bool   `json:"status"`
}

type FacebookApp struct {
	AppId   string `json:"app_id"`
	AppName string `json:"app_name"`
	Token   string `json:"token"`
	Status  bool   `json:"status"`
}

func (m *AppRequest) Validate() error {
	if len(m.AppName) < 1 {
		return errors.New("app name is required")
	}
	if len(m.State) < 1 {
		return errors.New("state is required")
	}
	if !m.Status.Valid {
		return errors.New("status is required")
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
		if len(m.InfoApp.Zalo.AccessToken) < 1 {
			return errors.New("access token is required")
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
		if len(m.InfoApp.Facebook.Token) < 1 {
			return errors.New("token is required")
		}
		// countOk += 1
	}
	// if countOk > 1 {
	// 	return errors.New("only one app can be active")
	// }
	return nil
}
