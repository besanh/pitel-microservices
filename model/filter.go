package model

import "database/sql"

type BalanceConfigFilter struct {
	Weight      []string     `json:"weight"`
	BalanceType []string     `json:"balance_type"`
	Provider    []string     `json:"provider"`
	Priority    []string     `json:"priority"`
	Status      sql.NullBool `json:"status"`
}

type PluginConfigFilter struct {
	PluginName []string     `json:"plugin_name"`
	PluginType []string     `json:"plugin_type"`
	Status     sql.NullBool `json:"status"`
}

type RecipientConfigFilter struct {
	Recipient     []string     `json:"recipient"`
	RecipientType []string     `json:"recipient_type"`
	Priority      []string     `json:"priority"`
	Provider      []string     `json:"provider"`
	Status        sql.NullBool `json:"status"`
}

type TemplateBssFilter struct {
	TemplateName string       `json:"template_name"`
	TemplateCode []string     `json:"template_code"`
	TemplateType []string     `json:"template_type"`
	Content      string       `json:"content"`
	Status       sql.NullBool `json:"status"`
}
