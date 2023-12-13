package model

import (
	"errors"

	"github.com/uptrace/bun"
)

type BalanceConfig struct {
	*Base
	bun.BaseModel        `bun:"table:balance_config,alias:bc"`
	Weight               string                `json:"weight" bun:"weight,type:text"`
	BalanceType          string                `json:"balance_type" bun:"balance_type,type:text,notnull"`
	Priority             string                `json:"priority" bun:"priority,type:text,notnull"`
	Provider             string                `json:"provider" bun:"provider,type:text,notnull"`
	Status               bool                  `json:"status" bun:"status,type:boolean"`
	CreatedBy            string                `json:"created_by" bun:"created_by,type:text"`
	InboxMarketingConfig *InboxMarketingConfig `json:"inbox_marketing_config" bun:"rel:belongs-to,join:id=id"`
}

type BalanceConfigView struct {
	bun.BaseModel `bun:"table:balance_config,alias:bc"`
	Weight        string `json:"weight" bun:"weight,type:text"`
	BalanceType   string `json:"balance_type" bun:"balance_type,type:text,notnull"`
	Priority      string `json:"priority" bun:"priority,type:text"`
	Provider      string `json:"provider" bun:"provider,type:text"`
	Status        bool   `json:"status" bun:"status,type:boolean"`
	CreatedBy     string `json:"created_by" bun:"created_by,type:text"`
	// InboxMarketingConfig *InboxMarketingConfig `json:"inbox_marketing_config" bun:"rel:belongs-to,join:id=id"`
}

type BalanceConfigBodyRequest struct {
	Weight      string `json:"weight"`
	BalanceType string `json:"balance_type"`
	Priority    string `json:"priority"`
	Provider    string `json:"provider"`
	Status      bool   `json:"status"`
}

type BalanceConfigPutBodyRequest struct {
	Weight      string `json:"weight"`
	BalanceType string `json:"balance_type"`
	Priority    string `json:"priority"`
	Provider    string `json:"provider"`
	Status      bool   `json:"status"`
}

func (r *BalanceConfigBodyRequest) Validate() error {
	if len(r.Weight) < 1 {
		return errors.New("weight is missing")
	}

	if len(r.BalanceType) < 1 {
		return errors.New("balance type is missing")
	}

	if len(r.Provider) < 1 {
		return errors.New("provider is missing")
	}

	if len(r.Priority) < 1 {
		return errors.New("priority is missing")
	}

	return nil
}
