package model

import (
	"errors"

	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type FacebookPage struct {
	*Base
	bun.BaseModel `bun:"table:chat_facebook_page,alias:cfb"`
	TenantId      string `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	AppId         string `json:"app_id" bun:"app_id,type:text,notnull"`
	OaId          string `json:"oa_id" bun:"oa_id,type:text,notnull"`
	OaName        string `json:"oa_name" bun:"oa_name,type:text,notnull"`
	TokenType     string `json:"token_type" bun:"token_type,type:text,notnull"`
	AccessToken   string `json:"access_token" bun:"access_token,type:text,notnull"`
	Avatar        string `json:"avatar" bun:"avatar,type:text,notnull"`
	Status        string `json:"status" bun:"status,notnull"`
}

type FacebookPageInfo struct {
	OaId        string `json:"oa_id"`
	OaName      string `json:"oa_name"`
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	Avatar      string `json:"avatar"`
	Status      string `json:"status"`
}

func (m *FacebookPageInfo) Validate() (err error) {
	if len(m.OaId) < 1 {
		err = errors.New("oa id is required")
	}

	if len(m.OaName) < 1 {
		err = errors.New("oa name is required")
	}

	if len(m.TokenType) < 1 {
		err = errors.New("token type is required")
	}
	if !slices.Contains([]string{"long_lived_token"}, m.TokenType) {
		err = errors.New("token type " + m.TokenType + " is not supported")
	}

	if len(m.AccessToken) < 1 {
		err = errors.New("access_token is required")
	}

	if len(m.Avatar) < 1 {
		err = errors.New("avatar is required")
	}

	if len(m.Status) < 1 {
		err = errors.New("status is required")
	}

	return
}
