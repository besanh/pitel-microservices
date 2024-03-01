package model

import (
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type ChatConnectionApp struct {
	bun.BaseModel  `bun:"table:chat_connection_app,alias:cca"`
	Id             string    `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	TenantId       string    `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull"`
	UpdatedAt      time.Time `bun:"updated_at,notnull"`
	ConnectionName string    `json:"connection_name" bun:"connection_name,type:text,notnull"`
	ConnectionType string    `json:"connection_type" bun:"connection_type,type:text,notnull"`
	AppId          string    `json:"app_id" bun:"app_id,type:text,notnull"`
	QueueId        string    `json:"queue_id" bun:"queue_id,type:text,notnull"`
	OaInfo         OaInfo    `json:"oa_info" bun:"oa_info,type:jsonb,notnull"`
	Status         string    `json:"status" bun:"status,notnull"`
}

type OaInfo struct {
	Zalo []struct {
		AppId    string `json:"app_id"`
		OaId     string `json:"oa_id"`
		OaName   string `json:"oa_name"`
		UrlOa    string `json:"url_oa"`
		Avatar   string `json:"avatar"`
		Cover    string `json:"cover"`
		CateName string `json:"cate_name"`
		Status   string `json:"status"`
	} `json:"zalo"`
	Facebook []struct {
		AppId       string `json:"app_id"`
		OaId        string `json:"oa_id"`
		OaName      string `json:"oa_name"`
		UrlOa       string `json:"url_oa"`
		Avatar      string `json:"avatar"`
		Cover       string `json:"cover"`
		AccessToken string `json:"access_token"`
		Status      string `json:"status"`
	} `json:"facebook"`
}

type ChatConnectionAppRequest struct {
	Id             string  `json:"id"`
	ConnectionName string  `json:"connection_name"`
	ConnectionType string  `json:"connection_type"`
	QueueId        string  `json:"queue_id"`
	OaInfo         *OaInfo `json:"oa_info"`
	Status         string  `json:"status"`
	OaId           string  `json:"oa_id"`
	AppId          string  `json:"app_id"` // use for receive from message when user grant permission to ott

	// Recieve from ott for update connection zalo
	OaName   string `json:"oa_name"`
	Avatar   string `json:"avatar"`
	Cover    string `json:"cover"`
	CateName string `json:"cate_name"`
	Code     int64  `json:"code"`
	Message  string `json:"message"`
}

type AccessInfo struct {
	CallbackUrl   string `json:"callback_url"`
	ChallangeCode string `json:"challange_code"`
	State         string `json:"state"`
}

func (m *ChatConnectionAppRequest) Validate() error {
	if len(m.ConnectionName) < 1 {
		return errors.New("connection name is required")
	}
	if len(m.ConnectionType) < 1 {
		return errors.New("connection type is required")
	}
	if !slices.Contains[[]string](variables.CONNECTION_TYPE, m.ConnectionType) {
		return errors.New("connection type " + m.ConnectionType + " is not supported")
	}
	if len(m.Status) < 1 {
		return errors.New("status is required")
	}
	return nil
}
