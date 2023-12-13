package model

import (
	"errors"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type PluginConfig struct {
	*Base
	bun.BaseModel `bun:"table:plugin_config,alias:pc"`
	PluginName    string `json:"plugin_name" bun:"plugin_name,type:text"`
	PluginType    string `json:"plugin_type" bun:"plugin_type,type:text"`
	Status        bool   `json:"status" bun:"status,type:boolean"`
}

type PluginConfigView struct {
	bun.BaseModel `bun:"table:plugin_config,alias:pc"`
	Id            string `json:"id" bun:"id,type:uuid"`
	PluginName    string `json:"plugin_name" bun:"plugin_name,type:text"`
	PluginType    string `json:"plugin_type" bun:"plugin_type,type:text"`
	Status        bool   `json:"status" bun:"status,type:boolean"`
}

type PluginConfigRequest struct {
	PluginName string `json:"plugin_name"`
	PluginType string `json:"plugin_type"`
	Status     bool   `json:"status"`
}

func (r *PluginConfigRequest) Validate() error {
	if len(r.PluginName) < 1 {
		return errors.New("plugin_name is missing")
	}
	if !slices.Contains[[]string](constants.CHANNEL, r.PluginType) {
		return errors.New("plugin_type is missing")
	}
	return nil
}
