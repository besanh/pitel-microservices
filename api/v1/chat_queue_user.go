package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type ChatQueueUser struct {
	chatQueueUser service.IChatQueueUser
}

func NewChatQueueUser(engine *gin.Engine, chatQueueUser service.IChatQueueUser) {
	handler := &ChatQueueUser{
		chatQueueUser: chatQueueUser,
	}
	Group := engine.Group("bss-message/v1/chat-queue-user")
	{
		Group.POST("", handler.InsertChatQueueUser)
		Group.PUT(":id", handler.UpdateChatQueueUserById)
	}
}

func (h *ChatQueueUser) InsertChatQueueUser(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ChatQueueUserRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	data.Source = res.Data.Source

	log.Info("insert chat queue user payload -> ", data)

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err := h.chatQueueUser.InsertChatQueueUser(c, res.Data, data)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (h *ChatQueueUser) UpdateChatQueueUserById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ChatQueueUserRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat queue user payload -> ", data)

	result, err := h.chatQueueUser.UpdateChatQueueUserById(c, res.Data, data)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(result))
}
