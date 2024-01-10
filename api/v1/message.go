package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/api"
	"github.com/tel4vn/fins-microservices/service"
)

type Message struct {
	messageService service.IMessage
}

func NewMessage(r *gin.Engine, messageService service.IMessage) {
	handler := &Message{
		messageService: messageService,
	}

	Group := r.Group("chat/v1/message")
	{
		Group.POST("send", api.MoveTokenToHeader(), handler.SendMessage)
	}
}

func (m *Message) SendMessage(c *gin.Context) {

}
