package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/model"
)

type (
	IWebhook interface {
		AbenlaWebhook(ctx context.Context, data model.WebhookReceiveSmsStatus) (int, any)
		FptWebhook(ctx context.Context, data model.FptWebhook) (int, any)
		IncomWebhook(ctx context.Context, data model.WebhookIncom) (int, any)
	}
	Webhook struct {
	}
)

func NewWebhook() IWebhook {
	return &Webhook{}
}
