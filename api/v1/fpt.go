package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Fpt struct {
	fpt service.IFpt
}

var FptHandler *Fpt

func NewFptWebhook(r *gin.Engine, fptService service.IFpt) {
	handler := &Fpt{
		fpt: fptService,
	}
	Group := r.Group("bss/v1/fpt")
	{
		Group.POST("webhook", handler.WebhookData)
	}
}

func (handler *Fpt) WebhookData(ctx *gin.Context) {
	routingConfigUuid := ctx.GetHeader("Authorization")
	log.Info(routingConfigUuid)

	fpt := model.FptWebhook{}
	if err := ctx.ShouldBindJSON(&fpt); err != nil {
		ctx.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	code, result := handler.fpt.FptWebhook(ctx, routingConfigUuid, fpt)

	ctx.JSON(code, result)
}
