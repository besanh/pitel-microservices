package v1

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type Message struct {
	messageService service.IMessage
}

func NewMessage(r *gin.Engine, messageService service.IMessage) {
	handler := &Message{
		messageService: messageService,
	}

	Group := r.Group("bss-message/v1/message")
	{
		Group.POST("send", handler.SendMessage)
		Group.GET("", handler.GetMessages)
		Group.POST("read", handler.MarkReadMessages)
		Group.POST("share-info", handler.ShareInfo)
	}
}

func (h *Message) SendMessage(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	message := model.MessageRequest{}
	var file *multipart.FileHeader

	if c.Query("event_name") == "form" {
		messageForm := model.MessageFormRequest{}
		if err := c.ShouldBind(&messageForm); err != nil {
			log.Error(err)
			c.JSON(response.BadRequestMsg(err))
			return
		}

		log.Info("send message body form: ", messageForm)
		message.EventName = messageForm.EventName
		message.AppId = messageForm.AppId
		message.ConversationId = messageForm.ConversationId
	} else {
		jsonBody := make(map[string]any, 0)
		if err := c.ShouldBind(&jsonBody); err != nil {
			log.Error(err)
			c.JSON(response.BadRequestMsg(err))
			return
		}
		log.Info("send message body: ", jsonBody)

		appId, _ := jsonBody["app_id"].(string)
		conversationId, _ := jsonBody["conversation_id"].(string)
		content, _ := jsonBody["content"].(string)
		message = model.MessageRequest{
			EventName:      "text",
			AppId:          appId,
			ConversationId: conversationId,
			Content:        content,
		}
	}

	code, result := h.messageService.SendMessageToOTT(c, res.Data, message, file)
	c.JSON(code, result)
}

func (h *Message) GetMessages(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	filter := model.MessageFilter{
		ConversationId: c.Query("conversation_id"),
		ExternalUserId: c.Query("external_user_id"),
	}

	code, result := h.messageService.GetMessages(c, res.Data, filter, limit, offset)
	c.JSON(code, result)
}

func (h *Message) MarkReadMessages(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	markReadMessages := model.MessageMarkRead{}
	if err := c.ShouldBindJSON(&markReadMessages); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("mark read message body: ", markReadMessages)

	if err := markReadMessages.ValidateMarkRead(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	code, result := h.messageService.MarkReadMessages(c, res.Data, markReadMessages)
	c.JSON(code, result)
}

func (h *Message) ShareInfo(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	shareInfo := model.ShareInfo{}
	if err := c.ShouldBindJSON(&shareInfo); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("share info body: ", shareInfo)

	code, result := h.messageService.ShareInfo(c, res.Data, shareInfo)
	c.JSON(code, result)
}
