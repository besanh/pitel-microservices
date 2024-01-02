package model

import (
	"errors"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type RecipientConfig struct {
	*Base
	bun.BaseModel `bun:"table:recipient_config,alias:rc"`
	Recipient     string `json:"recipient" bun:"recipient,type:text,notnull"`
	RecipientType string `json:"recipient_type" bun:"recipient_type,type:text,notnull"`
	Priority      string `json:"priority" bun:"priority,type:text"`
	Provider      string `json:"provider" bun:"provider,type:text,notnull"`
	Status        bool   `json:"status" bun:"status,type:boolean"`
	CreatedBy     string `json:"created_by" bun:"created_by,type:text"`
}

type RecipientConfigView struct {
	bun.BaseModel `bun:"table:recipient_config,alias:rc"`
	Recipient     string `json:"recipient" bun:"recipient,type:text,notnull"`
	RecipientType string `json:"recipient_type" bun:"recipient_type,type:text,notnull"`
	Priority      string `json:"priority" bun:"priority,type:text"`
	Provider      string `json:"provider" bun:"provider,type:text,notnull"`
	Status        bool   `json:"status" bun:"status,type:boolean"`
	CreatedBy     string `json:"created_by" bun:"created_by,type:text"`
}

type RecipientConfigRequest struct {
	Recipient     string `json:"recipient"`
	RecipientType string `json:"recipient_type"`
	Priority      string `json:"priority"`
	Provider      string `json:"provider"`
	Status        bool   `json:"status"`
}

type RecipientConfigPutRequest struct {
	Recipient     string `json:"recipient"`
	RecipientType string `json:"recipient_type"`
	Priority      string `json:"priority"`
	Provider      string `json:"provider"`
	Status        bool   `json:"status"`
}

func (r *RecipientConfigRequest) Validate() error {
	if len(r.Recipient) < 1 {
		return errors.New("recipient is missing")
	}
	if len(r.Recipient) > 0 {
		if !slices.Contains[[]string](constants.RECIPIENT, r.Recipient) {
			return errors.New("recipient " + r.Recipient + " not support")
		}
	}
	if len(r.Provider) < 1 {
		return errors.New("provider is missing")
	}

	if len(r.RecipientType) < 1 {
		return errors.New("recipient type is missing")
	}

	if len(r.Priority) < 1 {
		return errors.New("priority is missing")
	}

	return nil
}
