package apiv1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Incom struct {
	incom service.IIncom
}

var IncomHandler *Incom

func NewIncomWebhook(r *gin.Engine, incomService service.IIncom) {
	handler := &Incom{
		incom: incomService,
	}
	Group := r.Group("bss/v1/incom")
	{
		Group.POST("webhook/:routing_config", handler.WebhookData)
	}
}

func (handler *Incom) WebhookData(ctx *gin.Context) {
	routingConfigUuid := ctx.Param("routing_config")

	jsonBody := make(map[string]any)
	if err := ctx.ShouldBindJSON(&jsonBody); err != nil {
		ctx.JSON(response.BadRequestMsg(err))
		return
	}

	log.Info("webhook incom", jsonBody)

	idOmniMess, _ := jsonBody["IdOmniMess"].(string)
	status, _ := jsonBody["Status"].(string)
	channel, _ := jsonBody["Channel"].(string)
	errorCode, _ := jsonBody["ErrorCode"].(string)
	quantityTmp, _ := jsonBody["MtCount"].(float64)
	quantity := int(quantityTmp)
	telcoIdTmp, _ := jsonBody["TelcoId"].(float64)
	telcoId := int(telcoIdTmp)
	isChargedZnsTmp, _ := jsonBody["IsCharged"].(string)
	isChargedZns, _ := strconv.ParseBool(isChargedZnsTmp)

	incomData := model.WebhookIncom{
		IdOmniMess:   idOmniMess,
		Status:       status,
		Channel:      channel,
		ErrorCode:    errorCode,
		Quantity:     quantity,
		TelcoId:      telcoId,
		IsChargedZns: isChargedZns,
	}

	code, result := handler.incom.IncomWebhook(ctx, routingConfigUuid, incomData)

	ctx.JSON(code, result)
}
