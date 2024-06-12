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

type ChatCommand struct {
	chatCommandService service.IChatCommand
}

func NewChatCommand(engine *gin.Engine, chatCommandService service.IChatCommand) {
	handler := ChatCommand{
		chatCommandService: chatCommandService,
	}

	group := engine.Group("bss-message/v1/chat-command")
	{
		group.GET("", handler.GetChatCommands)
		group.GET(":id", handler.GetChatCommandById)
		group.POST("", handler.InsertChatCommand)
		group.PUT(":id", handler.UpdateChatCommand)
		group.DELETE(":id", handler.DeleteChatCommandById)
	}
}

func (handler *ChatCommand) GetChatCommands(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	total, result, err := handler.chatCommandService.GetChatCommands(c, res.Data, limit, offset)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatCommand) GetChatCommandById(c *gin.Context) {
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

	chatCommand, err := handler.chatCommandService.GetChatCommandById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(chatCommand))
}

func (handler *ChatCommand) InsertChatCommand(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("invalid token"))
		return
	}

	var chatCmd model.ChatCommandRequest
	err := c.ShouldBindJSON(&chatCmd)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatCmd.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatCommandService.InsertChatCommand(c, res.Data, chatCmd)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatCommand) UpdateChatCommand(c *gin.Context) {
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

	var chatCmd model.ChatCommandRequest
	err := c.ShouldBindJSON(&chatCmd)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	if err := chatCmd.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err = handler.chatCommandService.UpdateChatCommandById(c, res.Data, id, chatCmd)
	if err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (handler *ChatCommand) DeleteChatCommandById(c *gin.Context) {
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

	err := handler.chatCommandService.DeleteChatCommandById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
