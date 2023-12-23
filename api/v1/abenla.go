package apiv1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Abenla struct{}

var AbenlaHandler *Abenla

func NewAbenlaWebhook(r *gin.Engine) {
	handler := &Abenla{}

	Group := r.Group("/bss/v1/abenla")
	{
		Group.POST("/webhook/:routing_config", handler.WebhookAbenla)
	}
}

func (handler *Abenla) WebhookAbenla(ctx *gin.Context) {
	routingConfigUuid := ctx.Param("routing_config")
	// metadata, _ := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	code, result := response.ResponseXml("Status", "1")
	// 	ctx.XML(code, result)
	// 	return
	// }

	// authUser := model.AuthUser{
	// 	TenantId:           metadata.Get("tenant_id")[0],
	// 	BusinessUnitId:     metadata.Get("business_unit_id")[0],
	// 	UserId:             metadata.Get("user_id")[0],
	// 	Username:           metadata.Get("username")[0],
	// 	Services:           metadata.Get("services"),
	// 	DatabaseName:       metadata.Get("database_name")[0],
	// 	DatabaseHost:       metadata.Get("database_host")[0],
	// 	DatabasePort:       util.ParseInt(metadata.Get("database_port")[0]),
	// 	DatabaseUser:       metadata.Get("database_user")[0],
	// 	DatabasePassword:   metadata.Get("database_password")[0],
	// 	DatabaseEsHost:     metadata.Get("database_es_host")[0],
	// 	DatabaseEsUser:     metadata.Get("database_es_user")[0],
	// 	DatabaseEsPassword: metadata.Get("database_es_password")[0],
	// 	DatabaseEsIndex:    metadata.Get("database_es_index")[0],
	// }
	smsGuid := ctx.Query("SMSGUID")
	smsStatusQuery := ctx.Query("SMSSTATUS")
	smsStatus, _ := strconv.Atoi(smsStatusQuery)
	secretSign := ctx.Query("SECRECTSIGN")

	data := model.WebhookReceiveSmsStatus{
		SmsGuid:     smsGuid,
		Status:      smsStatus,
		SercretSign: secretSign,
	}

	code, result := service.NewAbenla().AbenlaWebhook(ctx, routingConfigUuid, data)
	ctx.XML(code, result)
}
