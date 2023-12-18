package model

import (
	"errors"
	"strings"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type ExternalPluginConnect struct {
	*Base
	bun.BaseModel `bun:"table:external_plugin_connect,alias:epc"`
	PluginName    string  `json:"plugin_name" bun:"plugin_name,type:text"`
	PluginType    string  `json:"plugin_type" bun:"plugin_type,type:text"`
	Config        *Config `json:"config" bun:"config,type:text"`
}

type Config struct {
	IncomConfig  IncomConfig  `json:"incom"`
	AbenlaConfig AbenlaConfig `json:"abenla"`
	FptConfig    FptConfig    `json:"fpt"`
}

type IncomConfig struct {
	Username string `json:"username" bun:"username"`
	Password string `json:"password" bun:"password"`
	Api      string `json:"api" bun:"api"`
	Status   bool   `json:"stataus" bun:"status"`
}

type FptConfig struct {
	GrantType    string `json:"grant_type" bun:"grant_type"`
	ClientId     string `json:"client_id" bun:"client_id"`
	ClientSercet string `json:"client_sercet" bun:"client_sercet"`
	Scope        string `json:"scope" bun:"scope,type:text"`
	SessionId    string `json:"session_id" bun:"session_id"`
	Api          string `json:"api" bun:"api"`
	BrandName    string `json:"brand_name" bun:"brand_name"`
	Status       bool   `json:"stataus" bun:"status"`
}

type AbenlaConfig struct {
	Username string `json:"username" bun:"username"`
	Password string `json:"password" bun:"password"`
	Api      string `json:"api" bun:"api"`
	Status   bool   `json:"stataus" bun:"status"`
}

func (r *ExternalPluginConnect) Validate() error {
	if len(r.PluginName) < 1 {
		return errors.New("plugin_name is missing")
	}
	if len(r.PluginType) < 1 {
		return errors.New("plugin_type is missing")
	}
	if !slices.Contains[[]string](constants.EXTERNAL_PLUGIN_CONNECT_TYPE, r.PluginType) {
		return errors.New("plugin_type is invalid")
	}
	if r.PluginType == "incom" {
	} else if r.PluginType == "abenla" {
		if len(r.Config.AbenlaConfig.Username) < 1 {
			return errors.New("username is missing")
		}
		if len(r.Config.AbenlaConfig.Password) < 1 {
			return errors.New("password is missing")
		}
		if len(r.Config.AbenlaConfig.Api) < 1 {
			return errors.New("api is missing")
		}
	} else if r.PluginType == "fpt" {
		if len(r.Config.FptConfig.ClientId) < 1 {
			return errors.New("client_id is missing")
		}
		if len(r.Config.FptConfig.ClientSercet) < 1 {
			return errors.New("client_sercet is missing")
		}
		if len(r.Config.FptConfig.Scope) < 1 {
			return errors.New("scope is missing")
		}
		scopes := strings.Split(r.Config.FptConfig.Scope, " ")
		for _, item := range scopes {
			if !slices.Contains[[]string](constants.SCOPE_FPT, item) {
				return errors.New("scope is invalid")
			}
		}
		if len(r.Config.FptConfig.Api) < 1 {
			return errors.New("api is missing")
		}
	}
	return nil
}
