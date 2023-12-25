package apiv1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Abenla struct {
	abenla service.IAbenla
}

var AbenlaHandler *Abenla

func NewAbenlaWebhook(r *gin.Engine, abenlaService service.IAbenla) {
	handler := &Abenla{
		abenla: abenlaService,
	}

	Group := r.Group("/bss/v1/abenla")
	{
		Group.POST("/webhook/:routing_config", handler.WebhookAbenla)
	}
}

func (handler *Abenla) WebhookAbenla(ctx *gin.Context) {
	routingConfigUuid := ctx.Param("routing_config")
	smsGuid := ctx.Query("SMSGUID")
	smsStatusQuery := ctx.Query("SMSSTATUS")
	smsStatus, _ := strconv.Atoi(smsStatusQuery)
	secretSign := ctx.Query("SECRECTSIGN")

	data := model.WebhookReceiveSmsStatus{
		SmsGuid:     smsGuid,
		Status:      smsStatus,
		SercretSign: secretSign,
	}

	code, result := handler.abenla.AbenlaWebhook(ctx, routingConfigUuid, data)
	ctx.XML(code, result)
}
