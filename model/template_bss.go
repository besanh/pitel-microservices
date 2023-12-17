package model

import (
	"errors"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type TemplateBss struct {
	*Base
	bun.BaseModel `bun:"table:template_bss,alias:tbss"`
	TemplateName  string `json:"template_name" bun:"template_name,type:text,notnull"`
	TemplateCode  string `json:"template_code" bun:"template_code,type:text,notnull"`
	TemplateType  string `json:"template_type" bun:"template_type,type:text,notnull"`
	Content       string `json:"content" bun:"content,type:text,notnull"`
	Partition     string `json:"partition" bun:"partition,type:text"`
	Status        bool   `json:"status" bun:"status,type:boolean"`
}

type TemplateBssBodyRequest struct {
	TemplateName string `json:"template_name"`
	TemplateCode string `json:"template_code"`
	TemplateType string `json:"template_type"`
	Content      string `json:"content"`
	Status       bool   `json:"status"`
}

func (r *TemplateBssBodyRequest) Validate() error {
	if len(r.TemplateName) < 1 {
		return errors.New("template name is missing")
	}
	if len(r.TemplateCode) < 1 {
		return errors.New("template code is missing")
	}
	if len(r.TemplateType) < 1 {
		return errors.New("template type is missing")
	}
	if !slices.Contains[[]string](constants.CHANNEL, r.TemplateType) {
		return errors.New("template type is invalid")
	}
	if len(r.Content) < 1 {
		return errors.New("content is missing")
	}
	return nil
}
