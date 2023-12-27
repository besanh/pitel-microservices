package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/model"
)

type (
	IWebhook interface {
		AbenlaWebhook(ctx context.Context, routingConfigUuid string, data model.WebhookReceiveSmsStatus) (int, any)
		FptWebhook(ctx context.Context, routingConfigUuid string, data model.FptWebhook) (int, any)
		IncomWebhook(ctx context.Context, routingConfigUuid string, data model.WebhookIncom) (int, any)
	}
	Webhook struct {
		Index string
	}
)

func NewWebhook(index string) IWebhook {
	return &Webhook{
		Index: index,
	}
}
