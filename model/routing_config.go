package model

import (
	"errors"

	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type RoutingConfig struct {
	*Base
	bun.BaseModel `bun:"table:routing_config,alias:brc"`
	RoutingName   string        `json:"routing_name" bun:"routing_name,type:text,notnull"`
	RoutingType   string        `json:"routing_type" bun:"routing_type,type:text,notnull"` // sms,zns,...
	RoutingFlow   RoutingFlow   `json:"routing_flow" bun:"routing_flow,type:text,notnull"`
	RoutingOption RoutingOption `json:"routing_option" bun:"routing_option,type:text,notnull"`
	Status        bool          `json:"status" bun:"status,type:boolean"`
}

type RoutingConfigView struct {
	bun.BaseModel `bun:"table:routing_config,alias:brc"`
	RoutingName   string        `json:"routing_name" bun:"routing_name,type:text,notnull"`
	RoutingType   string        `json:"routing_type" bun:"routing_type,type:text,notnull"` // sms,zns,...
	RoutingFlow   RoutingFlow   `json:"routing_flow" bun:"routing_flow,type:text,notnull"`
	RoutingOption RoutingOption `json:"routing_option" bun:"routing_option,type:text,notnull"`
	Status        bool          `json:"status" bun:"status,type:boolean"`
}

// Link to table recipient_Routing or balance Routing to control flow send data
type RoutingFlow struct {
	FlowType string   `json:"flow_type"` // table recipient or balance
	FlowUuid []string `json:"flow_uuid"`
}

// Include account info connected with external plugin
type RoutingOption struct {
	Incom  Incom  `json:"incom" bun:"incom,type:text"`
	Abenla Abenla `json:"abenla" bun:"abenla,type:text"`
	Fpt    Fpt    `json:"fpt" bun:"fpt,type:text"`
}

type Incom struct {
	Username    string   `json:"username" bun:"username"`
	Password    string   `json:"password" bun:"password"`
	ApiUrl      string   `json:"api_url" bun:"api_url"`
	WebhookUrl  []string `json:"webhook_url" bun:"webhook_url,type:text"`
	MaxAttempts int      `json:"max_attempts" bun:"max_attempts,type:text"`
	Signature   string   `json:"signature" bun:"signature,type:text"`
	Status      bool     `json:"status" bun:"status"`
}

type Abenla struct {
	Username      string   `json:"username" bun:"username"`
	Password      string   `json:"password" bun:"password"`
	ApiUrl        string   `json:"api_url" bun:"api_url"`
	ServiceTypeId string   `json:"service_type_id" bun:"service_type_id"`
	WebhookUrl    []string `json:"webhook_url" bun:"webhook_url,type:text"`
	MaxAttempts   int      `json:"max_attempts" bun:"max_attempts,type:text"`
	Signature     string   `json:"signature" bun:"signature,type:text"`
	Brandname     string   `json:"brand_name" bun:"brand_name,type:text"`
	Status        bool     `json:"status" bun:"status"`
}

type Fpt struct {
	ClientId     string   `json:"client_id" bun:"client_id"`
	ClientSecret string   `json:"client_secret" bun:"client_secret"`
	ApiUrl       string   `json:"api_url" bun:"api_url"`
	WebhookUrl   []string `json:"webhook_url" bun:"webhook_url,type:text"`
	MaxAttempts  int      `json:"max_attempts" bun:"max_attempts,type:text"`
	Signature    string   `json:"signature" bun:"signature,type:text"`
	Status       bool     `json:"status" bun:"status"`
}

func (r *RoutingConfig) Validate() (err error) {
	if len(r.RoutingName) < 1 {
		return errors.New("routing name is required")
	}
	if len(r.RoutingType) < 1 {
		return errors.New("routing type is required")
	}

	if !slices.Contains[[]string]([]string{"recipient", "balance"}, r.RoutingFlow.FlowType) {
		return errors.New("routing flow is invalid")
	}
	if len(r.RoutingFlow.FlowUuid) < 1 {
		return errors.New("routing flow uuid is required")
	}

	// isCheckRoutingOptionEnable := 0

	// if r.RoutingOption.Incom.Status {
	// 	isCheckRoutingOptionEnable += 1
	// }
	// if r.RoutingOption.Abenla.Status {
	// 	isCheckRoutingOptionEnable += 1
	// }
	// if r.RoutingOption.Fpt.Status {
	// 	isCheckRoutingOptionEnable += 1
	// }

	// if isCheckRoutingOptionEnable > 1 {
	// 	return errors.New("only one routing option is enable")
	// } else if isCheckRoutingOptionEnable == 0 {
	// 	return errors.New("routing option is required")
	// }

	if r.RoutingOption.Incom.Status {
		if len(r.RoutingOption.Incom.Username) < 1 {
			return errors.New("incom username is required")
		}
		if len(r.RoutingOption.Incom.Password) < 1 {
			return errors.New("incom password is required")
		}
		if len(r.RoutingOption.Incom.ApiUrl) < 1 {
			return errors.New("incom api url is required")
		}
	}

	if r.RoutingOption.Abenla.Status {
		if len(r.RoutingOption.Abenla.Username) < 1 {
			return errors.New("abenla username is required")
		}
		if len(r.RoutingOption.Abenla.Password) < 1 {
			return errors.New("abenla password is required")
		}
		if len(r.RoutingOption.Abenla.ApiUrl) < 1 {
			return errors.New("abenla api url is required")
		}
		if len(r.RoutingOption.Abenla.ServiceTypeId) < 1 {
			return errors.New("abenla service type id is required")
		}
	}

	if r.RoutingOption.Fpt.Status {
		if len(r.RoutingOption.Fpt.ClientId) < 1 {
			return errors.New("fpt client id is required")
		}
		if len(r.RoutingOption.Fpt.ClientSecret) < 1 {
			return errors.New("fpt client secret is required")
		}
		if len(r.RoutingOption.Fpt.ApiUrl) < 1 {
			return errors.New("fpt api url is required")
		}
	}

	return
}
