package v1

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type (
	IAPIChatMessageSample interface {
		InsertChatMsgSample(c *gin.Context)
		HandlePutChatMessageSampleUpload(c *gin.Context)
	}

	ChatMsgSample struct {
		chatMsgSampleService service.IChatMsgSample
	}
)

var APIChatMessageSampleHandler IAPIChatMessageSample

func NewChatMessageSample() IAPIChatMessageSample {
	return &ChatMsgSample{
		chatMsgSampleService: service.NewChatMsgSample(),
	}
}

func NewChatMsgSample(engine *gin.Engine, chatMsgSampleService service.IChatMsgSample) {
	handler := ChatMsgSample{
		chatMsgSampleService: chatMsgSampleService,
	}

	group := engine.Group("bss-message/v1/chat-sample")
	{
		group.GET("", handler.GetChatMsgSamples)
		group.GET(":id", handler.GetChatMsgSampleById)
		group.POST("", handler.InsertChatMsgSample)
		group.PUT(":id", handler.UpdateChatMsgSample)
		group.DELETE(":id", handler.DeleteChatMsgSampleById)
	}
}

func (handler *ChatMsgSample) GetChatMsgSamples(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	filter := model.ChatMsgSampleFilter{
		ConnectionId: c.Query("connection_id"),
		Channel:      c.Query("channel"),
		OaId:         c.Query("oa_id"),
		Keyword:      c.Query("keyword"),
	}

	total, result, err := handler.chatMsgSampleService.GetChatMsgSamples(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatMsgSample) GetChatMsgSampleById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	chatMsgSample, err := handler.chatMsgSampleService.GetChatMsgSampleById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(chatMsgSample))
}

func (handler *ChatMsgSample) InsertChatMsgSample(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	var chatMsgSampleRequest model.ChatMsgSampleRequest
	err := c.ShouldBind(&chatMsgSampleRequest)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatMsgSampleRequest.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatMsgSampleService.InsertChatMsgSample(c, res.Data, chatMsgSampleRequest, chatMsgSampleRequest.File)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatMsgSample) UpdateChatMsgSample(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	var chatMsgSampleRequest model.ChatMsgSampleRequest
	err := c.ShouldBind(&chatMsgSampleRequest)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatMsgSampleRequest.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err = handler.chatMsgSampleService.UpdateChatMsgSampleById(c, res.Data, id, chatMsgSampleRequest, chatMsgSampleRequest.File)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (handler *ChatMsgSample) DeleteChatMsgSampleById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	err := handler.chatMsgSampleService.DeleteChatMsgSampleById(c, res.Data, id)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *ChatMsgSample) HandlePutChatMessageSampleUpload(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		return
	}
	id := strings.TrimPrefix(c.Request.RequestURI, "/bss-chat/v1/chat-sample/upload/")
	if id == "" {
		c.JSON(response.BadRequestMsg("id is empty"))
		return
	}

	var chatMsgSampleRequest model.ChatMsgSampleRequest
	err := c.ShouldBind(&chatMsgSampleRequest)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatMsgSampleRequest.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err = handler.chatMsgSampleService.UpdateChatMsgSampleById(c, res.Data, id, chatMsgSampleRequest, chatMsgSampleRequest.File)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}
