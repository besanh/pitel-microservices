package apiv1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"golang.org/x/exp/slices"
)

type Webhook struct {
	webhook service.IWebhook
}

var WebhookHandler *Webhook

func NewWebhook(r *gin.Engine, webhookService service.IWebhook) {
	handler := &Webhook{
		webhook: webhookService,
	}
	Group := r.Group("bss/v1/webhook")
	{
		Group.POST(":plugin", api.ValidHeader(), handler.WebhookData)
	}
}

func (handler *Webhook) WebhookData(ctx *gin.Context) {
	if !slices.Contains[[]string](constants.EXTERNAL_PLUGIN_CONNECT_TYPE, ctx.Param("plugin")) {
		ctx.JSON(response.BadRequestMsg("plugin not support"))
		return
	}
	if ctx.Param("plugin") == "fpt" {
		fpt := model.FptWebhook{}
		if err := ctx.ShouldBindJSON(&fpt); err != nil {
			ctx.JSON(response.BadRequestMsg(err.Error()))
			return
		}
		log.Info("webhook fpt", fpt)

		code, result := handler.webhook.FptWebhook(ctx, fpt)

		ctx.JSON(code, result)
	} else if ctx.Param("plugin") == "incom" {
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

		code, result := handler.webhook.IncomWebhook(ctx, incomData)

		ctx.JSON(code, result)
	} else if ctx.Param("plugin") == "abenla" {
		smsGuid := ctx.Query("SMSGUID")
		smsStatusQuery := ctx.Query("SMSSTATUS")
		smsStatus, _ := strconv.Atoi(smsStatusQuery)
		secretSign := ctx.Query("SECRECTSIGN")

		data := model.WebhookReceiveSmsStatus{
			SmsGuid:     smsGuid,
			Status:      smsStatus,
			SercretSign: secretSign,
		}

		code, result := handler.webhook.AbenlaWebhook(ctx, data)
		ctx.XML(code, result)
	}
}
