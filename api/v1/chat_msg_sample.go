package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type ChatMsgSample struct {
	chatMsgSampleService service.IChatMsgSample
}

func NewChatMsgSample(engine *gin.Engine, chatMsgSampleService service.IChatMsgSample) {
	handler := ChatMsgSample{
		chatMsgSampleService: chatMsgSampleService,
	}

	group := engine.Group("bss-message/v1/chat-sample")
	{
		group.GET("", handler.GetChatMsgSamples)
		group.GET("/personalization", handler.GetPersonalizationValues)
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

	total, result, err := handler.chatMsgSampleService.GetChatMsgSamples(c, res.Data, limit, offset)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatMsgSample) GetPersonalizationValues(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	total, result, err := handler.chatMsgSampleService.GetChatPersonalizationValues(c, res.Data, limit, offset)
	if err != nil {
		log.Error(err)
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

	var chatCmd model.ChatMsgSampleRequest
	err := c.ShouldBindJSON(&chatCmd)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatCmd.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatMsgSampleService.InsertChatMsgSample(c, res.Data, chatCmd, chatCmd.File)
	if err != nil {
		log.Error(err)
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

	var chatCmd model.ChatMsgSampleRequest
	err := c.ShouldBindJSON(&chatCmd)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatCmd.Validate(); err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err = handler.chatMsgSampleService.UpdateChatMsgSampleById(c, res.Data, id, chatCmd, chatCmd.File)
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
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
