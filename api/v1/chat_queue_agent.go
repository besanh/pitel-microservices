package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
)

type ChatQueueAgent struct {
	chatQueueAgent service.IChatQueueAgent
}

func NewChatQueueAgent(engine *gin.Engine, chatQueueAgent service.IChatQueueAgent) {
	handler := &ChatQueueAgent{
		chatQueueAgent: chatQueueAgent,
	}
	Group := engine.Group("bss-message/v1/chat-queue-agent")
	{
		Group.POST("", handler.InsertChatQueueAgent)
		Group.PUT(":id", handler.UpdateChatQueueAgentById)
	}
}

func (h *ChatQueueAgent) InsertChatQueueAgent(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ChatQueueAgentRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}
	data.Source = res.Data.Source

	log.Info("insert chat queue agent payload -> ", data)

	if err := data.Validate(); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	err := h.chatQueueAgent.InsertChatQueueAgent(c, res.Data, data)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OKResponse())
}

func (h *ChatQueueAgent) UpdateChatQueueAgentById(c *gin.Context) {
	res := api.AuthMiddleware(c)
	if res == nil {
		c.JSON(response.ServiceUnavailableMsg("token is invalid"))
		return
	}

	var data model.ChatQueueAgentRequest
	if err := c.ShouldBind(&data); err != nil {
		log.Error(err)
		c.JSON(response.BadRequestMsg(err.Error()))
		return
	}

	log.Info("update chat queue agent payload -> ", data)

	result, err := h.chatQueueAgent.UpdateChatQueueAgentById(c, res.Data, data)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(response.OK(result))
}
