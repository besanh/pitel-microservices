package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/service"
)

type ChatConnectionQueue struct {
	chatConnectionQueueService service.IChatConnectionQueue
}

func NewChatConnectionQueue(engine *gin.Engine, chatConnectionQueueService service.IChatConnectionQueue) {
	handler := &ChatConnectionQueue{
		chatConnectionQueueService: chatConnectionQueueService,
	}

	Group := engine.Group("bss-message/v1/chat-connection-queue")
	{
		Group.GET(":id", handler.GetChatConnectionQueueById)
	}
}

func (handler *ChatConnectionQueue) GetChatConnectionQueueById(c *gin.Context) {
	res := api.AuthMiddleware(c)

	id := c.Param("id")
	if len(id) < 1 {
		c.JSON(response.BadRequestMsg("id is required"))
		return
	}

	chatConnectionApp, err := handler.chatConnectionQueueService.GetChatConnectionQueueById(c, res.Data, id)
	if err != nil {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	c.JSON(response.OK(map[string]any{
		"data": chatConnectionApp,
	}))
}
