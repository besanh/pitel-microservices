package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"golang.org/x/exp/slices"
)

type OttMessage struct {
	ottMessageService    service.IOttMessage
	connectionAppService service.IChatConnectionApp
}

func NewOttMessage(r *gin.Engine, messageService service.IOttMessage, connectionApp service.IChatConnectionApp) {
	handler := &OttMessage{
		ottMessageService:    messageService,
		connectionAppService: connectionApp,
	}

	Group := r.Group("bss-message/v1/ott")
	{
		Group.POST("", handler.GetOttMessage)
	}
}

func (h *OttMessage) GetOttMessage(c *gin.Context) {
	jsonBody := make(map[string]any, 0)
	if err := c.ShouldBindJSON(&jsonBody); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}
	log.Info(jsonBody)
	messageType, _ := jsonBody["type"].(string)
	eventName, _ := jsonBody["event_name"].(string)
	appId, _ := jsonBody["app_id"].(string)
	appName, _ := jsonBody["app_name"].(string)
	oaId, _ := jsonBody["oa_id"].(string)
	userIdByApp, _ := jsonBody["user_id_by_app"].(string)
	externalUserId, _ := jsonBody["uid"].(string)
	username, _ := jsonBody["username"].(string)
	avatar, _ := jsonBody["avatar"].(string)
	timestampTmp, _ := jsonBody["timestamp"].(string)
	timestamp, _ := strconv.ParseInt(timestampTmp, 10, 64)
	msgId, _ := jsonBody["msg_id"].(string)
	content, _ := jsonBody["text"].(string)
	connectionId, _ := jsonBody["connection_id"].(string)
	attachmentsTmp, _ := jsonBody["attachments"].([]any)
	attachmentsAny := make([]any, 0)
	for item := range attachmentsTmp {
		tmp := attachmentsTmp[item].(map[string]any)
		attType, _ := tmp["att_type"].(string)
		if slices.Contains[[]string](variables.EVENT_NAME_SEND_MESSAGE, attType) {
			attachment := map[string]any{
				"att_type": attType,
				"payload":  tmp["payload"],
			}
			attachmentsAny = append(attachmentsAny, attachment)
		}
	}
	attachments := make([]model.OttAttachments, 0)
	if err := util.ParseAnyToAny(attachmentsAny, &attachments); err != nil {
		c.JSON(response.BadRequestMsg(err))
		return
	}

	var message model.OttMessage
	if eventName == "" {
		connectionAppRequest := model.ChatConnectionAppRequest{
			OaId: oaId,
		}
		if err := h.connectionAppService.UpdateChatConnectionAppById(c, nil, connectionId, connectionAppRequest); err != nil {
			c.JSON(response.BadRequestMsg(err))
			return
		}
		c.JSON(response.OKResponse())
	} else {
		message = model.OttMessage{
			MessageType:    messageType,
			EventName:      eventName,
			AppId:          appId,
			AppName:        appName,
			OaId:           oaId,
			UserIdByApp:    userIdByApp,
			ExternalUserId: externalUserId,
			Username:       username,
			Avatar:         avatar,
			Timestamp:      timestamp,
			MsgId:          msgId,
			Content:        content,
			Attachments:    &attachments,
		}
		code, result := h.ottMessageService.GetOttMessage(c, message)
		c.JSON(code, result)
	}
}
