package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"golang.org/x/exp/slices"
)

type OttMessage struct {
	ottMessageService service.IOttMessage
}

func NewOttMessage(r *gin.Engine, messageService service.IOttMessage) {
	handler := &OttMessage{
		ottMessageService: messageService,
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
	userId, _ := jsonBody["uid"].(string)
	username, _ := jsonBody["username"].(string)
	avatar, _ := jsonBody["avatar"].(string)
	timestampTmp, _ := jsonBody["timestamp"].(string)
	timestamp, _ := strconv.ParseInt(timestampTmp, 10, 64)
	msgId, _ := jsonBody["msg_id"].(string)
	content, _ := jsonBody["text"].(string)
	attachmentsTmp, _ := jsonBody["attachments"].([]any)
	attachments := make([]model.OttAttachments, 0)
	for item := range attachmentsTmp {
		tmp := attachmentsTmp[item].(map[string]any)
		attType, _ := tmp["att_type"].(string)
		if slices.Contains[[]string](variables.EVENT_NAME_SEND_MESSAGE, attType) {
			attachments = append(attachments, model.OttAttachments{
				AttType: attType,
				Payload: tmp,
			})
		}
	}

	message := model.OttMessage{
		MessageType: messageType,
		EventName:   eventName,
		AppId:       appId,
		AppName:     appName,
		OaId:        oaId,
		UserIdByApp: userIdByApp,
		UserId:      userId,
		Username:    username,
		Avatar:      avatar,
		Timestamp:   timestamp,
		MsgId:       msgId,
		Content:     content,
		Attachments: &attachments,
	}
	code, result := h.ottMessageService.GetOttMessage(c, message)
	c.JSON(code, result)
}
