package model

import (
	"errors"
	"net/mail"

	"github.com/uptrace/bun"
	"golang.org/x/exp/slices"
)

type ChatEmail struct {
	*Base
	bun.BaseModel    `bun:"table:chat_email,alias:ce"`
	TenantId         string   `json:"tenant_id" bun:"tenant_id,type:uuid,notnull"`
	OaId             string   `json:"oa_id" bun:"oa_id,type:text,notnull"`
	EmailSubject     string   `json:"email_subject" bun:"email_subject,type:text,notnull"`
	EmailRecipient   []string `json:"email_recipient" bun:"email_recipient,type:text,notnull"`
	EmailContent     string   `json:"email_content" bun:"email_content,type:text,notnull"`
	EmailServer      string   `json:"email_server" bun:"email_server,type:text,notnull"`
	EmailUsername    string   `json:"email_username" bun:"email_username,type:text,notnull"`
	EmailPassword    string   `json:"email_password" bun:"email_password,type:text,notnull"`
	EmailPort        string   `json:"email_port" bun:"email_port,type:text,notnull"`
	EmailEncryptType string   `json:"email_encrypt_type" bun:"email_encrypt_type,type:text,notnull"`
	EmailStatus      bool     `json:"email_status" bun:"email_status,notnull"`
}

type ChatEmailCustom struct {
	*Base
	bun.BaseModel    `bun:"table:chat_email,alias:ce"`
	TenantId         string   `json:"tenant_id" bun:"tenant_id,type:uuid"`
	OaId             string   `json:"oa_id" bun:"oa_id,type:text"`
	EmailSubject     string   `json:"email_subject" bun:"email_subject"`
	EmailRecipient   []string `json:"email_recipient" bun:"email_recipient"`
	EmailContent     string   `json:"email_content" bun:"email_content"`
	EmailServer      string   `json:"email_server" bun:"email_server"`
	EmailUsername    string   `json:"email_username" bun:"email_username"`
	EmailPassword    string   `json:"email_password" bun:"email_password"`
	EmailPort        string   `json:"email_port" bun:"email_port"`
	EmailEncryptType string   `json:"email_encrypt_type" bun:"email_encrypt_type"`
	EmailStatus      bool     `json:"email_status" bun:"email_status"`

	// From connection
	ConnectionName string `json:"connection_name" bun:"connection_name"`
	ConnectionType string `json:"connection_type" bun:"connection_type"`
	OaInfo         OaInfo `json:"oa_info" bun:"oa_info"`
}

type ChatEmailRequest struct {
	OaId             string   `json:"oa_id" bun:"oa_id"`
	EmailSubject     string   `json:"email_subject" bun:"email_subject"`
	EmailRecipient   []string `json:"email_recipient" bun:"email_recipient"`
	EmailContent     string   `json:"email_content" bun:"email_content"`
	EmailServer      string   `json:"email_server" bun:"email_server"`
	EmailUsername    string   `json:"email_username" bun:"email_username"`
	EmailPassword    string   `json:"email_password" bun:"email_password"`
	EmailPort        string   `json:"email_port" bun:"email_port"`
	EmailEncryptType string   `json:"email_encrypt_type" bun:"email_encrypt_type"`
	EmailStatus      bool     `json:"email_status" bun:"email_status"`
	EmailRequestType string   `json:"email_request_type" bun:"email_request_type"`
}

func (m *ChatEmailRequest) Validate() error {
	if len(m.OaId) < 1 {
		return errors.New("oa id is required")
	}

	if len(m.EmailSubject) < 1 {
		return errors.New("email subject is required")
	}

	if len(m.EmailRecipient) < 1 {
		return errors.New("email recipient is required")
	}

	// pattern := `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	for _, email := range m.EmailRecipient {
		if len(email) < 1 {
			return errors.New("email recipient is required")
		}

		_, err := mail.ParseAddress(email)
		if err != nil {
			return errors.New("invalid email recipient: " + email)
		}
	}

	if len(m.EmailContent) < 1 {
		return errors.New("email content is required")
	}

	if !slices.Contains([]string{"manual", "auto"}, m.EmailRequestType) {
		return errors.New("email status is required")
	}

	if m.EmailRequestType == "manual" {
		if len(m.EmailServer) < 1 {
			return errors.New("email server is required")
		}

		if len(m.EmailUsername) < 1 {
			return errors.New("email username is required")
		}

		if len(m.EmailPassword) < 1 {
			return errors.New("email password is required")
		}

		if len(m.EmailPort) < 1 {
			return errors.New("email port is required")
		}

		if !slices.Contains([]string{"ssl", "tls"}, m.EmailEncryptType) {
			return errors.New("email encrypt type " + m.EmailEncryptType + " is not supported")
		}
	}

	return nil
}
