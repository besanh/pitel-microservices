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

type (
	IAPIMessage interface {
		HandlePostSendMessage(c *gin.Context)
	}

	Message struct {
		messageService service.IMessage
	}
)

var APIMessage IAPIMessage

func NewAPIMessage() IAPIMessage {
	return &Message{
		messageService: service.NewMessage(),
	}
}

func NewMessage(engine *gin.Engine, messageService service.IMessage) {
	handler := &Message{
		messageService: messageService,
	}

	Group := engine.Group("bss-message/v1/message")
	{
		Group.POST("send", handler.SendMessage)
		Group.GET("", handler.GetMessages)
		Group.POST("read", handler.MarkReadMessages)
		Group.POST("share-info", handler.ShareInfo)
		Group.GET("scroll", handler.GetMessagesWithScrollAPI)
	}
}

func (h *Message) SendMessage(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
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

		if err := messageForm.ValidateMessageForm(); err != nil {
			log.Error(err)
			c.JSON(response.BadRequestMsg(err))
			return
		}

		message.EventName = messageForm.EventName
		message.AppId = messageForm.AppId
		message.OaId = messageForm.OaId
		message.ConversationId = messageForm.ConversationId
		message.Url = messageForm.Url
		file = messageForm.File
	} else {
		jsonBody := make(map[string]any, 0)
		if err := c.ShouldBind(&jsonBody); err != nil {
			log.Error(err)
			c.JSON(response.BadRequestMsg(err))
			return
		}
		log.Info("send message body: ", jsonBody)

		appId, _ := jsonBody["app_id"].(string)
		oaId, _ := jsonBody["oa_id"].(string)
		conversationId, _ := jsonBody["conversation_id"].(string)
		content, _ := jsonBody["content"].(string)
		message = model.MessageRequest{
			EventName:      "text",
			AppId:          appId,
			OaId:           oaId,
			ConversationId: conversationId,
			Content:        content,
		}
	}

	data, err := h.messageService.SendMessageToOTT(c, res.Data, message, file)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Created(data))
}

func (h *Message) GetMessages(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	filter := model.MessageFilter{
		ConversationId: c.Query("conversation_id"),
		ExternalUserId: c.Query("external_user_id"),
	}

	total, data, err := h.messageService.GetMessages(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.Pagination(data, total, limit, offset))
}

func (h *Message) MarkReadMessages(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
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

	data, err := h.messageService.MarkReadMessages(c, res.Data, markReadMessages)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(data))
}

func (h *Message) ShareInfo(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	shareInfo := model.ShareInfo{}
	if err := c.ShouldBindJSON(&shareInfo); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("share info body: ", shareInfo)

	data, err := h.messageService.ShareInfo(c, res.Data, shareInfo)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(data))
}

func (h *Message) GetMessagesWithScrollAPI(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.Unauthorized())
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	scrollId := c.Query("scroll_id")

	filter := model.MessageFilter{
		ConversationId: c.Query("conversation_id"),
		ExternalUserId: c.Query("external_user_id"),
	}

	total, data, respScrollId, err := h.messageService.GetMessagesWithScrollAPI(c, res.Data, filter, limit, scrollId)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	result := map[string]any{
		"messages":  data,
		"scroll_id": respScrollId,
	}
	c.JSON(response.Pagination(result, total, limit, 0))
}

func (h *Message) HandlePostSendMessage(c *gin.Context) {
	res := api.AuthMiddlewareNewVersion(c)
	if res == nil {
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

		if err := messageForm.ValidateMessageForm(); err != nil {
			log.Error(err)
			c.JSON(response.BadRequestMsg(err))
			return
		}

		message.EventName = messageForm.EventName
		message.AppId = messageForm.AppId
		message.OaId = messageForm.OaId
		message.ConversationId = messageForm.ConversationId
		message.Url = messageForm.Url
		file = messageForm.File
	} else {
		jsonBody := make(map[string]any, 0)
		if err := c.ShouldBindJSON(&jsonBody); err != nil {
			log.Error(err)
			c.JSON(response.BadRequestMsg(err))
			return
		}
		log.Info("send message body: ", jsonBody)

		appId, _ := jsonBody["app_id"].(string)
		oaId, _ := jsonBody["oa_id"].(string)
		conversationId, _ := jsonBody["conversation_id"].(string)
		content, _ := jsonBody["content"].(string)
		message = model.MessageRequest{
			EventName:      "text",
			AppId:          appId,
			OaId:           oaId,
			ConversationId: conversationId,
			Content:        content,
		}
	}

	data, err := h.messageService.SendMessageToOTT(c, res.Data, message, file)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Created(data))
}
