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

type ChatConnectionApp struct {
	chatConnectionAppService service.IChatConnectionApp
}

func NewChatConnectionApp(engine *gin.Engine, chatConnectionAppService service.IChatConnectionApp) {
	handler := &ChatConnectionApp{
		chatConnectionAppService: chatConnectionAppService,
	}

	Group := engine.Group("bss-message/v1/chat-connection-app")
	{
		Group.GET("", handler.GetChatConnectionApp)
		Group.POST("", handler.InsertChatConnectionApp)
		Group.PUT(":id", handler.UpdateChatConnectionAppById)
		Group.GET(":id", handler.GetChatConnectionAppById)
		Group.DELETE(":id", handler.DeleteChatConnectionAppById)
	}
}

func (handler *ChatConnectionApp) GetChatConnectionApp(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}
	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	filter := model.ChatConnectionAppFilter{
		ConnectionName: c.Query("connection_name"),
		ConnectionType: c.Query("connection_type"),
		QueueId:        c.Query("queue_id"),
		Status:         c.Query("status"),
	}

	total, result, err := handler.chatConnectionAppService.GetChatConnectionApp(c, res.Data, filter, limit, offset)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.Pagination(result, total, limit, offset))
}

func (handler *ChatConnectionApp) InsertChatConnectionApp(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ChatConnectionAppRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("insert chat connection app payload -> ", &data)

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatConnectionAppService.InsertChatConnectionApp(c, res.Data, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatConnectionApp) GetChatConnectionAppById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	chatConnectionApp, err := handler.chatConnectionAppService.GetChatConnectionAppById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"data": chatConnectionApp,
	}))
}

func (handler *ChatConnectionApp) UpdateChatConnectionAppById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	var data model.ChatConnectionAppRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat connection app payload -> ", data)

	err := handler.chatConnectionAppService.UpdateChatConnectionAppById(c, res.Data, id, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}

func (handler *ChatConnectionApp) DeleteChatConnectionAppById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	err := handler.chatConnectionAppService.DeleteChatConnectionAppById(c, res.Data, id)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
