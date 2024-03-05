package v1

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type ChatRouting struct {
	chatRoutingService service.IChatRouting
}

func NewChatRouting(engine *gin.Engine, chatRoutingService service.IChatRouting) {
	handler := &ChatRouting{
		chatRoutingService: service.NewChatRouting(),
	}
	Group := engine.Group("bss-message/v1/chat-routing")
	{
		Group.POST("", handler.InsertChatRouting)
		Group.GET("", handler.GetChatRoutings)
		Group.GET(":id", handler.GetChatRoutingById)
		Group.PUT(":id", handler.UpdateChatRoutingById)
		Group.DELETE(":id", handler.DeleteChatRoutingById)
	}
}

func (handler *ChatRouting) InsertChatRouting(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ChatRoutingRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("insert chat routing payload -> ", &data)

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	id, err := handler.chatRoutingService.InsertChatRouting(c, res.Data, &data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OK(map[string]any{
		"id": id,
	}))
}

func (handler *ChatRouting) GetChatRoutings(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	limit := util.ParseLimit(c.Query("limit"))
	offset := util.ParseOffset(c.Query("offset"))

	var status sql.NullBool
	if len(c.Query("status")) > 0 {
		statusTmp, _ := strconv.ParseBool(c.Query("status"))
		status.Valid = true
		status.Bool = statusTmp
	}

	filter := model.ChatRoutingFilter{
		RoutingName:  c.Query("routing_name"),
		RoutingAlias: c.Query("routing_alias"),
		Status:       status,
	}

	total, chatRoutings, err := handler.chatRoutingService.GetChatRoutings(c, res.Data, filter, limit, offset)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.Pagination(chatRoutings, total, limit, offset))
}

func (handler *ChatRouting) GetChatRoutingById(c *gin.Context) {
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

	chatRouting, err := handler.chatRoutingService.GetChatRoutingById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"id": chatRouting.Id,
	}))
}

func (handler *ChatRouting) UpdateChatRoutingById(c *gin.Context) {
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

	var data model.ChatRoutingRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat routing payload -> ", &data)

	err := handler.chatRoutingService.UpdateChatRoutingById(c, res.Data, id, data)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (handler *ChatRouting) DeleteChatRoutingById(c *gin.Context) {
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

	err := handler.chatRoutingService.DeleteChatRoutingById(c, res.Data, id)
	if err != nil {
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	c.JSON(response.OKResponse())
}
